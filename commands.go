package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/andersfylling/disgord"
)

// This file implements all the functions that directly reply to the commands

// helpers

func getDivision(msg *disgord.Message) Division {
	if msg.GuildID != 0 {
		return NewDivision('G', msg.GuildID)
	}
	// if it is not a guild, use the author's ID as the ID
	if msg.Author.ID != 0 {
		return NewDivision('U', msg.Author.ID)
	}
	return NewDivision('T', 0) // test/temporary, discard
}

func getPerms(msg *disgord.Message) (disgord.PermissionBit, error) {
	if msg.GuildID == 0 { // DMs
		return disgord.PermissionBit(math.MaxUint64), nil // every permission feasible
	}
	bit, err := BotClient.Guild(msg.GuildID).Member(msg.Author.ID).GetPermissions()
	if err != nil {
		return 0, err
	}
	return bit, nil
}

func hasPerm(bit disgord.PermissionBit, perm disgord.PermissionBit) bool {
	// easily account for admin permissions
	return bit.Contains(perm) || bit.Contains(disgord.PermissionAdministrator) || bit.Contains(disgord.PermissionAll)
}

type WikiRes struct {
	BatchComplete bool `json:"batchcomplete"`
	Query         struct {
		Pages []struct {
			PageID     int    `json:"pageid"`
			Ns         int    `json:"ns"`
			Title      string `json:"title"`
			Missing    bool   `json:"missing"`
			Categories []struct {
				Ns    int    `json:"ns"`
				Title string `json:"title"`
			} `json:"categories"`
			Extract string `json:"extract"`
		} `json:"pages"`
	} `json:"query"`
}

func queryWiktionary(word string) (*WikiRes, error) {
	// query
	fmted := strings.ReplaceAll(word, " ", "_")
	fmted = url.QueryEscape(fmted)
	url := "https://simple.wiktionary.org/w/api.php?action=query&format=json&prop=categories%7Cextracts&titles=" + fmted + "&formatversion=2&explaintext=1"
	res, err := http.Get(url)

	// handle
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, errors.New(res.Status)
	}

	// read
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// format
	var resp WikiRes
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// templates

func msgerr(err error, msg *disgord.Message, s *disgord.Session) {
	if err != nil {
		msg.Reply(context.Background(), *s, "An error occured. Please report this as a bug.```prolog\n"+err.Error()+"```")
		fmt.Printf("\033[31mError handling message\nAuthor: %s (%d)\nContent: \"%s\"\nError: ", msg.Author.Tag(), msg.Author.ID, msg.Content)
		fmt.Println(err, "\033[0m")
	} else {
		fmt.Printf("Responded to \"%s\" from %d\n", msg.Content, msg.Author.ID) // logging
	}
}

func baseReply(msg *disgord.Message, s *disgord.Session, reply string) {
	_, err := msg.Reply(context.Background(), *s, disgord.Message{
		Content: reply,
		MessageReference: &disgord.MessageReference{ // "reply" client feature
			MessageID: msg.ID,
			ChannelID: msg.ChannelID,
			GuildID:   msg.GuildID,
		},
	})
	msgerr(err, msg, s)
}

func baseEmbedReply(msg *disgord.Message, s *disgord.Session, embed *disgord.Embed) {
	_, err := msg.Reply(context.Background(), *s, disgord.Message{
		Embeds: []*disgord.Embed{embed},
		MessageReference: &disgord.MessageReference{ // "reply" client feature
			MessageID: msg.ID,
			ChannelID: msg.ChannelID,
			GuildID:   msg.GuildID,
		},
	})
	msgerr(err, msg, s)
}

func baseEmbedDMReply(msg *disgord.Message, s *disgord.Session, embed *disgord.Embed, errorMessage string) {
	_, _, err := msg.Author.SendMsg(context.Background(), *s, &disgord.Message{ // DM feature
		Embeds: []*disgord.Embed{embed},
	})
	if err != nil && errorMessage != "" { // typical error is user having DMs disabled
		baseReply(msg, s, errorMessage) // this covers network errors because it also handles errors
	}
}

func baseTextFileReply(msg *disgord.Message, s *disgord.Session, content string, fileName string, fileContents string) {
	_, err := msg.Reply(context.Background(), *s, disgord.CreateMessage{
		Content: content,
		Files: []disgord.CreateMessageFile{
			{
				FileName: fileName,
				Reader:   strings.NewReader(fileContents),
			},
		},
	})
	msgerr(err, msg, s)
}

// -- handlers --

// simple response

func defaultTODOResponse(msg *disgord.Message, s *disgord.Session) {
	baseReply(msg, s, "This feature is incomplete. Don't worry, it's coming!")
}

func defaultResponse(msg *disgord.Message, s *disgord.Session) {
	baseReply(msg, s, defaultResponses[GRand.Intn(len(defaultResponses))])
}

func helpResponse(msg *disgord.Message, s *disgord.Session) {
	eFirst := &disgord.Embed{
		Color: 0x0096ff,
		Author: &disgord.EmbedAuthor{
			Name:    "TRAS Command List",
			IconURL: BotPFP,
		},
		Description: "**------------------------**\n",
		Fields: []*disgord.EmbedField{
			{
				Name:  "_ _\n@TRAS help",
				Value: "Summons this help list.",
			},
			{
				Name:  "_ _\n@TRAS about",
				Value: "Gives information about the bot. Add \"NoCB\" for devices that don't support links with command blocks.",
			},
		},
	}
	eSecond := &disgord.Embed{
		Color: 0x0096ff,
		Author: &disgord.EmbedAuthor{
			Name: "--Primary Commands--",
		},
		Fields: []*disgord.EmbedField{
			{
				Name:  "_ _\n@TRAS oof",
				Value: "Mega OOF",
			},
			{
				Name:  "_ _\n@TRAS f",
				Value: "Mega F",
			},
			{
				Name:  "_ _\n@TRAS pi",
				Value: "First 1 million digits of Pi",
			},
			{
				Name:  "_ _\n@TRAS big",
				Value: "Make a larger version of word/text made of the word. Starts getting wonky with emojis. Becomes file over 520 characters. You can enable thin letters with -t or --thin.\n*Format: @TRAS big (-t/--thin) [letter] [text]*",
			},
			{
				Name:  "_ _\n@TRAS jumble",
				Value: "Jumbles the words in a sentence so it's confusing to read.\n*Format: @TRAS jumble [text]*",
			},
			{
				Name:  "_ _\n@TRAS emojify",
				Value: "Turn all characters into emojis.\n*Format: @TRAS emojify [text]*",
			},
			{
				Name:  "_ _\n@TRAS flagify",
				Value: "Turn all letters into regional indicators, creating flags.\n*Format: @TRAS flagify [text]*",
			},
			{
				Name:  "_ _\n@TRAS superscript",
				Value: "Turn all numbers and letters plus a few math symbols into superscript. Some letters are always lowercase or replaced with something similar due to Unicode limitations.\n*Format: @TRAS superscript [text]*",
			},
			{
				Name:  "_ _\n@TRAS unicodify",
				Value: "Turn all numbers and letters into a non-Latin equivalent.\n*Format: @TRAS unicodify [text]*",
			},
			{
				Name:  "_ _\n@TRAS bold",
				Value: "Bolds all Latin letters and numbers using Unicode.\n*Format: @TRAS bold [text]*",
			},
			{
				Name:  "_ _\n@TRAS replace",
				Value: "Replaces every appearance of a set item with a set replacement.\n*Format: @TRAS replace [item] [replacement] [text]*",
			},
			{
				Name:  "_ _\n@TRAS overcomplicate",
				Value: "Replaces all words with synonyms of the word.\n*Format: @TRAS overcomplicate [text]*",
			},
			{
				Name:  "_ _\n@TRAS word info",
				Value: "Get the definition or Part-of-Speech of a word.\n*Format: @TRAS word info [definition/pos] [word]*",
			},
			{
				Name:  "_ _\n@TRAS ascii art",
				Value: "Generate ascii art. Over 15 characters responds with a file.\n*Format: @TRAS ascii art [text/{font:[Font (use \"\\ \" as space)]}/{getFonts}] [(font)text]*",
			},
			{
				Name:  "_ _\n@TRAS commands",
				Value: "View and manage custom server commands, managing requires 'Manage Messages' perms. Scheduled commands feature requires TRAS Deluxe TBD.\n*Format:@TRAS commands [manage/view] [(manage)...]*\n*Format (manage): @TRAS commands manage [set/delete/schedule] [(set/delete)trigger//(schedule)time of day (hh:mm:ss)] [(set/schedule)reply]*",
			},
			{
				Name:  "_ _\n@TRAS rank",
				Value: "Shows your rank, lets your reset your rank, and allows you to roll dice for a new rank if it's enabled. Admins get other commands as well. Dice rolling disabled by default.\n*Format: @TRAS rank [info|checkDice|dice|set(admin)|reset(part admin)|diceToggle(admin)] [user(4resetORset,admin)|amount(4set,admin)|-real(4info)] [amount(4set,admin)]*",
			},
			{
				Name:  "_ _\n@TRAS set nickname",
				Value: "Set the bot's Nickname on the server. Reset with '{RESET}'. Requires 'Manage Messages' or 'Change Nicknames'.\n*Format: @TRAS set nickname [nickname/{RESET}]*",
			},
			{
				Name:  "_ _\n@TRAS speak",
				Value: "Generate a sentence, repeat messages (requires send perms), and toggle and get status of fallback generated messages. Toggling requires 'Manage Messages' perms. Fallback messages off by default.\n*Format: @TRAS speak [generate/randomspeak] [(randomspeak)on/off/status//(generate)starter]*",
			},
			{
				Name:  "_ _\n@TRAS combinations",
				Value: "Sends file with all possible combinations of the units you have selected and given.\n*Format: @TRAS combinations [words/characters] [items]*",
			},
			{
				Name:  "_ _\n@TRAS ping",
				Value: "Check if the bot is alive. Add 'info' or 'information' for latency data.",
			},
		},
	}
	eThird := &disgord.Embed{
		Color: 0x0096ff,
		Author: &disgord.EmbedAuthor{
			Name: "--Alternatively Triggered Commands--",
		},
		Fields: []*disgord.EmbedField{
			{
				Name:  "Default fallback (mention with no valid command)",
				Value: "I reply, \"What's up?\"",
			},
			{
				Name:  "Generated messages",
				Value: "Fully generated messages *(not an AI so they're completely nonsensical)* can be toggled as the fallback instead of the default response.",
			},
		},
	}

	ccField, err := getGuildCustomCommandsFields(getDivision(msg)) // compatible with user and guild
	if err != nil {
		ccField = []*disgord.EmbedField{
			{
				Name:  "_ _\nError fetching custom commands",
				Value: err.Error(),
			},
		}
	}
	eFourth := &disgord.Embed{

		Color: 0x0096ff,
		Author: &disgord.EmbedAuthor{
			Name: "--Server-Specific Commands--",
		},
		Description: "*For the server this message was activated from*",
		Fields:      ccField,
	}

	// Has to be several messages due to embed size limitations
	baseReply(msg, s, helpCommandResponses[GRand.Intn(len(helpCommandResponses))]) // random help command response
	baseEmbedDMReply(msg, s, eFirst, "Your DMs are not open! Feel free to check out the commmands on https://tras.almostd.one.")
	baseEmbedDMReply(msg, s, eSecond, "")
	baseEmbedDMReply(msg, s, eThird, "")
	baseEmbedDMReply(msg, s, eFourth, "")
}

func aboutResponse(msg *disgord.Message, s *disgord.Session, nocb bool) {
	content := strings.ReplaceAll(BOT_ABOUT_INFO, "'", "`")
	if nocb {
		content = strings.ReplaceAll(content, "```md", "")
		content = strings.ReplaceAll(content, "```prolog", "")
		content = strings.ReplaceAll(content, "```py", "")
		content = strings.ReplaceAll(content, "```", "")
	}
	embed := &disgord.Embed{
		Color: 0x0096ff,
		Author: &disgord.EmbedAuthor{
			Name:    "About TRAS",
			IconURL: BotPFP,
		},
		Description: content,
		Thumbnail: &disgord.EmbedThumbnail{
			URL: "https://tras.almostd.one/img/traslogo.png",
		},
	}

	err := msg.React(context.Background(), *s, "👍")
	if err != nil {
		println(err.Error())
	}
	baseEmbedDMReply(msg, s, embed, "Your DMs are not open! Feel free to find the information on https://tras.almostd.one.")
}

func piResponse(msg *disgord.Message, s *disgord.Session) {
	embed := &disgord.Embed{
		Color:       0x0096ff,
		Title:       "Here's the first 1 million (10⁶) digits of Pi.",
		Description: "First 20: `3.1415926535897932384`\n\n[Download the rest](https://gist.githubusercontent.com/10Nates/95788a4abdd525d7d4dc15d3d45e32ae/raw/80987b58467d10353f0c2bc4ab2d1df8f127ca1c/pi-1mil.txt)",
	}
	// on TRAS 2, this was a file attachment, however
	// because of how text files appear now, it looks bad.
	// A link works the same regardless.

	baseEmbedReply(msg, s, embed)
}

func pingResponse(info bool, msg *disgord.Message, s *disgord.Session, procTimeStart time.Time) {
	if !info {
		baseReply(msg, s, "Pong!")
		return
	}

	hbTime, err := BotClient.AvgHeartbeatLatency()
	if err != nil {
		msgerr(err, msg, s)
		return
	}

	procTime := time.Since(procTimeStart)

	m, err := msg.Reply(context.Background(), *s, "Pong!") // end message
	if err != nil {
		msgerr(err, msg, s)
		return
	}

	resp := "Pong!\n" + // build response
		"`Average Heartbeat: " + hbTime.Truncate(time.Microsecond).String() + "`\n" +
		"`Processing Time:   " + procTime.String() + " `\n" +
		"`Response Latency:  " + m.Timestamp.Sub(msg.Timestamp.Time).String() + "`\n" +
		"*Response Latency is response msg date - initial msg date*"

	_, err = BotClient.Channel(msg.ChannelID).Message(m.ID).Update(&disgord.UpdateMessage{ // edit message
		Content: &resp,
	})
	msgerr(err, msg, s)
}

// simple replace

func emojifyResponse(text string, msg *disgord.Message, s *disgord.Session) {
	respText := text
	for k, v := range emojifyReplacements { // replace key with value
		respText = strings.ReplaceAll(respText, k, v)
	}

	//respond
	baseReply(msg, s, respText)
}

func flagifyResponse(text string, msg *disgord.Message, s *disgord.Session) {
	respText := text
	for k, v := range flagifyReplacements { // replace key with value
		respText = strings.ReplaceAll(respText, k, v)
	}

	//respond
	baseReply(msg, s, respText)
}

func superScriptResponse(text string, msg *disgord.Message, s *disgord.Session) {
	respText := text
	for k, v := range superScriptReplacements { // replace key with value
		respText = strings.ReplaceAll(respText, k, v)
	}

	//respond
	baseReply(msg, s, respText)
}

func unicodifyResponse(text string, msg *disgord.Message, s *disgord.Session) {
	respText := text
	for k, v := range unicodifyReplacements { // replace key with value
		respText = strings.ReplaceAll(respText, k, v)
	}

	//respond
	baseReply(msg, s, respText)
}

func boldResponse(text string, msg *disgord.Message, s *disgord.Session) {
	respText := text
	for k, v := range boldReplacements { // replace key with value
		respText = strings.ReplaceAll(respText, k, v)
	}

	//respond
	baseReply(msg, s, respText)
}

// complex replace

func replaceResponse(item string, replacement string, text string, msg *disgord.Message, s *disgord.Session) {
	respText := strings.ReplaceAll(text, item, replacement) // straight in, no need for filtering unless I'm mistaken
	if len(respText) > 2000 {                               // discord character limit
		baseTextFileReply(msg, s, "Your request didn't fit in a message, so I made it a file.", "replacement.txt", respText)
	} else {
		baseReply(msg, s, respText)
	}
}

func jumbleResponse(args []string, msg *disgord.Message, s *disgord.Session) {
	// This requires some level of explanation for what it's supposed to do. It's supposed
	// to sort of break down the grammar while maintaining the meaning. Basically it shifts
	// words around slightly but not significantly in a large body. Given enough passes
	// through entropy would take over though. This only does one pass.

	mangle := args
	base := strings.Join(args, " ")
	mod := strings.Join(args, " ")
	r := GRand

	// ensure it does something
	for mod == base && len(mangle) != 1 { // while loop
		for i := 0; i < len(mangle); i++ {
			if r.Intn(2) == 1 { // 50% chance

				if !(i < len(mangle)-1) {
					continue // prevent overflow
				}

				//swaps 0 with 1
				b := mangle[i+1]
				mangle[i+1] = mangle[i]
				mangle[i] = b
			} else if i > 0 { // prevent underflow

				//swaps 0 with -1
				b := mangle[i-1]
				mangle[i-1] = mangle[i]
				mangle[i] = b
			}
		}
		mod = strings.Join(mangle, " ")
	}
	// finished jumbling

	baseReply(msg, s, mod)
}

func overcompResponse(words []string, msg *disgord.Message, s *disgord.Session) {
	// for every word
	for i := 0; i < len(words); i++ {
		cleaned := punctuationSplitRegex.FindStringSubmatch(words[i]) // 0 is the whole match, so $1 = cleaned[1]
		// check if exists & get array of synonyms
		val, exists := ThesaurusLookup[strings.ToLower(cleaned[2])] // must be lowercase for lookup
		if exists {
			// replace with synonyms
			newval := val[GRand.Intn(len(val))]
			newvalarr := strings.Split(newval, "")

			// Find where capitals were
			capsLocations := capsRegex.FindAllStringIndex(cleaned[2], -1)
			if capsLocations != nil {
				locsA := map[int]bool{} // map to prevent index range issues
				for i := 0; i < len(capsLocations); i++ {
					locsA[capsLocations[i][0]] = true // list of where capitals are
				}

				for i := 0; i < len(newvalarr); i++ {
					val, exists := locsA[i]
					if exists && val {
						newvalarr[i] = strings.ToUpper(newvalarr[i]) // reapply capitals to new string
					}
				}
			}

			words[i] = cleaned[1] + strings.Join(newvalarr, "") + cleaned[3]
		}
	}

	resp := strings.Join(words, " ")
	if len(resp) > 2000 {
		baseTextFileReply(msg, s, "My literacy was too extravagant, so I made it a file.", "overcomplicate.txt", resp)
	} else {
		baseReply(msg, s, resp)
	}
}

// settings

func setNickResponse(newNick string, msg *disgord.Message, s *disgord.Session) {
	perms, err := getPerms(msg)
	if err != nil {
		msgerr(err, msg, s)
		return
	}

	if !(hasPerm(perms, disgord.PermissionManageNicknames) || hasPerm(perms, disgord.PermissionManageMessages)) {
		baseReply(msg, s, "You don't have permission \"Manage Nicknames\" or \"Manage Messages\". Sorry!")
		return
	}

	if msg.GuildID == 0 {
		baseReply(msg, s, "I can't change my nickname in DMs. Sorry!")
		return
	}

	re := ""
	if newNick == "{RESET}" {
		newNick = ""
		re = "re" // tell user the right thing
	}
	_, err = BotClient.Guild(msg.GuildID).SetCurrentUserNick(newNick)
	if err != nil {
		msgerr(err, msg, s)
		return
	}
	baseReply(msg, s, "Nickname "+re+"set!")
}

// generators

func bigTypeRespones(word string, text string, thin bool, msg *disgord.Message, s *disgord.Session) {
	textletters := strings.Split(text, "")
	for i := 0; i < len(textletters); i++ {
		_, exists := BigTypeLetters[textletters[i]]
		if !exists && textletters[i] != " " {
			textletters[i] = "⛝"
		}
	}

	wordLength := len(strings.Split(word, "")) // fixes unicode characters (kind of, they still aren't always monospace)

	inchar := strings.Repeat(" ", wordLength) // spacing within character
	space := strings.Repeat(" ", wordLength*2)
	midchar := "" // in between characters, gets added on 2nd char
	res := [5]string{}

	// construct 5 layers
	for i := 0; i < len(textletters); i++ {
		if i == 1 {
			if !thin {
				midchar = strings.Repeat(" ", wordLength+1)
			} else {
				midchar = strings.Repeat(" ", wordLength) // reduces readibility
			}
		}
		if textletters[i] != " " {
			if !thin {
				res[0] += midchar + BigTypeLetters[textletters[i]]["1"]
				res[1] += midchar + BigTypeLetters[textletters[i]]["2"]
				res[2] += midchar + BigTypeLetters[textletters[i]]["3"]
				res[3] += midchar + BigTypeLetters[textletters[i]]["4"]
				res[4] += midchar + BigTypeLetters[textletters[i]]["5"]
			} else {
				// this has to be done at this step since it shouldn't remove spaces between letters
				res[0] += midchar + strings.ReplaceAll(BigTypeLetters[textletters[i]]["1"], " ", "")
				res[1] += midchar + strings.ReplaceAll(BigTypeLetters[textletters[i]]["2"], " ", "")
				res[2] += midchar + strings.ReplaceAll(BigTypeLetters[textletters[i]]["3"], " ", "")
				res[3] += midchar + strings.ReplaceAll(BigTypeLetters[textletters[i]]["4"], " ", "")
				res[4] += midchar + strings.ReplaceAll(BigTypeLetters[textletters[i]]["5"], " ", "")
			}
		} else {
			res[0] += space
			res[1] += space
			res[2] += space
			res[3] += space
			res[4] += space
		}
	}
	bigString := res[0] + "\n" + res[1] + "\n" + res[2] + "\n" + res[3] + "\n" + res[4]
	// at this point, it's still all placeholder items. It still has to be replaced with the word and correct spacing

	bigString = strings.ReplaceAll(bigString, "_", inchar)
	bigString = strings.ReplaceAll(bigString, "c", word)

	if len(bigString) > 400 {
		baseTextFileReply(msg, s, "The result is over 400 characters, so I made it a file.", "big.txt", bigString)
	} else {
		baseReply(msg, s, "```\n"+bigString+"\n```")
	}
}

func wordInfoReply(info string, word string, msg *disgord.Message, s *disgord.Session) {
	res, err := queryWiktionary(word)
	if err != nil {
		msgerr(err, msg, s)
		return
	}
	q := res.Query.Pages[0]

	if q.Missing {
		baseReply(msg, s, "This word is not available on Wiktionary. Are you sure this is a real English word?\n\n*Note: This is case sensitive**")
		return
	}

	if info == "def" {
		// format
		extract := q.Extract
		extract = strings.ReplaceAll(extract, "\\n", "\n")
		extract = strings.ReplaceAll(extract, "== ", "**__")
		extract = strings.ReplaceAll(extract, " ==", "__**")
		sections := strings.Split(extract, "\n\n")

		// hyphens for subsections
		for i, v := range sections {
			section := strings.Split(v, "\n")
			if len(section) > 2 {
				for i2, v2 := range section[3:] {
					section[i2+3] = "- " + v2
				}
			}
			sections[i] = strings.Join(section, "\n")
		}

		// reply
		extfmted := strings.Join(sections[1:], "\n")
		baseReply(msg, s, extfmted)

	} else if info == "cat" {
		// format
		var categories []string
		for _, v := range q.Categories {
			categories = append(categories, strings.Replace(v.Title, "Category:", "", 1))
		}

		// reply
		catfmtd := "\n**__Categories__**\n- " + strings.Join(categories, "\n- ")
		catfmtd += "\n\n*Note: all Wiktionary categories, not just parts of speech*"

		baseReply(msg, s, catfmtd)

	} else {
		panic("Internal command incorrectly used")
	}
}
