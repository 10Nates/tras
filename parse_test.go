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

func TestCCDatabase(t *testing.T) {
	// COnnect
	conn := &db.Connection{
		Host:     "localhost",
		Port:     55001,
		Password: "Sulfur1-Capacity",
		DBName:   "tras",
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
