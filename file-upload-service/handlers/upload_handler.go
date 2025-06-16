package handlers

import (
	"bytes"
	"encoding/json"
	"file-upload-service/models"
	"file-upload-service/queue"
	"file-upload-service/repository"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func (h *UploadHandler) GetFiles(c *gin.Context) {
	userID, err := getUserIdFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	files, err := h.FileRepo.GetFilesByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	c.JSON(http.StatusOK, files)
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

type UploadHandler struct {
	FileRepo *repository.FileRepository
	Queue    *queue.RabbitMQ
}

func getUserIdFromToken(c *gin.Context) (uint, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return 0, fmt.Errorf("authorization header is missing")
	}

	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || strings.ToLower(tokenParts[0]) != "bearer" {
		return 0, fmt.Errorf("invalid authorization header format")
	}
	tokenString := tokenParts[1]

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})
	if err != nil || !token.Valid {
		return 0, fmt.Errorf("invalid token: %v", err)
	}

	claims := token.Claims.(jwt.MapClaims)
	userID := uint(claims["sub"].(float64))
	return userID, nil
}

func (h *UploadHandler) UploadFile(c *gin.Context) {
	log.Println("Attempting to upload file...")

	userID, err := getUserIdFromToken(c)
	if err != nil {
		log.Printf("Auth error: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		log.Printf("File error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	mimeType := fileHeader.Header.Get("Content-Type")
	if mimeType == "application/octet-stream" {

		switch strings.ToLower(filepath.Ext(fileHeader.Filename)) {
		case ".xlsx":
			mimeType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
		case ".xls":
			mimeType = "application/vnd.ms-excel"
		case ".doc":
			mimeType = "application/msword"
		case ".pptx":
			mimeType = "application/vnd.openxmlformats-officedocument.presentationml.presentation"
		case ".ppt":
			mimeType = "application/vnd.ms-powerpoint"
		case ".rtf":
			mimeType = "text/rtf"
		case ".dcm", ".dicom":
			mimeType = "application/dicom"
		}
	}

	allowedTypes := []string{
		"image/png",
		"image/jpeg",
		"image/bmp",
		"image/x-bmp",
		"image/x-ms-bmp",
		"text/plain",
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		"application/vnd.ms-excel",
		"application/msword",
		"application/vnd.openxmlformats-officedocument.presentationml.presentation",
		"application/vnd.ms-powerpoint",
		"text/rtf",
		"application/rtf",
		"application/dicom",
		"application/octet-stream",
	}

	if !contains(allowedTypes, mimeType) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported file type"})
		return
	}

	log.Printf("Received file: %s (Size: %d)", fileHeader.Filename, fileHeader.Size)

	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
		return
	}
	defer file.Close()

	fileContent, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
		return
	}

	storageURL := "http://localhost:8082/upload"
	req, err := http.NewRequest("POST", storageURL, bytes.NewReader(fileContent))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}
	req.Header.Set("Content-Type", mimeType)

	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusCreated {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Storage service error"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("Storage service error (HTTP %d): %s", resp.StatusCode, string(body))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Storage service failed"})
		return
	}

	var result struct {
		Filename string `json:"filename"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid storage response"})
		return
	}
	storedName := result.Filename

	newFile := &models.File{
		UserID:       userID,
		OriginalName: fileHeader.Filename,
		StoredName:   storedName,
		Size:         fileHeader.Size,
		MimeType:     mimeType,
	}

	if err := h.FileRepo.CreateFile(newFile); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file info"})
		return
	}

	msg := fmt.Sprintf(`{"file_id": %d}`, newFile.ID)
	if err := h.Queue.Publish("conversion_queue", msg); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to queue conversion"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "File uploaded and queued for conversion",
		"file_id": newFile.ID,
	})
}
