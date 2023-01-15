package db

import (
	"errors"

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

func (c *Connection) SetCustomCommand(key string, value string, div Division) (*CustomCommand, error) {
	// start transaction
	tx, err := c.DB.Begin()
	if err != nil {
		return nil, err
	}

	// fetch division data
	divData := &DivisionData{
		Div: div,
	}
	_, err = tx.Model(divData).WherePK().SelectOrInsert()
	if err != nil {
		return nil, err
	}

	// update command if key already exists in division
	var cmdID uuid.UUID
	cmdIndex := 0
	exists := false
	for i, cc := range divData.Cmds {
		if cc.Key == key {
			exists = true
			cmdID = cc.ID
			cmdIndex = i
			break
		}
	}

	var cmd *CustomCommand
	if exists {
		cmd = &CustomCommand{
			ID:  cmdID,
			Key: key,
			Val: value,
		}

		// update preexisting command in database
		_, err := tx.Model(cmd).WherePK().Update()
		if err != nil {
			return nil, err
		}

		// This shouldn't be necessary, but it only
		// works if I do it, therefore it is necessary.
		divData.Cmds[cmdIndex] = cmd
	} else {
		// add command to database
		cmd = &CustomCommand{
			ID:  uuid.New(),
			Key: key,
			Val: value,
		}
		_, err = tx.Model(cmd).Insert()
		if err != nil {
			return nil, err
		}

		// connect command to division
		divData.Cmds = append(divData.Cmds, cmd)

	}

	// update divsion data
	_, err = tx.Model(divData).WherePK().Update()
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

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
	_, err = tx.Model(divData).WherePK().SelectOrInsert()
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return divData, nil
}

func (c *Connection) RemoveCustomCommand(key string, div Division) error {
	// start transaction
	tx, err := c.DB.Begin()
	if err != nil {
		return err
	}

	// fetch division data
	divData := &DivisionData{
		Div: div,
	}
	_, err = tx.Model(divData).WherePK().SelectOrInsert()
	if err != nil {
		return err
	}

	// update command if key already exists in division
	exists := false
	var id uuid.UUID
	for i, cc := range divData.Cmds {
		if cc.Key == key {
			exists = true
			id = cc.ID
			divData.Cmds = removeSliceItem(divData.Cmds, i)
			break
		}
	}

	var cmd *CustomCommand
	if !exists {
		return errors.New("does not exist")
	}

	// add command to database
	cmd = &CustomCommand{
		ID: id,
	}
	_, err = tx.Model(cmd).WherePK().Delete()
	if err != nil {
		return err
	}

	// command was deleted from division data in for loop

	// update divsion data
	_, err = tx.Model(divData).WherePK().Update()
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// https://stackoverflow.com/questions/37334119/how-to-delete-an-element-from-a-slice-in-golang
func removeSliceItem(s []*CustomCommand, i int) []*CustomCommand {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func (c *Connection) AddMemberRankProgress(user snowflake.Snowflake, div Division, progress int64) error {
	return nil // TODO: Implement member ranks
}

func (c *Connection) GetMemberRankProgress(user snowflake.Snowflake, div Division) (int64, error) {
	return 0, nil // TODO: Implement member ranks
}
