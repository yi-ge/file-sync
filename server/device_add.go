package main

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (handler *DeviceAddHandler) postXhr(c *gin.Context) {
	var req DeviceAddRequest
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(400, gin.H{
			"status": -1,
			"msg":    "Missing required parameters",
			"result": nil,
		})
		return
	}

	// Check if the user exists
	var user User
	result := handler.DB.First(&user, "email = ?", req.Email)

	// If the user does not exist, create a new user and device
	if result.Error == gorm.ErrRecordNotFound {
		// Create new user
		newUser := User{
			Email:      req.Email,
			EmailSha1:  sha1Hash(req.Email),
			Verify:     req.Verify,
			PublicKey:  req.PublicKey,
			PrivateKey: req.PrivateKey,
			CreatedAt:  time.Now(),
		}
		handler.DB.Create(&newUser)

		// Create new device
		machineKey := uuid.New().String()
		newDevice := Device{
			Email:       req.Email,
			MachineId:   req.MachineId,
			MachineName: req.MachineName,
			MachineKey:  machineKey,
			CreatedAt:   time.Now(),
		}
		handler.DB.Create(&newDevice)

		// Encrypt the public and private keys
		encryptedPublicKey, _ := encrypt(req.PublicKey, req.Verify)
		encryptedPrivateKey, _ := encrypt(req.PrivateKey, req.Verify)

		c.JSON(200, gin.H{
			"status": 1,
			"msg":    "New user added",
			"result": gin.H{
				"publicKey":  encryptedPublicKey,
				"privateKey": encryptedPrivateKey,
				"machineKey": machineKey,
			},
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

	// Check if the device already exists
	var device Device
	result = handler.DB.First(&device, "email = ? AND machine_id = ?", req.Email, req.MachineId)

	// If the device exists, return an error
	if result.Error == nil {
		c.JSON(200, gin.H{
			"status": -2,
			"msg":    "Device already exists",
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

	// If the verification fails, return an error
	if user.Verify != req.Verify {
		c.JSON(200, gin.H{
			"status": -3,
			"msg":    "Verification Rejected",
			"result": nil,
		})
		return
	}

	// Add the new device
	machineKey := uuid.New().String()
	newDevice := Device{
		Email:       req.Email,
		MachineId:   req.MachineId,
		MachineName: req.MachineName,
		MachineKey:  machineKey,
		CreatedAt:   time.Now(),
	}
	handler.DB.Create(&newDevice)

	// Encrypt the public and private keys
	encryptedPublicKey, _ := encrypt(user.PublicKey, req.Verify)
	encryptedPrivateKey, _ := encrypt(user.PrivateKey, req.Verify)

	c.JSON(200, gin.H{
		"status": 2,
		"msg":    "Device added",
		"result": gin.H{
			"publicKey":  encryptedPublicKey,
			"privateKey": encryptedPrivateKey,
			"machineKey": machineKey,
		},
	})
}
