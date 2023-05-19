package main

import (
	"context"
	"log"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/session"
)

type DiscordMessage *gateway.MessageCreateEvent

const clydeID = 1081004946872352958

var (
	Session       *session.Session
	ClydeChannel  discord.ChannelID
	CurrentUserID string
)

var Ready chan bool = make(chan bool)

func RunDiscordSession(token string) {
	Session = session.New(token)
	Session.AddHandler(func(c *gateway.MessageCreateEvent) {
		if c.Author.ID != clydeID || c.ChannelID != ClydeChannel {
			return
		}

		if mode == TUI {
			tui.Send(DiscordMessage(c))
		} else if mode == CLI {
			CLIChan <- c.Content
		}

	})

	if err := Session.Open(context.Background()); err != nil {
		if mode == TUI {
			tui.Send(logMsg{Msg: "Unable to establish discord connection", Type: Error})
		} else if mode == CLI {
			log.Fatalln("Unable to establish discord connection")
		}
	}
	defer Session.Close()

	u, err := Session.Me()
	if err != nil {
		if mode == TUI {
			tui.Send(logMsg{Msg: "Unable to get user", Type: Error})
		} else if mode == CLI {
			log.Fatalln("Unable to get user")
		}
	}

	CurrentUserID = u.ID.String()

	if mode == TUI {
		tui.Send(logMsg{Msg: "Logged in as " + u.Username, Type: Info})
	} else if mode == CLI {
		Ready <- true
	}

	select {}
}
