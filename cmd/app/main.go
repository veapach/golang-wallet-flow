package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/veapach/golang-wallet-flow/internal/app"
	"github.com/veapach/golang-wallet-flow/internal/config"
)

func main() {
	_ = godotenv.Load()
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	if err = app.New(cfg).Run(); err != nil {
		log.Fatal(err)
	}
}
