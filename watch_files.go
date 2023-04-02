package main

import (
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
	jsoniter "github.com/json-iterator/go"
	"github.com/yi-ge/file-sync/utils"
)

var debounced = utils.Debounce(100 * time.Millisecond)

func checkAndUploadFile(filePath string, configs jsoniter.Any, privateKey []byte, data Data) {
	debounced(func() {
		sha256, err := utils.FileSHA256(filePath)
		if err != nil {
			logger.Errorf(err.Error())
		}
		for i := 0; i < configs.Size(); i++ {
			machineId := configs.Get(i, "machineId").ToString()
			if machineId == data.MachineId {
				actionPathEncrypted := configs.Get(i, "path").ToString()
				actionPathBase64, err := base64.URLEncoding.DecodeString(actionPathEncrypted)
				if err != nil {
					logger.Errorf(err.Error())
					break
				}
				actionPath := string(utils.RsaDecrypt([]byte(actionPathBase64), privateKey))
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

func watchFiles(data Data) {
	privateKeyEncrypted, err := getPrivateKey()
	if err != nil {
		logger.Errorf(err.Error())
		return
	}

	privateKeyHex, err := base64.RawURLEncoding.DecodeString(string(privateKeyEncrypted))
	if err != nil {
		logger.Errorf(err.Error())
		return
	}

	decrypted, privateKey, err := utils.AESMACDecryptBytes(privateKeyHex, data.RsaPrivateKeyPassword)

	if err != nil || !decrypted {
		logger.Errorf((errors.New("secret decrypt error: " + err.Error())).Error())
		return
	}

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
					log.Println("File system event:", event)
					if event.Has(fsnotify.Write) {
						log.Println("Modified file:", event.Name)
						filePath := event.Name
						checkAndUploadFile(filePath, configs, privateKey, data)
					}

					if event.Has(fsnotify.Remove) {
						filePath := event.Name
						go func() {
							for range time.Tick(time.Second) {
								fileExists, _ := utils.FileExists(filePath)
								if fileExists {
									err = watcher.Add(filePath)
									if err != nil {
										logger.Error(err)
									}

									checkAndUploadFile(filePath, configs, privateKey, data)
									break
								}
							}
						}()
					}
				case err, ok := <-watcher.Errors:
					if !ok {
						return
					}
					log.Println("File system error:", err)
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
		logger.Errorf("Watch file error: %s\n", err.Error())
		logger.Info("Retry watcher")
		time.Sleep(time.Duration(2) * time.Second)
		watchFiles(data)
	}
}
