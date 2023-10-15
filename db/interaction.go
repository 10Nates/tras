package db

import (
	"errors"

	"github.com/andersfylling/disgord"
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
		tx.Rollback()
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
			tx.Rollback()
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
			tx.Rollback()
			return nil, err
		}

		// connect command to division
		divData.Cmds = append(divData.Cmds, cmd)

	}

	// update divsion data
	_, err = tx.Model(divData).WherePK().UpdateNotZero()
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
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
		tx.Rollback()
		return nil, err
	}

	// need to commit because DivisionData is created
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
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
		tx.Rollback()
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
		tx.Rollback()
		return errors.New("does not exist")
	}

	// add command to database
	cmd = &CustomCommand{
		ID: id,
	}
	_, err = tx.Model(cmd).WherePK().Delete()
	if err != nil {
		tx.Rollback()
		return err
	}

	// command was deleted from division data in for loop

	// update divsion data
	_, err = tx.Model(divData).WherePK().UpdateNotZero()
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

// https://stackoverflow.com/questions/37334119/how-to-delete-an-element-from-a-slice-in-golang
func removeSliceItem(s []*CustomCommand, i int) []*CustomCommand {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func (c *Connection) SetRankMemberProgress(msg *disgord.Message, uID disgord.Snowflake, div Division, progress int64) error {
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
		tx.Rollback()
		return err
	}

	// find member
	var mem *RankMember
	for _, rm := range divData.RankMems {
		if rm.UserID == uint64(uID) {
			mem = rm
			break
		}
	}

	if mem == nil {
		// add member to database
		mem = &RankMember{
			ID:         uuid.New(),
			UserID:     uint64(uID),
			Progress:   progress,
			LastMsgTs:  msg.Timestamp.Time,
			LastChanID: uint64(msg.ChannelID), // Although this remains inaccurate for force set (uID mismatch) it doesn't really matter
		}

		_, err = tx.Model(mem).Insert()
		if err != nil {
			tx.Rollback()
			return err
		}

		divData.RankMems = append(divData.RankMems, mem)

	} else {
		// member does exist, update member
		mem.Progress = progress
		mem.LastMsgTs = msg.Timestamp.Time
		mem.LastChanID = uint64(msg.ChannelID) // refer to previous comment

		// update member in database
		_, err := tx.Model(mem).WherePK().Update()
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	// update divsion data
	_, err = tx.Model(divData).WherePK().UpdateNotZero()
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (c *Connection) GetRankMember(uID disgord.Snowflake, div Division) (*RankMember, error) {
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
		tx.Rollback()
		return nil, err
	}

	// find member
	var mem *RankMember
	for _, rm := range divData.RankMems {
		if rm.UserID == uint64(uID) {
			mem = rm
			break
		}
	}

	// need to commit because DivisionData is created
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if mem == nil {
		return new(RankMember), nil // if member does not exist, it is 0.
	} else {
		return mem, nil
	}
}

func (c *Connection) SetDiceAvailability(div Division, enabled bool) error {
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
		tx.Rollback()
		return err
	}

	// modify
	divData.Dice = enabled

	// update divsion data
	_, err = tx.Model(divData).WherePK().UpdateNotZero()
	if err != nil {
		tx.Rollback()
		return err
	}

	// need to commit because DivisionData is created
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

// Copy of SetDiceAvailability with different boolean to change
func (c *Connection) SetRandomSpeakAvailability(div Division, enabled bool) error {
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
		tx.Rollback()
		return err
	}

	// modify
	divData.RandSpeak = enabled

	// update divsion data
	_, err = tx.Model(divData).WherePK().UpdateNotZero()
	if err != nil {
		tx.Rollback()
		return err
	}

	// need to commit because DivisionData is created
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
