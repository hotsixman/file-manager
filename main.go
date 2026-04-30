package main

import (
	"file-manager/module/LM"
	"file-manager/module/connect"
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
	client := NewClient(lm)

	client.Connect()
	fmt.Println("Connected.")

	select {}
}

func NewClient(lm *LM.LockManager) *ndj.UDSClient {
	client := ndj.NewUDSClient(ndj.UDSClientOption{
		Path: os.Getenv("UDS_PATH"),
		ClientOption: ndj.ClientOption{
			Name: "foo",
			Key:  "foo",
		},
	})
	onReceive := connect.NewReceiver(lm, client.Client)
	client.Client.OnReceive = onReceive
	client.Client.OnError = func(err error) {
		fmt.Println("error:", err)
	}
	return client
}
