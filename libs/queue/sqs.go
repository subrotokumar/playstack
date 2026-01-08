package queue

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"gitlab.com/subrotokumar/glitchr/libs/core"
)

type Queue struct {
	SqsClient *sqs.Client
	log       core.Logger
}

func NewMessageQueue(region string, log *core.Logger) *Queue {
	ctx := context.Background()
	sdkConfig, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		log.Info("Couldn't load default configuration. Have you set up your AWS account?")
		log.Error(err.Error())
	}
	sqsClient := sqs.NewFromConfig(sdkConfig)
	return &Queue{
		SqsClient: sqsClient,
	}
}
