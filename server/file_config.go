package main

import (
	"crypto"
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (handler *FileConfigHandler) postXhr(c *gin.Context) {
	var req FileConfigRequest
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

	var user User
	result := handler.DB.First(&user, "email_sha1 = ?", emailSha1)

	if result.Error == gorm.ErrRecordNotFound {
		c.JSON(200, gin.H{
			"status": -2,
			"msg":    "Invalid email",
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

	// Verify the token
	sign := fmt.Sprintf("email=%s&machineId=%s&timestamp=%s&fileId=%s&action=%s&actionMachineId=%s&attribute=%s&%s", req.Email, req.MachineId, req.Timestamp, req.FileId, req.Action, req.ActionMachineId, req.Attribute, user.Verify)

	token, err := base64.URLEncoding.DecodeString(req.Token)
	if err != nil {
		c.JSON(200, gin.H{
			"status": -3,
			"msg":    "Invalid token",
			"result": nil,
		})
		return
	}

	publicKey, err := getPublicKeyFromPem(user.PublicKey)
	if err != nil {
		c.JSON(200, gin.H{
			"status": -3,
			"msg":    "Invalid token",
			"result": nil,
		})
		return
	}

	err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA1, []byte(sha1Hash(sign)), token)
	if err != nil {
		c.JSON(200, gin.H{
			"status": -3,
			"msg":    "Invalid token",
			"result": nil,
		})
		return
	}

	if req.Action == "add" {
		var config Config
		result := handler.DB.Where("email = ? AND machine_id = ? AND file_id = ? AND deleted_at IS NULL", user.Email, req.ActionMachineId, req.FileId).First(&config)

		if result.Error == nil {
			c.JSON(200, gin.H{
				"status": -5,
				"msg":    "This machine has been synchronized with the project of binding this file.",
				"result": nil,
			})
			return
		} else if result.Error != gorm.ErrRecordNotFound {
			c.JSON(500, gin.H{
				"status": -99,
				"msg":    "Unknown error.",
				"result": nil,
			})
			return
		}

		newConfig := Config{
			Email:     user.Email,
			MachineId: req.ActionMachineId,
			FileId:    req.FileId,
			FileName:  req.FileName,
			Path:      req.Path,
			Attribute: req.Attribute,
		}

		result = handler.DB.Create(&newConfig)

		if result.Error != nil {
			c.JSON(500, gin.H{
				"status": -99,
				"msg":    "Unknown error.",
				"result": nil,
			})
			return
		}

		c.JSON(200, gin.H{
			"status": 1,
			"msg":    "OK",
			"result": gin.H{
				"lastId": newConfig.ID,
				"fileId": req.FileId,
			},
		})
	} else if req.Action == "remove" {
		result := handler.DB.Model(&Config{}).Where("email = ? AND file_id = ? AND machine_id = ? AND deleted_at IS NULL", user.Email, req.FileId, req.ActionMachineId).Update("deleted_at", time.Now())

		if result.RowsAffected == 1 {
			c.JSON(200, gin.H{
				"status": 1,
				"msg":    "OK",
				"result": nil,
			})
		} else {
			c.JSON(200, gin.H{
				"status": -4,
				"msg":    "Delete fail.",
				"result": nil,
			})
		}
	} else {
		c.JSON(200, gin.H{
			"status": -2,
			"msg":    "Invalid action",
			"result": nil,
		})
	}
}
