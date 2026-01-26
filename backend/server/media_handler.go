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
	POST_PRESIGNED_URL_TTL = 1 * 60
	PUT_PRESIGNED_URL_TTL  = 1 * 1
)

const (
	ErrFailedToGeneratePresignedURL = "failed to generate presigned url"
	ErrInvalidVideoID               = "invalid video ID"
	ErrFailedToCreateVideoRecord    = "failed to create video record"
	ErrFailedToFetchVideo           = "failed to fetch video"
	ErrNoPermission                 = "you do not have permission to perform this action"
	ErrFailedToUpdateMetadata       = "failed to update video metadata"
)

const (
	MsgPresignedURLGenerated = "presigned URL generated successfully"
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
	ThumbnailAssetsRequest struct {
		Name        string `json:"name" validate:"required"`
		Size        int64  `json:"size" validate:"required"`
		ContentType string `json:"content_type" validate:"required"`
	}
	ThumbnailAssetsResponseData struct {
		UploadUrl string `json:"upload_url"`
	}
	ThumbnailAssetsResponse struct {
		Data    ThumbnailAssetsResponseData `json:"data,omitempty"`
		Message string                      `json:"message,omitempty"`
		Error   any                         `json:"error,omitempty"`
	}

	UpdateMetadataRequest struct {
		UserID      uuid.UUID       `json:"user_id"`
		Title       *string         `json:"title"`
		Status      *db.VideoStatus `json:"status" validate:"omitempty,oneof='PREUPLOAD' 'UPLOADED' 'PROCESSING' 'READY' 'FAILED'"`
		DurationSec *int32          `json:"duration_sec"`
	}

	GetVideoResponse struct {
		Data    []db.Video `json:"data"`
		Message string     `json:"message,omitempty"`
		Error   any        `json:"error,omitempty"`
	}
)

// VideoAssetsHandler godoc
//
// @Summary      Create presigned URL for video upload
// @Description Creates a video record and returns a presigned POST URL for uploading raw media
// @Tags         media
// @Accept       json
// @Produce      json
// @Param        body  body      AssetsRequest  true  "Video asset metadata"
// @Success      200   {object}  AssetsResponse
// @Failure      400   {object}  AssetsResponse
// @Failure      500   {object}  AssetsResponse
// @Security     BearerAuth
// @Router       /media/videos/signed-url [post]
func (s *Server) VideoAssetsHandler(c echo.Context) error {
	body := AssetsRequest{}
	if err := RequestBody(c, &body); err != nil {
		return c.JSON(http.StatusInternalServerError, AssetsResponse{Error: err.Error()})
	}

	videoId := uuid.Must(uuid.NewV7())
	userId := c.Get("sub").(uuid.UUID)
	key := fmt.Sprintf("/%s/%s/%s", userId.String(), videoId.String(), body.Name)
	_, err := s.store.CreateVideo(c.Request().Context(), db.CreateVideoParams{
		ID:          videoId,
		UserID:      userId,
		Title:       body.Name,
		Status:      db.VideoStatusPREUPLOAD,
		DurationSec: pgtype.Int4{Valid: false},
	})
	if err != nil {
		s.log.Error(ErrFailedToCreateVideoRecord, "err", err)
		return c.JSON(http.StatusInternalServerError, AssetsResponse{Error: ErrFailedToCreateVideoRecord})
	}
	presignedUrl, err := s.storage.PresignedClient().PresignPostObject(c.Request().Context(), &s3.PutObjectInput{
		Bucket:        aws.String(s.cfg.S3.RawMediaBucket),
		Key:           aws.String(key),
		ContentType:   aws.String(body.ContentType),
		ContentLength: aws.Int64(int64(body.Size)),
		Metadata:      map[string]string{},
	}, func(options *s3.PresignPostOptions) {
		options.Expires = time.Duration(POST_PRESIGNED_URL_TTL) * time.Second
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
		Message: MsgPresignedURLGenerated,
	})
}

// ThumbnailSignedUrlHandler godoc
//
// @Summary      Create presigned URL for thumbnail upload
// @Description Returns a presigned PUT URL for uploading a video thumbnail
// @Tags         media
// @Accept       json
// @Produce      json
// @Param        videoId  path      string                 true  "Video ID"
// @Param        body     body      ThumbnailAssetsRequest true  "Thumbnail metadata"
// @Success      200      {object}  ThumbnailAssetsResponse
// @Failure      400      {object}  AssetsResponse
// @Failure      403      {object}  AssetsResponse
// @Failure      500      {object}  AssetsResponse
// @Security     BearerAuth
// @Router       /media/videos/{videoId}/thumbnail/signed-url [put]
func (s *Server) ThumbnailSignedUrlHandler(c echo.Context) error {
	userId := c.Get("sub").(uuid.UUID)
	videoId, err := uuid.Parse(c.Param("videoId"))
	body := ThumbnailAssetsRequest{}
	if err := RequestBody(c, &body); err != nil {
		return c.JSON(http.StatusInternalServerError, AssetsResponse{Error: err.Error()})
	}

	if err != nil {
		return c.JSON(http.StatusBadRequest, AssetsResponse{Error: ErrInvalidVideoID})
	}

	key := fmt.Sprintf("/%s/%s/thumbnail", userId.String(), videoId.String())
	video, err := s.store.GetVideoByID(c.Request().Context(), videoId)
	if err != nil {
		s.log.Error(ErrFailedToFetchVideo, "err", err)
		return c.JSON(http.StatusInternalServerError, AssetsResponse{Error: ErrFailedToFetchVideo})
	}
	if video.UserID != userId {
		return c.JSON(http.StatusForbidden, AssetsResponse{Error: ErrNoPermission})
	}

	presignedUrl, err := s.storage.PresignedClient().PresignPutObject(c.Request().Context(), &s3.PutObjectInput{
		Bucket:        aws.String(s.cfg.S3.MediaBucket),
		Key:           aws.String(key),
		ContentType:   aws.String(body.ContentType),
		ContentLength: aws.Int64(int64(body.Size)),
		Metadata:      map[string]string{},
	}, func(options *s3.PresignOptions) {
		options.Expires = time.Duration(PUT_PRESIGNED_URL_TTL) * time.Second
	})
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, ThumbnailAssetsResponse{
		Data: struct {
			UploadUrl string "json:\"upload_url\""
		}{
			UploadUrl: presignedUrl.URL,
		},
		Message: MsgPresignedURLGenerated,
	})
}

// UpdateMediaInternalHandler godoc
//
// @Summary      Update video metadata (internal)
// @Description Updates title, status, or duration of a video
// @Tags         media-internal
// @Accept       json
// @Produce      json
// @Param        videoId  path      string                 true  "Video ID"
// @Param        body     body      UpdateMetadataRequest  true  "Metadata update payload"
// @Success      200      {string}  string "OK"
// @Failure      400      {object}  AssetsResponse
// @Failure      500      {object}  AssetsResponse
// @Security     BasicAuth
// @Router       /internal/media/videos/{videoId} [patch]
func (s *Server) UpdateMediaInternalHandler(c echo.Context) error {
	videoID, err := uuid.Parse(c.Param("videoId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, AssetsResponse{Error: ErrInvalidVideoID})
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
		s.log.Error(ErrFailedToUpdateMetadata, "err", err)
		return c.JSON(http.StatusInternalServerError, AssetsResponse{Error: ErrFailedToUpdateMetadata})
	}
	return c.NoContent(http.StatusOK)
}

// GetVideoHandler godoc
//
// @Summary      List ready videos
// @Description Returns videos with READY status
// @Tags         media
// @Produce      json
// @Success      200   {object}   GetVideoResponse
// @Failure      500   {object}  GetVideoResponse
// @Security     BearerAuth
// @Router       /media/videos [get]
func (s *Server) GetVideoHandler(c echo.Context) error {
	resp, err := s.store.SearchVideo(c.Request().Context(), db.SearchVideoParams{
		Status: db.NullVideoStatus{VideoStatus: db.VideoStatusREADY},
	})

	if err != nil {
		s.log.Error(ErrFailedToFetchVideo, "err", err)
		return c.JSON(http.StatusInternalServerError, GetVideoResponse{Error: ErrFailedToFetchVideo})
	}

	return c.JSON(http.StatusOK, GetVideoResponse{
		Data: resp,
	})
}
