package data

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

var SEA_DATABASE_REPORT_PREFIX = "report"
var SEA_DATABASE_REPORT_POSTFIX = "events"

var JOURNAL_DEFAULT_EVENTS_NAME = "default"

type Events interface {
	Initialize(topics []string)
	Default() *mongo.Collection
	GetStoreByTopic(topic string) *mongo.Collection
}

func getCollectionWithUnacknowledgedWrite(database string, name string) *mongo.Collection {
	wc := writeconcern.Unacknowledged()
	opts := options.Collection().SetWriteConcern(wc)
	return GetDatabase(database).Collection(name, opts)
}

type EventsApp struct {
	database_name string
	topics        []string
}

func (e *EventsApp) Initialize(topics []string) {
	e.topics = topics

	model := mongo.IndexModel{
		Keys: bson.M{
			"event_at": 1,
		},
		Options: options.Index().SetExpireAfterSeconds(12 * 30 * 24 * 60 * 60),
	}

	e.Default().Indexes().CreateOne(context.Background(), model)

	for _, topic := range topics {
		e.GetStoreByTopic(topic).Indexes().CreateOne(context.Background(), model)
	}
}

func (e *EventsApp) Default() *mongo.Collection {
	return getCollectionWithUnacknowledgedWrite(e.database_name, JOURNAL_DEFAULT_EVENTS_NAME)
}

func (e *EventsApp) GetStoreByTopic(topic string) *mongo.Collection {
	for _, v := range e.topics {
		if v == topic {
			return getCollectionWithUnacknowledgedWrite(e.database_name, topic)
		}
	}
	return e.Default()
}

func EventsDataStore(database_name string) Events {
	return &EventsApp{
		database_name: database_name,
		topics:        []string{},
	}
}
