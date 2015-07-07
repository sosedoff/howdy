package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
)

var (
	ConfigsPath string
)

func init() {
	flag.StringVar(&ConfigsPath, "c", "", "Path to all configs")
	flag.Parse()

	if ConfigsPath == "" {
		log.Fatalln("Please specify path to configs directory")
	}
}

func main() {
	file, err := os.Open(ConfigsPath)
	if err != nil {
		log.Fatalln(err)
	}

	stat, err := file.Stat()
	if err != nil {
		log.Fatalln(err)
	}

	if !stat.IsDir() {
		log.Fatalln("Path is not a directory")
	}

	files, err := ioutil.ReadDir(ConfigsPath)
	if err != nil {
		log.Fatalln(err)
	}

	for _, file := range files {
		config, err := ReadConfig(ConfigsPath + "/" + file.Name())
		if err != nil {
			log.Fatalln(err)
		}

		RunConfig(config)
	}
}
