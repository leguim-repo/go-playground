package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"

	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func loadAwsCredentials() (accessKey string, secretKey string, region string) {
	err := godotenv.Load("./internal/justforfun/aws/.env")
	if err != nil {
		log.Printf("Error loading .env: %v\n", err)
		return
	}

	accessKey = os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey = os.Getenv("AWS_SECRET_ACCESS_KEY")
	region = os.Getenv("AWS_DEFAULT_REGION")

	fmt.Printf("API Key: %s secret: %s region: %s\n", accessKey, secretKey, region)
	return accessKey, secretKey, region
}

func main() {
	accessKey, secretKey, region := loadAwsCredentials()

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
	)
	if err != nil {
		log.Fatalf("Error configuring AWS client: %v", err)
	}

	client := s3.NewFromConfig(cfg)

	result, err := client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	if err != nil {
		log.Fatalf("Error listing buckets: %v", err)
	}

	fmt.Println("Buckets found:")
	for _, bucket := range result.Buckets {
		fmt.Printf("- %s (created: %s)\n", *bucket.Name, bucket.CreationDate)
	}
}
