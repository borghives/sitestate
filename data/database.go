package data

import (
	"context"
	"log"
	"net"
	"net/url"
	"os"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"golang.org/x/net/proxy"
)

type M bson.M
type A bson.A
type E bson.E
type D bson.D

var (
	client         *mongo.Client
	isDisconnected bool
	once           sync.Once
	err            error
)

func GetDbConnectionUriFromEnv() string {
	mongoDBUriFmt := os.Getenv("MONGODB_URI")
	if mongoDBUriFmt == "" {
		log.Fatal("MONGODB_URI environment variable must be set")
	}

	return os.ExpandEnv(mongoDBUriFmt)

}

type mongoDialerWrapper struct {
	dialer proxy.Dialer
}

func (m *mongoDialerWrapper) DialContext(ctx context.Context, network, addr string) (net.Conn, error) {
	if cd, ok := m.dialer.(interface {
		DialContext(context.Context, string, string) (net.Conn, error)
	}); ok {
		return cd.DialContext(ctx, network, addr)
	}
	return m.dialer.Dial(network, addr)
}

func connectDbClient(connectionString string, proxyAddress string) (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(connectionString)
	if proxyAddress != "" {
		log.Println("Using proxy: ", proxyAddress)
		proxyUrl, err := url.Parse(proxyAddress)
		if err != nil {
			return nil, err
		}

		dialer, err := proxy.FromURL(proxyUrl, proxy.Direct)
		if err != nil {
			return nil, err
		}

		clientOptions = clientOptions.SetDialer(&mongoDialerWrapper{dialer: dialer})
	}
	for i := range 5 {
		client, err = mongo.Connect(clientOptions)
		if err == nil {
			err = client.Ping(context.Background(), nil)
		}
		if err == nil {
			break // Success!
		}
		log.Printf("MongoDb Ping Failed.  Waiting for MongoDB... (attempt %d): %v", i+1, err)
		time.Sleep(5 * time.Second)
	}
	return client, err
}

func InitDbConnection() {
	once.Do(func() {
		connectionString := GetDbConnectionUriFromEnv()
		proxyAddress := os.Getenv("ALL_PROXY")
		client, err = connectDbClient(connectionString, proxyAddress)
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
