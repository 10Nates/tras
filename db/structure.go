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
	LastMsgTs  time.Time
	LastChanID uint64 // guildID
}

type DivisionData struct {
	Div       Division `pg:",pk"`
	RandSpeak bool
	Cmds      []CustomCommand `pg:"rel:has-many"`
	RankMems  []RankMember    `pg:"rel:has-many"`
}

func createSchema(db *pg.DB) error {
	models := []interface{}{
		(*CustomCommand)(nil),
		(*RankMember)(nil),
		(*DivisionData)(nil),
	}

	for _, model := range models {
		err := db.Model(model).CreateTable(&orm.CreateTableOptions{
			Temp: true,
		})
		if err != nil {
			return err
		}
	}
	return nil
}
