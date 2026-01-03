package main

import (
	"github.com/labstack/gommon/log"
	"gitlab.com/subrotokumar/glitchr/internal/server"
)

func main() {
	svc := server.NewHTTPServer()
	err := svc.Run()
	if err != nil {
		log.Errorf("%s", err.Error())
	}
}
