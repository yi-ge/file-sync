package main

import (
	"crypto"
	"crypto/rsa"
	"encoding/base64"
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (handler *DeviceListHandler) postXhr(c *gin.Context) {
	var req DeviceListRequest
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

	var machineList []Device
	handler.DB.Find(&machineList, "email = ?", user.Email)

	c.JSON(200, gin.H{
		"status": 1,
		"msg":    "OK",
		"result": machineList,
	})
}
