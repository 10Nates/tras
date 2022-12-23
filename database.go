package main

import (
	"strings"

	"github.com/andersfylling/snowflake/v5"
)

// This file implements all the interactions between the database and the internal handlers.
// The user should never interact directly with the database except in extraneous circumstances.

// integrated method for storing information without worrying about
// overlap of IDs between users and guilds, used in various other files
type Division string

func (d *Division) Snowflake() snowflake.Snowflake {
	return snowflake.ParseSnowflakeString(strings.Split(string(*d), "-")[1])
}

func (d *Division) Type() byte {
	return strings.Split(string(*d), "-")[0][0]
}

// custom commands, also used in custom.go
type customCommand struct {
	key     string
	val     string
	divType byte
	divID   uint64
}
