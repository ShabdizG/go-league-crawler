package main

import (
	"fmt"
	"math"
)

// Default Values
const (
	// Allowed Values
	MaxAllowedTotalNumberOfMatches = 10000000
	MaxAllowedTotalNumberOfPlayers = 10000000
	// Default Values
	DefaultTotalNumberOfMatches          = 10000
	DefaultTotalNumberOfPlayers          = 10
	DefaultTotalNumberOfMatchesPerPlayer = math.MaxInt64
)

type CrawlerOption func(*Crawler) error

func WithMinNumberOfMatches(numberOfMatches int) func(*Crawler) error {
	return func(c *Crawler) error {
		if numberOfMatches > MaxAllowedTotalNumberOfMatches {
			return fmt.Errorf("Number of Matches exceeds allowed limit (%v)\n", MaxAllowedTotalNumberOfMatches)
		}
		if numberOfMatches > 0 {
			c.MinNumberOfMatches = numberOfMatches
		}
		return nil
	}
}

func WithMinNumberOfPlayers(numberOfPlayers int) func(*Crawler) error {
	return func(c *Crawler) error {
		if numberOfPlayers > MaxAllowedTotalNumberOfPlayers {
			return fmt.Errorf("Number of Players exceeds allowed limit (%v)\n", MaxAllowedTotalNumberOfPlayers)
		}
		if numberOfPlayers > 0 {
			c.MinNumberOfPlayers = numberOfPlayers
		}
		return nil
	}
}
