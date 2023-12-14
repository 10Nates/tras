package main

import (
	"db"
	"encoding/json"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/andersfylling/disgord"
	"github.com/andersfylling/snowflake/v5"
)

var BotID string // loaded on init
var BotPFP string
var BigTypeLetters map[string]map[string]string // this is way easier than the alternative
var ThesaurusLookup map[string][]string
var GRand = rand.New(rand.NewSource(time.Now().UnixNano()))
var DBConn *db.Connection

var slashCommands = []*disgord.CreateApplicationCommand{
	{
		Name:        "help",
		Description: "Alternative to @TRAS help, view help command without DMs.",
	},
	{
		Name:        "about",
		Description: "Alternative to @TRAS about. view about command without a DM.",
	},
}

type MessagePassthrough struct {
	Message     *disgord.Message
	Interaction *disgord.InteractionCreate
}

var OldData struct {
	Cmds    map[string]map[string]string `json:"cmds"`
	RandSpk map[string]bool              `json:"randSpk"`
}

func main() {
	// load bigtype letters
	bigtypejson, err := os.ReadFile("src/bigtype.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(bigtypejson, &BigTypeLetters)
	if err != nil {
		panic(err)
	}

	// load thesaurus
	thesarusjson, err := os.ReadFile("src/thesaurus.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(thesarusjson, &ThesaurusLookup)
	if err != nil {
		panic(err)
	}

	// connect to DB
	DBConn = &db.Connection{
		Host:     DB_HOST,
		Port:     DB_PORT,
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   DB_NAME,
	}
	err = DBConn.Connect()
	if err != nil {
		panic(err)
	}
	DBConn.CloseOnInterrupt()

	olddatajson, err := os.ReadFile("src/data.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(olddatajson, &OldData)
	if err != nil {
		panic(err)
	}

	for guild, guildcmds := range OldData.Cmds {
		guildID, err := strconv.ParseUint(guild, 10, 64)
		if err != nil {
			panic(err)
		}
		for key, resp := range guildcmds {
			_, err := DBConn.SetCustomCommand(key, resp, db.NewDivision('G', snowflake.NewSnowflake(guildID)))
			if err != nil {
				panic(err)
			}
		}
	}

	for guild, val := range OldData.RandSpk {
		guildID, err := strconv.ParseUint(guild, 10, 64)
		if err != nil {
			panic(err)
		}
		DBConn.SetRandomSpeakAvailability(db.NewDivision('G', snowflake.NewSnowflake(guildID)), val)
	}
}
