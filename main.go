package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const (
	VERSION = "0.3.1"
)

var (
	RunId             int64
	DbPath            string
	StartServer       bool
	ConfigsPath       string
	SingleConfigPath  string
	LogPath           string
	ShowVersion       bool
	SendNotifications bool
	TestMode          bool
	LogFile           *os.File
	LogWriter         *bufio.Writer
)

func init() {
	RunId = time.Now().Unix()

	flag.StringVar(&ConfigsPath, "c", "", "Path to all configs")
	flag.StringVar(&LogPath, "l", "", "Path to log")
	flag.BoolVar(&ShowVersion, "v", false, "Show version")
	flag.BoolVar(&SendNotifications, "n", true, "Send notifications")
	flag.BoolVar(&TestMode, "t", false, "Test mode")
	flag.BoolVar(&StartServer, "s", false, "Start server")
	flag.StringVar(&DbPath, "d", "howdy.sqlite3", "Path to database")
	flag.Parse()

	if ShowVersion {
		fmt.Println("Version:", VERSION)
		os.Exit(0)
	}

	if LogPath != "" {
		fd, err := os.OpenFile(LogPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalln(err)
		}

		LogFile = fd
		LogWriter = bufio.NewWriter(LogFile)

		log.SetOutput(LogWriter)
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

	// Setup database connection
	err := dbSetup()
	if err != nil {
		log.Fatal(err)
	}
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

	wg := &sync.WaitGroup{}
	wg.Add(len(files))

	for _, file := range files {
		if filepath.Ext(file.Name()) != ".yml" {
			wg.Done()
			continue
		}

		config, err := ReadConfig(ConfigsPath + "/" + file.Name())
		if err != nil {
			log.Println("Error reading config:", err)
			continue
		}

		if TestMode || config.Enabled {
			go RunConfig(config, wg)
		} else {
			log.Println("Skipping config:", config.Name)
			wg.Done()
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

func cleanup() {
	if LogWriter != nil {
		LogWriter.Flush()
	}

	if LogFile != nil {
		LogFile.Close()
	}
}

func main() {
	defer cleanup()

	if StartServer {
		runHttpServer()
		return
	}

	if ConfigsPath != "" {
		runAll()
	} else {
		runSingle()
	}
}
