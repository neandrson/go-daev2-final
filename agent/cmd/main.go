package main

import (
	"log"

	"github.com/joho/godotenv"

	agent "github.com/neandrson/go-daev2-final/agent/internal"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	agent.RunAgent()
}
