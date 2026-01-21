package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
	"gitlab.com/subrotokumar/glitchr/libs/db"
)

const (
	PUT_PRESIGNED_URL_TTL = 1 * 60
)

type (
	Asset struct {
		Id           uuid.UUID `json:"id"`
		Name         string    `json:"name"`
		Size         int       `json:"size"`
		ContentType  string    `json:"content_type"`
		Href         string    `json:"href"`
		OriginalName string    `json:"original_name"`
	}
	AssetsRequest struct {
		Name        string `json:"name" validate:"required"`
		Size        int64  `json:"size" validate:"required"`
		ContentType string `json:"content_type" validate:"required"`
	}
	AssetsResponseData struct {
		UploadUrl string            `json:"upload_url"`
		Header    map[string]string `json:"header"`
		Asset     Asset             `json:"asset"`
		Form      map[string]string `json:"form"`
	}
	AssetsResponse struct {
		Data    *AssetsResponseData `json:"data,omitempty"`
		Message string              `json:"message,omitempty"`
		Error   any                 `json:"error,omitempty"`
	}

	UpdateMetadataRequest struct {
		UserID      uuid.UUID       `json:"user_id"`
		Title       *string         `json:"title"`
		Status      *db.VideoStatus `json:"Status" validate:"omitempty,oneof='PREUPLOAD' 'UPLOADED' 'PROCESSING' 'READY' 'FAILED'"`
		DurationSec *int32          `json:"duration_sec"`
	}
)

func (s *Server) AssetsHandler(c echo.Context) error {
	body := AssetsRequest{}
	if err := RequestBody(c, &body); err != nil {
		return c.JSON(http.StatusInternalServerError, AssetsResponse{Error: err.Error()})
	}

	videoId := uuid.Must(uuid.NewV7())
	userId := c.Get("sub").(uuid.UUID)
	key := fmt.Sprintf("/%d/%s", userId, videoId)
	_, err := s.store.CreateVideo(c.Request().Context(), db.CreateVideoParams{
		ID:          videoId,
		UserID:      userId,
		Title:       body.Name,
		Status:      db.VideoStatusPREUPLOAD,
		DurationSec: pgtype.Int4{Valid: false},
	})
	if err != nil {
		s.log.Error("Failed to create video record", "err", err)
		return c.JSON(http.StatusInternalServerError, AssetsResponse{Error: "Failed to create video record"})
	}
	presignedUrl, err := s.storage.PresignedClient().PresignPostObject(c.Request().Context(), &s3.PutObjectInput{
		Bucket:        aws.String(s.cfg.S3.RawMediaBucket),
		Key:           aws.String(key),
		ContentType:   aws.String(body.ContentType),
		ContentLength: aws.Int64(int64(body.Size)),
		Metadata:      map[string]string{},
	}, func(options *s3.PresignPostOptions) {
		options.Expires = time.Duration(PUT_PRESIGNED_URL_TTL) * time.Second
	})
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, AssetsResponse{
		Data: &AssetsResponseData{
			UploadUrl: presignedUrl.URL,
			Header:    map[string]string{},
			Asset: Asset{
				Id:           videoId,
				Name:         body.Name,
				Size:         int(body.Size),
				ContentType:  body.ContentType,
				Href:         "",
				OriginalName: body.Name,
			},
			Form: presignedUrl.Values,
		},
		Message: "Presigned URL generated successfully",
	})
}

func (s *Server) UpdateMediaInternalHandler(c echo.Context) error {
	videoID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, AssetsResponse{Error: "Invalid video ID"})
	}
	body := UpdateMetadataRequest{}
	if err := RequestBody(c, &body); err != nil {
		return c.JSON(http.StatusBadRequest, AssetsResponse{Error: err.Error()})
	}

	params := db.PatchVideosParams{
		ID:     videoID,
		UserID: body.UserID,
	}
	if body.Title != nil {
		params.Title = pgtype.Text{
			String: *body.Title,
			Valid:  true,
		}
	}

	if body.Status != nil {
		params.Status = db.NullVideoStatus{
			VideoStatus: *body.Status,
			Valid:       true,
		}
	}

	if body.DurationSec != nil {
		params.DurationSec = pgtype.Int4{
			Int32: *body.DurationSec,
			Valid: true,
		}
	}
	if err := s.store.PatchVideos(c.Request().Context(), params); err != nil {
		s.log.Error("Failed to update video metadata", "err", err)
		return c.JSON(http.StatusInternalServerError, AssetsResponse{Error: "Failed to update video metadata"})
	}
	return c.NoContent(http.StatusOK)
}
