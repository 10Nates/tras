package discordless

import (
	"fmt"
	"time"

	"github.com/andersfylling/disgord"
)

type ParseCommand func(msg *disgord.Message, s *disgord.Session)

func CreateHeadlessMessage(content string, identifier string) (*disgord.Message, *disgord.Session) {
	var s *disgord.Session

	newmsg := &disgord.Message{
		Author: &disgord.User{
			Email: identifier, // allows source reference later down the line, not necessarily an actual email
		},
		Member:          &disgord.Member{},
		Content:         content,
		Timestamp:       disgord.Time{Time: time.Now()},
		EditedTimestamp: disgord.Time{Time: time.Now()},
	}

	return newmsg, s
}

func HeadlessReply(content string, identifier string) {
	if identifier == "TEST" {
		return // No need to send anywhere - only checking for errors
	}
	fmt.Println(content) // TODO: Custom API reply
}

func HeadlessReact(emoji interface{}, identifier string) {
	if identifier == "TEST" {
		return // No need to send anywhere - only checking for errors
	}
	fmt.Println("Reacted: ", emoji) // TODO: Custom API react
}
