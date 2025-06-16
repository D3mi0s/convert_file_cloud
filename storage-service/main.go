package main

import (
	"log"
	"net/http"

	"storage-service/handlers"
	"storage-service/storage"

	"github.com/gorilla/mux"
)

func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	minioClient, err := storage.NewMinIOClient()
	if err != nil {
		log.Fatal("MinIO init failed:", err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/upload", handlers.UploadHandler(minioClient))

	loggedRouter := handlers.LoggingMiddleware(CORSMiddleware(r))

	log.Println("Storage Service (Upload) started on :8082")
	log.Fatal(http.ListenAndServe(":8082", loggedRouter))
}
