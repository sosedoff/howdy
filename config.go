package main

import (
	"fmt"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

type Check interface {
	Name() string
	Perform() error
}

type Notifier interface {
	Perform(string) error
}

type Config struct {
	Name      string     `yaml:"name"`
	Enabled   bool       `yaml:"enabled"`
	Checks    []Check    `yaml:"-"`
	Notifiers []Notifier `yaml:"-"`
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
		Checks    map[string][]map[string]interface{} `yaml:"checks"`
		Notifiers map[string]map[string]interface{}   `yaml:"notify"`
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
		case "port":
			for _, item := range items {
				config.Checks = append(config.Checks, ParsePortCheck(item))
			}
		default:
			return nil, fmt.Errorf("Invalid check type:", t)
		}
	}

	for t, items := range checks.Notifiers {
		switch t {
		case "slack":
			config.Notifiers = append(config.Notifiers, ParseSlackNotifier(items))
		default:
			return nil, fmt.Errorf("Invalid notifier type:", t)
		}
	}

	return &config, nil
}
