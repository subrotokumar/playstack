package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gitlab.com/subrotokumar/glitchr/consumer/config"
	"gitlab.com/subrotokumar/glitchr/libs/core"
	"gitlab.com/subrotokumar/glitchr/libs/queue"
)

func main() {
	cfg := config.Config{}
	err := core.ConfigFromEnv(&cfg)
	if err != nil {
		panic(err)
	}
	log := core.NewLogger(cfg.App.Env, cfg.App.Name, cfg.Log.Level)
	log.Info("Raw uploaded video consumer started")

	q := queue.NewMessageQueue(cfg.Aws.Region, log)
	queueUrl := cfg.Sqs.Url
	maxMessages := cfg.Sqs.MaxMessages
	maxWait := cfg.Sqs.MaxWait
	emptyQueueSleep := cfg.Sqs.EmptyQueueSleep

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			log.Info("Shutting down consumer gracefully")
			return
		default:
			messages, err := q.GetMessages(ctx, queueUrl, maxMessages, maxWait)
			if err != nil {
				log.Error("Failed to get messages", "error", err)
				time.Sleep(2 * time.Second)
				continue
			}

			if len(messages) == 0 {
				time.Sleep(time.Second * time.Duration(emptyQueueSleep))
				continue
			}

			for _, msg := range messages {
				log.Info("Processing message", "id", *msg.MessageId, "body", *msg.Body)

				err := q.DeleteMessage(ctx, queueUrl, *msg.ReceiptHandle)
				if err != nil {
					log.Error("Failed to delete message", "id", *msg.MessageId, "error", err)
				}
			}
		}
	}
}
