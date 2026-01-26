package server

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type Status string

const (
	StatusUp   Status = "UP"
	StatusDown Status = "DOWN"
)

type (
	DatabaseHealthStatus struct {
		Status        Status `json:"status"`
		TotalConns    int32  `json:"total_conns"`
		IdleConns     int32  `json:"idle_conns"`
		AcquiredConns int32  `json:"acquired_conns"`
		MaxConns      int32  `json:"max_conns"`
	}
	HealthResponse struct {
		Status   Status                `json:"status"`
		Database *DatabaseHealthStatus `json:"db,omitempty"`
	}
)

// LivenessHandler godoc
//
//	@Summary		Liveness probe
//	@Description	Indicates whether the application process is alive
//	@Tags			health
//	@Produce		json
//	@Success		200	{object}	HealthResponse	"Service is alive"
//	@Router			/health/liveness [get]
func (s *Server) LivenessHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, HealthResponse{
		Status: StatusUp,
	})
}

// ReadinessHandler godoc
//
//	@Summary		Readiness probe
//	@Description	Indicates whether the application is ready to receive traffic
//	@Tags			health
//	@Produce		json
//	@Success		200	{object}	HealthResponse	"Service is ready or not"
//	@Router			/health/readiness [get]
func (s *Server) ReadinessHandler(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), 500*time.Millisecond)
	defer cancel()

	dbStatus := StatusDown
	httpStatus := http.StatusServiceUnavailable

	if _, err := s.store.GetTimestamp(ctx); err == nil {
		dbStatus = StatusUp
		httpStatus = http.StatusOK
	}
	stat := s.store.Stat()

	db := &DatabaseHealthStatus{
		Status:        dbStatus,
		TotalConns:    stat.TotalConns(),
		IdleConns:     stat.IdleConns(),
		AcquiredConns: stat.AcquiredConns(),
		MaxConns:      stat.MaxConns(),
	}

	return c.JSON(httpStatus, HealthResponse{
		Status:   dbStatus,
		Database: db,
	})
}
