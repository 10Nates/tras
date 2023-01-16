package db

import (
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/google/uuid"
)

// custom command data, also used in customcmds.go
type CustomCommand struct {
	ID  uuid.UUID `pg:",pk,type:uuid,default:uuid_generate_v4()"`
	Key string
	Val string
}

type RankMember struct {
	ID         uuid.UUID `pg:",pk,type:uuid,default:uuid_generate_v4()"`
	UserID     uint64
	Progress   int64
	LastMsgTs  time.Time // used for base attentiveness score
	LastChanID uint64    // score booster
}

type DivisionData struct {
	Div           Division `pg:",pk"`
	RandSpeak     bool
	LastRandSpeak time.Time // used to manage randSpeak interval
	Cmds          []*CustomCommand
	RankMems      []*RankMember
}

func createSchema(db *pg.DB) error {
	// I had to do this to get uuid to work
	_, err := db.ExecOne(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`)
	if err != nil {
		return err
	}

	models := []interface{}{
		(*CustomCommand)(nil),
		(*RankMember)(nil),
		(*DivisionData)(nil),
	}

	for _, model := range models {
		err := db.Model(model).CreateTable(&orm.CreateTableOptions{
			IfNotExists: true,
		})
		if err != nil {
			return err
		}
	}
	return nil
}
