package config

import (
	"gitlab.com/subrotokumar/glitchr/pkg/core"
)

type Config struct {
	App struct {
		Name string   `yaml:"name" envconfig:"SERVICE_NAME" default:"glitchr_consumer"`
		Env  core.Env `yaml:"env" envconfig:"SERVICE_ENV" default:"dev"`
	} `yaml:"app"`
	Log struct {
		Level *string `yaml:"level" envconfig:"LOG_LEVEL" default:"INFO"`
	} `yaml:"log"`
	Aws struct {
		Region string `yaml:"region" envconfig:"AWS_REGION" default:"ap-south-1"`
	} `yaml:"aws"`
	Sqs struct {
		Url             string `yaml:"url" envconfig:"SQS_QUEUE_URL"`
		MaxMessages     int32  `yaml:"max_messages" envconfig:"SQS_MAX_MESSAGES" default:"1"`
		MaxWait         int32  `yaml:"max_wait" envconfig:"SQS_MAX_WAIT_TIME" default:"1"`
		EmptyQueueSleep int32  `yaml:"max_wait" envconfig:"SQS_MAX_WAIT_TIME" default:"1"`
	} `yaml:"sqs"`
}
