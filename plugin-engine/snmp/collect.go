package snmp

import (
	"NMS-Lite/consts"
	"NMS-Lite/snmpclient"
	"NMS-Lite/utils"
	"encoding/json"
	"fmt"
	"strconv"
)

func Collect(context map[string]interface{}) map[string]interface{} {

	logger := utils.NewLogger("snmp", "Collect")

	var errors []map[string]interface{}

	result := map[string]interface{}{}

	client, err := snmpclient.Init(context)

	defer client.Close()

	if err != nil {

		errors = append(errors, map[string]interface{}{

			consts.ErrorName: "Error Initializing snmp client",

			consts.ErrorMessage: err.Error(),
		})

		logger.Error(fmt.Sprintf("Error initializing SNMP client: %s", err.Error()))

		result[consts.Error] = errors

		result[consts.Status] = consts.FailedStatus

		return result
	}

	collectionResult, err := client.Walk(consts.TabularOids) //result

	if err != nil {

		errors = append(errors, map[string]interface{}{

			consts.ErrorName: "Error fetching interface details",

			consts.ErrorMessage: err.Error(),
		})

		logger.Error(fmt.Sprintf("Error fetching tabular OIDs: %s", err.Error()))

		result[consts.Error] = errors

		result[consts.Status] = consts.FailedStatus

		return result
	}

	for oidName, oidResult := range collectionResult {

		oidResultMap, ok := oidResult.(map[string]interface{})

		if !ok {

			errors = append(errors, map[string]interface{}{

				consts.ErrorName: "Error asserting oid Type",

				consts.ErrorMessage: err.Error(),
			})

			logger.Error(fmt.Sprintf("Error asserting type for OID %v", oidName))

			result[consts.Error] = errors

			result[consts.Status] = consts.FailedStatus

			continue
		}

		result[strconv.Itoa(oidName)] = oidResultMap
	}

	jsonResult, err := json.Marshal(result)

	if err != nil {

		errors = append(errors, map[string]interface{}{

			consts.ErrorName: "Error marshaling result",

			consts.ErrorMessage: err.Error(),
		})

		logger.Error(fmt.Sprintf("Error marshaling collection result: %s", err.Error()))

		result[consts.Error] = errors

		result[consts.Status] = consts.FailedStatus

		return result
	}

	resultLogger := utils.NewLogger("snmp", "ResultData")

	resultLogger.Info(string(jsonResult))

	result[consts.Error] = "[]"

	result[consts.Status] = consts.SuccessStatus

	return result
}
