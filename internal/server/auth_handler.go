package server

import (
	api "gofr.dev/pkg/gofr"
)

type SignUpRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (s *Server) SignupHandler(c *api.Context) (any, error) {
	ctx := Ctx{c}
	body := SignUpRequest{}
	ctx.Body(&body)
	return ctx.Json(map[string]string{
		"message": "Hello World",
	})
}
