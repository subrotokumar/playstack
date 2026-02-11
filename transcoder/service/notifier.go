package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"gitlab.com/subrotokumar/playstack/libs/db"
)

type (
	UpdateMetadataRequest struct {
		Title       *string        `json:"title"`
		Status      db.VideoStatus `json:"status" validate:"omitempty,oneof='PREUPLOAD' 'UPLOADED' 'PROCESSING' 'READY' 'FAILED'"`
		DurationSec *int32         `json:"duration_sec"`
	}
)

func (s *Service) UpdateMetadata(ctx context.Context, request UpdateMetadataRequest) error {
	s.log.Info("Updating video metadata in database")

	userID, videoID := s.cfg.UserAndVideoID()

	url := s.cfg.NotifierService.URL + "/internal/media/videos/" + videoID
	payload := make(map[string]any)
	payload["user_id"] = userID
	payload["status"] = string(request.Status)

	if request.Title != nil {
		payload["title"] = *request.Title
	}
	if request.DurationSec != nil {
		payload["duration_sec"] = *request.DurationSec
	}

	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		s.log.Error("failed to marshal update payload", "err", err)
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "PATCH", url, bytes.NewReader(bodyBytes))
	if err != nil {
		s.log.Error("failed to build request", "err", err)
		return err
	}

	username := s.cfg.NotifierService.Username
	password := s.cfg.NotifierService.PASSWORD
	req.Header.Set("Content-Type", "application/json")
	if username != "" || password != "" {
		req.SetBasicAuth(username, password)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		s.log.Error("notifier request failed", "err", err)
		return err
	}
	defer res.Body.Close()

	respBody, _ := io.ReadAll(res.Body)
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		s.log.Error("notifier returned non-2xx", "status", res.StatusCode, "body", string(respBody))
		return fmt.Errorf("notifier returned status: %s", res.Status)
	}

	s.log.Info("Notifier updated metadata", "status", res.StatusCode)
	return nil
}
