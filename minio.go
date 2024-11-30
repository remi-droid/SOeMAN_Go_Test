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

// To upload a document to the minio storage
func UploadDocument(fileName string, fileData []byte) error {
	ctx := context.Background()

	// Create a reader from the file data.
	fileReader := bytes.NewReader(fileData)

	_, err := minioClient.PutObject(ctx, bucketName, fileName, fileReader, int64(len(fileData)), minio.PutObjectOptions{
		ContentType: "application/octet-stream",
	})
	if err != nil {
		return err
	}

	log.Printf("Document %s uploaded successfully", fileName)
	return nil
}

// To download a document from the minio storage
func DownloadDocument(fileName string) (*minio.Object, error) {
	ctx := context.Background()
	return minioClient.GetObject(ctx, bucketName, fileName, minio.GetObjectOptions{})
}

// Remove all the documents from the bucket
func ClearMinioStorage() (int, error) {

	deletedCount := 0
	ctx := context.Background()

	// Get the list of all the documents in the bucket
	documents := minioClient.ListObjects(ctx, bucketName, minio.ListObjectsOptions{
		Recursive: true,
	})

	// Delete all the documents one by one
	for document := range documents {
		if document.Err != nil {
			return deletedCount, fmt.Errorf("failed to list document %s: %w", document.Key, document.Err)
		}
		if err := minioClient.RemoveObject(ctx, bucketName, document.Key, minio.RemoveObjectOptions{}); err != nil {
			return deletedCount, fmt.Errorf("failed to delete document %s: %w", document.Key, err)
		}
		deletedCount++
	}

	return deletedCount, nil
}
