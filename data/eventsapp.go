package data

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

var SEA_DATABASE_REPORT = "pierianseareport"

var DB_JOURNAL_DEFAULT_EVENTS_COLLECTION_NAME = "journal_events"
var DB_JOURNAL_PUT_PAGE_EVENTS_COLLECTION_NAME = "put_page_events"
var DB_JOURNAL_GET_PAGE_EVENTS_COLLECTION_NAME = "get_page_events"
var DB_JOURNAL_PUT_STANZA_EVENTS_COLLECTION_NAME = "put_stanza_events"
var DB_JOURNAL_GET_STANZA_EVENTS_COLLECTION_NAME = "get_stanza_events"

type EventType int

const (
	EventDefault EventType = iota
	EventPutPage
	EventGetPage
	EventPutStanza
	EventGetStanza
)

func (e EventType) String() string {
	return [...]string{
		"Default",
		"PutPage",
		"GetPage",
		"PutStanza",
		"GetStanza",
	}[e]
}

type Events interface {
	Initialize()
	Default() *mongo.Collection
	PutPage() *mongo.Collection
	GetPage() *mongo.Collection
	PutStanza() *mongo.Collection
	GetStanza() *mongo.Collection
	GetEventStoreByType(eType EventType) *mongo.Collection
}

func getCollectionWithUnacknowledgedWrite(name string) *mongo.Collection {
	wc := writeconcern.Unacknowledged()
	opts := options.Collection().SetWriteConcern(wc)
	return GetDatabase(SEA_DATABASE).Collection(name, opts)
}

type EventsApp struct{}

func (e *EventsApp) Initialize() {
	model := mongo.IndexModel{
		Keys: bson.M{
			"event_at": 1,
		},
		Options: options.Index().SetExpireAfterSeconds(12 * 30 * 24 * 60 * 60),
	}

	_, err := e.Default().Indexes().CreateOne(context.Background(), model)
	if err != nil {
		log.Printf("error creating indexes default event: %s", err)
	}

	_, err = e.PutPage().Indexes().CreateOne(context.Background(), model)
	if err != nil {
		log.Printf("error creating indexes put page event: %s", err)
	}

	_, err = e.GetPage().Indexes().CreateOne(context.Background(), model)
	if err != nil {
		log.Printf("error creating indexes get page event: %s", err)
	}

	_, err = e.PutStanza().Indexes().CreateOne(context.Background(), model)
	if err != nil {
		log.Printf("error creating indexes put stanza event: %s", err)
	}

	_, err = e.GetStanza().Indexes().CreateOne(context.Background(), model)
	if err != nil {
		log.Printf("error creating indexes get stanza event: %s", err)
	}
}

func (e *EventsApp) Default() *mongo.Collection {
	return getCollectionWithUnacknowledgedWrite(DB_JOURNAL_DEFAULT_EVENTS_COLLECTION_NAME)
}

func (e *EventsApp) PutPage() *mongo.Collection {
	return getCollectionWithUnacknowledgedWrite(DB_JOURNAL_PUT_PAGE_EVENTS_COLLECTION_NAME)
}

func (e *EventsApp) GetPage() *mongo.Collection {
	return getCollectionWithUnacknowledgedWrite(DB_JOURNAL_GET_PAGE_EVENTS_COLLECTION_NAME)
}

func (e *EventsApp) PutStanza() *mongo.Collection {
	return getCollectionWithUnacknowledgedWrite(DB_JOURNAL_PUT_STANZA_EVENTS_COLLECTION_NAME)
}

func (e *EventsApp) GetStanza() *mongo.Collection {
	return getCollectionWithUnacknowledgedWrite(DB_JOURNAL_GET_STANZA_EVENTS_COLLECTION_NAME)
}

func (e *EventsApp) GetEventStoreByType(eType EventType) *mongo.Collection {
	switch eType {
	case EventPutPage:
		return e.PutPage()
	case EventGetPage:
		return e.GetPage()
	case EventPutStanza:
		return e.PutStanza()
	case EventGetStanza:
		return e.GetStanza()
	default:
		return e.Default()
	}
}

func EventsDataStore() Events {
	return &EventsApp{}
}
