package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
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
			debounced := utils.Debounce(100 * time.Millisecond)

			for {
				select {
				case event, ok := <-watcher.Events:
					if !ok {
						return
					}
					log.Println("event:", event)
					if event.Has(fsnotify.Write) {
						log.Println("modified file:", event.Name)
						debounced(func() {
							filePath := event.Name
							sha256, err := utils.FileSHA256(filePath)
							if err != nil {
								logger.Errorf(err.Error())
							}
							for i := 0; i < configs.Size(); i++ {
								machineId := configs.Get(i, "machineId").ToString()
								if machineId == data.MachineId {
									actionPath := configs.Get(i, "path").ToString()
									if actionPath == filePath {
										fileId := configs.Get(i, "fileId").ToString()
										if fileId != "" {
											fileName := filepath.Base(filePath)
											f, err := os.ReadFile(filePath)
											if err != nil {
												fmt.Println("read fail", err)
											}
											fileContent := string(f)
											timestamp := time.Now().UnixNano() / 1e6
											err = fileUpload(fileId, fileName, sha256, fileContent, timestamp, data)
											if err != nil {
												logger.Errorf(err.Error())
											}
										}
										break
									}
								}
							}
						})
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
