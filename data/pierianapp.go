package data

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var SEA_DATABASE = "pieriansea"

var DB_PAGE_COLLECTION_NAME = "page"
var DB_STANZA_COLLECTION_NAME = "stanza"
var DB_BUNDLE_COLLECTION_NAME = "bundle"

type Pierian interface {
	Initialize()
	Page() *mongo.Collection
	Stanza() *mongo.Collection
	Bundle() *mongo.Collection
}

type PierianApp struct{}

func (p *PierianApp) Initialize() {
	models := []mongo.IndexModel{
		{
			Keys: bson.M{
				"link": 1,
			},
		},
		{
			Keys: bson.M{
				"updated_time": 1,
			},
		},
		{
			Keys: bson.M{
				"created_time": -1,
			},
		},
	}
	_, err := p.Page().Indexes().CreateMany(context.Background(), models)
	if err != nil {
		log.Printf("error creating indexes 0: %s", err)
	}
}

func (*PierianApp) Page() *mongo.Collection {
	return GetDatabase(SEA_DATABASE).Collection(DB_PAGE_COLLECTION_NAME)
}

func (*PierianApp) Stanza() *mongo.Collection {
	return GetDatabase(SEA_DATABASE).Collection(DB_STANZA_COLLECTION_NAME)
}

func (*PierianApp) Bundle() *mongo.Collection {
	return GetDatabase(SEA_DATABASE).Collection(DB_BUNDLE_COLLECTION_NAME)
}

func PierianDataStore() Pierian {
	return &PierianApp{}
}
