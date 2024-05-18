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

	decodedContext, err := base64.StdEncoding.DecodeString(os.Args[1]) //getting the encoded json object

	if err != nil {

		Logger.Error(fmt.Sprintf("Error in decoding context %s", err.Error()))

		return
	}

	contexts := make([]map[string]interface{}, 0)

	if err = json.Unmarshal(decodedContext, &contexts); err != nil {

		Logger.Error(fmt.Sprintf("Conversion of json to map failed: %s", err.Error()))

		return
	}

	// Create a slice to store results
	results := make([]map[string]interface{}, 0)

	for _, context := range contexts {

		Logger.Debug(fmt.Sprintf("Context: %s\n", context))

		wg.Add(1)

		go func(context map[string]interface{}) {

			defer wg.Done()

			defer func() {

				if r := recover(); r != nil {

					Logger.Error(fmt.Sprintf("Panic occurred: %v\n%s", r, debug.Stack()))

					return
				}
			}()

			if pluginName, ok := context[consts.PluginName].(string); ok {

				var result map[string]interface{}

				switch pluginName {

				case consts.Discover:
					result = snmp.Discover(context)

				case consts.Collect:
					result = snmp.Collect(context)

				default:
					fmt.Println("Unsupported plugin type!")
				}

				context[consts.Result] = result

				results = append(results, context)
			}
		}(context)

		wg.Wait()
	}

	jsonResult, err := json.Marshal(results)

	if err != nil {

		Logger.Error(fmt.Sprintf("Error marshaling collection result: %s", err.Error()))

	}

	fmt.Println(base64.StdEncoding.EncodeToString(jsonResult))

	DataLogger := utils.NewLogger("data", "BootstrapResult")

	DataLogger.Info(string(jsonResult))
}
