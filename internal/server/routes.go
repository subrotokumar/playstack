package server

import (
	"log"

	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
)

func (s *Server) Mux() *echo.Echo {
	e := echo.New()
	s.resisterMetricsRoutes(e)
	s.registerRoutes(e)
	return e
}

func (s *Server) resisterMetricsRoutes(e *echo.Echo) {
	customRegistry := prometheus.NewRegistry()
	customCounter := prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "custom_requests_total",
			Help: "How many HTTP requests processed, partitioned by status code and HTTP method.",
		},
	)
	if err := customRegistry.Register(customCounter); err != nil {
		log.Fatal(err)
	}

	e.Use(echoprometheus.NewMiddlewareWithConfig(echoprometheus.MiddlewareConfig{
		AfterNext: func(c echo.Context, err error) {
			customCounter.Inc()
		},
		Registerer: customRegistry,
	}))
	e.GET("/metrics", echoprometheus.NewHandlerWithConfig(echoprometheus.HandlerConfig{Gatherer: customRegistry}))

}

func (s *Server) registerRoutes(e *echo.Echo) {
	e.GET("/health/liveness", s.LivenessHandler)
	e.GET("/health/readiness", s.ReadinessHandler)

	e.POST("/auth/signup", s.SignupHandler)
	e.POST("/auth/login", s.LoginHandler)
	e.POST("/auth/refresh", s.RefreshTokenHandler)
	e.POST("/auth/confirm-signup", s.ConfirmSignupHandler)
}
