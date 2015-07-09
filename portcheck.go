package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

type PortCheck struct {
	Network string
	Host    string
	Port    int
}

func ParsePortCheck(data map[string]interface{}) PortCheck {
	check := PortCheck{Network: "tcp"}

	if data["net"] != nil {
		check.Network = data["net"].(string)
	}

	if data["host"] != nil {
		check.Host = data["host"].(string)
	}

	if data["port"] != nil {
		check.Port = data["port"].(int)
	}

	return check
}

func (check PortCheck) Name() string {
	return "PORT"
}

func (check PortCheck) Perform() error {
	log.Printf(
		"Performing %v check for host=%v net=%v port=%v\n",
		check.Name(), check.Host, check.Network, check.Port,
	)

	if check.Host == "" {
		return fmt.Errorf("Host should not be empty")
	}

	if check.Port <= 0 {
		return fmt.Errorf("Invalid port: %v", check.Port)
	}

	dialer := net.Dialer{
		Timeout:   time.Second * 10,
		KeepAlive: 0,
	}

	conn, err := dialer.Dial(check.Network, fmt.Sprintf("%v:%v", check.Host, check.Port))
	if err != nil {
		return err
	}

	conn.Close()
	return nil
}
