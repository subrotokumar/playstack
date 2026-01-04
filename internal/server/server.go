package server

import (
	"log"
	"log/slog"
	"net/http"

	validation "github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"gitlab.com/subrotokumar/glitchr/internal/core"
	"gitlab.com/subrotokumar/glitchr/internal/service"
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
		cfg     core.Config
		idp     service.IdentityProvider
		handler *http.Server
		log     *slog.Logger
	}
	Ctx struct {
		echo.Context
	}
)

func NewHTTPServer() *Server {
	validator = validation.New(validation.WithRequiredStructEnabled())

	cfg, err := core.ConfigFromEnv()
	if err != nil {
		log.Fatalln(err.Error())
	}

	logger := core.NewLogger(cfg.App.Env, cfg.App.Name, nil)

	srv := &Server{
		cfg: cfg,
		idp: service.NewIndentityProvider(cfg.Aws.Region, cfg.Cognito.ClientID, cfg.Cognito.ClientSecret),
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
