package main

import (
	"fmt"
	"log"
	"net"
	"time"

	fastping "github.com/tatsushid/go-fastping"
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
	log.Printf("Performing PING check for ip=%v\n", check.Host)
	if check.Host == "" {
		return fmt.Errorf("Host should not be empty")
	}

	pongCount := 0

	p := fastping.NewPinger()
	p.Network("udp")

	ra, err := net.ResolveIPAddr("ip4:icmp", check.Host)
	if err != nil {
		return err
	}
	p.AddIPAddr(ra)

	p.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
		pongCount += 1
	}

	err = p.Run()
	if err != nil {
		return err
	}

	return nil
}
