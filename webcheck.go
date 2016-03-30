package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

type WebCheck struct {
	Url    string `yaml:"url"`
	Status int    `yaml:"status"`
	Format string `yaml:"format"`
}

func init() {
	// TODO: Make sure SSL verification works. This is just a workaround.
	cfg := &tls.Config{InsecureSkipVerify: true}

	http.DefaultClient.Transport = &http.Transport{
		TLSClientConfig: cfg,
	}

	http.DefaultClient.Timeout = time.Second * 10
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

func (check WebCheck) Name() string {
	return "WEB"
}

func (check WebCheck) Perform() error {
	log.Printf(
		"[%v] url=%v status=%v format=%v\n",
		check.Name(), check.Url, check.Status, check.Format,
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
		return fmt.Errorf(
			"Expected HTTP status %v for %v, got: %v",
			check.Status, check.Url, resp.StatusCode,
		)
	}

	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, check.Format) {
		return fmt.Errorf(
			"Expected HTTP format '%v' for %v to include '%v'",
			contentType, check.Url, check.Format,
		)
	}

	return nil
}
