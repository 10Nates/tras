package main

import (
	"db"
	"discordless"
	"fmt"
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
		Port:     55000,
		Password: "Sulfur1-Capacity",
		DBName:   "tras",
	}

	err := conn.Connect()
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	conn.CloseOnInterrupt()

	_, err = conn.AddCustomCommand("test", "this is a test", db.NewDivision('H', 0))
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	_, err = conn.AddCustomCommand("test2", "this is a test 2", db.NewDivision('H', 0))
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	div, err := conn.GetDivsion(db.NewDivision('H', 0))
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	fmt.Printf("div: %v\n", div)
}
