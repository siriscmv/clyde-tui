package main

import (
	"log"
	"os"

	godotenv "github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")
	token := os.Getenv("TOKEN")

	if token == "" {
		log.Fatalln("No $TOKEN given.")
		os.Exit(1)
	} else if os.Getenv("CLYDE_CHANNEL_ID") == "" {
		log.Fatalln("No $CLYDE_CHANNEL_ID given.")
		os.Exit(1)
	}

	go RunDiscordSession(token)
	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
}
