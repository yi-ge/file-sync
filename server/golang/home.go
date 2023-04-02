package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HomeHandler(c *gin.Context) {
	c.AbortWithStatus(http.StatusForbidden)
	c.Writer.WriteString("403 Unauthorized")
}
