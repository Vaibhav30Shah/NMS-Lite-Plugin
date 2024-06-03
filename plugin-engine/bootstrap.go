package main

import (
	"NMS-Lite/server"
	"NMS-Lite/utils"
	"fmt"
	"sync"
)

func main() {
	logger := utils.NewLogger("bootstrap", "Bootstrap")

	logger.Info("Starting Plugin Engine")

	wg := new(sync.WaitGroup)

	wg.Add(1)

	context, recvSocket, sendSocket, err := server.Init()

	if err != nil {

		logger.Error(fmt.Sprintf("Error initializing the ZMQ sockets: %v", err.Error()))
	}

	defer context.Term()

	defer recvSocket.Close()

	defer sendSocket.Close()

	result := make(chan []map[string]interface{})

	go server.ReceiveRequests(recvSocket, result)

	go server.SendResults(sendSocket, result)

	wg.Wait()
}
