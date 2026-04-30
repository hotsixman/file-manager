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
		case "movePre":
			movePre(message, lm, client)
		case "movePost":
			movePost(message, lm, client)
		case "copyPre":
			copyPre(message, lm, client)
		case "copyPost":
			copyPost(message, lm, client)
		case "renamePre":
			renamePre(message, lm, client)
		case "renamePost":
			renamePost(message, lm, client)
		case "deletePre":
			deletePre(message, lm, client)
		case "deletePost":
			deletePost(message, lm, client)
		case "readPre":
			readPre(message, lm, client)
		case "readPost":
			readPost(message, lm, client)
		case "uploadPre":
			uploadPre(message, lm, client)
		case "uploadPost":
			uploadPost(message, lm, client)
		case "downloadPre":
			downloadPre(message, lm, client)
		case "downloadPost":
			downloadPost(message, lm, client)
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
