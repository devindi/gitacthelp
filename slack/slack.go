package slack

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type SlackClient struct {
	HookUrl string
}

func (slack SlackClient) SendMessage(msg string) error {
	message := map[string]interface{}{
		"text": msg,
	}

	bytesRepresentation, err := json.Marshal(message)
	if err != nil {
		return err
	}

	_, err = http.Post(slack.HookUrl, "application/json", bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		return err
	}
	return nil
}

func (slack SlackClient) Write(p []byte) (n int, err error) {
	message := string(p)
	err = slack.SendMessage(message)
	if err != nil {
		return 0, err
	} else {
		return len(p), nil
	}
}
