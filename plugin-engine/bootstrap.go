package main

import (
	"NMS-Lite/consts"
	"NMS-Lite/snmp"
	"NMS-Lite/utils"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/pebbe/zmq4"
	"runtime/debug"
	"sync"
)

func main() {

	wg := new(sync.WaitGroup)

	logger := utils.NewLogger("bootstrap", "Bootstrap")

	logger.Info("Starting Plugin Engine")

	context, err := zmq4.NewContext()

	if err != nil {

		logger.Error(fmt.Sprintf("Error creating ZMQ context: %v", err.Error()))

		return
	}

	defer context.Term()

	recvSocket, err := context.NewSocket(zmq4.PULL)

	sendSocket, err := context.NewSocket(zmq4.PUSH)

	if err != nil {

		logger.Error(fmt.Sprintf("Error creating ZMQ recvSocket: %v", err.Error()))

		return
	}

	defer recvSocket.Close()

	defer sendSocket.Close()

	err = recvSocket.Connect("tcp://*:5555")

	err = sendSocket.Connect("tcp://*:5556")

	if err != nil {

		logger.Error(fmt.Sprintf("Error binding ZMQ recvSocket: %v", err.Error()))

		return
	}

	for {
		recdContext, err := recvSocket.Recv(0)

		decodedContext, err := base64.StdEncoding.DecodeString(recdContext)

		if err != nil {

			logger.Error(fmt.Sprintf("Error receiving request: %v", err.Error()))

			continue
		}

		if err != nil {

			logger.Error(fmt.Sprintf("Error in decoding context %s", err.Error()))

			return
		}

		contexts := make([]map[string]interface{}, 0)

		logger.Info(fmt.Sprintf("Context: %v", decodedContext))

		if err = json.Unmarshal([]byte(decodedContext), &contexts); err != nil {

			logger.Error(fmt.Sprintf("Conversion of json to map failed: %s", err.Error()))

			return
		}

		// Create a slice to store results
		results := make([]map[string]interface{}, 0)

		for _, context := range contexts {

			logger.Debug(fmt.Sprintf("Context: %s\n", context))

			wg.Add(1)

			go func(context map[string]interface{}) {

				defer wg.Done()

				defer func() {

					if r := recover(); r != nil {

						logger.Error(fmt.Sprintf("Panic occurred: %v\n%s", r, debug.Stack()))

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
		}

		wg.Wait()

		jsonResult, err := json.Marshal(results)

		if err != nil {

			logger.Error(fmt.Sprintf("Error marshaling collection result: %s", err.Error()))

		}

		encodedResult := base64.StdEncoding.EncodeToString(jsonResult)

		sendSocket.Send(encodedResult, zmq4.DONTWAIT)

		logger.Info("Result sent Via ZMQ.")

		//fmt.Println(base64.StdEncoding.EncodeToString(jsonResult))

		dataLogger := utils.NewLogger("data", "BootstrapResult")

		dataLogger.Info(string(jsonResult))
	}
}
