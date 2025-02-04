package snmp

import (
	"NMS-Lite/consts"
	"NMS-Lite/snmpclient"
	"NMS-Lite/utils"
	"encoding/json"
	"fmt"
	"strings"
)

func Discover(context map[string]interface{}) map[string]interface{} {

	logger := utils.NewLogger("snmp", "Discover") //logger

	var errors []map[string]interface{}

	var systemName string

	result := map[string]interface{}{
		consts.SystemName: systemName,
	}

	var newContext = make(map[string]interface{})

	credentialsArray := context[consts.SnmpCredential].([]interface{})

	for _, credentials := range credentialsArray {

		credentialMap := credentials.(map[string]interface{})

		newContext[consts.SnmpCommunity] = credentialMap[consts.SnmpCommunity].(string)

		newContext[consts.SnmpVersion] = credentialMap[consts.SnmpVersion].(string)

		newContext[consts.SnmpPort] = uint16(context[consts.SnmpPort].(float64))

		newContext[consts.ObjectIp] = context[consts.ObjectIp]
	}

	client, err := snmpclient.Init(newContext)

	defer client.Close()

	if err != nil {

		errors = append(errors, map[string]interface{}{

			consts.ErrorName: "Error Initializing snmp client",

			consts.ErrorMessage: err.Error(),
		})

		logger.Error(fmt.Sprintf("Error initializing SNMP client: %s", err.Error()))

		context[consts.Error] = errors

		context[consts.Status] = consts.FailedStatus

		return result
	}

	discoveryResult, err := client.Get([]string{consts.ScalerOids[consts.SystemName]})

	if err != nil {

		errors = append(errors, map[string]interface{}{

			consts.ErrorName: "Error fetching system name",

			consts.ErrorMessage: err.Error(),
		})

		logger.Error(fmt.Sprintf("Error fetching system name: %s", err.Error()))

		context[consts.Error] = errors

		context[consts.Status] = consts.FailedStatus

		return result
	}

	if len(result) == 0 {

		errors = append(errors, map[string]interface{}{

			consts.ErrorName: "No response received from SNMP device",

			consts.ErrorMessage: err.Error(),
		})

		logger.Error("No response received from SNMP device")

		context[consts.Error] = errors

		context[consts.Status] = consts.FailedStatus

		return result
	}

	switch value := discoveryResult[0].Value.(type) {

	case string:
		systemName = strings.TrimPrefix(strings.TrimSuffix(value, `"`), `"`)

	case []byte:
		systemName = string(value)

	default:
		logger.Error("Unsupported value type for system name")

		context[consts.Error] = errors

		context[consts.Status] = consts.FailedStatus

		return result
	}

	result[consts.SystemName] = systemName

	jsonResult, err := json.Marshal(result)

	resultLogger := utils.NewLogger("snmp", "ResultData")

	resultLogger.Info(string(jsonResult))

	context[consts.Error] = "[]"

	context[consts.Status] = consts.SuccessStatus

	return result
}
