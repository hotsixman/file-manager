package socket

import (
	"errors"
	"file-lock-manager/module/LM"
	"file-lock-manager/module/types"

	ndj "github.com/hotsixman/ndj-flow-client-go"
	"github.com/mitchellh/mapstructure"
)

func NewReceiver(lm *LM.LockManager, client *ndj.Client) func(message ndj.Message) {
	onReceive := func(message ndj.Message) {
		command := message.Header.Metadata["command"]

		switch command {
		case "lock":
			{
				var data types.LockUnlockCommand
				err := mapstructure.Decode(<-message.Body, &data)
				if err != nil {
					client.Send(ndj.SendHeader{
						To: message.Header.From,
						ID: message.Header.ID,
						Metadata: map[string]string{
							"status": "400",
						},
					}, nil)
					return
				}

				resultStatus, err := lm.Lock(data.Path, data.Status)
				if err != nil {
					status := "500"
					if exception, ok := errors.AsType[*types.LockManagerException](err); ok {
						switch exception.Code {
						case types.ExceptionStatus.PATH_IS_NOT_ABS, types.ExceptionStatus.INVALID_STATUS:
							status = "400"
						case types.ExceptionStatus.ALREADY_LOCKED:
							status = "429"
						case types.ExceptionStatus.ANCESTOR_LOCKED, types.ExceptionStatus.DECENDENT_LOCKED:
							status = "422"
						}
					}
					client.Send(ndj.SendHeader{
						To: message.Header.From,
						ID: message.Header.ID,
						Metadata: map[string]string{
							"status": status,
						},
					}, nil)
					return
				}

				body := make([]any, 0)
				sendData := map[string]int{
					"status": resultStatus,
				}
				body = append(body, sendData)
				client.Send(ndj.SendHeader{
					To: message.Header.From,
					ID: message.Header.ID,
					Metadata: map[string]string{
						"status": "200",
					},
				}, body)
			}
		case "unlock":
			{
				var data types.LockUnlockCommand
				err := mapstructure.Decode(<-message.Body, &data)
				if err != nil {
					client.Send(ndj.SendHeader{
						To: message.Header.From,
						ID: message.Header.ID,
						Metadata: map[string]string{
							"status": "400",
						},
					}, nil)
					return
				}

				err = lm.Unlock(data.Path)
				if err != nil {
					status := "500"
					if exception, ok := errors.AsType[*types.LockManagerException](err); ok {
						switch exception.Code {
						case types.ExceptionStatus.PATH_IS_NOT_ABS:
							status = "400"
						case types.ExceptionStatus.ALREADY_UNLOCKED:
							status = "204"
						case types.ExceptionStatus.ANCESTOR_LOCKED, types.ExceptionStatus.DECENDENT_LOCKED:
							status = "422"
						}
					}
					client.Send(ndj.SendHeader{
						To: message.Header.From,
						ID: message.Header.ID,
						Metadata: map[string]string{
							"status": status,
						},
					}, nil)
					return
				}

				client.Send(ndj.SendHeader{
					To: message.Header.From,
					ID: message.Header.ID,
					Metadata: map[string]string{
						"status": "200",
					},
				}, nil)
			}
		case "check":
			{
				var data types.CheckCommand
				err := mapstructure.Decode(<-message.Body, &data)
				if err != nil {
					client.Send(ndj.SendHeader{
						To: message.Header.From,
						ID: message.Header.ID,
						Metadata: map[string]string{
							"status": "400",
						},
					}, nil)
					return
				}

				result, err := lm.CheckLocked(data.Path)
				if err != nil {
					status := "500"
					if exception, ok := errors.AsType[*types.LockManagerException](err); ok {
						switch exception.Code {
						case types.ExceptionStatus.PATH_IS_NOT_ABS:
							status = "400"
						}
						client.Send(ndj.SendHeader{
							To: message.Header.From,
							ID: message.Header.ID,
							Metadata: map[string]string{
								"status": status,
							},
						}, nil)
						return
					}
				}

				body := make([]any, 0)
				body = append(body, result)
				client.Send(ndj.SendHeader{
					To: message.Header.From,
					ID: message.Header.ID,
					Metadata: map[string]string{
						"status": "200",
					},
				}, body)
			}
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
