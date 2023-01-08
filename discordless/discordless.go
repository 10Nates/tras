package discordless

import (
	"fmt"
	"time"

	"github.com/andersfylling/disgord"
)

type ParseCommand func(msg *disgord.Message, s *disgord.Session)

func Test(handler ParseCommand) {
	CreateHeadlessMessage("<@462051981863682048> about", "test", handler)
}

func CreateHeadlessMessage(content string, identifier string, handler ParseCommand) *disgord.Message {
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

	handler(newmsg, s)

	return newmsg
}

func HeadlessReact(emoji interface{}) {
	fmt.Println("Responded with emoji: ", emoji)
}
