package main

import (
	"context"
	"fmt"

	"github.com/andersfylling/disgord"
)

// This file implements all the functions that directly reply to the commands

const BOT_ABOUT_INFO = `<Created_by Nathan Hedge>

[Website](https://tras.almostd.one/)
[Add Link](https://bit.ly/gotras)
[Top.gg Page](https://top.gg/bot/494273862427738113)
[Git Repo](https://github.com/10Nates/tras)

[Legal]()
[ ](TRAS operates under the MIT license)
[ ](https://github.com/10Nates/tras/LICENSE)`

// templates
func msgerr(err error, msg *disgord.Message, s *disgord.Session) {
	if err != nil {
		if msg != nil {
			msg.Reply(context.Background(), *s, err.Error())
		}
		fmt.Printf("\033[31mError handling message\nAuthor: %s (%d)\nContent: \"%s\"\nError: ", msg.Author.Tag(), msg.Author.ID, msg.Content)
		fmt.Println(err, "\033[0m")
	} else {
		fmt.Printf("Responded to \"%s\" from %d\n", msg.Content, msg.Author.ID)
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
	_, _, err := msg.Author.SendMsg(context.Background(), *s, &disgord.Message{
		Embeds: []*disgord.Embed{embed},
	})
	if err != nil && errorMessage != "" {
		baseReply(msg, s, errorMessage)
	}
}

// handlers
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
				Value: "Gives information about the bot.",
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
				Value: "Make a larger version of word/text made of the word. Starts getting wonky with emojis. Becomes file over 520 characters. You can enable thin letters with -t or --thin.\n*Format: txt.big [-t|--thin (pushes other vars forward)] [word] [text (optional)]*",
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
				Value: "Get the definition or Part-of-Speech of a word.\n*Format: @TRAS word info [def|pos] [word]*",
			},
			{
				Name:  "_ _\n@TRAS ascii art",
				Value: "Generate ascii art. Over 15 characters responds with a file.\n*Format: @TRAS ascii art [text|{Font:[Font (use '_' as space)]}|{getFonts}] [text]*",
			},
			{
				Name:  "_ _\n@TRAS commands",
				Value: "View and manage custom server commands, managing requires 'Manage Messages' perms.\n*Format: @TRAS commands [manage|view] [set|delete] [activator] [reply (multiword)]*",
			},
			{
				Name:  "_ _\n@TRAS rank",
				Value: "Shows your rank, lets your reset your rank, and allows you to roll dice for a new rank if it's enabled. Admins get other commands as well. Dice rolling disabled by default.\n*Format: @TRAS rank [info|checkDice|dice|set(admin)|reset(part admin)|diceToggle(admin)] [user(4resetORset,admin)|amount(4set,admin)|-real(4info)] [amount(4set,admin)]*",
			},
			{
				Name:  "_ _\n@TRAS set nickname",
				Value: "Set the bot's Nickname on the server. Reset with '{RESET}'. Requires 'Manage Messages' or 'Change Nicknames'.\n*Format: @TRAS set nickname [nickname|{RESET}]*",
			},
			{
				Name:  "_ _\n@TRAS speak",
				Value: "Generate a sentence, repeat messages (requires send perms), and toggle and get status of random generated messages. Toggling requires 'Manage Messages' perms. Random messages off by default.\n*Format: @TRAS speak [generate|repeat|toggleRandSpeak|randSpeakStatus] [channel ID or channel tag] [message]*",
			},
			{
				Name:  "_ _\n@TRAS combinations",
				Value: "Sends file with all possible combinations of the units you have selected and given.\n*Format: @TRAS combinations [words|characters] [items]*",
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
	eFourth := &disgord.Embed{

		Color: 0x0096ff,
		Author: &disgord.EmbedAuthor{
			Name: "--Server-Specific Commands--",
		},
		Description: "*For the server this message was activated from*",
		// Fields:      getGuildCustomCommandsFields(msg.GuildID),
		// TODO: add fields for custom commands

	}

	// Has to be several messages due to embed size limitations
	baseReply(msg, s, "I slid a command list into your DMs.")
	baseEmbedDMReply(msg, s, eFirst, "Your DMs are not open! Feel free to check out the commands in  any of the linked sites at @TRAS about.")
	baseEmbedDMReply(msg, s, eSecond, "")
	baseEmbedDMReply(msg, s, eThird, "")
	baseEmbedDMReply(msg, s, eFourth, "")
}

func aboutResponse(msg *disgord.Message, s *disgord.Session) {
	embed := &disgord.Embed{
		Color: 0x0096ff,
		Author: &disgord.EmbedAuthor{
			Name:    "About TRAS",
			IconURL: BotPFP,
		},
		Title:       "Text Response Automation System",
		Description: "```md\n<Version " + BOT_VERSION + ">\n" + BOT_ABOUT_INFO + "```",
		Thumbnail: &disgord.EmbedThumbnail{
			URL: "https://tras.almostd.one/img/traslogo.png",
		},
	}
	baseEmbedReply(msg, s, embed)
}
