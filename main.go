package main

import (
	"context"
	"flag"
	"fmt"
	"go-league-crawler/pkg/logging"
	"go-league-crawler/pkg/storage"
	"time"

	log "github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
)

// TODOs:
// command line parameters
// Properly implement the "worker paradigm"
// Add more Default/Optional Parameteres for Crawler and DBManager

var (
	// Reads as: At most 15,000 Requests every 10 Minutes
	// However, golang's rate package needs to have the rate given based on (Milli)Seconds.
	// 	e.g.: 10 minutes/ 15000 rq  = 600 seconds / 15000 rq  = 0.040 seconds/rq
	rt                    = rate.Every(40 * time.Millisecond)
	limiter *rate.Limiter = rate.NewLimiter(rt, 1)
	// Number of concurrent threads the Crawler will utilize

	// DB Manager Properties
	host             string = "127.0.0.1"
	dbname           string = "go-league-crawler-test"
	matchCollection  string = "matches"
	playerCollection string = "players"

	// Crawler Properties
	platform           string = "EUW"
	startPlayer        string = "dwaynehart"
	concurrency        int    = 6
	minNumberOfMatches int    = 100
	minNumberOfPlayers int    = 0

	// Command Line Flag Pointers
	hostPtr               *string = flag.String("host", host, "Host of the Target DB")
	dbnamePtr             *string = flag.String("db", dbname, "Name of the Target DB")
	matchCollectionPtr    *string = flag.String("mc", matchCollection, "Collection where to ingest the match data into")
	playerCollectionPtr   *string = flag.String("pc", playerCollection, "Collection where to ingest the player data into")
	platformPtr           *string = flag.String("pl", platform, "Region to crawl data from")
	startPlayerPtr        *string = flag.String("s", startPlayer, "Player with whom to begin to crawl data from")
	concurrencyPtr        *int    = flag.Int("con", concurrency, "Degree of Concurrency (No. of Threads)")
	minNumberOfMatchesPtr *int    = flag.Int("m", minNumberOfMatches, "Minimum Number of Matches to Crawl before terminating")
	minNumberOfPlayersPtr *int    = flag.Int("p", minNumberOfPlayers, "Minimum Players of Matches to Crawl before terminating")
)

func main() {
	// Init Logger
	logging.InitLogger("./log/logfile.log")

	// Parse Command Line Flags and log them
	flag.Parse()
	log.WithFields(log.Fields{
		"Host":                     *hostPtr,
		"DB":                       *dbnamePtr,
		"Matches":                  *matchCollectionPtr,
		"Players":                  *playerCollectionPtr,
		"Platform":                 *platformPtr,
		"Starting Player":          *startPlayerPtr,
		"Concurrency Level":        *concurrencyPtr,
		"Minimum Matches to Crawl": *minNumberOfMatchesPtr,
		"Minimum Players to Crawl": *minNumberOfPlayersPtr,
	}).Info("Started Crawler with the following parameters")
	now := time.Now()

	// Init DB Manager
	mm := storage.NewMManager(host, dbname, matchCollection, playerCollection)
	defer mm.Client.Disconnect(context.Background())
	err := mm.Init()
	if err != nil {
		panic(err)
	}
	// Init Crawler
	EUWCrawler, err := NewCrawler(
		// Mandatory Parameters
		mm, platform, startPlayer, limiter, concurrency,
		// Optional Parameters
		WithMinNumberOfMatches(*minNumberOfMatchesPtr),
		WithMinNumberOfPlayers(*minNumberOfPlayersPtr),
	)
	// Start Crawling Matches
	EUWCrawler.Start()
	then := time.Now()
	fmt.Println("Finished in ", then.Sub(now))
}
