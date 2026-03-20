package data

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"sync"

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
		client, err = mongo.Connect(options.Client().ApplyURI(connectionString))
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

// InitDbProxyClient initializes the MongoDB client with a SOCKS5 proxy.
func InitDbProxyClient(connectionString string) {
	once.Do(func() {
		// 1. Create a dialer for the SOCKS5 proxy (your Tailscale sidecar)
		// In Cloud Run sidecars, this is always localhost:1055
		dialer, err := proxy.SOCKS5("tcp", "localhost:1055", nil, proxy.Direct)
		if err != nil {
			log.Fatalf("Failed to connect to SOCKS5 proxy: %v", err)
		}

		clientOptions := options.Client().ApplyURI(connectionString).
			SetDialer(&mongoDialerWrapper{dialer: dialer})
		client, err = mongo.Connect(clientOptions)
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
