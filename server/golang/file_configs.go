package main

import (
	"crypto"
	"crypto/rsa"
	"encoding/base64"
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (handler *FileConfigsHandler) postXhr(c *gin.Context) {
	var req FileConfigsRequest
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
	sign := fmt.Sprintf("email=%s&machineId=%s&timestamp=%s&%s", req.Email, req.MachineId, req.Timestamp, user.Verify)

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

	var configs []ConfigResponse
	err = handler.DB.Raw(`
		SELECT
			config.id as id,
			config.file_name as fileName,
			config.file_id as fileId,
			file.update_at as updateAt,
			config.machine_id as machineId,
			device.machine_name as machineName,
			config.path as path,
			config.attribute as attribute,
			config.created_at as createdAt
		FROM
			config
		LEFT JOIN file on file.id = (
			SELECT f.id FROM file AS f
			WHERE config.file_id = f.file_id
			ORDER BY f.update_at DESC
			LIMIT 1
		)
		LEFT JOIN device on config.machine_id = device.machine_id
		WHERE
			config.email = ? AND
			deleted_at IS NULL
		ORDER BY
			config.file_id,
			config.created_at,
			config.id
	`, user.Email).Scan(&configs).Error

	if err != nil {
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
		"result": configs,
	})
}
