package main

import (
	"context"

	_ "github.com/joho/godotenv/autoload"
	"gitlab.com/subrotokumar/playstack/transcoder/service"
)

func main() {
	worker := service.New()
	worker.Run(context.Background())
}
