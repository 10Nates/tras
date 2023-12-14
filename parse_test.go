package main

import (
	"db"
	"discordless"
	"fmt"
	"testing"

	"github.com/andersfylling/disgord"
)

// passthrough to discordless testing file

func parseCommandPassthrough(msg *disgord.Message, s *disgord.Session) {
	parseCommand(MessagePassthrough{Message: msg}, s)
}

func BenchmarkParser(b *testing.B) {
	for i := 0; i < b.N; i++ {
		discordless.Test(parseCommandPassthrough)
	}
}

func TestParser(t *testing.T) {
	// Connect
	DBConn = &db.Connection{
		Host:     DB_HOST,
		Port:     DB_PORT,
		Password: "Sulfur1-Capacity", // dev only password
		DBName:   DB_NAME,
	}

	err := DBConn.Connect()
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	DBConn.CloseOnInterrupt()

	res := discordless.Test(parseCommandPassthrough)

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

func TestCCDatabase(t *testing.T) {
	// Connect
	conn := &db.Connection{
		Host:     DB_HOST,
		Port:     DB_PORT,
		Password: "Sulfur1-Capacity", // dev only password
		DBName:   DB_NAME,
	}

	err := conn.Connect()
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	conn.CloseOnInterrupt()

	d := db.NewDivision('H', 0)

	// Add commands
	_, err = conn.SetCustomCommand("test", "this is a test", d)
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	_, err = conn.SetCustomCommand("test2", "this is a test 2", d)
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	// check
	div, err := conn.GetDivsion(d)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	fmt.Printf("div: %v\n", div)
	for _, cc := range div.Cmds {
		fmt.Printf("%v\n", cc)
	}

	// Alter commands
	_, err = conn.SetCustomCommand("test", "this is an alternate test", d)
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	_, err = conn.SetCustomCommand("test2", "this is another alt test", d)
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	// check
	div2, err := conn.GetDivsion(d)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	fmt.Printf("div2: %v\n", div2)
	for _, cc := range div2.Cmds {
		fmt.Printf("%v\n", cc)
	}

	// Remove commands
	err = conn.RemoveCustomCommand("test", d)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	err = conn.RemoveCustomCommand("test2", d)
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	// check
	div3, err := conn.GetDivsion(d)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	fmt.Printf("div3: %v\n", div3)
	for _, cc := range div3.Cmds {
		fmt.Printf("%v\n", cc)
	}
}
