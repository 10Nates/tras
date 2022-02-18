package main

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/andersfylling/disgord"
	"github.com/andersfylling/disgord/std"
)

func main() {
	//load client
	client := disgord.New(disgord.Config{
		BotToken: os.Getenv("Token"),
		Intents:  disgord.IntentGuildMessages | disgord.IntentDirectMessages | disgord.IntentGuildMessageReactions | disgord.IntentDirectMessageReactions,
	})
	defer client.Gateway().StayConnectedUntilInterrupted()

	//startup message
	client.Gateway().BotReady(func() {
		fmt.Println("Bot ready at " + time.Now().Local().Format(time.RFC1123))
		client.UpdateStatusString("@me help")
	})

	//filter out unwanted messages
	content, err := std.NewMsgFilter(context.Background(), client)
	if err != nil {
		panic(err)
	}
	content.NotByBot(client)
	content.ContainsBotMention(client)

	//on message with mention
	client.Gateway().
		WithMiddleware(content.NotByBot, content.ContainsBotMention).       // filter
		MessageCreate(func(s disgord.Session, evt *disgord.MessageCreate) { // on message

			go parseCommand(evt.Message, &s)
		})
}

func parseCommand(msg *disgord.Message, session *disgord.Session) {
	cstr := msg.Content
	rsplitstr := regexp.MustCompile(`([^\\])( )`).ReplaceAllString(cstr, "$1\n")
	carr := strings.Split(rsplitstr, "\n")

	fmt.Println(carr)
}
