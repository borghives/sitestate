package data

import (
	"context"
	"log"
	"net"
	"os"
	"sync"

	"github.com/borghives/go-cmd-tool/shared"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"golang.org/x/net/proxy"
)

type M bson.M
type A bson.A
type E bson.E
type D bson.D

var (
	config         shared.SiteConfig
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
	log.Println("Using MongoDB URI: ", mongoDBUriFmt)
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

func InitDbConnection() {
	once.Do(func() {
		var err error
		config, err = shared.LoadSiteConfig()
		if err != nil {
			log.Fatalf("Failed to load site config: %v\n", err)
		}

		log.Printf("Using MongoDB URI: %s", config.MongoDBUri)
		log.Printf("Using MongoDB Auth URI: %s", config.MongoDBAuthUri)
		log.Printf("Using Project ID: %s", config.ProjectID)
		log.Printf("Using Proxy Address: %s", config.ProxyAddress)
		log.Printf("Proxy ENV: %s", os.Getenv("ALL_PROXY"))
		client = shared.MustGetDbClient(&config)

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
