package main

import (
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
			err := check.Perform()
			if err != nil {
				log.Println("Check failed:", err)
			}
		}
	}
}
