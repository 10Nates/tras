package discordless

import "time"

// File for unit testing but instead of it being automated I look at it with my eyes
// Later implemented in another file with actual automation

var tests = []string{
	"help",
	"about",
	"oof",
	"f",
	"pi",

	"big hi hello",
	"big -t hi hello",

	"jumble The quick brown fox",
	"emojify The quick brown fox",
	"flagify The quick brown fox",
	"superscript The quick brown fox",
	"unicodify The quick brown fox",
	"bold The quick brown fox",
	"replace sly quick The sly brown fox",
	"overcomplicate The quick brown fox",

	"word info def fox",
	"word info cat fox",

	"ascii art getfonts",
	"ascii art big Hello!",

	"commands view",
	"commands manage",
	"commands manage set test works",
	"commands manage delete test",

	"rank info",
	"rank checkDice",
	"rank dice",
	"rank set",
	"rank set invalid",
	"rank set invalid 42",
	"rank set <@123456789123456789> invalid",
	"rank set <@123456789123456789> 42",
	"rank reset",
	"rank reset invalid",
	"rank reset <@123456789123456789>",

	"set nick TestName1",
	"rest nick",
	"set nickname TestName2",
	"rest nickname",

	"speak generate",
	"speak generate How is",
	"speak randomspeak status",
	"speak randomspeak on",
	"speak randomspeak off",
	"speak randomspeak invalid",

	"combinations words The quick brown fox",
	"combinations characters Hello",

	"ping",
	"ping info",
}

var testChannel = make(chan string)

func Test(handler ParseCommand) []string {
	prefix := "<@462051981863682048> "

	res := []string{}
	for _, test := range tests {
		handler(CreateHeadlessMessage(prefix+test, "TEST"))
		select {
		case v := <-testChannel:
			res = append(res, v)
		case <-time.After(5 * time.Second):
			res = append(res, "An error occurred. INTERNAL TIMEOUT")
		}
	}

	return res
}
