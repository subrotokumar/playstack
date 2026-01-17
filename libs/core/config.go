package core

import (
	"os"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v2"
)

func ConfigFromFile(cfg any, path string) (err error) {
	f, err := os.Open(path)
	if err != nil {
		return err
	}

	defer func() {
		if cerr := f.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(cfg)
	return
}

func ConfigFromEnv(cfg any) error {
	return envconfig.Process("", cfg)
}
