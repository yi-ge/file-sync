package main

import (
	"crypto"
	"crypto/rsa"
	"encoding/base64"
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (handler *DeviceRemoveHandler) postXhr(c *gin.Context) {
	var req DeviceRemoveRequest
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
	sign := fmt.Sprintf("email=%s&machineId=%s&timestamp=%s&removeMachineId=%s&machineKey=%s&%s", req.Email, req.MachineId, req.Timestamp, req.RemoveMachineId, req.MachineKey, user.Verify)

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

	var machine Device
	handler.DB.First(&machine, "email = ? AND machine_id = ?", user.Email, req.MachineId)

	if sha1Hash(machine.MachineKey) != req.MachineKey {
		c.JSON(200, gin.H{
			"status": -4,
			"msg":    "Invalid machineKey",
			"result": nil,
		})
		return
	}

	result = handler.DB.Delete(&Device{}, "machine_id = ?", req.RemoveMachineId)

	if result.RowsAffected == 1 {
		c.JSON(200, gin.H{
			"status": 1,
			"msg":    "OK",
			"result": gin.H{
				"removedMachineId": req.RemoveMachineId,
			},
		})
	} else {
		c.JSON(200, gin.H{
			"status": -5,
			"msg":    "Failed to remove",
			"result": nil,
		})
	}
}
