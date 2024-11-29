package main

import (
	"fmt"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var minioClient *minio.Client
var bucketName = "uploads"

func InitMinioStorage() error {
	var err error
	minioClient, err = minio.New("minio:9000", &minio.Options{
		Creds:  credentials.NewStaticV4("root", "password", ""),
		Secure: false, // Use true for HTTPS
	})
	if err != nil {
		return fmt.Errorf("Failed to initialize MinIO client: %v", err)
	}

	return nil
}
