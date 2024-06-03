package server

import (
	"NMS-Lite/consts"
	"NMS-Lite/snmp"
	"NMS-Lite/utils"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/pebbe/zmq4"
	"os"
	"runtime/debug"
	"sync"
)

var logger = utils.NewLogger("bootstrap", "Bootstrap")

func Init() (context *zmq4.Context, recvSocket *zmq4.Socket, sendSocket *zmq4.Socket, err error) {

	if context, err = zmq4.NewContext(); err != nil {

		logger.Error(fmt.Sprintf("Error creating ZMQ context: %v", err.Error()))

		return
	}

	if recvSocket, err = context.NewSocket(zmq4.PULL); err != nil {

		logger.Error(fmt.Sprintf("Error creating ZMQ recvSocket: %v", err.Error()))

		return
	}

	if sendSocket, err = context.NewSocket(zmq4.PUSH); err != nil {

		logger.Error(fmt.Sprintf("Error creating ZMQ sendSocket: %v", err.Error()))

		return
	}

	if err = recvSocket.Connect("tcp://localhost:5555"); err != nil {

		logger.Error(fmt.Sprintf("Error connecting ZMQ recvSocket: %v", err.Error()))

		return
	}

	if err = sendSocket.Connect("tcp://localhost:5556"); err != nil {

		logger.Error(fmt.Sprintf("Error connecting ZMQ sendSocket: %v", err.Error()))

		return
	}

	return
}

func ReceiveRequests(recvSocket *zmq4.Socket, result chan<- []map[string]interface{}) {

	//recdContext, err := recvSocket.Recv(0)

	decodedContext, err := base64.StdEncoding.DecodeString(os.Args[1])

	if err != nil {

		logger.Error(fmt.Sprintf("Error receiving request: %v", err.Error()))

		return
	}

	contexts := make([]map[string]interface{}, 0)

	if err := json.Unmarshal([]byte(decodedContext), &contexts); err != nil {

		logger.Error(fmt.Sprintf("Conversion of json to map failed: %s", err.Error()))

		return
	}

	logger.Info(fmt.Sprintf("Context: %v", string(decodedContext)))

	results := make([]map[string]interface{}, 0)

	wg := new(sync.WaitGroup)

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

			return
		}(context)
	}

	wg.Wait()

	logger.Info(fmt.Sprintf("Results: %v", results))

	result <- results

	return
}
