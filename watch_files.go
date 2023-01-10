package main

import (
	"fmt"
	"log"
	"time"

	"github.com/fsnotify/fsnotify"
)

func watchFiles(data Data) {
	// TODO: 网络异常后重新检查
	// TODO: Recheck after network anomaly
	configs, err := listConfigs(data)
	if err == nil {
		// TODO: 检查config是否已被从其他设备移除
		// TODO: Check if config has been removed from other devices
		fmt.Println(configs.ToString())

		// Create new watcher.
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			log.Fatal(err)
		}
		defer watcher.Close()

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
				err = watcher.Add(actionPath)
				if err != nil {
					logger.Error(err)
				}
			}
		}

	} else {
		logger.Error(err)
		logger.Info("retry watcher")
		time.Sleep(time.Duration(2) * time.Second)
		watchFiles(data)
	}
}
