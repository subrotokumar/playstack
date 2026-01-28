package service

import (
	"context"
)

func (s *Service) UpdateMetadata(ctx context.Context) error {
	s.log.Info("Updating video metadata in database")
	return nil
}
