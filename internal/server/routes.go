package server

import (
	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/http/response"
)

const (
	HeaderContextType = "ContextType"

	ContextTypeJson = "application/json"
)

func (s *Server) HttpServer() *gofr.App {
	e := gofr.New()
	// e.Use(middleware.Logger())
	// e.Use(middleware.Recover())

	// e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
	// 	AllowOrigins:     []string{"https://*", "http://*"},
	// 	AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
	// 	AllowHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
	// 	AllowCredentials: true,
	// 	MaxAge:           300,
	// }))

	e.GET("/", HelloWorldHandler)

	return e
}

func HelloWorldHandler(c *gofr.Context) (any, error) {
	headers := map[string]string{
		HeaderContextType: ContextTypeJson,
	}

	return response.Response{
		Data: map[string]string{
			"message": "Hello World",
		},
		Headers: headers,
	}, nil
}
