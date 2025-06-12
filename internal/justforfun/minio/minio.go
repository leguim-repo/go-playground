package main

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
)

func printListBuckets(client *minio.Client) {
	buckets, err := client.ListBuckets(context.Background())
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
	endpoint := "localhost:8090"
	accessKeyID := "admin"
	secretAccessKey := "password"
	useSSL := false

	ctx := context.Background()

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("%#v\n", minioClient) // minioClient is now set up

	printListBuckets(minioClient)

	// Make a new bucket called testbucket.
	bucketName := "testbucket"
	location := "eu-central-1"

	err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: location})
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := minioClient.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			log.Printf("We already own %s\n", bucketName)
		} else {
			log.Fatalln(err)
		}
	} else {
		log.Printf("Successfully created %s\n", bucketName)
	}

	printListBuckets(minioClient)
	printListObjects(minioClient, "raw")
}
