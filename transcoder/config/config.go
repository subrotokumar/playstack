package config

import (
	"gitlab.com/subrotokumar/glitchr/libs/core"
	"gitlab.com/subrotokumar/glitchr/libs/storage"
)

type Config struct {
	App struct {
		Name string   `yaml:"name" envconfig:"SERVICE_NAME" default:"glitchr"`
		Env  core.Env `yaml:"env" envconfig:"SERVICE_ENV" default:"dev"`
	} `yaml:"app"`
	Log struct {
		Level *string `yaml:"level" envconfig:"LOG_LEVEL" default:"INFO"`
	} `yaml:"log"`
	Database struct {
		Username string `yaml:"user" envconfig:"DB_USERNAME"`
		Password string `yaml:"pass" envconfig:"DB_PASSWORD"`
		Host     string `yaml:"host" envconfig:"DB_HOST"`
		Port     string `yaml:"port" envconfig:"DB_PORT"`
		DbName   string `yaml:"name" envconfig:"DB_NAME"`
	} `yaml:"database"`
	Aws struct {
		Region          string `yaml:"region" envconfig:"AWS_REGION" default:"ap-south-1"`
		AccessKeyID     string `yaml:"secret_id" envconfig:"AWS_ACCESS_KEY_ID"`
		SecretAccessKey string `yaml:"secret_key" envconfig:"AWS_SECRET_ACCESS_KEY"`
	} `yaml:"aws"`
	S3 struct {
		Bucket string `yaml:"bucket" envconfig:"S3_BUCKET"`
		Key    string `yaml:"key" envconfig:"S3_KEY"`
	} `yaml:"s3"`
	Events storage.S3Event `yaml:"events" envconfig:"SQS_MESSAGE"`
}
