package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/andersfylling/disgord"
	"github.com/andersfylling/disgord/std"
)

var BotID string // loaded on init
var BotPFP string
var BotClient *disgord.Client

func main() {
	//load client
	client := disgord.New(disgord.Config{
		BotToken: os.Getenv("Token"),
		Intents: disgord.IntentGuildMessages | disgord.IntentDirectMessages |
			disgord.IntentGuildMessageReactions | disgord.IntentDirectMessageReactions,
	})
	defer client.Gateway().StayConnectedUntilInterrupted()

	BotClient = client

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
		WithMiddleware(content.NotByBot, content.NotByWebhook, content.ContainsBotMention, content.HasBotMentionPrefix). // filter
		MessageCreate(func(s disgord.Session, evt *disgord.MessageCreate) {                                              // on message
			go parseCommand(evt.Message, &s)
		})
}

func parseCommand(msg *disgord.Message, s *disgord.Session) {
	procTimeStart := time.Now() // timer for ping info
	cstr := msg.Content
	rsplitstr := argumentSplitRegex.ReplaceAllString(cstr, "$1\n") // separate arguments with "\\" and "\"
	rfixspacestr := strings.ReplaceAll(rsplitstr, "\\ ", " ")      // fix "\ " to " " (this should only happen when that should happen)
	rfixslashstr := strings.ReplaceAll(rfixspacestr, "\\\\", "\\") // fix "\\" to "\" (this is somewhat wonky behavior because "\" doesn't do
	carr := strings.Split(rfixslashstr, "\n")                      // anything outside of the spaces so it can still work without disappearing but oh well)

	args := []string{}
	argsl := []string{}

	for i := 0; i < len(carr); i++ {
		if !strings.Contains(carr[i], BotID) { // ignore where bot is mentioned
			args = append(args, carr[i])
			argsl = append(argsl, strings.ToLower(carr[i]))
		}
	}

	if len(args) < 1 { // prevent error in switch case
		args = append(args, "")
		argsl = append(argsl, "")
	}

	switch argsl[0] {
	case "help":
		helpResponse(msg, s)
	case "about":
		if len(argsl) > 1 && argsl[1] == "nocb" { // remove code blocks so iOS can click the links
			aboutResponse(msg, s, true)
		} else {
			aboutResponse(msg, s, false)
		}
	case "oof":
		// big OOF
		baseReply(msg, s, "oof oof oof     oof oof oof     oof oof oof\noof        oof     oof        oof     oof\noof        oof     oof        oof     oof oof oof\noof        oof     oof        oof     oof\noof oof oof     oof oof oof     oof")
	case "f":
		// big F
		baseReply(msg, s, "F F F F F F\nF F \nF F F F F F\nF F\nF F")
	case "pi":
		piResponse(msg, s)
	case "big":

	case "jumble":

	case "emojify":
		if len(argsl) > 1 {
			text := strings.Join(argsl[1:], " ") // case insensitive
			emojifyResponse(text, msg, s)
		} else {
			baseReply(msg, s, "What would you like me to change to emojis?")
		}
	case "flagify":
		if len(argsl) > 1 {
			text := strings.Join(argsl[1:], " ") // case insensitive
			flagifyResponse(text, msg, s)
		} else {
			baseReply(msg, s, "Ya gotta tell me what to flagify!")
		}
	case "superscript":
		if len(argsl) > 1 {
			text := strings.Join(args[1:], " ") // case sensitive
			superScriptResponse(text, msg, s)
		} else {
			baseReply(msg, s, "What do you need to be superscripts?")
		}
	case "unicodify":
		if len(argsl) > 1 {
			text := strings.Join(args[1:], " ") // case sensitive
			unicodifyResponse(text, msg, s)
		} else {
			baseReply(msg, s, "You need to tell me what to unicodify.")
		}
	case "bold":
		if len(argsl) > 1 {
			text := strings.Join(args[1:], " ") // case sensitive
			boldResponse(text, msg, s)
		} else {
			baseReply(msg, s, "What needs bolding?")
		}
	case "replace":
		if len(argsl) > 3 {
			text := strings.Join(args[3:], " ") // case sensitive
			replaceResponse(args[1], args[2], text, msg, s)
		} else {
			baseReply(msg, s, "Tell me the [what to replace], the [replacement], and then provide the [body of text].")
		}
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
			text := strings.Join(args[2:], " ") // case sensitive
			if len(argsl) > 2 {
				setNickResponse(text, msg, s)
			} else {
				baseReply(msg, s, "What should by nickname be?")
			}
		} else {
			defaultResponse(msg, s)
		}
	case "speak":

	case "combinations":

	case "ping":
		if len(argsl) > 1 && (argsl[1] == "info" || argsl[1] == "information") {
			pingResponse(true, msg, s, procTimeStart)
		} else {
			pingResponse(false, msg, s, procTimeStart)
		}
	default:
		defaultResponse(msg, s)
	}
}
