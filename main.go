package main

import (
	"log"
	"os"
	"strings"

	godotenv "github.com/joho/godotenv"
)

type Mode string

const (
	CLI Mode = "CLI"
	TUI Mode = "TUI"
)

var mode Mode

func main() {
	godotenv.Load(".env")
	token := os.Getenv("CLYDE_DISCORD_USER_TOKEN")

	if token == "" {
		log.Fatalln("No $CLYDE_DISCORD_USER_TOKEN given.")
	} else if os.Getenv("CLYDE_CHANNEL_ID") == "" {
		log.Fatalln("No $CLYDE_CHANNEL_ID given.")
	}

	go InitClipboard()
	go RunDiscordSession(token)

	args := os.Args[1:]

	if len(args) > 0 {
		mode = CLI
		RunCLI(strings.Join(args, " "))
	} else {
		mode = TUI
		RunTUI()
	}
}
