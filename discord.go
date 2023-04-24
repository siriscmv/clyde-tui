package main

import (
	"context"
	"os"
	"strings"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/session"
	"golang.design/x/clipboard"
)

type DiscordMessage *gateway.MessageCreateEvent

const clydeID = 1081004946872352958

var (
	s             *session.Session
	clydeChannel  discord.ChannelID
	CurrentUserID string
)

func RunDiscordSession(token string) {
	s = session.New(token)
	s.AddHandler(func(c *gateway.MessageCreateEvent) {
		if c.Author.ID != clydeID || c.ChannelID != clydeChannel {
			return
		}
		p.Send(DiscordMessage(c))
	})

	if err := s.Open(context.Background()); err != nil {
		p.Send(logMsg{Msg: "Unable to establish discord connection", Type: Error})
	}
	defer s.Close()

	u, err := s.Me()
	if err != nil {
		p.Send(logMsg{Msg: "Unable to get user", Type: Error})
	}

	p.Send(logMsg{Msg: "Logged in as " + u.Username, Type: Info})
	CurrentUserID = u.ID.String()
	select {}
}

func AskClyde(msg string) {
	if clydeChannel == 0 {
		id, err := discord.ParseSnowflake(os.Getenv("CLYDE_CHANNEL_ID"))
		if err != nil {
			p.Send(logMsg{Msg: "Unable to parse channel id", Type: Error})
		}
		clydeChannel = discord.ChannelID(id)
	}

	if strings.Contains(msg, "@cb") {
		err := clipboard.Init()
		if err != nil {
			p.Send(Log{Msg: "Unable to initialize clipboard", Type: Error})
		} else {
			msg = strings.ReplaceAll(msg, "@cb", string(clipboard.Read(clipboard.FmtText)))
		}

	}

	s.SendMessage(clydeChannel, msg)
}
