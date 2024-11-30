package main

import (
	"bytes"
	"context"
	"fmt"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var minioClient *minio.Client
var bucketName = "uploads"

func InitMinioStorage() error {
	var err error
	minioClient, err = minio.New("minio:9000", &minio.Options{
		Creds:  credentials.NewStaticV4("root", "password", ""),
		Secure: false, // For HTTP only
	})
	if err != nil {
		return fmt.Errorf("failed to initialize MinIO client: %w", err)
	}

	// Check if the bucket exists or create it.
	ctx := context.Background()
	exists, errBucketExists := minioClient.BucketExists(ctx, bucketName)
	if errBucketExists != nil {
		return fmt.Errorf("failed to check bucket: %w", errBucketExists)
	}
	if !exists {
		err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return fmt.Errorf("failed to create bucket: %w", err)
		}
		log.Printf("Bucket %s created successfully", bucketName)
	}

	return nil
}

// To upload a file to the minio storage
func UploadFile(fileName string, fileData []byte) error {
	ctx := context.Background()

	// Create a reader from the file data.
	fileReader := bytes.NewReader(fileData)

	_, err := minioClient.PutObject(ctx, bucketName, fileName, fileReader, int64(len(fileData)), minio.PutObjectOptions{
		ContentType: "application/octet-stream",
	})
	if err != nil {
		return err
	}

	log.Printf("File %s uploaded successfully", fileName)
	return nil
}
