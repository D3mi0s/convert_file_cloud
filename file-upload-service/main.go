package main

import (
	"file-upload-service/handlers"
	"file-upload-service/models"
	"file-upload-service/queue"
	"file-upload-service/repository"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main() {

	checkDB, err := gorm.Open(postgres.Open("host=localhost user=postgres password=admin dbname=postgres port=5432 sslmode=disable"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to default database:", err)
	}

	var dbExists bool
	err = checkDB.Raw(
		"SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = 'upload_service')",
	).Scan(&dbExists).Error

	if err != nil {
		log.Fatal("Failed to check database existence:", err)
	}

	if !dbExists {
		log.Println("Database 'upload_service' not found, creating...")
		err = checkDB.Exec("CREATE DATABASE upload_service").Error
		if err != nil {
			log.Fatal("Failed to create database:", err)
		}
		log.Println("Database 'upload_service' created successfully")
	}

	sqlDB, _ := checkDB.DB()
	sqlDB.Close()

	dsn := "host=localhost user=postgres password=admin dbname=upload_service port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to upload_service database:", err)
	}

	if err := db.AutoMigrate(&models.File{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	rabbitMQ, err := queue.NewRabbitMQ("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ:", err)
	}
	defer rabbitMQ.Conn.Close()

	fileRepo := &repository.FileRepository{DB: db}
	uploadHandler := &handlers.UploadHandler{
		FileRepo: fileRepo,
		Queue:    rabbitMQ,
	}

	r := gin.Default()
	r.Use(CORSMiddleware())
	r.POST("/upload", uploadHandler.UploadFile)
	r.GET("/files", uploadHandler.GetFiles)
	log.Println("File Upload Service started on :8081")
	r.Run(":8081")
}
