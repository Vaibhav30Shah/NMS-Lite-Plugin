package snmp

import (
	"NMS-Lite/consts"
	"NMS-Lite/snmpclient"
	"NMS-Lite/utils"
	"encoding/json"
	"fmt"
)

func Collect(context map[string]interface{}, errors *[]map[string]interface{}) {
	Logger := utils.NewLogger("snmp", "Collect")

	ip, port, community, timeOut, version := utils.ValidateContext(context)

	client, err := snmpclient.Init(ip, community, uint16(port), version, int(timeOut))
	if err != nil {
		Logger.Error(fmt.Sprintf("Error initializing SNMP client: %s", err.Error()))
		return
	}
	defer client.Close()

	collectionResult, err := client.Walk(consts.TabularOids)
	if err != nil {
		Logger.Error(fmt.Sprintf("Error fetching tabular OIDs: %s", err.Error()))
		return
	}

	var flattenedResult []map[string]interface{}

	for oidName, oidResult := range collectionResult {
		oidResultMap, ok := oidResult.(map[string]interface{})
		if !ok {
			Logger.Error(fmt.Sprintf("Error asserting type for OID %s", oidName))
			continue
		}

		for interfaceIndex, interfaceData := range oidResultMap {
			newData := make(map[string]interface{})
			newData["interface.index"] = interfaceIndex
			newData[oidName] = interfaceData
			flattenedResult = append(flattenedResult, newData)
		}
	}

	jsonResult, err := json.Marshal(flattenedResult)
	if err != nil {
		Logger.Error(fmt.Sprintf("Error marshaling collection result: %s", err.Error()))
		return
	}

	fmt.Println("Result of collect successfully appended to result.json")
	Logger2 := utils.NewLogger("snmp", "ResultData")
	Logger2.Info(string(jsonResult))
}
