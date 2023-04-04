package main

import "time"

type User struct {
	ID         uint   `gorm:"primaryKey"`
	Email      string `gorm:"unique"`
	EmailSha1  string
	Verify     string
	PublicKey  string
	PrivateKey string
	CreatedAt  time.Time
}

type Device struct {
	ID          uint `gorm:"primaryKey"`
	Email       string
	MachineId   string `gorm:"unique"`
	MachineName string
	MachineKey  string
	CreatedAt   time.Time
}

type Config struct {
	ID        uint64     `gorm:"primaryKey;autoIncrement;column:id"`
	Email     string     `gorm:"column:email"`
	MachineId string     `gorm:"column:machineId"`
	FileId    string     `gorm:"column:fileId"`
	FileName  string     `gorm:"column:fileName"`
	Path      string     `gorm:"column:path"`
	Attribute string     `gorm:"column:attribute"`
	DeletedAt *time.Time `gorm:"column:deletedAt"`
	CreatedAt time.Time  `gorm:"column:createdAt"`
}

type File struct {
	ID            uint64    `gorm:"primaryKey;autoIncrement;column:id"`
	Email         string    `gorm:"column:email"`
	EmailSha1     string    `gorm:"column:emailSha1"`
	FileId        string    `gorm:"column:fileId"`
	FileName      string    `gorm:"column:fileName"`
	Content       string    `gorm:"column:content;type:longtext"`
	SHA256        string    `gorm:"column:sha256"`
	FromMachineId string    `gorm:"column:fromMachineId"`
	UpdateAt      time.Time `gorm:"column:updateAt"`
}

type Log struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement;column:id"`
	Email     string    `gorm:"column:email"`
	MachineId string    `gorm:"column:machineId"`
	Action    string    `gorm:"column:action"`
	Content   string    `gorm:"column:content"`
	CreatedAt time.Time `gorm:"column:createdAt"`
}
