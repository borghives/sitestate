package data

import (
	"context"
	"log"
	"os"
	"sync"

	"github.com/borghives/kosmos-go"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type M bson.M
type A bson.A
type E bson.E
type D bson.D

var (
	client         *mongo.Client
	isDisconnected bool
	once           sync.Once
)

func GetDbConnectionUriFromEnv() string {
	mongoDBUriFmt := os.Getenv("MONGODB_URI")
	if mongoDBUriFmt == "" {
		log.Fatal("MONGODB_URI environment variable must be set")
	}
	log.Println("Using MongoDB URI: ", mongoDBUriFmt)
	return os.ExpandEnv(mongoDBUriFmt)

}

func InitDbConnection() {
	once.Do(func() {

		client = kosmos.MustHaveObserverClient()

		log.Println("Connected to MongoDB!")
		isDisconnected = false
	})
}

// DisconnectMongoClient handles the disconnection of the MongoDB client.
func DisconnectDbClient() {
	if client != nil {
		if err := client.Disconnect(context.Background()); err != nil {
			log.Fatalf("Failed to disconnect MongoDB client: %v", err)
		}
		client = nil
		isDisconnected = true
		log.Println("Disconnected from MongoDB!")
	}
}

// GetDatabase returns the database with the given name.
// If the database does not exist, it will be created.
func GetDatabase(name string) *mongo.Database {
	return getMongoClient().Database(name)
}

// getMongoClient returns the instance of the MongoDB client.
func getMongoClient() *mongo.Client {
	if client == nil {
		if isDisconnected {
			log.Println("MongoDB client was previously disconnected")
			return nil
		}
		log.Fatal("MongoDB client is not initialized")
	}
	return client
}
