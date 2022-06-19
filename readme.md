# Go Leage Crawler

This project entails a Web Crawler for League of Legends Matches and Players utilizing the Riot API.
Up to this date, it makes use of the Match-V5 and the Summoner-V4 Endpoints. 
It can essentially be used as a command line application.

## Requirements

The Crawler is based on the following requirements:
- Go 1.14
- MongoDB 4.2.6

The required packages used in the golang source code can be obtained from the corresponding `go.mod` file

### API Key

To run the Crawler it is necessary to include your Riot API Key. 
At this current stage the source code assumes your **API key** to be included as an environmental variable `DEV_KEY` 

## Usage

In the project directory create the binary via 

`go build -o ./bin/go-league-crawler`

After creating the executable file you can run Crawler by including a set of specfics

`go run go-league-crawler -s "ben trades" -m 100`

The Crawler will then crawl at least 100 Matches beginning with the player "ben trades". In case the given player has less matches played, the next player's matchlist will be crawled.


## Parameters

	-s       Player with whom to begin to crawl data from
    -m       Minimum Number of Matches to Crawl before terminating
	-p       Minimum Players of Matches to Crawl before terminatinghost    
    -host    Host of the Target DB
	-dbname  Name of the Target DB
	-mc      Collection where to ingest the match data into
	-pc      Collection where to ingest the player data into
	-pl      Region to crawl data from
	-con     Degree of Concurrency (No. of Threads)
