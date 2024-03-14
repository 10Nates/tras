package main

import (
	"db"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"time"

	"github.com/andersfylling/disgord"
	"github.com/andersfylling/snowflake/v5"
)

// This file implements special features, such as custom commands and ranking

// -- Custom commands --

// helpers

func getGuildCustomCommandsFields(DID db.Division) ([]*disgord.EmbedField, error) {
	cmds, err := getCustomCommands(DID)
	if err != nil {
		return nil, err
	}

	if len(cmds) > 25 { // discord embed cap
		return []*disgord.EmbedField{
			{
				Name:  "_ _\nThe server's number of commands exceeds discord's embed field cap of 25. Use \"@TRAS commands view\" instead.",
				Value: "_ _",
			},
		}, nil
	}

	newEmbedFields := []*disgord.EmbedField{}
	for i := 0; i < len(cmds); i++ { // generative embeds
		newEmbedFields = append(newEmbedFields, &disgord.EmbedField{
			Name:  "_ _\n@TRAS " + cmds[i].Key,
			Value: "I respond " + cmds[i].Val,
		})
	}
	if len(newEmbedFields) == 0 { // no embeds on server
		newEmbedFields = append(newEmbedFields, &disgord.EmbedField{
			Name:  "_ _\nNo custom commands are currently on this server",
			Value: "_ _",
		})
	}

	return newEmbedFields, nil
}

func getCustomCommands(div db.Division) ([]*db.CustomCommand, error) {
	divData, err := DBConn.GetDivsion(div)
	if err != nil {
		return nil, err
	}

	return divData.Cmds, nil
}

// handlers

func handleViewCustomCommands(msg *disgord.Message, s *disgord.Session) {
	div := getDivision(msg)
	cmds, err := getCustomCommands(div)
	if err != nil {
		msgerr(err, msg, s)
		return
	}

	respArr := []string{"**__ Commands List:__** \n"}
	c := 0

	for _, cc := range cmds {
		entry := "- \"" + cc.Key + "\", returns: \"" + cc.Val + "\""
		if len(respArr[c]+entry+"\n") > 2000 { // if too large to fit in single message
			c++
			respArr = append(respArr, entry) // expand array
		} else {
			respArr[c] += entry + "\n"
		}
	}

	baseReact(msg, s, "👍")
	for _, v := range respArr {
		baseDMReply(msg, s, v, nil)
	}
}

func handleSetCustomCommand(msg *disgord.Message, s *disgord.Session, key string, value string) {
	div := getDivision(msg)

	_, err := DBConn.SetCustomCommand(key, value, div)
	if err != nil {
		msgerr(err, msg, s)
		return
	}

	baseReply(msg, s, "Command \""+key+"\" set successfully!")
}

func handleDeleteCustomCommand(msg *disgord.Message, s *disgord.Session, key string) {
	div := getDivision(msg)
	err := DBConn.RemoveCustomCommand(key, div)
	if err != nil {
		msgerr(err, msg, s)
		return
	}

	baseReply(msg, s, "Command \""+key+"\" removed successfully!")
}

// parser

func parseCustomCommand(msg *disgord.Message, s *disgord.Session, arg string) bool {
	div := getDivision(msg)

	cmds, err := getCustomCommands(div)
	if err != nil {
		msgerr(err, msg, s) // msgerr is warranted here because we know that they at least pinged the bot
		return false
	}

	for _, cc := range cmds {
		if arg == cc.Key {
			baseReply(msg, s, cc.Val)
			return true
			// break
		}
	}

	return false
}

// -- Ranking --

// helpers

// func baseAttentionScore(timeDiff time.Duration) float64 {
// 	x := timeDiff.Seconds()
// 	score := math.Max(0, (-1.0/125.0)*(600-180*x+math.Pow(x, 2))*math.Min(1, 10.0/math.Abs(-45+x)))
// 	return score
// }

func calcNewMemberProgress(msg *disgord.Message) (int64, error) {
	if msg.MentionEveryone {
		// never adds score if it mentions everyone
		return 0, nil
	}

	div := getDivision(msg)

	rankMem, err := DBConn.GetRankMember(msg.Author.ID, div)
	if err != nil {
		return 0, err
	}

	// base attention score is based on time between messages.
	// I played around in desmos for a while and found the equation I liked,
	// then I asked wolframalpha to simplify it.
	tsDiff := msg.Timestamp.Time.Sub(rankMem.LastMsgTs).Seconds()
	score := math.Max(0, (-1.0/125.0)*(600-180*tsDiff+math.Pow(tsDiff, 2))*math.Min(1, 10.0/math.Abs(-45+tsDiff)))

	if msg.MessageReference != nil { // replying to someone else

		score *= 2

		if tsDiff > 3.3 { // prevent gaming spam filter with message reference

			// this inflates base score since you are
			// extrememly likely to be "attentive" to what
			// you responded to regardless of time difference
			score += 10
		}
	}

	if msg.ChannelID != snowflake.Snowflake(rankMem.LastChanID) { // not on the same channel
		score *= 0.5
	}

	if 3 > len(msg.Mentions) && len(msg.Mentions) > 0 {
		if msg.Mentions[0].ID != msg.Author.ID && ((len(msg.Mentions) > 1 && msg.Mentions[1].ID != msg.Author.ID) || len(msg.Mentions) < 2) { // not mentioning self
			score *= 1.1 // if there is 1 or 2 mentions, increase the score slightly
		}
	}

	newProg := int64(score) + rankMem.Progress

	return newProg, nil
}

func updateMemberProgress(msg *disgord.Message) error {
	//calculate
	newProg, err := calcNewMemberProgress(msg)
	if err != nil {
		// Because this runs on every message, returning an error would be a nuisance in the event
		// of a repeating failure. As such, it is only logged.
		logmsgerr(msg, err)
		return err
	}
	div := getDivision(msg)
	DBConn.SetRankMemberProgress(msg, msg.Author.ID, div, newProg)
	return nil
}

func getDiceStatus(msg *disgord.Message) (bool, error) {
	div := getDivision(msg)
	data, err := DBConn.GetDivsion(div)
	if err != nil {
		return false, err
	}

	return data.Dice, nil
}

func toggleDiceStatus(msg *disgord.Message) (bool, error) {
	curStat, err := getDiceStatus(msg)
	if err != nil {
		return false, err
	}

	err = DBConn.SetDiceAvailability(getDivision(msg), !curStat) // flip status
	if err != nil {
		return false, err
	}

	return !curStat, nil
}

func forceSetUserRank(msg *disgord.Message, uID disgord.Snowflake, newProgress int64) error {
	err := DBConn.SetRankMemberProgress(msg, uID, getDivision(msg), newProgress)
	return err
}

// handlers

func diceRollResponse(msg *disgord.Message, s *disgord.Session) {
	// Sets your progress to a random value within 100 levels
	rand.Seed(time.Now().UnixNano())
	newLevel := rand.Float64() * 100
	newProgress := int64(math.Pow(float64(newLevel), 2))
	err := forceSetUserRank(msg, msg.Author.ID, newProgress)
	if err != nil {
		msgerr(err, msg, s)
		return
	}

	// Modified from commands.go/getUserRankInfo
	levelStr := strconv.Itoa(int(newLevel))
	progStr := strconv.Itoa(int(newProgress))
	nextMilestone := strconv.Itoa(int(math.Pow(math.Floor(newLevel)+1, 2)))

	baseReply(msg, s, "Dice rolled! Your stats are now:\n"+
		"Level:"+levelStr+"\n"+"Progress:"+progStr+"/"+nextMilestone)
}

// -- Random Speak --

// helpers

type RandSpeakData struct {
	status        bool
	LastRandSpeak time.Time
}

func getRandSpeakInfo(msg *disgord.Message) (*RandSpeakData, error) {
	div := getDivision(msg)
	data, err := DBConn.GetDivsion(div)
	if err != nil {
		return nil, err
	}

	return &RandSpeakData{
		status:        data.RandSpeak,
		LastRandSpeak: data.LastRandSpeak,
	}, nil
}

func executeRandSpeakRoll(msg *disgord.Message, s *disgord.Session) error {
	rsdata, err := getRandSpeakInfo(msg)
	if err != nil {
		return err
	}

	if !rsdata.status {
		// randomSpeak disabled
		return nil
	}

	probabilityWeight := -math.Pow(math.E, float64(time.Now().Unix()-rsdata.LastRandSpeak.Unix())*(-1/60.0)) + 1
	if GRand.Float64()*25 < probabilityWeight { // max odds 1 in 25, min odds 0 (immediately after last randSpeak)
		DBConn.SetLastRandomSpeakTime(getDivision(msg), time.Now())
		randSpeakGenerateResponse(msg, s, "")
	}

	return nil
}

// -- Data Management --

// helpers

type dataMeDownloadRes struct {
	Data []*db.RankMemberExport
}

type deleteNonce struct {
	UserID disgord.Snowflake
	DivID  uint64
	Exp    time.Time
	Val    string
}

var activeNonces = []deleteNonce{}

func randomString(length int) string {
	const charset = "ABCDEFGHIJKLMONPQRSTUVWXYZ"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

// handlers

func dataMangementHandler(msg *disgord.Message, h *disgord.InteractionCreate, s *disgord.Session) {
	target := h.Data.Options[0].Value.(string)
	action := h.Data.Options[1].Value.(string)

	switch target {
	case "me":
		// action branch
		if action == "download" {
			unformatted, err := DBConn.FetchAllUserData(msg.Author.ID)
			if err != nil {
				msgerr(err, msg, s)
				return
			}
			// convert to json string
			j, err := json.MarshalIndent(dataMeDownloadRes{Data: unformatted}, "", "  ")
			if err != nil {
				msgerr(err, msg, s)
				return
			}
			baseTextFileDMReply(msg, s, "Here's your data, compiled fresh.", msg.Author.ID.String()+"-data.json", string(j), h)

		} else if action == "delete" {
			// trim expired nonces
			for i, a := range activeNonces {
				if a.Exp.Before(time.Now()) {
					// https://stackoverflow.com/questions/37334119/how-to-delete-an-element-from-a-slice-in-golang
					// order doesn't matter
					activeNonces[i] = activeNonces[len(activeNonces)-1]
					activeNonces = activeNonces[:len(activeNonces)-1]
				}
			}

			if len(h.Data.Options) <= 2 {
				// delete nonce generation
				nonce := randomString(12)
				activeNonces = append(activeNonces, deleteNonce{
					UserID: msg.Author.ID,
					DivID:  getDivision(msg).DivID,
					Exp:    msg.Timestamp.Time.Add(2 * time.Minute),
					Val:    nonce,
				})
				baseDMReply(msg, s, fmt.Sprintf("## Are you *sure* you want to delete your TRAS data? **This action is PERMANENT** and applies to **EVERY server you are in.**\n\n*Note: This does not include DM pseudo-server data, which must be removed separately.*\n\n> Run the command `/mydata target:me action:delete confirmdelete:%s` within the next 2 minutes to confirm.", nonce), h)
				return
			}

			// delete user data if confirmdelete matches nonce
			for i, a := range activeNonces {
				// time preverified with previous trim
				// 				verify nonce							verify user						verify location
				if h.Data.Options[2].Value.(string) == a.Val && msg.Author.ID == a.UserID && getDivision(msg).DivID == a.DivID {

					err := DBConn.RemoveAllUserData(a.UserID)
					if err != nil {
						msgerr(err, msg, s)
						return
					}

					activeNonces[i] = activeNonces[len(activeNonces)-1]
					activeNonces = activeNonces[:len(activeNonces)-1]

					baseDMReply(msg, s, "Your data has been deleted. You are now a clean slate in TRAS-land.", h)
					return
				}
			}

			baseDMReply(msg, s, "Your confirmation could not be verified. Are you sure you typed it correctly?", h)
		}

	case "server":
		// check for perms
		perms, err := getPerms(msg, s)
		if err != nil {
			msgerr(err, msg, s)
			return
		}
		if !hasPerm(perms, disgord.PermissionAdministrator) {
			baseDMReply(msg, s, "You don't have Administrator permissions.", h)
			return
		}

		// action branch
		if action == "download" {
			unformatted, err := DBConn.FetchAllDivisionData(getDivision(msg))
			if err != nil {
				msgerr(err, msg, s)
				return
			}
			// convert to json string
			j, err := json.MarshalIndent(unformatted, "", "  ")
			if err != nil {
				msgerr(err, msg, s)
			}
			baseTextFileDMReply(msg, s, "Here's your server's data, compiled fresh.", h.GuildID.String()+"-data.json", string(j), h)

		} else if action == "delete" {
			// trim expired nonces
			for i, a := range activeNonces {
				if a.Exp.Before(time.Now()) {
					// https://stackoverflow.com/questions/37334119/how-to-delete-an-element-from-a-slice-in-golang
					// order doesn't matter
					activeNonces[i] = activeNonces[len(activeNonces)-1]
					activeNonces = activeNonces[:len(activeNonces)-1]
				}
			}

			if len(h.Data.Options) <= 2 {
				// delete nonce generation
				nonce := randomString(12)
				activeNonces = append(activeNonces, deleteNonce{
					UserID: msg.Author.ID,
					DivID:  getDivision(msg).DivID,
					Exp:    time.Now().Add(2 * time.Minute),
					Val:    nonce,
				})
				baseDMReply(msg, s, fmt.Sprintf("## Are you *sure* you want to delete this server's TRAS data? **This action is PERMANENT.**\n\n> Run the command `/mydata target:server action:delete confirmdelete:%s` within the next 2 minutes to confirm.", nonce), h)
				return
			}

			// delete server data if confirmdelete matches nonce
			for i, a := range activeNonces {
				// time preverified with previous trim
				// 				verify nonce							verify user						verify location
				if h.Data.Options[2].Value.(string) == a.Val && msg.Author.ID == a.UserID && getDivision(msg).DivID == a.DivID {
					err := DBConn.RemoveAllDivisionData(getDivision(msg))
					if err != nil {
						msgerr(err, msg, s)
						return
					}

					activeNonces[i] = activeNonces[len(activeNonces)-1]
					activeNonces = activeNonces[:len(activeNonces)-1]

					baseDMReply(msg, s, "The server's data has been deleted. It is now a clean slate in TRAS-land.", h)
					return
				}
			}

			baseDMReply(msg, s, "Your confirmation could not be verified. Are you sure you typed it correctly?", h)
		}

	default:
		msgerr(errors.New("invalid command structure"), msg, s)
	}
}
