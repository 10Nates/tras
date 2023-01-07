package main

import (
	"github.com/andersfylling/snowflake/v5"
)

// This file implements all the interactions between the database and the internal handlers.
// The user should never interact directly with the database except in extraneous circumstances.

// integrated method for storing information without worrying about
// overlap of IDs between users and guilds, used in various other files
type Division string

func (d *Division) Snowflake() snowflake.Snowflake {
	return snowflake.ParseSnowflakeString(string(*d)[1:]) //faster than strings split + allows removal of weird - in between
}

func (d *Division) Type() byte {
	return string(*d)[0] // first byte of string is div type
}

func NewDivision(divType byte, divID snowflake.Snowflake) Division {
	return Division(string(divType) + divID.HexString())
}

// custom commands, also used in custom.go
type CustomCommand struct {
	key     string
	val     string
	divType byte
	divID   uint64
}
