package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

type WebCheck struct {
	Url    string `yaml:"url"`
	Status int    `yaml:"status"`
	Format string `yaml:"format"`
}

func ParseWebCheck(data map[string]interface{}) WebCheck {
	check := WebCheck{}

	if data["url"] != nil {
		check.Url = data["url"].(string)
	}

	if data["status"] != nil {
		check.Status = data["status"].(int)
	} else {
		check.Status = 200
	}

	if data["format"] != nil {
		check.Format = data["format"].(string)
	} else {
		check.Format = "html"
	}

	return check
}

func (check WebCheck) Perform() error {
	log.Printf(
		"Performing WEB check for url=%v status=%v format=%v\n",
		check.Url, check.Status, check.Format,
	)

	if check.Url == "" {
		return fmt.Errorf("URL should not be empty")
	}

	resp, err := http.Get(check.Url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != check.Status {
		return fmt.Errorf("Expected HTTP status: %v, got: %v", check.Status, resp.StatusCode)
	}

	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, check.Format) {
		return fmt.Errorf("Expected HTTP format '%v' to include '%v'", contentType, check.Format)
	}

	return nil
}
