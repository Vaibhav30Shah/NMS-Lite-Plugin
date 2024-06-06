package server

import (
	"NMS-Lite/utils"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/pebbe/zmq4"
	"runtime/debug"
)

var logger = utils.NewLogger("bootstrap", "Bootstrap")

var context *zmq4.Context

var err error

var recvSocket *zmq4.Socket

var sendSocket *zmq4.Socket

func Init() error {

	if context, err = zmq4.NewContext(); err != nil {

		logger.Error(fmt.Sprintf("Error creating ZMQ context: %v", err.Error()))

		return err
	}

	if recvSocket, err = context.NewSocket(zmq4.PULL); err != nil {

		logger.Error(fmt.Sprintf("Error creating ZMQ recvSocket: %v", err.Error()))

		return err
	}

	if sendSocket, err = context.NewSocket(zmq4.PUSH); err != nil {

		logger.Error(fmt.Sprintf("Error creating ZMQ sendSocket: %v", err.Error()))

		return err
	}

	if err = recvSocket.Connect("tcp://localhost:5555"); err != nil {

		logger.Error(fmt.Sprintf("Error connecting ZMQ recvSocket: %v", err.Error()))

		return err
	}

	if err = sendSocket.Connect("tcp://localhost:5556"); err != nil {

		logger.Error(fmt.Sprintf("Error connecting ZMQ sendSocket: %v", err.Error()))

		return err
	}

	return nil
}

func ReceiveRequests(result chan<- []map[string]interface{}) {

	defer func() {

		if r := recover(); r != nil {

			logger.Error(fmt.Sprintf("Panic occurred: %v\n%s", r, debug.Stack()))

			go ReceiveRequests(result)
		}
	}()

	for {
		recdContext, err := recvSocket.Recv(0)

		logger.Trace(recdContext)

		decodedContext, err := base64.StdEncoding.DecodeString(recdContext)

		//decodedContext, err := base64.StdEncoding.DecodeString(os.Args[1])

		logger.Trace("received context:=" + string(decodedContext))

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

		result <- contexts
	}
}

func SendResults(result <-chan []map[string]interface{}) {

	logger := utils.NewLogger("server", "SendResults")

	defer func() {

		if r := recover(); r != nil {

			logger.Error(fmt.Sprintf("Panic occurred: %v\n%s", r, debug.Stack()))

			go SendResults(result)
		}
	}()

	for {
		fmt.Sprintf("%v", result)

		for results := range result {

			jsonResult, err := json.Marshal(results)

			logger.Info("json result: " + string(jsonResult))

			if err != nil {

				logger.Error(fmt.Sprintf("Error marshaling collection result: %s", err.Error()))

				continue
			}

			encodedResult := base64.StdEncoding.EncodeToString(jsonResult)

			sendSocket.Send(encodedResult, zmq4.DONTWAIT)

			logger.Info("Result sent Via ZMQ.")

			logger.Info(string(jsonResult))

			dataLogger := utils.NewLogger("data", "BootstrapResult")

			dataLogger.Info(string(jsonResult))
		}
	}
}

func Close() {

	defer func() {

		if r := recover(); r != nil {

			logger.Error(fmt.Sprintf("Panic occurred: %v\n%s", r, debug.Stack()))
		}
	}()

	if err := context.Term(); err != nil {

		logger.Error(fmt.Sprintf("Error closing ZMQ socket: %v", err.Error()))
	}

	if err := recvSocket.Close(); err != nil {

		logger.Error(fmt.Sprintf("Error closing ZMQ socket: %v", err.Error()))
	}

	if err := sendSocket.Close(); err != nil {

		logger.Error(fmt.Sprintf("Error closing ZMQ socket: %v", err.Error()))
	}
}
