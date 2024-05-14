package main

import (
	"NMS-Lite/snmp"
	"NMS-Lite/utils"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"runtime/debug"
	"sync"
)

const (
	Discover   = "Discover"
	PluginName = "plugin_type"
	Collect    = "Collect"
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

			if pluginName, ok := context[PluginName].(string); ok {
				var result map[string]interface{}
				switch pluginName {

				case Discover:
					result = snmp.Discover(context, &errors)

				case Collect:
					result = snmp.Collect(context, &errors)

				default:
					fmt.Println("Unsupported plugin type!")
				}
				if result != nil {
					context["result"] = result
				}

				jsonResult, err := json.Marshal(result)
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
