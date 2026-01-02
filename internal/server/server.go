package server

import (
	"os"
	"strconv"

	validation "github.com/go-playground/validator/v10"
	"gitlab.com/subrotokumar/glitchr/internal/service"
	"gofr.dev/pkg/gofr"
	api "gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/http"
	"gofr.dev/pkg/gofr/http/response"
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
		port int
		idp  service.IdentityProvider
	}
	Ctx struct {
		*gofr.Context
	}
)

func NewServer() *api.App {
	_ = http.ErrorPanicRecovery{}
	validator = validation.New()
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	NewServer := &Server{
		port: port,
		idp:  service.NewIndentityProvider(),
	}
	server := NewServer.HttpServer()
	return server
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

func (ctx *Ctx) Json(res any) (any, error) {
	return response.Response{
		Data: res,
	}, nil
}

func (ctx *Ctx) JsonWithErr(res any, httpError HttpError) (any, error) {
	return response.Response{
		Data: res,
	}, nil
}
