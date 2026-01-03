package server

import (
	"log"
	"net/http"
	"os"

	validation "github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"gitlab.com/subrotokumar/glitchr/internal/config"
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
	ResponseEntity struct {
		Data    any    `json:"data,omitempty"`
		Message string `json:"message,omitempty"`
		Error   any    `json:"error,omitempty"`
	}
	Server struct {
		cfg     config.Config
		idp     service.IdentityProvider
		handler *http.Server
	}
	Ctx struct {
		echo.Context
	}
)

func NewHTTPServer() *Server {
	validator = validation.New()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	cfg, err := config.ConfigFromEnv()
	if err != nil {
		log.Fatalln(err.Error())
	}

	srv := &Server{
		cfg: cfg,
		idp: service.NewIndentityProvider(),
	}
	srv.handler = &http.Server{
		Addr:    "0.0.0.0:" + port,
		Handler: srv.Mux(),
	}
	return srv
}

func (s *Server) Run() error {
	return s.handler.ListenAndServe()
}

func (ctx *Ctx) Body(v any) error {
	if err := ctx.Bind(v); err != nil {
		return err
	}
	if err := validator.Struct(v); err != nil {
		return err
	}
	return nil
}
