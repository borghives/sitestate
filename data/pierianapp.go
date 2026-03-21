package data

import (
	"context"
	"log"

	"github.com/borghives/sitepages"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var SEA_DATABASE = "pieriansea"

var DB_PAGE_COLLECTION_NAME = "page"
var DB_STANZA_COLLECTION_NAME = "stanza"
var DB_BUNDLE_COLLECTION_NAME = "bundle"
var DB_USER_PAGE_COLLECTION_NAME = "user_page"

type Pierian interface {
	Initialize()
	Page() *mongo.Collection
	Stanza() *mongo.Collection
	Bundle() *mongo.Collection
	Relation(graphType sitepages.RelationGraphType) *mongo.Collection
}

type PierianApp struct{}

func (p *PierianApp) Initialize() {
	pageIndex := []mongo.IndexModel{
		{
			Keys: D{{Key: "link", Value: 1}},
		},
		{
			Keys: D{{Key: "updated_time", Value: 1}},
		},
		{
			Keys: D{{Key: "created_time", Value: -1}},
		},
	}
	_, err := p.Page().Indexes().CreateMany(context.Background(), pageIndex)
	if err != nil {
		log.Printf("error creating indexes 0: %s", err)
	}

	relationIndex := []mongo.IndexModel{
		{
			Keys: D{{Key: "subjid", Value: 1}},
		},
		{
			Keys: D{{Key: "objid", Value: 1}},
		},
		{
			Keys: D{{Key: "event_at", Value: -1}},
		},
		{
			Keys: D{
				{Key: "subjid", Value: 1},
				{Key: "objid", Value: 1},
			},
		},
		{
			Keys: D{
				{Key: "subjid", Value: 1},
				{Key: "relation", Value: 1},
			},
		},
		{
			Keys: D{
				{Key: "subjid", Value: 1},
				{Key: "objid", Value: 1},
				{Key: "relation", Value: 1},
			},
			Options: options.Index().SetUnique(true),
		},
	}

	_, err = p.Relation(sitepages.RelationGraphType_UserPage).Indexes().CreateMany(context.Background(), relationIndex)
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

func (*PierianApp) Relation(graphType sitepages.RelationGraphType) *mongo.Collection {
	database := GetDatabase(SEA_DATABASE)
	if graphType == sitepages.RelationGraphType_UserPage {
		return database.Collection(DB_USER_PAGE_COLLECTION_NAME)
	}

	return database.Collection("rel_" + graphType.String())
}

func PierianDataStore() Pierian {
	return &PierianApp{}
}
