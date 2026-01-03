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

func (s *Server) LivenessHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, HealthResponse{
		Status: StatusUp,
	})
}

func (s *Server) ReadinessHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, HealthResponse{
		Status: StatusDown,
	})
}
