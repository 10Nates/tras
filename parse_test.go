package main

import (
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
