package main

import (
	"fmt"
	"log"
	"os/exec"
)

type PingCheck struct {
	Host string
}

func ParsePingCheck(data map[string]interface{}) PingCheck {
	check := PingCheck{}

	if data["host"] != nil {
		check.Host = data["host"].(string)
	}

	return check
}

func (check PingCheck) Name() string {
	return "PING"
}

func (check PingCheck) Perform() error {
	log.Printf("Performing %v check for ip=%v\n", check.Name(), check.Host)
	if check.Host == "" {
		return fmt.Errorf("Host should not be empty")
	}

	return exec.Command("ping", "-c", "1", check.Host).Run()
}
