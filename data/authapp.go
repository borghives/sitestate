package data

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var SEA_DATABASE_AUTH = "pierianauth"

var DB_USER_COLLECTION_NAME = "user"
var DB_AUTH_SESSION_COLLECTION_NAME = "auth_session"

type Authentication interface {
	Initialize()
	User() *mongo.Collection
	AuthSession() *mongo.Collection
}

type AuthenticationApp struct {
	database_name string
}

func (a *AuthenticationApp) Initialize() {
	usernameIndex := mongo.IndexModel{
		Keys: bson.M{
			"name": 1,
		},
		Options: options.Index().SetUnique(true),
	}

	expireIndex := mongo.IndexModel{
		Keys: bson.M{
			"created_time": 1,
		},
		Options: options.Index().SetExpireAfterSeconds(30 * 24 * 60 * 60),
	}

	_, err := a.User().Indexes().CreateOne(context.Background(), usernameIndex)
	if err != nil {
		log.Printf("error creating indexes for user: %s", err)
	}

	_, err = a.AuthSession().Indexes().CreateOne(context.Background(), expireIndex)
	if err != nil {
		log.Printf("error creating indexes for user: %s", err)
	}
}

func (a AuthenticationApp) GetDatabase() *mongo.Database {
	return GetDatabase(a.database_name)
}

func (a AuthenticationApp) User() *mongo.Collection {
	return a.GetDatabase().Collection(DB_USER_COLLECTION_NAME)
}

func (a AuthenticationApp) AuthSession() *mongo.Collection {
	return a.GetDatabase().Collection(DB_AUTH_SESSION_COLLECTION_NAME)
}

func AuthenticationDataStore(namespace string) Authentication {
	if namespace == "" {
		return &AuthenticationApp{
			database_name: SEA_DATABASE_AUTH,
		}
	} else {
		return &AuthenticationApp{
			database_name: namespace + "_auth",
		}
	}
}
