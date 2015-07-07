package main

import (
	"fmt"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

type Check interface {
	Perform() error
}

type Config struct {
	Name    string  `yaml:"name"`
	Enabled bool    `yaml:"enabled"`
	Checks  []Check `yaml:"-"`
}

func ReadConfig(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	config := Config{}
	if err = yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	checks := struct {
		Checks map[string][]map[string]interface{} `yaml:"checks"`
	}{}
	if err = yaml.Unmarshal(data, &checks); err != nil {
		return nil, err
	}

	for t, items := range checks.Checks {
		switch t {
		case "web":
			for _, item := range items {
				config.Checks = append(config.Checks, ParseWebCheck(item))
			}
		case "ping":
			for _, item := range items {
				config.Checks = append(config.Checks, ParsePingCheck(item))
			}
		case "dns":
			for _, item := range items {
				config.Checks = append(config.Checks, ParseDnsCheck(item))
			}
		default:
			return nil, fmt.Errorf("Invalid check type:", t)
		}
	}

	return &config, nil
}
