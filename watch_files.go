package main

import (
	"log"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/yi-ge/file-sync/utils"
)

func watchFiles(data Data) {
	configs, err := listConfigs(data)
	if err == nil {
		if watcher != nil {
			watcher.Close()
			watcher = nil
		}

		// Create new watcher.
		watcher, err = fsnotify.NewWatcher()
		if err != nil {
			log.Fatal(err)
		}
		// defer watcher.Close()

		// Start listening for events.
		go func() {
			for {
				select {
				case event, ok := <-watcher.Events:
					if !ok {
						return
					}
					log.Println("event:", event)
					if event.Has(fsnotify.Write) {
						log.Println("modified file:", event.Name)
					}

					if event.Has(fsnotify.Remove) {
						actionPath := event.Name
						go func() {
							for range time.Tick(time.Second) {
								fileExists, _ := utils.FileExists(actionPath)
								if fileExists {
									err = watcher.Add(actionPath)
									// TODO：检查一次该文件
									if err != nil {
										logger.Error(err)
									}
									break
								}
							}
						}()
					}
				case err, ok := <-watcher.Errors:
					if !ok {
						return
					}
					log.Println("error:", err)
				}
			}
		}()

		for i := 0; i < configs.Size(); i++ {
			machineId := configs.Get(i, "machineId").ToString()
			if machineId == data.MachineId {
				actionPath := configs.Get(i, "path").ToString()
				// Add a path.
				fileExists, _ := utils.FileExists(actionPath)
				if fileExists {
					err = watcher.Add(actionPath)
					if err != nil {
						logger.Error(err)
					}
				} else {
					go func() {
						for range time.Tick(time.Second) {
							fileExists, _ := utils.FileExists(actionPath)
							if fileExists {
								err = watcher.Add(actionPath)
								if err != nil {
									logger.Error(err)
								}
								break
							}
						}
					}()
				}
			}
		}

		logger.Infof("Watcher is Working.")
	} else {
		logger.Error(err)
		logger.Info("retry watcher")
		time.Sleep(time.Duration(2) * time.Second)
		watchFiles(data)
	}
}
