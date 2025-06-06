package main

import (
	"fmt"
	"log"

	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func main() {
	accessKey := "pepe"
	secretKey := "pepa"
	region := "eu-central-1"

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
	)
	if err != nil {
		log.Fatalf("Error al configurar el cliente AWS: %v", err)
	}

	client := s3.NewFromConfig(cfg)

	result, err := client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	if err != nil {
		log.Fatalf("Error al listar buckets: %v", err)
	}

	fmt.Println("Buckets encontrados:")
	for _, bucket := range result.Buckets {
		fmt.Printf("- %s (creado: %s)\n", *bucket.Name, bucket.CreationDate)
	}
}
