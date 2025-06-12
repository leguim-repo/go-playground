package main

import (
	"context"
	"fmt"
	"go-playground/internal/justforfun/minio_demo"
	"log"
)

func main() {
	// Define MinIO configuration
	config := minio_demo.MinioConfig{
		Endpoint:        "localhost:8090",
		AccessKeyID:     "admin",
		SecretAccessKey: "password",
		UseSSL:          false,
	}

	// Initialize MinIO client
	minioClient, err := minio_demo.InitMinioClient(config)
	if err != nil {
		log.Fatalf("Error initializing MinIO client: %v", err)
	}

	ctx := context.Background()
	bucketName := "raw"
	objectName := "navy_seals_inventory.json"

	minio_demo.PrintListBuckets(ctx, minioClient)
	minio_demo.PrintListObjects(minioClient, bucketName)

	// Get object content example
	content, err := minio_demo.GetObjectContent(ctx, minioClient, bucketName, objectName)
	if err != nil {
		log.Fatalf("Error getting object content: %v", err)
	}
	fmt.Println("Object content:", content)

	// Download object example
	localPath := "./datalake/downloaded-navy_seals_inventory.json"
	err = minio_demo.DownloadObject(ctx, minioClient, bucketName, objectName, localPath)
	if err != nil {
		log.Fatalf("Error downloading object: %v", err)
	}
	fmt.Printf("Object downloaded successfully to %s\n", localPath)

}
