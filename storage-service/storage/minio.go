package storage

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func NewMinIOClient() (*minio.Client, error) {
	client, err := minio.New("localhost:9000", &minio.Options{
		Creds:  credentials.NewStaticV4("minioadmin", "minioadmin", ""),
		Secure: false,
	})
	if err != nil {
		return nil, fmt.Errorf("connection error: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = client.ListBuckets(ctx)
	if err != nil {
		return nil, fmt.Errorf("ping failed: %w", err)
	}

	log.Println("Successfully connected to MinIO")

	exists, err := client.BucketExists(ctx, "uploads")
	if !exists && err == nil {
		err = client.MakeBucket(ctx, "uploads", minio.MakeBucketOptions{})
		if err != nil {
			return nil, fmt.Errorf("bucket creation failed: %w", err)
		}
		log.Println("Created bucket 'uploads'")
	}

	return client, nil
}
