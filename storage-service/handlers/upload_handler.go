package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/minio/minio-go/v7"
)

func UploadHandler(client *minio.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Upload request received")

		if r.Method != "POST" {
			log.Printf("Invalid method: %s", r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		contentType := r.Header.Get("Content-Type")
		log.Printf("Content-Type: %s", contentType)

		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading body: %v", err)
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		log.Printf("Received %d bytes", len(body))

		storedName := fmt.Sprintf("%d_%d", time.Now().UnixNano(), len(body))
		log.Printf("Generated filename: %s", storedName)

		_, err = client.PutObject(
			r.Context(),
			"uploads",
			storedName,
			bytes.NewReader(body),
			int64(len(body)),
			minio.PutObjectOptions{ContentType: contentType},
		)

		if err != nil {
			log.Printf("MinIO upload error: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		log.Printf("File %s uploaded successfully (%d bytes)", storedName, len(body))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{
			"status":   "success",
			"filename": storedName,
		})
	}
}
