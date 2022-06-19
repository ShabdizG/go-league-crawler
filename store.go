package main

import (
	"container/list"
	"errors"
	"fmt"
	"sync"
)

type CacheType struct {
	mux   sync.RWMutex
	Cache map[string]struct{}
}

func NewCache() *CacheType {
	return &CacheType{
		Cache: make(map[string]struct{}),
	}
}

func (c CacheType) String() string {
	c.mux.Lock()
	defer c.mux.Unlock()
	keys := []string{}
	for k, _ := range c.Cache {
		keys = append(keys, k)
	}
	return fmt.Sprintf("%v", keys)
}

// Store depicts a Structure to cache IDs that have already been traversed through
type Store struct {
	Match       *CacheType
	PlayerKnown *CacheType
	mux         sync.Mutex
	PlayerQueue *list.List
}

// NewStore returns an reference to a new sink object
func NewStore() *Store {
	return &Store{
		Match:       NewCache(),
		PlayerKnown: NewCache(),
		PlayerQueue: list.New(),
	}
}

// ConfirmMatch inserts GameId into Sink
func (s *Store) ConfirmMatch(id string) {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.Match.Cache[id] = struct{}{}
}

// ConfirmPlayer inserts AccoundId into SInk
func (s *Store) ConfirmPlayer(id string) {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.PlayerKnown.Cache[id] = struct{}{}
}

// MatchExists checks if a match had been inserted to the sink
func (s *Store) MatchExists(id string) bool {
	s.mux.Lock()
	defer s.mux.Unlock()
	_, ok := s.Match.Cache[id]
	return ok
}

// IsPlayerKnown checks if a match had been inserted to the sink
func (s *Store) IsPlayerKnown(id string) bool {
	s.mux.Lock()
	defer s.mux.Unlock()
	_, ok := s.PlayerKnown.Cache[id]
	return ok
}

// NumMatches returns the number of Matches crawled
func (s *Store) NumMatches() int {
	s.mux.Lock()
	defer s.mux.Unlock()
	return len(s.Match.Cache)
}

// NumPlayers returns the players of Matches crawled
func (s *Store) NumPlayers() int {
	s.mux.Lock()
	defer s.mux.Unlock()
	return len(s.PlayerKnown.Cache)
}

// NextPlayer returns the next player in the queue to process
func (s *Store) NextPlayer() (string, error) {
	s.mux.Lock()
	defer s.mux.Unlock()
	if s.PlayerQueue.Len() == 0 {
		return "", errors.New("No player to process")
	}
	next := s.PlayerQueue.Front()
	s.PlayerQueue.Remove(next)
	return next.Value.(string), nil
}

// AddToQueue inserts a player at the back of the queue of players to process
func (s *Store) AddToQueue(player string) {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.PlayerQueue.PushBack(player)
}
