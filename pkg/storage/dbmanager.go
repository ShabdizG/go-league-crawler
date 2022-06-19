package storage

import types "go-league-crawler/pkg/types/lol"

type DBManager interface {
	InsertMatch(match types.Match) error
	InsertPlayer(player types.Summoner) error
}

type DB struct {
	Host          string
	Database      string
	MatchStorage  string
	PlayerStorage string
}
