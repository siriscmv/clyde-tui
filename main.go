package main

import (
	"log"
	"os"

	godotenv "github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")
	token := os.Getenv("CLYDE_DISCORD_USER_TOKEN")

	if token == "" {
		log.Fatalln("No $CLYDE_DISCORD_USER_TOKEN given.")
		os.Exit(1)
	} else if os.Getenv("CLYDE_CHANNEL_ID") == "" {
		log.Fatalln("No $CLYDE_CHANNEL_ID given.")
		os.Exit(1)
	}

	go InitClipboard()
	go RunDiscordSession(token)
	RunTUI()
}
