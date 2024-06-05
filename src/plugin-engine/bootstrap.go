package main

import (
	"NMS-Lite/consts"
	"NMS-Lite/server"
	"NMS-Lite/snmp"
	"NMS-Lite/utils"
	"fmt"
	"runtime/debug"
	"sync"
)

func main() {
	logger := utils.NewLogger("bootstrap", "Bootstrap")

	logger.Info("Starting Plugin Engine")

	err := server.Init()

	defer server.Close()

	if err != nil {

		logger.Error(fmt.Sprintf("Error initializing the ZMQ sockets: %v", err.Error()))
	}

	channel := make(chan []map[string]interface{})

	go server.ReceiveRequests(channel)

	results := make([]map[string]interface{}, 0)

	go server.SendResults(channel)

	for {

		contexts := <-channel

		wg := sync.WaitGroup{}

		for _, context := range contexts {

			wg.Add(1)

			go func(context map[string]interface{}) {

				defer wg.Done()

				defer func() {

					if r := recover(); r != nil {

						logger.Error(fmt.Sprintf("Panic occurred: %v\n%s", r, debug.Stack()))
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
						logger.Error("Unsupported plugin type!")
					}

					context[consts.Result] = result

					results = append(results, context)
				}

				channel <- contexts

				return
			}(context)
		}

		wg.Wait()

		logger.Trace(fmt.Sprintf("Results: %v", results))
	}
}
