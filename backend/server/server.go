package server

import (
	"log/slog"
	"net/http"
	"os"

	validation "github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"gitlab.com/subrotokumar/playstack/backend/config"
	"gitlab.com/subrotokumar/playstack/libs/core"
	"gitlab.com/subrotokumar/playstack/libs/db"
	idp "gitlab.com/subrotokumar/playstack/libs/idp"
	"gitlab.com/subrotokumar/playstack/libs/storage"
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
		log     *core.Logger
		store   *db.SQLStore
		storage *storage.Storage
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
		slog.Error("Failed to load config: %s", "err", err.Error())
		os.Exit(0)
	}

	logger := core.NewLogger(cfg.App.Env, cfg.App.Name, nil)

	pgxpool, err := db.NewPgxPool(cfg.ConnectionUrl(), int32(cfg.Database.MinConn), int32(cfg.Database.MinConn))
	if err != nil {
		core.LogFatal("Failed to get pgxpool", "err", err.Error())
	}
	dbStore := db.NewSQLStore(pgxpool)

	storage := storage.NewStorageProvider(cfg.Aws.Region)

	srv := &Server{
		cfg:     cfg,
		idp:     idp.NewIndentityProvider(cfg.Aws.Region, cfg.Cognito.ClientID, cfg.Cognito.ClientSecret),
		log:     logger,
		store:   dbStore,
		storage: storage,
	}
	srv.handler = &http.Server{
		Addr:    cfg.App.Host + ":" + cfg.App.Port,
		Handler: srv.Mux(),
	}

	return srv
}

func (s *Server) Run() error {
	defer s.store.Close()
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
