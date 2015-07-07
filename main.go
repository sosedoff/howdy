package main

import (
	"fmt"
	"io/ioutil"
	"log"
)

func main() {
	base := "./configs"

	files, err := ioutil.ReadDir(base)
	if err != nil {
		log.Fatalln(err)
	}

	for _, file := range files {
		config, err := ReadConfig(base + "/" + file.Name())
		if err != nil {
			log.Fatalln(err)
		}

		for _, check := range config.Checks {
			err = check.Perform()
			if err != nil {
				log.Println("Check failed:", err)
				msg := fmt.Sprintf(
					"%v check failed for %v: %v",
					check.Name(), config.Name, err.Error(),
				)

				for _, notifier := range config.Notifiers {
					err = notifier.Perform(msg)
					if err != nil {
						log.Println("Notifier failed:", err)
					}
				}
			}
		}
	}
}
