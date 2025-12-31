package main

import (
	"gitlab.com/subrotokumar/glitchr/internal/server"
)

func main() {
	svc := server.NewServer()
	svc.Run()
}
