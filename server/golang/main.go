package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	r := gin.Default()

	// CORS
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Next()
	})

	// Connect to the database
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("MYSQL_USER"),
		os.Getenv("MYSQL_PASS"),
		os.Getenv("MYSQL_HOST"),
		os.Getenv("MYSQL_PORT"),
		os.Getenv("MYSQL_NAME"),
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		return
	}

	// Migrate the schema
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Device{})
	db.AutoMigrate(&Config{})
	db.AutoMigrate(&File{})
	db.AutoMigrate(&Log{})

	// Create the DeviceAddHandler
	deviceAddHandler := DeviceAddHandler{DB: db}

	// Create the DeviceListHandler
	deviceListHandler := DeviceListHandler{DB: db}

	// Create the DeviceRemoveHandler
	deviceRemoveHandler := DeviceRemoveHandler{DB: db}

	// Create the FileConfigHandler
	fileConfigHandler := FileConfigHandler{DB: db}

	// Create the FileConfigsHandler
	fileConfigsHandler := FileConfigsHandler{DB: db}

	// Create the FileCheckHandler
	fileCheckHandler := FileCheckHandler{DB: db}

	// Create the FileSyncHandler
	fileSyncHandler := FileSyncHandler{DB: db}

	r.GET("/", HomeHandler)
	r.POST("/device/add", deviceAddHandler.postXhr)
	r.POST("/device/list", deviceListHandler.postXhr)
	r.POST("/device/remove", deviceRemoveHandler.postXhr)
	r.POST("/file/configs", fileConfigsHandler.postXhr)
	r.POST("/file/config", fileConfigHandler.postXhr)
	r.POST("/file/check", fileCheckHandler.postXhr)
	r.POST("/file/sync", fileSyncHandler.postXhr)

	r.Run(":8080")
}
