package main

import (
	"db"
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
