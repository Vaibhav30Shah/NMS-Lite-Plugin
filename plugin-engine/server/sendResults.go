package server

import (
	"NMS-Lite/utils"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/pebbe/zmq4"
)

func SendResults(sendSocket *zmq4.Socket, result <-chan []map[string]interface{}) {

	logger := utils.NewLogger("engine", "SendResults")

	for results := range result {

		jsonResult, err := json.Marshal(results)

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
	return
}
