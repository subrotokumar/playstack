package server

import (
	api "gofr.dev/pkg/gofr"
)

func (s *Server) HttpServer() *api.App {
	app := api.New()
	app.POST("/auth/signup", s.SignupHandler)
	app.POST("/auth/login", s.SignupHandler)

	return app
}
