package main

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"

	"github.com/andersfylling/disgord"
)

// This file implements all the functions that directly reply to the commands

// helpers
func getDivision(msg *disgord.Message) Division {
	if msg.GuildID != 0 {
		return Division("G-" + msg.GuildID.HexString())
	}
	return Division("U-" + msg.Author.ID.HexString()) // if it is not a guild, use the author's ID as the ID
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

// handlers

// simple response
func defaultResponse(msg *disgord.Message, s *disgord.Session) {
	baseReply(msg, s, "What's up?")
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
	baseReply(msg, s, HELP_COMMAND_RESPONSES[rand.Intn(len(HELP_COMMAND_RESPONSES))]) // random help command response
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

func jumbleResponse(msg *disgord.Message, s *disgord.Session) {

}

func overcompResponse(msg *disgord.Message, s *disgord.Session) {

}

// manipulation without database modification

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
