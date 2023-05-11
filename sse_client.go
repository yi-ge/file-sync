package main

import (
	"crypto/tls"
	"net/http"
	"strconv"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/yi-ge/file-sync/utils"
	sse "github.com/yi-ge/sse/v2"
)

func StartSSEClient(data Data) {
	timestamp := time.Now().UnixNano() / 1e6
	emailSha1 := utils.GetSha1Str(data.Email)
	eventURL := apiURL + "/events?email=" + emailSha1 + "&machineId=" + data.MachineId + "&timestamp=" + strconv.FormatInt(timestamp, 10)
	// fmt.Println(eventURL)
	client := sse.NewClient(eventURL)
	client.AutoReconnect = true

	// disabling ssl verification for self signed certs
	client.Connection.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client.OnConnect(func(c *sse.Client) {
		logger.Infof("Connected!")
	})

	client.OnDisconnect(func(c *sse.Client) {
		logger.Infof("Disconnected!")

		if watcher != nil {
			watcher.Close()
			watcher = nil
		}
	})

	logger.Infof("Registered Server Events!")

	err := client.Subscribe("messages", func(msg *sse.Event) {
		if string(msg.Event) == "connected" {
			logger.Infof("Event server connected.")
			go watchFiles(data) // Recheck after network anomaly
		} else if string(msg.Event) == "file" {
			logger.Infof("File Event: %s", string(msg.Data))
			var json []Jobs
			err := jsoniter.Unmarshal(msg.Data, &json)
			if err != nil {
				logger.Errorf(err.Error())
				return
			}
			go job(json, emailSha1, data)
		} else if string(msg.Event) == "config" {
			// Check if config has been removed from other devices
			go watchFiles(data)
		} else if string(msg.Event) == "heartbeat" {
			logger.Infof("Receive heartbeat package.")
		}
	})

	if err != nil {
		logger.Infof(err.Error())
	}
}
