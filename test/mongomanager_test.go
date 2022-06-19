package storage

import (
	"encoding/json"
	"go-league-crawler/pkg/storage"
	types "go-league-crawler/pkg/types/lol"
	"io/ioutil"
	"runtime"
	"testing"
	"time"

	"golang.org/x/time/rate"
)

var (
	// Reads as: At most 15,000 Requests every 10 Minutes
	// However, golang's rate package needs to have the rate given based on (Milli)Seconds.
	// 	e.g.: 10 minutes/ 15000 rq  = 600 seconds / 15000 rq  = 0.040 seconds/rq
	rt                    = rate.Every(40 * time.Millisecond)
	limiter *rate.Limiter = rate.NewLimiter(rt, 1)
	// Number of concurrent threads the Crawler will utilize
	numCpus int = runtime.NumCPU()

	// DB credentials
	localhost        string = "127.0.0.1"
	dbname                  = "go-league-crawler-test"
	matchCollection         = "matches"
	playerCollection        = "player"
	// DB
	mm = storage.NewMManager(localhost, dbname, matchCollection, playerCollection)
)

func TestInsertMatch(t *testing.T) {
	// Reading and decoding test file (match)
	jsonMatch, err := ioutil.ReadFile("./data/match/EUW1_5413144108.json")
	if err != nil {
		t.Fatalf("Reading test file (match) failed!")
	}
	match := types.Match{}
	err = json.Unmarshal([]byte(jsonMatch), &match)
	if err != nil {
		t.Fatalf("Error at Decoding test file (match)!")
	}

	//Conducting test
	err = mm.InsertMatch(match)
	if err != nil {
		t.Fatalf("Error at inserting test match!")
	}
}

func TestInsertPlayer(t *testing.T) {
	// Reading and decoding test file (match)
	jsonMatch, err := ioutil.ReadFile("./data/player/dwaynehart.json")
	if err != nil {
		t.Fatalf("Reading test file (match) failed!")
	}
	sum := types.Summoner{}
	err = json.Unmarshal([]byte(jsonMatch), &sum)
	if err != nil {
		t.Fatalf("Error at Decoding test file (match)!")
	}

	//Conducting test
	err = mm.InsertPlayer(sum)
	if err != nil {
		t.Fatalf("Error at inserting test match!")
	}
}
