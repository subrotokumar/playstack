package storage

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Storage struct {
	client        *s3.Client
	presignClient *s3.PresignClient
}

func NewStorageProvider(region string) *Storage {
	ctx := context.Background()
	sdkConfig, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		fmt.Println("Couldn't load default configuration. Have you set up your AWS account?")
		log.Fatal(err)
	}

	client := s3.NewFromConfig(sdkConfig)
	presignClient := s3.NewPresignClient(client)
	return &Storage{
		client:        client,
		presignClient: presignClient,
	}
}
