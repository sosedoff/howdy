package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
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

func (notifier SlackNotifier) Perform(messages []string) error {
	payload := map[string]interface{}{
		"unfurl_links": false,
		"username":     "service-health",
		"attachments": []map[string]string{
			map[string]string{
				"color": "danger",
				"text":  strings.Join(messages, "\n"),
			},
		},
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
		respData, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("Failed to deliver payload: %s\n", respData)
	}

	return nil
}
