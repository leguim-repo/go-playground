package main

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"io"
	"log"
)

// MinioConfig holds the configuration for MinIO client
type MinioConfig struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	UseSSL          bool
}

// initMinioClient initializes and returns a new MinIO client
func initMinioClient(config MinioConfig) (*minio.Client, error) {
	// Initialize MinIO client object
	minioClient, err := minio.New(config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccessKeyID, config.SecretAccessKey, ""),
		Secure: config.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize MinIO client: %w", err)
	}

	return minioClient, nil
}

// getObjectContent retrieves the content of an object from MinIO bucket and returns it as a string
func getObjectContent(ctx context.Context, client *minio.Client, bucketName, objectName string) (string, error) {
	// Get object from bucket
	object, err := client.GetObject(ctx, bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to get object %s from bucket %s: %w", objectName, bucketName, err)
	}
	defer object.Close()

	// Read object content
	content, err := io.ReadAll(object)
	if err != nil {
		return "", fmt.Errorf("failed to read object content: %w", err)
	}

	return string(content), nil
}

// downloadObject downloads an object from MinIO bucket to a local file
func downloadObject(ctx context.Context, client *minio.Client, bucketName, objectName, localPath string) error {
	// Download object to local file
	err := client.FGetObject(ctx, bucketName, objectName, localPath, minio.GetObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to download object %s to %s: %w", objectName, localPath, err)
	}

	return nil
}

func printListBuckets(ctx context.Context, client *minio.Client) {
	buckets, err := client.ListBuckets(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	for _, bucket := range buckets {
		log.Printf("Bucket: %s\n", bucket.Name)
	}
}

func printListObjects(client *minio.Client, bucketName string) {
	ctx, cancel := context.WithCancel(context.Background())
	log.Println("Objects found in bucket: ", bucketName)
	defer cancel()

	objectCh := client.ListObjects(ctx, bucketName, minio.ListObjectsOptions{
		Recursive: true,
	})
	for object := range objectCh {
		if object.Err != nil {
			log.Println(object.Err)
			return
		}
		log.Println("* ", object.Key)
	}
}

func main() {
	// Define MinIO configuration
	config := MinioConfig{
		Endpoint:        "localhost:8090",
		AccessKeyID:     "admin",
		SecretAccessKey: "password",
		UseSSL:          false,
	}

	// Initialize MinIO client
	minioClient, err := initMinioClient(config)
	if err != nil {
		log.Fatalf("Error initializing MinIO client: %v", err)
	}

	ctx := context.Background()
	bucketName := "raw"
	objectName := "navy_seals_inventory.json"

	printListBuckets(ctx, minioClient)
	printListObjects(minioClient, bucketName)

	// Get object content example
	content, err := getObjectContent(ctx, minioClient, bucketName, objectName)
	if err != nil {
		log.Fatalf("Error getting object content: %v", err)
	}
	fmt.Println("Object content:", content)

	// Download object example
	localPath := "./datalake/downloaded-navy_seals_inventory.json"
	err = downloadObject(ctx, minioClient, bucketName, objectName, localPath)
	if err != nil {
		log.Fatalf("Error downloading object: %v", err)
	}
	fmt.Printf("Object downloaded successfully to %s\n", localPath)

}
