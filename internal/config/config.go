package config

import (
	"os"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v2"
)

type Env string

const (
	EnvDevelopment Env = "development"
	EnvStaging     Env = "staging"
	EnvProduction  Env = "production"
)

type Config struct {
	Service struct {
		Name string `yaml:"name" envconfig:"SERVICE_NAME" default:"glitchr"`
		Port string `yaml:"port" envconfig:"SERVICE_PORT" default:"8080"`
		Host string `yaml:"host" envconfig:"SERVICE_HOST" default:"0.0.0.0"`
		Env  Env    `yaml:"env" envconfig:"SERVICE_ENV" default:"development"`
	} `yaml:"service"`
	Database struct {
		Username string `yaml:"user" envconfig:"DB_USERNAME"`
		Password string `yaml:"pass" envconfig:"DB_PASSWORD"`
		Host     string `yaml:"host" envconfig:"DB_HOST"`
		Port     string `yaml:"port" envconfig:"DB_PORT"`
		DbName   string `yaml:"name" envconfig:"DB_NAME"`
	} `yaml:"database"`
}

func ConfigFromFile(path string) (cfg Config, err error) {
	f, err := os.Open(path)
	if err != nil {
		return cfg, err
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	if err = decoder.Decode(&cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}

func ConfigFromEnv() (cfg Config, err error) {
	if err = envconfig.Process("", &cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}
