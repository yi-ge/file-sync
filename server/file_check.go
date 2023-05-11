package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (handler *FileCheckHandler) postXhr(c *gin.Context) {
	var req FileCheckRequest
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(400, gin.H{
			"status": -1,
			"msg":    "Missing required parameters",
			"result": nil,
		})
		return
	}

	emailSha1 := sha1Hash(req.Email)

	var file File
	result := handler.DB.Order("update_at DESC").First(&file, "email_sha1 = ? AND file_id = ?", emailSha1, req.FileId)

	if result.Error == gorm.ErrRecordNotFound {
		c.JSON(200, gin.H{
			"status": 0,
			"msg":    "File not found",
			"result": nil,
		})
		return
	} else if result.Error != nil {
		c.JSON(500, gin.H{
			"status": -99,
			"msg":    "Unknown error.",
			"result": nil,
		})
		return
	}

	if file.SHA256 != req.Sha256 {
		c.JSON(200, gin.H{
			"status": 1,
		})
	} else {
		c.JSON(200, gin.H{
			"status": 2,
		})
	}
}
