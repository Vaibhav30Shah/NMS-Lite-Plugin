package snmp

import (
	"NMS-Lite/consts"
	"NMS-Lite/snmpclient"
	"NMS-Lite/utils"
	"encoding/json"
	"fmt"
)

func Collect(context map[string]interface{}, errors *[]map[string]interface{}) (map[string]interface{}, *[]map[string]interface{}) {

	Logger := utils.NewLogger("snmp", "Collect")

	ip, port, community, timeOut, version := utils.ValidateContext(context)

	client, err := snmpclient.Init(ip, community, uint16(port), version, int(timeOut))

	if err != nil {

		*errors = append(*errors, map[string]interface{}{

			consts.ErrorName: "Error Initializing snmp client",

			consts.ErrorMessage: err.Error(),
		})

		Logger.Error(fmt.Sprintf("Error initializing SNMP client: %s", err.Error()))

		return nil, errors
	}

	defer client.Close()

	collectionResult, err := client.Walk(consts.TabularOids)

	if err != nil {

		*errors = append(*errors, map[string]interface{}{

			consts.ErrorName: "Error fetching tabular oids",

			consts.ErrorMessage: err.Error(),
		})

		Logger.Error(fmt.Sprintf("Error fetching tabular OIDs: %s", err.Error()))

		return nil, errors
	}

	result := map[string]interface{}{
		"interface": make([]map[string]interface{}, 0),
	}

	for oidName, oidResult := range collectionResult {

		oidResultMap, ok := oidResult.(map[string]interface{})

		if !ok {

			*errors = append(*errors, map[string]interface{}{

				consts.ErrorName: "Error asserting oid Type",

				consts.ErrorMessage: err.Error(),
			})

			Logger.Error(fmt.Sprintf("Error asserting type for OID %s", oidName))

			continue
		}

		result[oidName] = oidResultMap
	}

	jsonResult, err := json.Marshal(result)

	if err != nil {

		*errors = append(*errors, map[string]interface{}{

			consts.ErrorName: "Error marshaling result",

			consts.ErrorMessage: err.Error(),
		})

		Logger.Error(fmt.Sprintf("Error marshaling collection result: %s", err.Error()))
	}

	Logger2 := utils.NewLogger("snmp", "ResultData")

	Logger2.Info(string(jsonResult))

	return result, errors
}
