package NMS_Lite

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

func Bootstrap() {

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
				switch pluginName {

				case Discover:
					snmp.Discover(context, &errors)

				case Collect:
					snmp.Collect(context, &errors)

				default:
					fmt.Println("Unsupported plugin type!")
				}
			}
		}(context)

		wg.Wait()
	}
}
