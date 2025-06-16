package handlers

import (
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/minio/minio-go/v7"
)

func DownloadHandler(client *minio.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		filename := vars["filename"]

		obj, err := client.GetObject(
			r.Context(),
			"uploads",
			filename,
			minio.GetObjectOptions{},
		)
		if err != nil {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}
		defer obj.Close()

		stat, err := obj.Stat()
		if err == nil {
			w.Header().Set("Content-Type", stat.ContentType)
		}

		if _, err := io.Copy(w, obj); err != nil {
			http.Error(w, "Download failed", http.StatusInternalServerError)
		}
	}
}
