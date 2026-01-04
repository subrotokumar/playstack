package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Status string

const (
	StatusUp   Status = "UP"
	StatusDown Status = "DOWN"
)

type HealthResponse struct {
	Status Status `json:"status"`
}

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
	return c.JSON(http.StatusOK, HealthResponse{
		Status: StatusDown,
	})
}
