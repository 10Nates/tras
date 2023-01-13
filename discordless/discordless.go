package discordless

import (
	"fmt"
	"time"

	"github.com/andersfylling/disgord"
)

// This file creates discordless/headless messages and assists in parsing their results.
// This allows for tests and other external calls

type ParseCommand func(msg *disgord.Message, s *disgord.Session)

func CreateHeadlessMessage(content string, identifier string) (*disgord.Message, *disgord.Session) {
	var s *disgord.Session

	newmsg := &disgord.Message{
		Author: &disgord.User{
			Email: identifier, // allows source reference later down the line, not necessarily an actual email
		},
		Content:         content,
		Timestamp:       disgord.Time{Time: time.Now()},
		EditedTimestamp: disgord.Time{Time: time.Now()},
	}

	return newmsg, s
}

func HeadlessReply(content string, identifier string) {
	if identifier == "TEST" {
		go func() { // non-blocking
			testChannel <- content // send for error checking & printing
		}()
		return
	}

	fmt.Println(content) // TODO: Custom API reply
}

func HeadlessReact(emoji interface{}, identifier string) {
	if identifier == "TEST" {
		go func() { // non-blocking
			testChannel <- fmt.Sprint("(reaction) ", emoji) // send for error checking & printing
		}()
		return
	}

	fmt.Println("Reacted: ", emoji) // TODO: Custom API react
}
