package main

import (
	"NMS-Lite/consts"
	"NMS-Lite/snmp"
	"NMS-Lite/utils"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"runtime/debug"
	"sync"
)

func main() {

	wg := new(sync.WaitGroup)

	Logger := utils.NewLogger("bootstrap", "Bootstrap")

	Logger.Info("Starting Plugin Engine")

	decodedContext, err := base64.StdEncoding.DecodeString(os.Args[1])

	if err != nil {

		Logger.Error(fmt.Sprintf("Error in decoding context %s", err.Error()))

		return
	}

	contexts := make([]map[string]interface{}, 0)

	if err = json.Unmarshal(decodedContext, &contexts); err != nil {

		Logger.Error(fmt.Sprintf("Conversion of json to map failed: %s", err.Error()))

		return
	}

	//fmt.Println(contexts)

	//fmt.Println(string(decodedContext))

	for _, context := range contexts {
		Logger.Debug(fmt.Sprintf("Context: %s\n", context))

		wg.Add(1)

		go func(context map[string]interface{}) {

			defer wg.Done()

			errors := make([]map[string]interface{}, 0)

			defer func(context map[string]interface{}, contexts []map[string]interface{}) {

				if r := recover(); r != nil {

					Logger.Error(fmt.Sprintf("Panic occurred: %v\n%s", r, debug.Stack()))

					return
				}
			}(context, contexts)

			if pluginName, ok := context[consts.PluginName].(string); ok {

				var result map[string]interface{}

				var err1 *[]map[string]interface{}

				switch pluginName {

				case consts.Discover:
					result, err1 = snmp.Discover(context, &errors)

				case consts.Collect:
					result, err1 = snmp.Collect(context, &errors)

				default:

					fmt.Println("Unsupported plugin type!")
				}
				if result != nil {

					context[consts.Result] = result
				}

				if err1 != nil {

					context[consts.Error] = err1

					context[consts.Status] = "Failed"

				}
				if err1 == nil {

					context[consts.Error] = "[]"

					context[consts.Status] = "Success"
				}

				jsonResult, err := json.Marshal(context)

				if err != nil {

					Logger.Error(fmt.Sprintf("Error marshaling collection result: %s", err.Error()))

				}

				DataLogger := utils.NewLogger("data", "BootstrapResult")

				DataLogger.Info(string(jsonResult))
			}
		}(context)

		wg.Wait()
	}
}
