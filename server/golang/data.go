package main

import (
	"gorm.io/gorm"
)

type DeviceAddHandler struct {
	DB *gorm.DB
}

type DeviceListHandler struct {
	DB *gorm.DB
}

type DeviceRemoveHandler struct {
	DB *gorm.DB
}

type FileConfigHandler struct {
	DB *gorm.DB
}

type FileConfigsHandler struct {
	DB *gorm.DB
}

type FileCheckHandler struct {
	DB *gorm.DB
}

type FileSyncHandler struct {
	DB *gorm.DB
}

type EventStreamHandler struct {
	DB *gorm.DB
}

type DeviceAddRequest struct {
	Email       string `json:"email"`
	MachineId   string `json:"machineId"`
	MachineName string `json:"machineName"`
	Verify      string `json:"verify"`
	PublicKey   string `json:"publicKey"`
	PrivateKey  string `json:"privateKey"`
}

type DeviceListRequest struct {
	Email     string `json:"email"`
	MachineId string `json:"machineId"`
	Timestamp string `json:"timestamp"`
	Token     string `json:"token"`
}

type DeviceRemoveRequest struct {
	Email           string `json:"email"`
	MachineId       string `json:"machineId"`
	Timestamp       string `json:"timestamp"`
	RemoveMachineId string `json:"removeMachineId"`
	MachineKey      string `json:"machineKey"`
	Token           string `json:"token"`
}

type FileConfigRequest struct {
	Email           string `json:"email"`
	MachineId       string `json:"machineId"`
	Timestamp       string `json:"timestamp"`
	FileId          string `json:"fileId"`
	Action          string `json:"action"`
	ActionMachineId string `json:"actionMachineId"`
	Attribute       string `json:"attribute"`
	Token           string `json:"token"`
	Path            string `json:"path"`
	FileName        string `json:"fileName"`
}

type FileConfigsRequest struct {
	Email     string `json:"email"`
	MachineId string `json:"machineId"`
	Timestamp string `json:"timestamp"`
	Token     string `json:"token"`
}

type ConfigResponse struct {
	ID          uint64 `json:"id"`
	FileName    string `json:"fileName"`
	FileId      string `json:"fileId"`
	UpdateAt    string `json:"updateAt"`
	MachineId   string `json:"machineId"`
	MachineName string `json:"machineName"`
	Path        string `json:"path"`
	Attribute   string `json:"attribute"`
	CreatedAt   string `json:"createdAt"`
}

type FileCheckRequest struct {
	Email  string `json:"email"`
	FileId string `json:"fileId"`
	Sha256 string `json:"sha256"`
}

type FileSyncRequest struct {
	Email     string `json:"email"`
	MachineId string `json:"machineId"`
	Timestamp string `json:"timestamp"`
	FileId    string `json:"fileId"`
	Token     string `json:"token"`
	UpdateAt  int64  `json:"updateAt,omitempty"`
	Content   string `json:"content,omitempty"`
	Sha256    string `json:"sha256,omitempty"`
	FileName  string `json:"fileName,omitempty"`
}
