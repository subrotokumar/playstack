package server

import (
	"os"
	"strconv"

	_ "github.com/joho/godotenv/autoload"
	"gofr.dev/pkg/gofr"
)

type Server struct {
	port int
}

func NewServer() *gofr.App {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	NewServer := &Server{
		port: port,
	}

	// Declare Server config
	server := NewServer.HttpServer()

	return server
}
