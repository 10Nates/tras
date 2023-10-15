package main

import (
	"context"
	"db"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/andersfylling/disgord"
	"github.com/andersfylling/disgord/std"
)

var BotID string // loaded on init
var BotPFP string
var BigTypeLetters map[string]map[string]string // this is way easier than the alternative
var ThesaurusLookup map[string][]string
var GRand = rand.New(rand.NewSource(time.Now().UnixNano()))
var DBConn *db.Connection

func main() {
	// load bigtype letters
	bigtypejson, err := os.ReadFile("src/bigtype.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(bigtypejson, &BigTypeLetters)
	if err != nil {
		panic(err)
	}

	// load thesaurus
	thesarusjson, err := os.ReadFile("src/thesaurus.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(thesarusjson, &ThesaurusLookup)
	if err != nil {
		panic(err)
	}

	// connect to DB
	DBConn = &db.Connection{
		Host:     DB_HOST,
		Port:     DB_PORT,
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   DB_NAME,
	}
	err = DBConn.Connect()
	if err != nil {
		panic(err)
	}
	DBConn.CloseOnInterrupt()

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
		WithMiddleware(content.NotByBot, content.NotByWebhook, content.ContainsBotMention, content.HasBotMentionPrefix). // filter
		MessageCreate(func(s disgord.Session, evt *disgord.MessageCreate) {                                              // on message
			// used for standard message parsing
			go parseCommand(evt.Message, &s)
		})

	client.Gateway().
		WithMiddleware(content.NotByBot, content.NotByWebhook).
		MessageCreate(func(s disgord.Session, evt *disgord.MessageCreate) { // on message (any)
			if content.ContainsBotMention(evt) != nil { // middleware !content.ContainsBotMention
				return
			}
			// used for ranking and randomspeak

			updateMemberProgress(evt.Message)
			// TODO: randomspeak
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

	for i := 1; i < len(carr); i++ { // first argument should always be bot mention
		args = append(args, carr[i])
		argsl = append(argsl, strings.ToLower(carr[i]))
	}

	if len(args) < 1 { // prevent error in switch case
		args = append(args, "")
		argsl = append(argsl, "")
	}

	// custom commands never override standard
	// commands to prevent deadlock
	successful_cc := parseCustomCommand(msg, s, argsl[0])

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
		if len(argsl) > 1 && (argsl[1] == "-t" || argsl[1] == "--thin") {
			if len(argsl) > 3 {
				text := strings.Join(argsl[3:], " ") // case insensitive
				bigTypeRespones(args[2], text, true, msg, s)
			} else if len(argsl) == 3 {
				bigTypeRespones(args[2], argsl[2], true, msg, s) // word is case sensitive but text is not
			} else {
				baseReply(msg, s, "I need to know the [word] to enlarge, OR the [word] and [text] to enlarge with it.")
			}
		} else {
			if len(argsl) > 2 {
				text := strings.Join(argsl[2:], " ") // case insensitive
				bigTypeRespones(args[1], text, false, msg, s)
			} else if len(argsl) == 2 {
				bigTypeRespones(args[1], argsl[1], false, msg, s) // word is case sensitive but text is not
			} else {
				baseReply(msg, s, "I need to know the [word] to enlarge, OR the [word] and [text] to enlarge with it.")
			}
		}
	case "jumble":
		if len(argsl) > 1 {
			jumbleResponse(args[1:], msg, s) // case sensitive, presplit
		} else {
			baseReply(msg, s, "What should I jumble?")
		}
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
	case "replace", "rep":
		if len(argsl) > 3 {
			text := strings.Join(args[3:], " ") // case sensitive
			replaceResponse(args[1], args[2], text, msg, s)
		} else {
			baseReply(msg, s, "Tell me the [what to replace], the [replacement], and then provide the [body of text].")
		}
	case "overcomplicate", "overcomp":
		if len(argsl) > 1 {
			overcompResponse(args[1:], msg, s) // case sensitive
		} else {
			baseReply(msg, s, "Whichever lexical constructions shall I reform?")
		}
	case "word":
		if len(argsl) > 1 && argsl[1] == "info" {
			if len(argsl) > 2 { // type
				if len(argsl) > 3 { // word
					word := strings.Join(args[3:], " ")

					switch argsl[2] {
					case "definition", "definitions", "define", "def", "defs":
						wordInfoReply("def", word, msg, s)

					case "categories", "category", "cat", "cats", "partofspeech", "pos":
						wordInfoReply("cat", word, msg, s)

					default:
						baseReply(msg, s, "That's not an info type I can provide.")
					}

				} else {
					baseReply(msg, s, "What word do you want info on?")
				}
			} else {
				baseReply(msg, s, "What type of info do you want? Defintion, or categories?")
			}
		} else {
			defaultResponse(msg, s, successful_cc)
		}
	case "ascii":
		if len(argsl) > 1 && argsl[1] == "art" {
			if len(argsl) > 2 {
				if argsl[2] == "getfonts" {
					asciiGetFonts(msg, s)
				} else {
					if len(argsl) > 3 {
						text := strings.Join(args[3:], " ")
						asciiResponse(msg, s, args[2], text, 100) // max width 100, maybe add option in future
					} else {
						baseReply(msg, s, "What text do you want to be generated?")
					}
				}
			} else {
				baseReply(msg, s, "I need to know the [font] and the [text] you want me to use.")
			}
		} else {
			defaultResponse(msg, s, successful_cc)
		}
	case "commands", "cmds":
		if len(argsl) > 1 && (argsl[1] == "view" || argsl[1] == "list") {
			handleViewCustomCommands(msg, s)
		} else if len(argsl) > 1 && argsl[1] == "manage" {
			// check for permissions
			perms, err := getPerms(msg, s)
			if err != nil {
				msgerr(err, msg, s)
				return
			}
			if !hasPerm(perms, disgord.PermissionAdministrator) {
				baseReply(msg, s, "You don't have administrator permission. Sorry!")
				return
			}

			// restricted cases
			word := ""
			if len(argsl) > 2 {
				word = argsl[2]
			}
			switch word {
			case "set":
				if len(argsl) > 4 {
					text := strings.Join(args[4:], " ")
					handleSetCustomCommand(msg, s, argsl[3], text)
				} else {
					baseReply(msg, s, "You need to tell me the [trigger] and [what I should respond with].")
				}
			case "delete", "del", "remove", "rem", "reset":
				if len(argsl) > 3 {
					handleDeleteCustomCommand(msg, s, argsl[3])
				} else {
					baseReply(msg, s, "You need to tell me the [trigger] to delete.")
				}
			case "schedule":
				defaultTODOResponse(msg, s) // TODO: schedule feature
			default:
				baseReply(msg, s, "The format for manage is `@TRAS commands manage [set/delete/schedule] [(set/delete)trigger//(schedule)time of day (hh:mm:ss)] [(set/schedule)reply]`")
			}
		} else {
			baseReply(msg, s, "Would you like to [view] the commands, or [manage] them as the admin?")
		}
	case "rank":
		if len(argsl) > 1 {
			switch argsl[1] {
			case "info":
				baseReply(msg, s, "TRAS' \"progress\" meter takes various elements of your messages' metadata into account when valuing them.\n"+
					"This is a differnet approach than TRAS 2 due to the API changes between the sunset of v2 and the creation of v3.\n"+
					"Levels are the logarithm of your \"progress\" to base 2, meaning you require 2 times the \"progress\" per level."+
					"\nI included the \"dice roll\" feature as a fun gimmick to portray my thoughts about levels in general - pointless and silly.")
			case "checkdice":
				statusStr := "OFF"
				status, err := getDiceStatus(msg)
				if err != nil {
					msgerr(err, msg, s)
					return
				}
				if status {
					statusStr = "ON"
				}
				baseReply(msg, s, "Dice rolls are currently "+statusStr)
			case "dice":
				defaultTODOResponse(msg, s) // TODO: rank dice
			case "set", "reset", "toggledice":
				// check for permissions
				perms, err := getPerms(msg, s)
				if err != nil {
					msgerr(err, msg, s)
					return
				}
				if !hasPerm(perms, disgord.PermissionAdministrator) {
					baseReply(msg, s, "You don't have administrator permission. Sorry!")
					return
				}
				switch argsl[1] {
				case "set":
					if len(argsl) > 3 {
						user, validMention := extractSnowflake(argsl[2])
						if !validMention {
							baseReply(msg, s, "That was not a valid user mention.")
							return
						}

						num, err := strconv.Atoi(argsl[3])
						if err != nil {
							baseReply(msg, s, "Invalid number!")
							return
						}

						err = forceSetUserRank(msg, user, int64(num))
						if err != nil {
							msgerr(err, msg, s)
							return
						}
					} else {
						baseReply(msg, s, "You need to tell me the [user] and the [value] to set the progress to.")
					}
				case "reset":
					if len(argsl) > 2 {
						user, validMention := extractSnowflake(argsl[2])
						if !validMention {
							baseReply(msg, s, "That was not a valid user mention.")
							return
						}
						err := forceSetUserRank(msg, user, 0)
						if err != nil {
							msgerr(err, msg, s)
							return
						}
					} else {
						baseReply(msg, s, "You need to tell me the [user] to reset the progress of.")
					}
				case "toggledice":
					toggleDiceResponse(msg, s)
				}
			default:
				user, validMention := extractSnowflake(argsl[1])
				if validMention {
					getUserRankInfo(msg, s, user)
					return
				}

				// check for permissions
				perms, err := getPerms(msg, s)
				if err != nil {
					msgerr(err, msg, s)
					return
				}
				helpContent := "The format is `@TRAS rank [info/checkDice/dice] [(info)-real]`"
				if hasPerm(perms, disgord.PermissionAdministrator) {
					helpContent += "\nThe format for admin controls is `@TRAS rank [set/reset/toggleDice] [(set/reset)user] [value]`"
				}
				baseReply(msg, s, helpContent)
			}
		} else {
			getUserRankInfo(msg, s, msg.Author.ID)
		}
	case "set":
		if len(argsl) > 1 && (argsl[1] == "nickname" || argsl[1] == "nick") {
			text := strings.Join(args[2:], " ") // case sensitive
			if len(argsl) > 2 {
				setNickResponse(text, msg, s)
			} else {
				baseReply(msg, s, "What should by nickname be?")
			}
		} else {
			defaultResponse(msg, s, successful_cc)
		}
	case "reset":
		if len(argsl) > 1 && (argsl[1] == "nickname" || argsl[1] == "nick") {
			// A more natural way of resetting nickname
			text := "{RESET}" // case sensitive
			setNickResponse(text, msg, s)
		} else {
			defaultResponse(msg, s, successful_cc)
		}
	case "speak":
		defaultTODOResponse(msg, s) // TODO: speak
	case "combinations", "combos", "powerset":
		if len(argsl) > 1 { // option
			if len(argsl) > 2 { // text
				switch argsl[1] {
				case "words", "w":
					combosResponse(args[2:], msg, s)

				case "characters", "chars", "c":
					text := strings.Join(args[2:], " ")
					ltrs := strings.Split(text, "")

					combosResponse(ltrs, msg, s)
				default:
					baseReply(msg, s, "That's not an option.")
				}
			} else {
				baseReply(msg, s, "What do you want the combinations for?")
			}
		} else {
			baseReply(msg, s, "Which combinations do you want? Words, or characters?")
		}
	case "ping":
		if len(argsl) > 1 && (argsl[1] == "info" || argsl[1] == "information") {
			pingResponse(true, msg, s, procTimeStart)
		} else {
			pingResponse(false, msg, s, procTimeStart)
		}
	default:
		defaultResponse(msg, s, successful_cc)
	}
}
