package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type SlackNotifier struct {
	Webhook string
	Channel string
}

func ParseSlackNotifier(data map[string]interface{}) SlackNotifier {
	notifier := SlackNotifier{}

	if data["webhook"] != nil {
		notifier.Webhook = data["webhook"].(string)
	}

	if data["channel"] != nil {
		notifier.Channel = data["channel"].(string)
	}

	return notifier
}

func (notifier SlackNotifier) Perform(message string) error {
	payload := map[string]string{
		"text":     message,
		"username": "service-health",
	}

	if notifier.Channel != "" {
		payload["channel"] = notifier.Channel
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	reader := bytes.NewReader(data)
	resp, err := http.Post(notifier.Webhook, "application/json", reader)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("Failed to deliver payload")
	}

	return nil
}
