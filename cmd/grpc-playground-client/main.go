package main

import (
	"context"

	"github.com/SKF/go-utility/log"
	"github.com/lonnblad/grpc-playground/cmd/grpc-playground-client/api"
	"github.com/lonnblad/grpc-playground/cmd/grpc-playground-client/conn"
)

func main() {
	host := "localhost"
	port := "50051"

	client := api.Client{}
	err := conn.Connect(&client, host, port)
	if err != nil {
		log.WithError(err).Infof("Failed connecet to host")
		return
	}

	defer client.Close()

	// go func() {
	// 	n := make([]byte, 30)
	// 	for range n {
	// 		if err := client.DeepPing(); err != nil {
	// 			log.WithError(err).Infof("Failed to Deep Ping")
	// 		}
	// 		time.Sleep(2 * time.Second)
	// 	}
	// }()

	// if err = client.InvokePanic(); err == nil {
	// 	log.WithError(err).Infof("Failed to Invoke Panic")
	// }

	recChannel := make(chan string)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := client.GetReceiveStream(ctx, recChannel); err != nil {
			log.WithError(err).Infof("Failed to GetStream")
			return
		}
	}()

	// sendChannel := make(chan string)
	// go func() {
	// 	if err := client.GetSendStream(sendChannel); err != nil {
	// 		log.WithError(err).Infof("Failed to GetStream")
	// 	}
	// }()

	for msg := range recChannel {
		// select {
		// case msg := <-recChannel:
		log.WithField("msg", msg).
			Info("Got message from stream")
		// case <-time.After(time.Second):
		// 	sendChannel <- "Ping"
		// }
		// cancel()
	}
}
