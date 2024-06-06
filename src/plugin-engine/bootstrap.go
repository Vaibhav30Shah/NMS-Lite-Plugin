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

var logger = utils.NewLogger("bootstrap", "Bootstrap")

func main() {

	logger.Info("Starting Plugin Engine")

	if err := server.Init(); err != nil {

		logger.Error(fmt.Sprintf("Error initializing the ZMQ sockets: %v", err.Error()))
	}

	defer server.Close()

	channel := make(chan []map[string]interface{})

	go server.ReceiveRequests(channel)

	go server.SendResults(channel)

	executeTasks(channel)

	select {}
}

func executeTasks(channel chan []map[string]interface{}) {

	defer func() {

		if r := recover(); r != nil {

			logger.Error(fmt.Sprintf("Panic occurred: %v\n%s", r, debug.Stack()))

			executeTasks(channel)
		}
	}()

	for {
		contexts := <-channel

		var wg sync.WaitGroup

		results := make([]map[string]interface{}, 0, len(contexts))

		for _, context := range contexts {

			wg.Add(1)

			if pluginName, ok := context[consts.PluginName]; ok {

				switch pluginName {

				case consts.Discover:

					go func(context map[string]interface{}) {

						defer wg.Done()

						defer func() {

							if r := recover(); r != nil {

								logger.Error(fmt.Sprintf("Panic occurred: %v\n%s", r, debug.Stack()))
							}
						}()

						result := snmp.Discover(context)

						context[consts.Result] = result

						results = append(results, context)

						return

					}(context)

				case consts.Collect:

					go func(context map[string]interface{}) {

						defer wg.Done()

						defer func() {

							if r := recover(); r != nil {

								logger.Error(fmt.Sprintf("Panic occurred: %v\n%s", r, debug.Stack()))
							}
						}()

						result := snmp.Collect(context)

						context[consts.Result] = result

						results = append(results, context)

						return

					}(context)

				default:
					logger.Error("Unsupported plugin type!")
				}
			}
		}

		wg.Wait()

		logger.Trace(fmt.Sprintf("Results: %v", results))

		channel <- contexts
	}
}
