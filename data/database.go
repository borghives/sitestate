package data

import (
	"context"
	"log"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

var SEA_DATABASE = "pieriansea"
var SEA_DATABASE_REPORT = "pierianseareport"
var SEA_DATABASE_AUTH = "pierianauth"
var (
	client                *mongo.Client
	client_unacknowledged *mongo.Client
	isDisconnected        bool
	once                  sync.Once
	err                   error
)

// InitMongoClient initializes the MongoDB client. It will only create one instance.
func InitMongoClient(connectionString string) {
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

		// Unacknowledged client
		wc := writeconcern.Unacknowledged()
		clientOptions.SetWriteConcern(wc)
		client_unacknowledged, err = mongo.Connect(context.Background(), clientOptions)
		if err != nil {
			log.Fatalf("Failed to connect to client_unacknowledged MongoDB: %v", err)
		}

		// Optional: Ping the database to verify connection
		if err = client_unacknowledged.Ping(context.Background(), nil); err != nil {
			log.Fatalf("Failed to ping client_unacknowledged MongoDB: %v", err)
		}
	})
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

// getMongoQuickClient returns the unacknowledged instance of the MongoDB client.
func getMongoQuickClient() *mongo.Client {
	if client_unacknowledged == nil {
		if isDisconnected {
			log.Println("MongoDB client was previously disconnected")
			return nil
		}
		log.Fatal("MongoDB client is not initialized")
	}
	return client_unacknowledged
}

// DisconnectMongoClient handles the disconnection of the MongoDB client.
func DisconnectMongoClient() {
	if client != nil {
		if err := client.Disconnect(context.Background()); err != nil {
			log.Fatalf("Failed to disconnect MongoDB client: %v", err)
		}
		client = nil
		isDisconnected = true
		log.Println("Disconnected from MongoDB!")
	}
}

func GetDB() *mongo.Database {
	return getMongoClient().Database(SEA_DATABASE)
}

func GetDB_Unacknowledged() *mongo.Database {
	return getMongoQuickClient().Database(SEA_DATABASE)
}

func GetReportDB() *mongo.Database {
	return getMongoClient().Database(SEA_DATABASE_REPORT)
}

func GetAuthDB() *mongo.Database {
	return getMongoClient().Database(SEA_DATABASE_REPORT)
}
