package config

import (
	"fmt"

	"gitlab.com/subrotokumar/playstack/libs/core"
)

type Config struct {
	App struct {
		Name string   `yaml:"name" envconfig:"SERVICE_NAME" default:"playstack-backend"`
		Port string   `yaml:"port" envconfig:"SERVICE_PORT" default:"8080"`
		Host string   `yaml:"host" envconfig:"SERVICE_HOST" default:"0.0.0.0"`
		Env  core.Env `yaml:"env" envconfig:"SERVICE_ENV" default:"dev"`
	} `yaml:"app"`
	BasicAuth struct {
		Username string `yaml:"username" envconfig:"BASIC_AUTH_USERNAME" default:"admin"`
		Password string `yaml:"password" envconfig:"BASIC_AUTH_PASSWORD" default:"password"`
	} `yaml:"basic_auth"`
	Log struct {
		Level *string `yaml:"level" envconfig:"LOG_LEVEL" default:"INFO"`
	} `yaml:"log"`
	Database struct {
		Username string `yaml:"user" envconfig:"DB_USERNAME" required:"true"`
		Password string `yaml:"pass" envconfig:"DB_PASSWORD" required:"true"`
		Host     string `yaml:"host" envconfig:"DB_HOST" required:"true"`
		Port     string `yaml:"port" envconfig:"DB_PORT" default:"5432" required:"true"`
		DbName   string `yaml:"name" envconfig:"DB_NAME" required:"true"`
		SslMode  string `yaml:"ssl_mode" envconfig:"DB_SSL_MODE"`
		MaxConn  int32  `yaml:"max_conn" envconfig:"DB_MAX_CONN" default:"10"`
		MinConn  int32  `yaml:"min_conn" envconfig:"DB_MIN_CONN" default:"2"`
	} `yaml:"database"`
	Aws struct {
		Region string `yaml:"region" envconfig:"AWS_REGION" required:"true"`
	} `yaml:"aws"`
	Cognito struct {
		ClientID     string `yaml:"client_id" envconfig:"COGNITO_CLIENT_ID" required:"true"`
		ClientSecret string `yaml:"client_secret" envconfig:"COGNITO_CLIENT_SECRET" required:"true"`
		UserPoolID   string `yaml:"user_pool_id" envconfig:"COGNITO_USER_POOL_ID" required:"true"`
	} `yaml:"cognito"`
	S3 struct {
		RawMediaBucket string `yaml:"raw_media_bucket" envconfig:"RAW_MEDIA_BUCKET" required:"true"`
		MediaBucket    string `yaml:"media_bucket" envconfig:"MEDIA_BUCKET"`
	} `yaml:"s3"`
}

func (cfg Config) ConnectionUrl() string {
	conn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.Database.Username, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.DbName)
	if cfg.Database.SslMode != "" {
		conn += "?sslmode=require"
	}
	return conn
}
