package main

import (
	"fmt"
	"log"
)

func RunCheck(config *Config, check Check) {
	err := check.Perform()

	if err != nil {
		log.Println("Check failed:", err)

		msg := fmt.Sprintf(
			"%v check failed for %v: %v",
			check.Name(), config.Name, err.Error(),
		)

		if SendNotifications {
			for _, notifier := range config.Notifiers {
				err = notifier.Perform(msg)
				if err != nil {
					log.Println("Notifier failed:", err)
				}
			}
		}
	}
}

func RunConfig(config *Config) {
	if TestMode || config.Enabled {
		for _, check := range config.Checks {
			RunCheck(config, check)
		}
	} else {
		log.Println("Skipping config:", config.Name)
	}
}
