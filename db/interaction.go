package db

import (
	"github.com/andersfylling/snowflake/v5"
	"github.com/google/uuid"
)

// This file implements all the interactions between the database and the internal handlers.
// The user should never interact directly with the database except in extraneous circumstances.

// integrated method for storing information without worrying about
// overlap of IDs between users and guilds, used in various other files
type Division struct {
	DivType byte
	DivID   uint64
}

func (d *Division) Snowflake() snowflake.Snowflake {
	return snowflake.Snowflake(d.DivID)
}

func NewDivision(divType byte, divID snowflake.Snowflake) Division {
	return Division{
		DivType: divType,
		DivID:   uint64(divID),
	}
}

func (c *Connection) AddCustomCommand(key string, value string, div Division) (*CustomCommand, error) {
	cmd := &CustomCommand{
		ID:  uuid.New(),
		Key: key,
		Val: value,
	}
	_, err := c.DB.Model(cmd).Insert()
	if err != nil {
		return nil, err
	}

	return cmd, nil
}
