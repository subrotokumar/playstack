package config

import (
	"gitlab.com/subrotokumar/glitchr/pkg/core"
)

type Config struct {
	App struct {
		Name string   `yaml:"name" envconfig:"SERVICE_NAME" default:"glitchr"`
		Port string   `yaml:"port" envconfig:"SERVICE_PORT" default:"8080"`
		Host string   `yaml:"host" envconfig:"SERVICE_HOST" default:"0.0.0.0"`
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
		Region string `yaml:"region" envconfig:"AWS_REGION" default:"ap-south-1"`
	} `yaml:"aws"`
	Cognito struct {
		ClientID     string `yaml:"client_id" envconfig:"COGNITO_CLIENT_ID"`
		ClientSecret string `yaml:"client_secret" envconfig:"COGNITO_CLIENT_SECRET"`
	} `yaml:"cognito"`
}
