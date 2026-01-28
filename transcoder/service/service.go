package service

import (
	"gitlab.com/subrotokumar/playstack/libs/core"
	"gitlab.com/subrotokumar/playstack/libs/storage"
	"gitlab.com/subrotokumar/playstack/transcoder/config"
)

type Service struct {
	cfg     config.Config
	log     *core.Logger
	storage *storage.Storage
	bucket  string
	path    string
}

func New() *Service {
	cfg := config.Config{}
	if err := core.ConfigFromEnv(&cfg); err != nil {
		panic(err)
	}
	log := core.NewLogger(cfg.App.Env, cfg.App.Name, cfg.Log.Level)
	s3Event, err := cfg.ParseS3Event()
	if err != nil {
		log.Fatal("failed to unmarshell SQS_MESSAGE")
	}
	cfg.S3Event = s3Event
	storage := storage.NewStorageProvider(cfg.Aws.Region)
	return &Service{
		cfg:     cfg,
		log:     log,
		storage: storage,
	}
}
