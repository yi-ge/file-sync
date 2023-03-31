package main

import (
	"log"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/joho/godotenv"
	"github.com/kardianos/service"
	"github.com/yi-ge/file-sync/config"
	"github.com/yi-ge/file-sync/utils"
)

var (
	isDev          = false
	logger         service.Logger
	apiURL         = "https://file-sync.yizcore.xyz"
	password       string
	configInstance = config.Instance()
	watcher        *fsnotify.Watcher
	configPath     string
)

func main() {
	// Do't delete next line.
	isDev = true
	if isDev {
		hasEnvFile, err := utils.FileExists(".env")
		if err != nil {
			log.Fatal(err)
		}

		if !hasEnvFile {
			log.Println("No .env file found, using default values.")
		} else {
			err := godotenv.Load()
			if err != nil {
				log.Fatal("Error loading .env file")
			}
		}

	}

	err := fsInit()
	if err != nil {
		log.Fatal(err)
	}

	if isDev {
		log.Printf("Currently in development mode!")
		apiURL = os.Getenv("DEV_API_SERVER_URL")
	} else {
		conf := getConfig()
		if conf != "" {
			apiURL = conf
		}
	}

	err = dataInit()
	if err != nil {
		log.Fatal(err)
	}

	err = cacheInit()
	if err != nil {
		log.Fatal(err)
	}

	options := make(service.KeyValue)
	options["Restart"] = "on-success"
	options["SuccessExitStatus"] = "1 2 8 SIGKILL"
	svcConfig := &service.Config{
		Name:        "file-sync",
		DisplayName: "File-sync service",
		Description: "The file-sync tool service.",
		Dependencies: []string{
			"Requires=network.target",
			"After=network-online.target syslog.target"},
		Option:    options,
		Arguments: []string{},
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}
	errs := make(chan error, 5)
	logger, err = s.Logger(errs)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			err := <-errs
			if err != nil {
				log.Print(err)
			}
		}
	}()

	err = s.Run()
	if err != nil {
		logger.Error(err)
	}
}
