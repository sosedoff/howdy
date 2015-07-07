package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

const (
	VERSION = "0.2.0"
)

var (
	ConfigsPath       string
	SingleConfigPath  string
	ShowVersion       bool
	SendNotifications bool
	TestMode          bool
)

func init() {
	flag.StringVar(&ConfigsPath, "c", "", "Path to all configs")
	flag.BoolVar(&ShowVersion, "v", false, "Show version")
	flag.BoolVar(&SendNotifications, "n", true, "Send notifications")
	flag.BoolVar(&TestMode, "t", false, "Test mode")
	flag.Parse()

	if ShowVersion {
		fmt.Println("Version:", VERSION)
		os.Exit(0)
	}

	if ConfigsPath == "" {
		if len(flag.Args()) == 0 {
			log.Fatalln("Please specify path to configs directory")
		} else {
			SingleConfigPath = flag.Args()[0]
		}
	}

	if !SendNotifications {
		log.Println("Will not send any notifications")
	}

	ConfigsPath = strings.Replace(ConfigsPath, "~", os.Getenv("HOME"), 1)
	SingleConfigPath = strings.Replace(SingleConfigPath, "~", os.Getenv("HOME"), 1)
}

func runAll() {
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

	wg := sync.WaitGroup{}
	wg.Add(len(files))

	for _, file := range files {
		if filepath.Ext(file.Name()) != ".yml" {
			wg.Done()
			continue
		}

		config, err := ReadConfig(ConfigsPath + "/" + file.Name())
		if err != nil {
			log.Println(err)
			continue
		}

		if TestMode || config.Enabled {
			go RunConfig(config, &wg)
		} else {
			log.Println("Skipping config:", config.Name)
		}
	}

	wg.Wait()
}

func runSingle() {
	config, err := ReadConfig(SingleConfigPath)
	if err != nil {
		log.Fatalln(err)
	}

	RunConfig(config, nil)
}

func main() {
	if ConfigsPath != "" {
		runAll()
	} else {
		runSingle()
	}
}
