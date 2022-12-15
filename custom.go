package main

import (
	"github.com/andersfylling/disgord"
	"github.com/andersfylling/snowflake/v5"
)

// This file implements all the functions for handling custom commands

type customCommand struct { // also used in database
	key   string
	val   string
	guild uint64 // snowflake
}

func getGuildCustomCommandsFields(guildID snowflake.Snowflake) ([]*disgord.EmbedField, error) {
	cmds, err := getCustomCommands(guildID)
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

func getCustomCommands(guildID snowflake.Snowflake) ([]*customCommand, error) {
	return []*customCommand{}, nil // TODO: implement custom commands
}
