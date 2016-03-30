package main

import (
	"fmt"
	"log"
	"net"
)

type DnsCheck struct {
	Domain string
	Ip     string
}

func ParseDnsCheck(data map[string]interface{}) DnsCheck {
	check := DnsCheck{}

	if data["domain"] != nil {
		check.Domain = data["domain"].(string)
	}

	if data["ip"] != nil {
		check.Ip = data["ip"].(string)
	}

	return check
}

func (check DnsCheck) Name() string {
	return "DNS"
}

func (check DnsCheck) Perform() error {
	if check.Domain == "" {
		return fmt.Errorf("Domain should not be empty")
	}

	if check.Ip == "" {
		return fmt.Errorf("IP should not be empty")
	}

	log.Printf("[%v] domain=%v ip=%v\n", check.Name(), check.Domain, check.Ip)

	addr, err := net.ResolveIPAddr("ip4:icmp", check.Domain)
	if err != nil {
		return err
	}

	if addr.String() != check.Ip {
		return fmt.Errorf("Expected IP: %v, got: %v", check.Ip, addr.String())
	}

	return nil
}
