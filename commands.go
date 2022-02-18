package main

import (
	"context"

	"github.com/andersfylling/disgord"
)

func baseReply(msg *disgord.Message, s *disgord.Session, reply string) {
	msg.Reply(context.Background(), *s, reply)
}
