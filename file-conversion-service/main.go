package main

import (
	"bytes"
	"encoding/json"
	"file-conversion-service/converter"
	"file-conversion-service/repository"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rabbitmq/amqp091-go"
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

	dsn := "host=localhost user=postgres password=admin dbname=upload_service port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	fileRepo := repository.NewFileRepository(db)

	conv := converter.NewConverter()

	rabbitMQ, err := NewRabbitMQ("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal("RabbitMQ connection failed:", err)
	}
	defer rabbitMQ.Close()

	log.Println("Waiting for conversion tasks...")
	rabbitMQ.Consume("conversion_queue", func(msg string) {
		HandleConversion(msg, fileRepo, conv)
	})
}

func HandleConversion(msg string, repo *repository.FileRepository, conv *converter.FileConverter) {
	var task struct {
		FileID uint `json:"file_id"`
	}

	if err := json.Unmarshal([]byte(msg), &task); err != nil {
		log.Printf("Error decoding message: %v", err)
		return
	}

	file, err := repo.GetFileByID(task.FileID)
	if err != nil {
		log.Printf("File not found: %d. Error: %v", task.FileID, err)
		return
	}

	downloadURL := fmt.Sprintf("http://localhost:8083/download/%s", file.StoredName)
	log.Printf("Downloading from: %s", downloadURL)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(downloadURL)
	if err != nil {
		log.Printf("Download failed: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Download error: HTTP %d", resp.StatusCode)
		return
	}

	fileData, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read file data: %v", err)
		return
	}

	pdfData, err := conv.ConvertFile(fileData, file.MimeType)
	if err != nil {
		log.Printf("Conversion failed: %v", err)
		return
	}

	uploadURL := "http://localhost:8082/upload"
	req, err := http.NewRequest("POST", uploadURL, bytes.NewReader(pdfData))
	if err != nil {
		log.Printf("Failed to create upload request: %v", err)
		return
	}
	req.Header.Set("Content-Type", "application/pdf")

	resp, err = client.Do(req)
	if err != nil {
		log.Printf("Upload failed: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		log.Printf("Upload error: HTTP %d", resp.StatusCode)
		return
	}

	var result struct {
		Filename string `json:"filename"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("Failed to parse response: %v", err)
		return
	}

	if err := repo.UpdateFileStatus(file.ID, "completed"); err != nil {
		log.Printf("Failed to update status: %v", err)
	}

	log.Printf("Successfully processed file ID %d", file.ID)

	if err := repo.UpdateConvertedName(file.ID, result.Filename); err != nil {
		log.Printf("Failed to update converted name: %v", err)
	}
}

type RabbitMQ struct {
	conn *amqp091.Connection
	ch   *amqp091.Channel
}

func NewRabbitMQ(uri string) (*RabbitMQ, error) {
	conn, err := amqp091.Dial(uri)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &RabbitMQ{conn: conn, ch: ch}, nil
}

func (rmq *RabbitMQ) Consume(queue string, handler func(string)) {
	_, err := rmq.ch.QueueDeclare(
		queue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal("Failed to declare queue:", err)
	}

	msgs, err := rmq.ch.Consume(
		queue,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal("Failed to register consumer:", err)
	}

	for msg := range msgs {
		handler(string(msg.Body))
	}
}

func (rmq *RabbitMQ) Close() {
	rmq.ch.Close()
	rmq.conn.Close()
}
