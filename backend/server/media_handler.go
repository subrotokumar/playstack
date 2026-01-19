package server

import (
	"net/http"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

const (
	PUT_PRESIGNED_URL_TTL = 1 * 60
)

type (
	Asset struct {
		Id           string `json:"id"`
		Name         string `json:"name"`
		Size         int    `json:"size"`
		ContentType  string `json:"content_type"`
		Href         string `json:"href"`
		OriginalName string `json:"original_name"`
	}
	AssetsResponse struct {
		UploadUrl string            `json:"upload_url"`
		Header    map[string]string `json:"header"`
		Asset     Asset             `json:"asset"`
		Form      map[string]string `json:"form"`
	}
)

func (s *Server) AssetsHandler(c echo.Context) error {
	name := c.FormValue("name")
	size, err := strconv.Atoi(c.FormValue("size"))
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	contentType := c.FormValue("content_type")

	presignedUrl, err := s.storage.PresignedClient().PresignPostObject(c.Request().Context(), &s3.PutObjectInput{
		Bucket:        aws.String(s.cfg.S3.RawMediaBucket),
		Key:           aws.String("/user"),
		ContentType:   aws.String(contentType),
		ContentLength: aws.Int64(int64(size)),
		Metadata:      map[string]string{},
	}, func(options *s3.PresignPostOptions) {
		options.Expires = time.Duration(PUT_PRESIGNED_URL_TTL) * time.Second
	})
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, AssetsResponse{
		UploadUrl: presignedUrl.URL,
		Header:    map[string]string{},
		Asset: Asset{
			Id:           uuid.Must(uuid.NewV7()).String(),
			Name:         name,
			Size:         size,
			ContentType:  contentType,
			Href:         "",
			OriginalName: name,
		},
		Form: presignedUrl.Values,
	})
}
