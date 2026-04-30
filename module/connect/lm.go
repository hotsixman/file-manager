package connect

import (
	"errors"
	"file-manager/module/LM"
	"file-manager/module/types"

	ndj "github.com/hotsixman/ndj-flow-client-go"
	"github.com/mitchellh/mapstructure"
)

func lock(message ndj.Message, lm *LM.LockManager, client *ndj.Client) {
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
		sendLMError(message, client, err)
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

func unlock(message ndj.Message, lm *LM.LockManager, client *ndj.Client) {
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
		sendLMError(message, client, err)
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

func check(message ndj.Message, lm *LM.LockManager, client *ndj.Client) {
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

	result, err := lm.Check(data.Path)
	if err != nil {
		sendLMError(message, client, err)
		return
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

func movePre(message ndj.Message, lm *LM.LockManager, client *ndj.Client) {
	var data types.TwoPathCommand
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

	err = lm.MovePre(data.Src, data.Dest)
	if err != nil {
		sendLMError(message, client, err)
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

func movePost(message ndj.Message, lm *LM.LockManager, client *ndj.Client) {
	var data types.TwoPathCommand
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

	lm.MovePost(data.Src, data.Dest)
	client.Send(ndj.SendHeader{
		To: message.Header.From,
		ID: message.Header.ID,
		Metadata: map[string]string{
			"status": "200",
		},
	}, nil)
}

func copyPre(message ndj.Message, lm *LM.LockManager, client *ndj.Client) {
	var data types.TwoPathCommand
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

	err = lm.CopyPre(data.Src, data.Dest)
	if err != nil {
		sendLMError(message, client, err)
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

func copyPost(message ndj.Message, lm *LM.LockManager, client *ndj.Client) {
	var data types.TwoPathCommand
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

	lm.CopyPost(data.Src, data.Dest)
	client.Send(ndj.SendHeader{
		To: message.Header.From,
		ID: message.Header.ID,
		Metadata: map[string]string{
			"status": "200",
		},
	}, nil)
}

func renamePre(message ndj.Message, lm *LM.LockManager, client *ndj.Client) {
	var data types.TwoPathCommand
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

	err = lm.RenamePre(data.Src, data.Dest)
	if err != nil {
		sendLMError(message, client, err)
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

func renamePost(message ndj.Message, lm *LM.LockManager, client *ndj.Client) {
	var data types.TwoPathCommand
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

	lm.RenamePost(data.Src, data.Dest)
	client.Send(ndj.SendHeader{
		To: message.Header.From,
		ID: message.Header.ID,
		Metadata: map[string]string{
			"status": "200",
		},
	}, nil)
}

func deletePre(message ndj.Message, lm *LM.LockManager, client *ndj.Client) {
	var data types.OnePathCommand
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

	err = lm.DeletePre(data.Path)
	if err != nil {
		sendLMError(message, client, err)
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

func deletePost(message ndj.Message, lm *LM.LockManager, client *ndj.Client) {
	var data types.OnePathCommand
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

	lm.DeletePost(data.Path)
	client.Send(ndj.SendHeader{
		To: message.Header.From,
		ID: message.Header.ID,
		Metadata: map[string]string{
			"status": "200",
		},
	}, nil)
}

func readPre(message ndj.Message, lm *LM.LockManager, client *ndj.Client) {
	var data types.OnePathCommand
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

	err = lm.ReadPre(data.Path)
	if err != nil {
		sendLMError(message, client, err)
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

func readPost(message ndj.Message, lm *LM.LockManager, client *ndj.Client) {
	var data types.OnePathCommand
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

	lm.ReadPost(data.Path)
	client.Send(ndj.SendHeader{
		To: message.Header.From,
		ID: message.Header.ID,
		Metadata: map[string]string{
			"status": "200",
		},
	}, nil)
}

func uploadPre(message ndj.Message, lm *LM.LockManager, client *ndj.Client) {
	var data types.OnePathCommand
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

	err = lm.UploadPre(data.Path)
	if err != nil {
		sendLMError(message, client, err)
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

func uploadPost(message ndj.Message, lm *LM.LockManager, client *ndj.Client) {
	var data types.OnePathCommand
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

	lm.UploadPost(data.Path)
	client.Send(ndj.SendHeader{
		To: message.Header.From,
		ID: message.Header.ID,
		Metadata: map[string]string{
			"status": "200",
		},
	}, nil)
}

func downloadPre(message ndj.Message, lm *LM.LockManager, client *ndj.Client) {
	var data types.OnePathCommand
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

	err = lm.DownloadPre(data.Path)
	if err != nil {
		sendLMError(message, client, err)
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

func downloadPost(message ndj.Message, lm *LM.LockManager, client *ndj.Client) {
	var data types.OnePathCommand
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

	lm.DownloadPost(data.Path)
	client.Send(ndj.SendHeader{
		To: message.Header.From,
		ID: message.Header.ID,
		Metadata: map[string]string{
			"status": "200",
		},
	}, nil)
}

func sendLMError(message ndj.Message, client *ndj.Client, err error) {
	status := "500"
	if exception, ok := errors.AsType[*types.LockManagerException](err); ok {
		switch exception.Code {
		case types.ExceptionStatus.PATH_IS_NOT_ABS, types.ExceptionStatus.INVALID_STATUS:
			status = "400"
		case types.ExceptionStatus.ALREADY_LOCKED:
			status = "429"
		case types.ExceptionStatus.ANCESTOR_LOCKED, types.ExceptionStatus.DECENDENT_LOCKED:
			status = "422"
		case types.ExceptionStatus.ALREADY_UNLOCKED:
			status = "204"
		}
	}
	client.Send(ndj.SendHeader{
		To: message.Header.From,
		ID: message.Header.ID,
		Metadata: map[string]string{
			"status": status,
		},
	}, nil)
}
