package main

import (
	"context"

	"github.com/andersfylling/disgord"
)

func baseReply(msg *disgord.Message, s *disgord.Session, reply string) {
	msg.Reply(context.Background(), *s, reply)
}

func baseEmbedReply(msg *disgord.Message, s *disgord.Session, embed *disgord.Embed) {
	msg.Reply(context.Background(), *s, embed)
}

func defaultResponse(msg *disgord.Message, s *disgord.Session) {
	baseReply(msg, s, "What's up?")
}