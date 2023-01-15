package main

import (
	"db"

	"github.com/andersfylling/disgord"
)

// This file implements all the functions for handling server-specific custom commands

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

// internal handlers

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
		baseDMReply(msg, s, v)
	}
}

func handleSetCustomCommand(msg *disgord.Message, s *disgord.Session, key string, value string) {
	div := getDivision(msg)

	_, err := DBConn.SetCustomCommand(key, value, div)
	if err != nil {
		msgerr(err, msg, s)
	}

	baseReply(msg, s, "Command \""+key+"\" set successfully!")
}

func handleDeleteCustomCommand(msg *disgord.Message, s *disgord.Session, key string) {
	div := getDivision(msg)
	err := DBConn.RemoveCustomCommand(key, div)
	if err != nil {
		msgerr(err, msg, s)
	}

	baseReply(msg, s, "Command \""+key+"\" removed successfully!")
}

// parser

func parseCustomCommand(msg *disgord.Message, s *disgord.Session, arg string) bool {
	div := getDivision(msg)

	cmds, err := getCustomCommands(div)
	if err != nil {
		msgerr(err, msg, s)
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
