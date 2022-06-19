package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go-league-crawler/pkg/storage"
	types "go-league-crawler/pkg/types/lol"
	"net/http"
	"os"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
)

// Go does not support constant maps
var (
	platformMap = map[string]string{
		"EUW": "EUW1",
		"KR":  "KR",
		"NA":  "NA1",
		"BR":  "BR1",
		"JP":  "JP1",
		"OC":  "OC1",
		"EUN": "EUN1",
	}
	regionMap = map[string]string{
		"EUW1": "europe",
	}
)

const (
	// SUMMONERS refers to the api resource for fetching information about players
	SUMMONERS = "summoner/v4/summoners/"

	// MATCHLIST_BY_PUUID refers to the api resource for fetching a player's match history
	MATCHLIST_BY_PUUID = "match/v5/matches/by-puuid/"

	// MATCH refers to the api resource for fetching a match dto
	MATCH = "match/v5/matches/"

	// RANKED refers to the queueId referencing Summoner's Rift - Ranked Games
	RANKED = "420"
)

// Void is a shortcut type for struct{} that is specifically used for maps
type Void struct{}

// RequestError encapsulates the Status Code from the HTTP Request, in case in was not 200
type RequestError struct {
	err        error
	StatusCode int
}

func (e *RequestError) Error() string {
	return e.err.Error()
}

// Crawler type that in essence consists of region, api key, an http client and a rate limiter
type Crawler struct {
	// Mandatory Parameters and internal Variables
	dbm         storage.DBManager
	platform    string
	region      string
	apiKey      string
	root        string
	maxAttempts int
	client      *http.Client
	ratelimit   *rate.Limiter
	startPlayer string
	store       *Store
	concurrency int
	// Channels to control flow of excecution between goroutines
	playerChan   chan string
	participants chan []string
	ready        chan Void
	quit         chan Void
	// Optional Parameters
	MinNumberOfMatches            int
	MinNumberOfPlayers            int
	TotalNumberOfMatchesPerPlayer int
	mux                           sync.Mutex
}

// NewCrawler initializes an instance of a Crawler for a specific region
func NewCrawler(dbm storage.DBManager, platf string, startPlayer string, rl *rate.Limiter, concurrency int, options ...CrawlerOption) (*Crawler, error) {
	httpClient := &http.Client{}
	nCrwl := Crawler{
		// Mandatory Parameters and internal Variables
		dbm:         dbm,
		platform:    GetPlatform(platf),
		region:      GetRegion(GetPlatform(platf)),
		apiKey:      os.Getenv("DEV_KEY"),
		root:        ".api.riotgames.com/lol/",
		maxAttempts: 10,
		client:      httpClient,
		ratelimit:   rl,
		startPlayer: startPlayer,
		store:       NewStore(),
		concurrency: concurrency,
		// Optional Parameters
		MinNumberOfMatches: DefaultTotalNumberOfMatches,
		MinNumberOfPlayers: DefaultTotalNumberOfPlayers,
		// TODO:
		// TotalNumberOfMatchesPerPlayer:
	}

	for _, opt := range options {
		if err := opt(&nCrwl); err != nil {
			return &nCrwl, err
		}
	}

	return &nCrwl, nil
}

// Start sets the crawler in motion
func (c *Crawler) Start() {
	playerChan := make(chan string)
	participants := make(chan []string, c.concurrency)
	ready := make(chan Void)
	defer close(ready)
	defer close(playerChan)
	defer close(participants)
	var wgWorker sync.WaitGroup
	var wgDispatcher sync.WaitGroup

	ctxWorker, cancelWorker := context.WithCancel(context.Background())
	defer cancelWorker()

	ctxDispatcher, cancelDispatcher := context.WithCancel(context.Background())
	defer cancelDispatcher()

	sp, err := c.GetPlayerByName(c.startPlayer)
	if err != nil {
		log.Infof("Erronous Start Player given!")
		return
	}
	// c.store.AddToQueue(sp.Puuid)
	// Store participants and Queue next Players to crawl Matches from.
	// Keep track of the number of matches crawled and terminate everything
	// when the termination criterion has been met

	go c.CheckTermination(cancelWorker, playerChan)

	go c.QueuePlayers(ctxDispatcher, ctxWorker, participants, playerChan, &wgDispatcher)

	// spawn worker and do work
	for i := 0; i < c.concurrency-1; i++ {
		go c.CrawlPlayer(ctxWorker, i+1, playerChan, ready, participants, &wgWorker)
		wgWorker.Add(1)
	}

	playerChan <- sp.Puuid

	wgWorker.Wait()
	cancelDispatcher()
	wgDispatcher.Wait()
	log.Infof("Finished")
	log.Infof("These are the matches that I have crawled (%d in total)", c.store.NumMatches())
	log.Infof("%v", c.store.Match)
	log.Infof(fmt.Sprintf("These are the players that I have crawled matches from (%d in total)", c.store.NumPlayers()))
	log.Infof("%v", c.store.PlayerKnown)
}

func (c *Crawler) Finished() bool {
	if c.store.NumMatches() >= c.MinNumberOfMatches || c.store.NumPlayers() >= c.MinNumberOfPlayers {
		return true
	}
	return false
}

func (c *Crawler) CheckTermination(cancelWorker func(), player chan string) {
	for {
		if c.Finished() {
			cancelWorker()
			return
		}
	}
}

func (c *Crawler) QueuePlayers(ctxDispatcher context.Context, ctxWorker context.Context, participants <-chan []string, player chan<- string, wgDispatcher *sync.WaitGroup) {
	wgDispatcher.Add(1)
OUTER:
	for {
		select {
		case <-ctxDispatcher.Done():
			log.Info("Dispatcher's Job has Finished")
			wgDispatcher.Done()
			return
		case p := <-participants:
		INNER_PARTICIPANTS:
			for _, player := range p {
				if c.store.IsPlayerKnown(player) {
					log.Infof("Player %s already known", player)
					continue INNER_PARTICIPANTS
				}
				if player == "" {
					log.Warnf("Empty Player about to be inserted")
				}
				c.store.AddToQueue(player)
				log.Infof("New Player %s pushed into Queue", player)
			}
		INNER_PLAYER:
			for i := 0; i < (c.concurrency - len(c.playerChan)); i++ {
				next, err := c.store.NextPlayer()
				if err != nil {
					log.Error(err)
					continue INNER_PLAYER
				}
				if next == "" {
					log.Warnf("Empty Player ID but programm has not continued...")
				}
				select {
				case <-ctxWorker.Done():
					continue OUTER
				default:
					log.Infof("Next player in line: %s", next)
					player <- next
				}
			}

		}
	}
}

// CrawlPlayer fetches a matchlist based on the given summoner and places the gameIds into the crawler's match channel
func (c *Crawler) CrawlPlayer(ctx context.Context, workerID int, playerChan <-chan string, ready chan<- Void, participants chan<- []string, wg *sync.WaitGroup) {
	log.Printf("Goroutine with WorkerID %v started", workerID)
OUTER:
	for {
		select {
		case <-ctx.Done():
			wg.Done()
			log.Infof("[WorkerID:%v]: Goroutine finished", workerID)
			return
		case <-time.After(1 * time.Minute):
			log.Infof("[WorkerID:%v]: 1 Minute passed without any new player coming in ...", workerID)
		case player := <-playerChan:
			log.Infof("[WorkerID:%v]: Received player %v", workerID, player)
			if player == "" {
				log.Warnf("Empty Player ID")
			}
			log.Infof("[WorkerID:%v]: Getting Matchlist of player %v ...", workerID, player)
			ml, err := c.GetMatchList(player)
			if err != nil {
				log.Infof("[WorkerID:%v]: Error occured for Player", workerID, player)
				continue OUTER
			}
			identifiedParticipants := []string{}
			// Process each match from matchlist
		INNER:
			for _, m := range *ml {
				select {
				case <-ctx.Done():
					continue OUTER
				default:
					if c.store.MatchExists(m) {
						log.Infof("[WorkerID:%v]: Match already crawled %v", workerID, m)
						continue INNER
					}
					match, err := c.GetMatch(m)
					if err != nil {
						log.Errorf("[WorkerID:%v] Error fetching match %v", workerID, m)
						continue INNER
					}
					log.Infof("[WorkerID:%v]: New Match: %v", workerID, match.Info.GameID)
					// Handle Match
					c.dbm.InsertMatch(*match)
					c.store.ConfirmMatch(match.MetaData.MatchID)
					log.Infof("[WorkerID:%v][Region: %v][Player: %v]: Total Number of Matches crawled so far: %v", workerID, c.platform, player, c.store.NumMatches())
					identifiedParticipants = append(identifiedParticipants, match.MetaData.Participants...)
				}
			}
			// Finish up current player and get next player to process
			participants <- identifiedParticipants
			summoner, err := c.GetPlayerByPUUID(player)
			c.dbm.InsertPlayer(*summoner)
			log.Infof("[WorkerID:%v] Finished working on player %s", workerID, player)
			c.store.ConfirmPlayer(player)
		}
	}
}

// _sendRequest represents a proxy function that includes the Rate Limit Management
// before requesting the ressource. It will be invoked by the function SendRequest
func (c *Crawler) _sendRequest(req *http.Request) (*http.Response, error) {
	if !c.ratelimit.Allow() {
		r := c.ratelimit.Reserve()
		log.Warnf("Rate Limit has been reached, need to wait for %v", r.Delay())
		time.Sleep(r.Delay())
	}
	return c.client.Do(req)
}

// SendRequest attempts (at most maxAttempts times) to get the http contents from the given url
// automatically re-attempts to request the url based on the settings of the Crawler in case of internal server issues or rate limit exceedances
// in case of neither: returns an error with the status code
func (c *Crawler) SendRequest(url string, maxAttempts ...int) (*http.Response, error) {
	maxAtt := c.maxAttempts
	if maxAttempts[0] != 0 {
		maxAtt = maxAttempts[0]
	}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("X-Riot-Token", c.apiKey)
	for i := 1; i <= maxAtt; i++ {
		response, e := c._sendRequest(req)
		if e != nil {
			return nil, e
		}
		// log.Infof("Status Code: %v", response.StatusCode)
		switch response.StatusCode {
		case http.StatusOK: // 200
			return response, nil
		case http.StatusForbidden:
			log.Warnf("Non-Temporary Error: %v for url %v", http.StatusText(response.StatusCode), url)
		case http.StatusTooManyRequests: // 429
			log.Warnf("Temporary Error: %v for url %v", http.StatusText(response.StatusCode), url)
			// see: https://stackoverflow.com/questions/17573190/how-to-multiply-duration-by-integer
			time.Sleep(time.Second * time.Duration(10*i))
		case http.StatusInternalServerError, http.StatusServiceUnavailable, http.StatusGatewayTimeout: // 500, 503, 504:
			log.Warnf("Temporary Error: %v occured at attempt No. %d for url %v", http.StatusText(response.StatusCode), i, url)
			// see: https://stackoverflow.com/questions/17573190/how-to-multiply-duration-by-integer
			time.Sleep(time.Second * time.Duration(10*i))
			continue
		default:
			log.Fatalf("Non-temporary Error occured for url %v", url)
			err := &RequestError{
				fmt.Errorf("Non-temporary Error: %s", http.StatusText(response.StatusCode)),
				response.StatusCode,
			}
			return nil, err
		}
	}
	log.Fatalf("Maximum Attempts (%d) reached. Could not retrieve Document. Continuing...", maxAtt)
	err := &RequestError{
		errors.New("Maximum Attempts Exceeded"),
		0,
	}
	return nil, err
}

// GetPlayerByName retrieves a SummonerDTO based on a given name
func (c *Crawler) GetPlayerByName(name string) (*types.Summoner, error) {
	// Example https://euw1.api.riotgames.com/lol/summoner/v4/summoners/by-name/dwaynehart
	SummonerDTO := types.Summoner{}
	url := "https://" + c.platform + c.root + SUMMONERS + "by-name/" + name
	log.Infof("Requesting URL: %v", url)
	response, err := c.SendRequest(url, c.maxAttempts)
	if err != nil {
		return nil, err
	}
	err = json.NewDecoder(response.Body).Decode(&SummonerDTO)
	if err != nil {
		log.Fatal("Error at Decoding")
		print(err)
	}
	return &SummonerDTO, nil
}

// GetPlayerByPUUID retrieves a SummonerDTO based on a given puuid
func (c *Crawler) GetPlayerByPUUID(puuid string) (*types.Summoner, error) {
	// Example https://euw1.api.riotgames.com/lol/summoner/v4/summoners/by-name/dwaynehart
	SummonerDTO := types.Summoner{}
	url := "https://" + c.platform + c.root + SUMMONERS + "by-puuid/" + puuid
	log.Infof("Requesting URL:", url)
	response, err := c.SendRequest(url, c.maxAttempts)
	if err != nil {
		return nil, err
	}
	err = json.NewDecoder(response.Body).Decode(&SummonerDTO)
	if err != nil {
		log.Fatal("Error at Decoding")
		print(err)
	}
	return &SummonerDTO, nil
}

// GetMatchList reveices the entire matchlist according to the queue type of a speciifc player
func (c *Crawler) GetMatchList(puuid string) (*[]string, error) {
	var (
		start     int = 0
		count     int = 100
		err       error
		matches   *[]string
		matchList []string = []string{}
	)
	for {
		matches, err = c.GetPaginatedMatchList(puuid, start, count)
		if err != nil {
			log.Error(err)
		}
		if len(*matches) == 0 {
			break
		}
		matchList = append(matchList, *matches...)
		start += count
	}
	return &matchList, nil
}

// GetPaginatedMatchList retrieves a paginated MatchListDTO based on an accountId and a queueID
func (c *Crawler) GetPaginatedMatchList(puuid string, start int, count int) (*[]string, error) { //(*types.MatchList, error) {
	matchList := []string{}
	url := fmt.Sprintf("https://%s%s%s%s", c.region, c.root, MATCHLIST_BY_PUUID, puuid)
	url += fmt.Sprintf("/ids?queue=%v&start=%v&count=%v", RANKED, start, count)
	log.Infof("Request URL: %v", url)
	response, err := c.SendRequest(url, c.maxAttempts)
	if err != nil {
		return nil, err
	}
	err = json.NewDecoder(response.Body).Decode(&matchList)
	if err != nil {
		log.Infof("Error at Decoding")
	}
	return &matchList, nil
}

// GetMatch retrieves a Match based on a gameId
func (c *Crawler) GetMatch(gameID string) (*types.Match, error) {
	MatchDTO := types.Match{}
	//url := "https://" + c.Region + c.root + MATCH + gameID
	url := fmt.Sprintf("https://%s%s%s%s", c.region, c.root, MATCH, gameID)
	log.Infof("Request URL: %v", url)
	response, err := c.SendRequest(url, c.maxAttempts)
	if err != nil {
		return nil, err
	}
	err = json.NewDecoder(response.Body).Decode(&MatchDTO)
	if err != nil {
		log.Infof("Error at Decoding")
	}
	return &MatchDTO, nil
}
