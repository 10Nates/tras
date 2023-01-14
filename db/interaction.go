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
	// start transaction
	tx, err := c.DB.Begin()
	if err != nil {
		return nil, err
	}

	// add command to database
	cmd := &CustomCommand{
		ID:  uuid.New(),
		Key: key,
		Val: value,
	}
	_, err = tx.Model(cmd).Insert()
	if err != nil {
		return nil, err
	}

	// fetch division data
	divData := &DivisionData{
		Div: div,
	}
	_, err = tx.Model(divData).SelectOrInsert()
	if err != nil {
		return nil, err
	}

	// connect command to division
	divData.Cmds = append(divData.Cmds, *cmd)

	_, err = tx.Model(divData).WherePK().Update()
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	// Check if insertion was successful (unneccessary)
	// divDataCheck := &DivisionData{
	// 	Div: div,
	// }
	// c.DB.Model(divDataCheck).WherePK().Select()
	// pass := false
	// for _, cc := range divDataCheck.Cmds {
	// 	if cc.ID == cmd.ID {
	// 		pass = true
	// 	}
	// }
	// if !pass {
	// 	return nil, errors.New("not inserted correctly")
	// }

	return cmd, nil
}

func (c *Connection) GetDivsion(div Division) (*DivisionData, error) {
	// start transaction
	tx, err := c.DB.Begin()
	if err != nil {
		return nil, err
	}

	// fetch division data
	divData := &DivisionData{
		Div: div,
	}
	_, err = tx.Model(divData).SelectOrInsert()
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return divData, nil
}
