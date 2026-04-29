package main

import (
	"file-lock-manager/module/LM"
	"file-lock-manager/module/socket"
	"fmt"
	"os"

	ndj "github.com/hotsixman/ndj-flow-client-go"
	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()
}

func main() {
	lm := LM.NewLockManager()

	client := ndj.NewUDSClient(ndj.UDSClientOption{
		Path: os.Getenv("UDS_PATH"),
		ClientOption: ndj.ClientOption{
			Name: "foo",
			Key:  "foo",
		},
	})

	onReceive := socket.NewReceiver(lm, client.Client)
	client.Client.OnReceive = onReceive
	client.Client.OnError = func(err error) {
		fmt.Println("error:", err)
	}
	client.Connect()
	fmt.Println("Connect")

	select {}
}
