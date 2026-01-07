package server

import (
	"log"
	"log/slog"
	"net/http"

	validation "github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"gitlab.com/subrotokumar/glitchr/backend/config"
	"gitlab.com/subrotokumar/glitchr/pkg/core"
	idp "gitlab.com/subrotokumar/glitchr/pkg/idp"
	"gitlab.com/subrotokumar/glitchr/pkg/logger"
)

const (
	HeaderContextType = "ContextType"
	ContextTypeJson   = "application/json"
)

var (
	validator         *validation.Validate
	HeaderJsonContent = map[string]string{
		HeaderContextType: ContextTypeJson,
	}
)

type (
	Server struct {
		cfg     config.Config
		idp     idp.IdentityProvider
		handler *http.Server
		log     *slog.Logger
	}
	Ctx struct {
		echo.Context
	}
)

func NewHTTPServer() *Server {
	validator = validation.New(validation.WithRequiredStructEnabled())

	cfg := config.Config{}
	err := core.ConfigFromEnv(&cfg)
	if err != nil {
		log.Fatalf("Failed to load config: %s", err.Error())
	}

	logger := logger.New(cfg.App.Env, cfg.App.Name, nil)

	srv := &Server{
		cfg: cfg,
		idp: idp.NewIndentityProvider(cfg.Aws.Region, cfg.Cognito.ClientID, cfg.Cognito.ClientSecret),
		log: logger,
	}
	srv.handler = &http.Server{
		Addr:    cfg.App.Host + ":" + cfg.App.Port,
		Handler: srv.Mux(),
	}

	return srv
}

func (s *Server) Run() error {
	s.log.Info("Server running at " + s.cfg.App.Host + ":" + s.cfg.App.Port)
	return s.handler.ListenAndServe()
}

func RequestBody(ctx echo.Context, v any) error {
	if err := ctx.Bind(v); err != nil {
		return err
	}
	if err := validator.Struct(v); err != nil {
		return err
	}
	return nil
}
