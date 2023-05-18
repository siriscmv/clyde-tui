package main

import (
	"context"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/session"
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
		tui.Send(DiscordMessage(c))
	})

	if err := s.Open(context.Background()); err != nil {
		tui.Send(logMsg{Msg: "Unable to establish discord connection", Type: Error})
	}
	defer s.Close()

	u, err := s.Me()
	if err != nil {
		tui.Send(logMsg{Msg: "Unable to get user", Type: Error})
	}

	tui.Send(logMsg{Msg: "Logged in as " + u.Username, Type: Info})
	CurrentUserID = u.ID.String()
	select {}
}
