package snmp

import (
	"NMS-Lite/consts"
	"NMS-Lite/snmpclient"
	"NMS-Lite/utils"
	"encoding/json"
	"fmt"
	"strings"
)

func Discover(context map[string]interface{}, errors *[]map[string]interface{}) (map[string]interface{}, *[]map[string]interface{}) {

	Logger := utils.NewLogger("snmp", "Discover")

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

	result, err := client.Get([]string{consts.ScalerOids[consts.SystemName]})

	if err != nil {

		*errors = append(*errors, map[string]interface{}{

			consts.ErrorName: "Error fetching system name",

			consts.ErrorMessage: err.Error(),
		})

		Logger.Error(fmt.Sprintf("Error fetching system name: %s", err.Error()))

		return nil, errors
	}

	if len(result) == 0 {

		*errors = append(*errors, map[string]interface{}{

			consts.ErrorName: "No response received from SNMP device",

			consts.ErrorMessage: err.Error(),
		})

		Logger.Error("No response received from SNMP device")

		return nil, errors
	}

	var systemName string

	switch value := result[0].Value.(type) {

	case string:
		systemName = strings.TrimPrefix(strings.TrimSuffix(value, `"`), `"`)

	case []byte:
		systemName = string(value)

	default:
		Logger.Error("Unsupported value type for system name")

		return nil, nil
	}

	discoveryResult := map[string]interface{}{

		"system_name": systemName,
	}

	jsonResult, err := json.Marshal(discoveryResult)

	Logger2 := utils.NewLogger("snmp", "ResultData")

	Logger2.Info(string(jsonResult))

	return discoveryResult, nil
}
