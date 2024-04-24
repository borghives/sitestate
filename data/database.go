package data

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type M primitive.M
type A primitive.A
type E primitive.E
type D primitive.D

var (
	client         *mongo.Client
	isDisconnected bool
	once           sync.Once
	err            error
)

func GetDbConnectionUriFromEnv() string {
	mongoDBPwd := os.Getenv("SECRET_MONGODB_PWD")

	mongoDBUriFmt := os.Getenv("MONGODB_URI")
	if mongoDBUriFmt == "" {
		log.Fatal("MONGODB_URI environment variable must be set")
	}

	return fmt.Sprintf(mongoDBUriFmt, mongoDBPwd)

}

// InitDbClient initializes the MongoDB client. It will only create one instance.
func InitDbClient(connectionString string) {
	once.Do(func() {
		clientOptions := options.Client().ApplyURI(connectionString)
		client, err = mongo.Connect(context.Background(), clientOptions)
		if err != nil {
			log.Fatalf("Failed to connect to MongoDB: %v", err)
		}

		// Optional: Ping the database to verify connection
		if err = client.Ping(context.Background(), nil); err != nil {
			log.Fatalf("Failed to ping MongoDB: %v", err)
		}

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
