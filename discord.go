package main

import (
	"context"
	"log"
	"os"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/session"
)

type DiscordMessage *gateway.MessageCreateEvent

var clydeChannel discord.ChannelID

const clydeID = 1081004946872352958

func RunDiscordSession(token string) {
	s := session.New(token)
	s.AddHandler(func(c *gateway.MessageCreateEvent) {
		if c.Author.ID != clydeID || c.ChannelID != clydeChannel {
			return
		}

		p.Send(DiscordMessage(c))
	})

	s.AddIntents(gateway.IntentGuildMessages)

	if err := s.Open(context.Background()); err != nil {
		p.Send(logMsg{Msg: "Unable to establish discord connection", Type: Error})
	}
	defer s.Close()

	u, err := s.Me()
	if err != nil {
		p.Send(logMsg{Msg: "Unable to get user", Type: Error})
	}

	p.Send(logMsg{Msg: "Logged in as " + u.Username, Type: Info})

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}

	select {}
}

var HttpClient = api.NewClient(os.Getenv("token"))

func SendDiscordMessage(msg string) {
	if clydeChannel == 0 {
		id, err := discord.ParseSnowflake(os.Getenv("CLYDE_CHANNEL_ID"))
		if err != nil {
			p.Send(logMsg{Msg: "Unable to parse channel id", Type: Error})
		}
		clydeChannel = discord.ChannelID(id)
	}

	HttpClient.SendMessage(clydeChannel, msg)
}
