package main

import (
	"github.com/andersfylling/disgord"
)

// This file implements all the functions for handling server-specific custom commands

func getGuildCustomCommandsFields(DID Division) ([]*disgord.EmbedField, error) {
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
			Name:  "_ _\n@TRAS " + cmds[i].key,
			Value: "I respond " + cmds[i].val,
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

func getCustomCommands(guildID Division) ([]*CustomCommand, error) {
	return []*CustomCommand{}, nil // TODO: implement custom commands
}

func newCustomCommand(key string, val string, div Division) *CustomCommand {
	return &CustomCommand{
		key:     key,
		val:     val,
		divType: div.Type(), // U for user, G for guild
		divID:   uint64(div.Snowflake()),
	}
}
