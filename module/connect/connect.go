package connect

import (
	"file-manager/module/LM"

	ndj "github.com/hotsixman/ndj-flow-client-go"
)

func NewReceiver(lm *LM.LockManager, client *ndj.Client) func(message ndj.Message) {
	onReceive := func(message ndj.Message) {
		command := message.Header.Metadata["command"]

		switch command {
		case "lock":
			lock(message, lm, client)
		case "unlock":
			unlock(message, lm, client)
		case "check":
			check(message, lm, client)
		default:
			{
				client.Send(ndj.SendHeader{
					To: message.Header.From,
					ID: message.Header.ID,
					Metadata: map[string]string{
						"status": "404",
					},
				}, nil)
			}
		}
	}
	return onReceive
}
