package config

import (
	"encoding/json"
	"strings"

	"gitlab.com/subrotokumar/playstack/libs/core"
	"gitlab.com/subrotokumar/playstack/libs/storage"
)

type Config struct {
	App struct {
		Name string   `yaml:"name" envconfig:"SERVICE_NAME" default:"playstack-transcoder"`
		Env  core.Env `yaml:"env" envconfig:"SERVICE_ENV" default:"dev"`
	} `yaml:"app"`
	Log struct {
		Level *string `yaml:"level" envconfig:"LOG_LEVEL" default:"INFO"`
	} `yaml:"log"`
	Aws struct {
		Region          string `yaml:"region" envconfig:"AWS_REGION" default:"ap-south-1"`
		AccessKeyID     string `yaml:"secret_id" envconfig:"AWS_ACCESS_KEY_ID"`
		SecretAccessKey string `yaml:"secret_key" envconfig:"AWS_SECRET_ACCESS_KEY"`
		MediaBucket     string `yaml:"media_bucket" envconfig:"MEDIA_BUCKET" required:"true"`
	} `yaml:"aws"`
	NotifierService struct {
		URL      string `yaml:"api" envconfig:"NOTIFIER_SERVICE_ENDPOINT" default:"http://localhost:8080"`
		Username string `yaml:"username" envconfig:"BASIC_AUTH_USERNAME"`
		PASSWORD string `yaml:"password" envconfig:"BASIC_AUTH_PASSWORD"`
	} `yaml:"notifier_service"`
	Event   string `yaml:"events" envconfig:"SQS_MESSAGE" required:"true"`
	S3Event storage.S3Event
}

func (cfg *Config) ParseS3Event() (event storage.S3Event, err error) {
	err = json.Unmarshal([]byte(cfg.Event), &event)
	return event, err
}

func (cfg *Config) Bucket() string {
	return cfg.S3Event.Records[0].S3.Bucket.Name
}

func (cfg *Config) Key() string {
	return cfg.S3Event.Records[0].S3.Object.Key
}

func (cfg *Config) UserAndVideoID() (string, string) {
	keys := strings.Split(cfg.Key(), "/")
	return keys[1], keys[2]
}

func (cfg *Config) ObjectSize() int64 {
	return cfg.S3Event.Records[0].S3.Object.Size
}
