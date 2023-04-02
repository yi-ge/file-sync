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

func (handler *FileSyncHandler) postXhr(c *gin.Context) {
	var req FileSyncRequest
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
	sign := fmt.Sprintf("email=%s&fileId=%s&machineId=%s&timestamp=%s&%s", req.Email, req.FileId, req.MachineId, req.Timestamp, user.Verify)

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

	if req.UpdateAt == 0 { // Download
		var file File
		result := handler.DB.First(&file, "email = ? AND file_id = ?", user.Email, req.FileId)

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

		c.JSON(200, gin.H{
			"status": 1,
			"msg":    "OK",
			"result": file,
		})
	} else { // Upload
		if req.Content == "" || req.Sha256 == "" || req.FileName == "" {
			c.JSON(200, gin.H{
				"status": -4,
				"msg":    "Missing required parameters",
				"result": nil,
			})
			return
		}

		updateAt := time.Unix(req.UpdateAt/1000, 0)

		file := File{
			Email:         user.Email,
			EmailSha1:     emailSha1,
			FileId:        req.FileId,
			FileName:      req.FileName,
			Content:       req.Content,
			SHA256:        req.Sha256,
			FromMachineId: req.MachineId,
			UpdateAt:      updateAt,
		}

		result := handler.DB.Create(&file)
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
				"lastId": file.ID,
				"fileId": req.FileId,
			},
		})
	}
}
