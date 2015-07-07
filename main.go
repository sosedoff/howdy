package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

const (
	VERSION = "0.1.0"
)

var (
	ConfigsPath string
	ShowVersion bool
)

func init() {
	flag.StringVar(&ConfigsPath, "c", "", "Path to all configs")
	flag.BoolVar(&ShowVersion, "v", false, "Show version")
	flag.Parse()

	if ShowVersion {
		fmt.Println("Version:", VERSION)
		os.Exit(0)
	}

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
		if filepath.Ext(file.Name()) != ".yml" {
			continue
		}

		config, err := ReadConfig(ConfigsPath + "/" + file.Name())
		if err != nil {
			log.Fatalln(err)
		}

		RunConfig(config)
	}
}
