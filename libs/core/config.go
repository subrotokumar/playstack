package core

import (
	"os"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v2"
)

func ConfigFromFile(cfg any, path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	if err = decoder.Decode(&cfg); err != nil {
		return err
	}
	return nil
}

func ConfigFromEnv(cfg any) error {
	return envconfig.Process("", cfg)
}
