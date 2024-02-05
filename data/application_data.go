package data

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB_PAGE_COLLECTION_NAME = "page"
var DB_STANZA_COLLECTION_NAME = "stanza"
var DB_HOST_INFO_COLLECTION_NAME = "hostinfo"
var DB_SESSION_INFO_COLLECTION_NAME = "session_info"
var DB_JOURNAL_DEFAULT_EVENTS_COLLECTION_NAME = "journal_events"
var DB_JOURNAL_AGGREGATE_EVENTS_COLLECTION_NAME = "agg_events"
var DB_JOURNAL_PUT_PAGE_EVENTS_COLLECTION_NAME = "put_page_events"
var DB_JOURNAL_GET_PAGE_EVENTS_COLLECTION_NAME = "get_page_events"
var DB_JOURNAL_PUT_STANZA_EVENTS_COLLECTION_NAME = "put_stanza_events"
var DB_JOURNAL_GET_STANZA_EVENTS_COLLECTION_NAME = "get_stanza_events"

func GetPageCollection() *mongo.Collection {
	return GetDB().Collection(DB_PAGE_COLLECTION_NAME)
}

func GetStanzaCollection() *mongo.Collection {
	return GetDB().Collection(DB_STANZA_COLLECTION_NAME)
}

func GetSessionInfoCollection() *mongo.Collection {
	return GetDB().Collection(DB_SESSION_INFO_COLLECTION_NAME)
}

func GetHostInfoCollection() *mongo.Collection {
	return GetDB().Collection(DB_HOST_INFO_COLLECTION_NAME)
}

func GetDefaultEventsCollection() *mongo.Collection {
	return GetDB_Unacknowledged().Collection(DB_JOURNAL_DEFAULT_EVENTS_COLLECTION_NAME)
}

func GetAggregateEventsCollection() *mongo.Collection {
	return GetDB_Unacknowledged().Collection(DB_JOURNAL_AGGREGATE_EVENTS_COLLECTION_NAME)
}

func GetPutPageEventsCollection() *mongo.Collection {
	return GetDB_Unacknowledged().Collection(DB_JOURNAL_PUT_PAGE_EVENTS_COLLECTION_NAME)
}

func GetGetPageEventsCollection() *mongo.Collection {
	return GetDB_Unacknowledged().Collection(DB_JOURNAL_GET_PAGE_EVENTS_COLLECTION_NAME)
}

func GetPutStanzaEventsCollection() *mongo.Collection {
	return GetDB_Unacknowledged().Collection(DB_JOURNAL_PUT_STANZA_EVENTS_COLLECTION_NAME)
}

func GetGetStanzaEventsCollection() *mongo.Collection {
	return GetDB_Unacknowledged().Collection(DB_JOURNAL_GET_STANZA_EVENTS_COLLECTION_NAME)
}

func GetDayGetPageReportsCollection() *mongo.Collection {
	return GetDB().Collection("day_getpage")
}

func GetDayGetStanzaReportsCollection() *mongo.Collection {
	return GetDB().Collection("day_getstanza")
}

func GetDayPutPageReportsCollection() *mongo.Collection {
	return GetDB().Collection("day_putpage")
}

func GetDayPutStanzaReportsCollection() *mongo.Collection {
	return GetDB().Collection("day_putstanza")
}

func EnsurePageIndexes() {
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
	}
	_, err := GetPageCollection().Indexes().CreateMany(context.Background(), models)
	if err != nil {
		log.Printf("error creating indexes 0: %s", err)
	}
}

func EnsureEventIndexes() {
	model := mongo.IndexModel{
		Keys: bson.M{
			"created_at": 1,
		},
		Options: options.Index().SetExpireAfterSeconds(12 * 30 * 24 * 60 * 60),
	}

	_, err := GetDefaultEventsCollection().Indexes().CreateOne(context.Background(), model)
	if err != nil {
		log.Printf("error creating indexes default event: %s", err)
	}

	_, err = GetPutPageEventsCollection().Indexes().CreateOne(context.Background(), model)
	if err != nil {
		log.Printf("error creating indexes put page event: %s", err)
	}

	_, err = GetGetPageEventsCollection().Indexes().CreateOne(context.Background(), model)
	if err != nil {
		log.Printf("error creating indexes get page event: %s", err)
	}

	_, err = GetPutStanzaEventsCollection().Indexes().CreateOne(context.Background(), model)
	if err != nil {
		log.Printf("error creating indexes put stanza event: %s", err)
	}

	_, err = GetGetStanzaEventsCollection().Indexes().CreateOne(context.Background(), model)
	if err != nil {
		log.Printf("error creating indexes get stanza event: %s", err)
	}
}

func EnsureReportsIndexes() {
	models := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "target_id", Value: 1},
				{Key: "floor", Value: -1},
				{Key: "ceiling", Value: -1},
			},
		},
		{
			Keys: bson.D{
				{Key: "target_id", Value: 1},
				{Key: "floor", Value: -1},
			},
		},
		{
			Keys: bson.D{
				{Key: "target_id", Value: 1},
			},
		},
		{
			Keys: bson.D{
				{Key: "floor", Value: -1},
			},
		},
		{
			Keys: bson.D{
				{Key: "floor", Value: -1},
			},
		},
	}

	_, err := GetDayGetPageReportsCollection().Indexes().CreateMany(context.Background(), models)
	if err != nil {
		log.Printf("error creating indexes day get page reports: %s", err)
	}

	_, err = GetDayPutPageReportsCollection().Indexes().CreateMany(context.Background(), models)
	if err != nil {
		log.Printf("error creating indexes day put page reports: %s", err)
	}

	_, err = GetDayGetStanzaReportsCollection().Indexes().CreateMany(context.Background(), models)
	if err != nil {
		log.Printf("error creating indexes day get stanza reports: %s", err)
	}

	_, err = GetDayPutStanzaReportsCollection().Indexes().CreateMany(context.Background(), models)
	if err != nil {
		log.Printf("error creating indexes day put stanza reports: %s", err)
	}

}

func EnsureIndexes() {
	EnsurePageIndexes()
	EnsureEventIndexes()
	EnsureReportsIndexes()
}
