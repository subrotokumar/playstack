package server

import (
	"context"
	"log/slog"
	"strings"

	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus"
	echoSwagger "github.com/swaggo/echo-swagger"
	_ "gitlab.com/subrotokumar/glitchr/backend/swagger"
	"gitlab.com/subrotokumar/glitchr/libs/idp"
)

func (s *Server) Mux() *echo.Echo {
	e := echo.New()
	s.registerMiddleware(e)
	s.registerOpenAPIRoutes(e)
	s.resisterMetricsRoutes(e)
	s.registerRoutes(e)
	return e
}

func (s *Server) registerMiddleware(e *echo.Echo) {
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowCredentials: true,
	}))

	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig(
		middleware.RequestLoggerConfig{
			Skipper: func(c echo.Context) bool {
				return strings.Contains(c.Request().URL.Path, "/swagger")
			},
			LogMethod:        true,
			LogURI:           true,
			LogStatus:        true,
			LogLatency:       true,
			LogHost:          true,
			LogRequestID:     true,
			LogContentLength: true,
			LogResponseSize:  true,
			LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
				if v.Error == nil {
					slog.LogAttrs(context.Background(), slog.LevelInfo, "REQUEST",
						slog.String("method", v.Method),
						slog.String("uri", v.URI),
						slog.Int("status", v.Status),
						slog.Duration("latency", v.Latency),
						slog.String("host", v.Host),
						slog.String("bytes_in", v.ContentLength),
						slog.Int64("bytes_out", v.ResponseSize),
						// slog.String("user_agent", v.UserAgent),
						slog.String("remote_ip", v.RemoteIP),
						slog.String("request_id", v.RequestID),
					)
				} else {
					slog.LogAttrs(context.Background(), slog.LevelError, "REQUEST_ERROR",
						slog.String("method", v.Method),
						slog.String("uri", v.URI),
						slog.Int("status", v.Status),
						slog.Duration("latency", v.Latency),
						slog.String("host", v.Host),
						slog.String("bytes_in", v.ContentLength),
						slog.Int64("bytes_out", v.ResponseSize),
						// slog.String("user_agent", v.UserAgent),
						slog.String("remote_ip", v.RemoteIP),
						slog.String("request_id", v.RequestID),
						slog.String("error", v.Error.Error()),
					)
				}
				return nil
			},
		},
	)))
}

func (s *Server) registerOpenAPIRoutes(e *echo.Echo) {
	// Swagger UI
	e.GET("/swagger/*", echoSwagger.WrapHandler)
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
		s.log.Fatal(err.Error())
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
	e.POST("/auth/profile", s.ProfileHandler)

	e.Use(idp.NewAuthMiddleware(s.cfg.Aws.Region, s.cfg.Cognito.UserPoolID, s.cfg.Cognito.ClientID).AuthMiddleware())
	e.POST("/upload/policies/assets", s.AssetsHandler)
}
