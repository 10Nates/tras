package main

import (
	"db"
	"discordless"
	"testing"
)

// passthrough to discordless testing file

func BenchmarkParser(b *testing.B) {
	for i := 0; i < b.N; i++ {
		discordless.Test(parseCommand)
	}
}

func TestParser(t *testing.T) {
	res := discordless.Test(parseCommand)

	nF := 0
	for _, v := range res {
		if len(v) > 17 && v[0:18] == "An error occurred." {
			t.Log(v)
			t.Fail()
			nF++
		}
	}
	if t.Failed() {
		t.Logf("Number failed: %d/%d", nF, len(res))
	}
}

func TestDatabase(t *testing.T) {
	conn := &db.Connection{
		Host:     "localhost",
		Port:     5432,
		Password: "Sulfur1-Capacity",
		DBName:   "tras",
	}

	conn.Connect()
	conn.CloseOnInterrupt()

	cmd, err := conn.AddCustomCommand("test", "this is a test", db.NewDivision('H', 0))
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	t.Log(cmd.ID, cmd.Key, cmd.Val)
}
