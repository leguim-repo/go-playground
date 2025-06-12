package minio_demo

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

// InitMinioClient initializes and returns a new MinIO client
func InitMinioClient(config MinioConfig) (*minio.Client, error) {
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

// GetObjectContent retrieves the content of an object from MinIO bucket and returns it as a string
func GetObjectContent(ctx context.Context, client *minio.Client, bucketName, objectName string) (string, error) {
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

// DownloadObject downloads an object from MinIO bucket to a local file
func DownloadObject(ctx context.Context, client *minio.Client, bucketName, objectName, localPath string) error {
	// Download object to local file
	err := client.FGetObject(ctx, bucketName, objectName, localPath, minio.GetObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to download object %s to %s: %w", objectName, localPath, err)
	}

	return nil
}

func PrintListBuckets(ctx context.Context, client *minio.Client) {
	buckets, err := client.ListBuckets(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	for _, bucket := range buckets {
		log.Printf("Bucket: %s\n", bucket.Name)
	}
}

func PrintListObjects(client *minio.Client, bucketName string) {
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
