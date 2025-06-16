package main

import (
	"log"
	"net/http"

	"download-service/handlers"
	"download-service/storage"

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
	r.HandleFunc("/download/{filename}", handlers.DownloadHandler(minioClient))

	finalRouter := CORSMiddleware(r)

	log.Println("Download Service started on :8083")
	log.Fatal(http.ListenAndServe(":8083", finalRouter))
}
