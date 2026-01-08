package service

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"gitlab.com/subrotokumar/glitchr/libs/core"
	"gitlab.com/subrotokumar/glitchr/libs/storage"
	"gitlab.com/subrotokumar/glitchr/transcoder/config"
	"gitlab.com/subrotokumar/glitchr/transcoder/ffmpeg"
)

type Service struct {
	cfg     config.Config
	log     *core.Logger
	storage *storage.Storage
}

func New() *Service {
	cfg := config.Config{}
	if err := core.ConfigFromEnv(&cfg); err != nil {
		panic(err)
	}
	log := core.NewLogger(cfg.App.Env, cfg.App.Name, cfg.Log.Level)
	storage := storage.NewStorageProvider(cfg.Aws.Region)
	return &Service{
		cfg:     cfg,
		log:     log,
		storage: storage,
	}
}

func (s *Service) Run() error {
	s.log.Info("Transcorder service started")
	// Implement the core logic of the transcorder service here
	return nil
}

func (s *Service) Download(ctx context.Context, destPath string) error {
	s.log.Info("Downloading file", "path", destPath)

	if err := os.MkdirAll(filepath.Dir(destPath), 0o755); err != nil {
		return fmt.Errorf("create dir: %w", err)
	}

	out, err := s.storage.Client().GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.cfg.S3.Bucket),
		Key:    aws.String(s.cfg.S3.Key),
	})
	if err != nil {
		return fmt.Errorf("get object failed: %w", err)
	}
	defer out.Body.Close()

	file, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer file.Close()

	_, err = io.Copy(file, out.Body)
	if err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	return nil
}

func (s *Service) Transcode(ctx context.Context, inputPath, outputDir string) error {
	cmdArgs := ffmpeg.DASH_CMD(inputPath, outputDir)
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	return cmd.Run()
}
