package snmp

import (
	"NMS-Lite/consts"
	"NMS-Lite/snmpclient"
	"NMS-Lite/utils"
	"encoding/json"
	"fmt"
	"strings"
)

func Discover(context map[string]interface{}, errors *[]map[string]interface{}) {
	Logger := utils.NewLogger("snmp", "Discover")

	ip, port, community, timeOut, version := utils.ValidateContext(context)

	client, err := snmpclient.Init(ip, community, uint16(port), version, int(timeOut))
	if err != nil {
		Logger.Error(fmt.Sprintf("Error initializing SNMP client: %s", err.Error()))
		return
	}
	defer client.Close()

	result, err := client.Get([]string{consts.ScalerOids["system.name"]})
	if err != nil {
		Logger.Error(fmt.Sprintf("Error fetching system name: %s", err.Error()))
		return
	}

	if len(result) == 0 {
		Logger.Error("No response received from SNMP device")
		return
	}

	var systemName string
	switch value := result[0].Value.(type) {
	case string:
		systemName = strings.TrimPrefix(strings.TrimSuffix(value, `"`), `"`)
	case []byte:
		systemName = string(value)
	default:
		Logger.Error("Unsupported value type for system name")
		return
	}

	discoveryResult := map[string]interface{}{
		"system_name": systemName,
	}

	jsonResult, err := json.Marshal(discoveryResult)
	if err != nil {
		Logger.Error(fmt.Sprintf("Error marshaling discovery result: %s", err.Error()))
		return
	}

	fmt.Println(string(jsonResult))
	Logger2 := utils.NewLogger("snmp", "ResultData")
	Logger2.Info(string(jsonResult))
}
