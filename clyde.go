package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/diamondburned/arikawa/v3/discord"
)

func AskClyde(prompt string, instructions string) {
	if clydeChannel == 0 {
		id, err := discord.ParseSnowflake(os.Getenv("CLYDE_CHANNEL_ID"))
		if err != nil {
			tui.Send(logMsg{Msg: "Unable to parse channel id", Type: Error})
		}
		clydeChannel = discord.ChannelID(id)
	}

	if strings.Contains(prompt, "@cb") {
		prompt = strings.ReplaceAll(prompt, "@cb", ReadClipboard())
	}

	if len(instructions) > 0 {
		prompt += ".\n" + instructions
	}

	s.SendMessage(clydeChannel, prompt)
}

func FormatClydeReponse(msg string) (string, string) {
	parsed := strings.ReplaceAll(msg, fmt.Sprintf("<@!%s>", CurrentUserID), "`@You`")
	var md string

	if os.Getenv("GLAMOUR_STYLE") != "" {
		md, _ = glamour.RenderWithEnvironmentConfig(parsed)
	} else {
		md, _ = glamour.Render(parsed, "dark")
	}

	trimmed := strings.Trim(md, "\n") + "\n"
	return parsed, trimmed
}
