package main

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/gin-gonic/gin"
)

func (handler *EventStreamHandler) eventStream(c *gin.Context) {
	email := c.Query("email")
	if email == "" {
		c.String(400, "Invalid request")
		return
	}

	emailSha1 := sha1Hash(email)

	// Set the headers for SSE
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")

	// Send the connected event
	c.Stream(func(w io.Writer) bool {
		fmt.Fprint(w, "event: connected\n")
		fmt.Fprint(w, "data: 1\n\n")
		return true
	})

	// Send events in a loop
	ticker := time.NewTicker(1 * time.Second)
	heartbeatTicker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()
	defer heartbeatTicker.Stop()

	for {
		select {
		case <-ticker.C:
			var files []File
			handler.DB.Where("email_sha1 = ?", emailSha1).Find(&files)

			if len(files) > 0 {
				data, _ := json.Marshal(files)
				c.Stream(func(w io.Writer) bool {
					fmt.Fprint(w, "event: file\n")
					fmt.Fprintf(w, "data: %s\n\n", data)
					return true
				})
			}
		case <-heartbeatTicker.C:
			c.Stream(func(w io.Writer) bool {
				fmt.Fprint(w, "event: heartbeat\n")
				fmt.Fprint(w, "data: 1\n\n")
				return true
			})
		}
	}
}
