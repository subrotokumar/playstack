package main

import (
	"github.com/labstack/gommon/log"
	"gitlab.com/subrotokumar/glitchr/internal/server"
)

//	@title			Glitchr
//	@version		1.0
//	@description	API Dccumentation for Glitch Backend.

//	@contact.name	Subroto Kumar
//	@contact.url	http://github.com/subrotokumar/glitchr
//	@contact.email	subrotokumar@outlook.in

//	@license.name	Apache 2.0

//	@host		localhost:8080
//	@BasePath	/

// @externalDocs.description	OpenAPI
// @externalDocs.url			https://swagger.io/resources/open-api/
func main() {
	svc := server.NewHTTPServer()
	err := svc.Run()
	if err != nil {
		log.Errorf("%s", err.Error())
	}
}
