package storage

import (
	"context"
	"fmt"
	types "go-league-crawler/pkg/types/lol"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	CONTEXT_TIMOUT    = 60 * time.Second
	MATCH_COLLECTION  = "matches"
	PLAYER_COLLECTION = "players"
)

type MongoManager struct {
	*DB
	Client *mongo.Client
}

// NewMManager returns a new MManager object in order to create a new mongo connection
func NewMManager(host string, database string, matchCollection string, playerCollection string) *MongoManager {
	db := &DB{
		Host:          host,
		Database:      database,
		MatchStorage:  matchCollection,
		PlayerStorage: playerCollection,
	}
	mm := &MongoManager{
		DB: db,
	}
	client, err := mm.connect()
	if err != nil {
		fmt.Print(err)
	}
	mm.Client = client
	return mm
}

func (mm *MongoManager) Init() error {
	client, err := mm.connect()
	mm.Client = client
	err = mm.ping()
	return err
}

func (mm *MongoManager) connect() (*mongo.Client, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(fmt.Sprintf("mongodb://%v:27017", mm.Host)))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), CONTEXT_TIMOUT)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	return client, err
}

// Ping Tests the Client Connection
func (mm *MongoManager) ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), CONTEXT_TIMOUT)
	defer cancel()
	return mm.Client.Ping(ctx, readpref.Primary())
}

func (mm *MongoManager) InsertMatch(match types.Match) error {
	res, err := mm.Client.Database(mm.Database).Collection(mm.MatchStorage).InsertOne(context.TODO(), match)
	fmt.Println(fmt.Sprintf("Stored Match with ID: %v\n", res.InsertedID))
	return err
}

func (mm *MongoManager) InsertPlayer(player types.Summoner) error {
	res, err := mm.Client.Database(mm.Database).Collection(mm.PlayerStorage).InsertOne(context.TODO(), player)
	fmt.Printf("Stored Player with ID: %v\n", res.InsertedID)
	return err
}
