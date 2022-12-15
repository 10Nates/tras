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

const BOT_VERSION = "3.0.0"

var BotID string
var BotPFP string

func main() {
	//load client
	client := disgord.New(disgord.Config{
		BotToken: os.Getenv("Token"),
		Intents: disgord.IntentGuildMessages | disgord.IntentDirectMessages |
			disgord.IntentGuildMessageReactions | disgord.IntentDirectMessageReactions,
	})
	defer client.Gateway().StayConnectedUntilInterrupted()

	//startup message
	client.Gateway().BotReady(func() {
		usr, err := client.CurrentUser().Get()
		if err != nil {
			panic(err) // Bot shouldn't start
		}
		BotID = usr.ID.String()
		BotPFP, err = usr.AvatarURL(256, false)
		if err != nil {
			panic(err) // Bot shouldn't start
		}
		fmt.Println("Bot started @ " + time.Now().Local().Format(time.RFC1123))
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
		WithMiddleware(content.NotByBot, content.ContainsBotMention, content.HasBotMentionPrefix). // filter
		MessageCreate(func(s disgord.Session, evt *disgord.MessageCreate) {                        // on message
			go parseCommand(evt.Message, &s)
		})
}

func parseCommand(msg *disgord.Message, s *disgord.Session) {
	cstr := msg.Content
	rsplitstr := regexp.MustCompile(`([^\\])( )`).ReplaceAllString(cstr, "$1\n")
	carr := strings.Split(rsplitstr, "\n")

	args := []string{}
	argsl := []string{}

	for i := 0; i < len(carr); i++ {
		if !strings.Contains(carr[i], BotID) {
			args = append(args, carr[i])
			argsl = append(argsl, strings.ToLower(carr[i]))
		}
	}

	if len(args) < 1 {
		args = append(args, "")
		argsl = append(argsl, "")
	}

	switch argsl[0] {
	case "help":
		helpResponse(msg, s)
	case "about":
		aboutResponse(msg, s)
	case "oof":
		// big OOF
		baseReply(msg, s, "oof oof oof     oof oof oof     oof oof oof\noof        oof     oof        oof     oof\noof        oof     oof        oof     oof oof oof\noof        oof     oof        oof     oof\noof oof oof     oof oof oof     oof")
	case "f":
		// big F
		baseReply(msg, s, "F F F F F F\nF F \nF F F F F F\nF F\nF F")
	case "pi":

	case "big":

	case "jumble":

	case "emojify":

	case "flagify":

	case "superscript":

	case "unicodify":

	case "bold":

	case "replace":

	case "overcomplicate":

	case "word":
		if len(argsl) > 1 && argsl[1] == "info" {

		} else {
			defaultResponse(msg, s)
		}
	case "ascii":
		if len(argsl) > 1 && argsl[1] == "art" {

		} else {
			defaultResponse(msg, s)
		}
	case "commands":

	case "rank":

	case "set":
		if len(argsl) > 1 && (argsl[1] == "nickname" || argsl[1] == "nick") {

		} else {
			defaultResponse(msg, s)
		}
	case "speak":

	case "combinations":

	default:
		defaultResponse(msg, s)
	}
}
