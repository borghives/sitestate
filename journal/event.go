package journal

import (
	"context"
	"log"
	"time"

	"github.com/borghives/sitepages"
	"github.com/borghives/sitestate/data"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type EventStat struct {
	Key  string `bson:"key"`
	Info string `bson:"info, omitempty"`
}

type EventDocument struct {
	ID           primitive.ObjectID     `bson:"_id,omitempty"`
	SessionId    primitive.ObjectID     `bson:"session_id"`
	TargetId     primitive.ObjectID     `bson:"target_id"`
	Statistics   []EventStat            `bson:"statistics,omitempty"`
	EventAt      time.Time              `bson:"event_at,omitempty"`
	Duration     time.Duration          `bson:"duration"`
	UpdateResult *mongo.UpdateResult    `bson:"update_result,omitempty"`
	BulkResult   *mongo.BulkWriteResult `bson:"bulk_result,omitempty"`
}

type JournalEvent struct {
	ID                    primitive.ObjectID
	Session               sitepages.WebSession
	TargetId              primitive.ObjectID
	EventCollectionClient *mongo.Collection
	Statistics            []EventStat
	EventAt               time.Time
	DoneAt                time.Time
	UpdateResult          *mongo.UpdateResult
	BulkResult            *mongo.BulkWriteResult
}

func CreateJournalEvent() *JournalEvent {
	return &JournalEvent{
		ID:                    primitive.NewObjectID(),
		EventCollectionClient: data.GetDefaultEventsCollection(),
		EventAt:               time.Now(),
		Statistics:            []EventStat{},
	}
}

func (e *JournalEvent) SetTargetId(targetId primitive.ObjectID) {
	e.TargetId = targetId
}

func (e *JournalEvent) SetEventCollection(dbCollection *mongo.Collection) {
	e.EventCollectionClient = dbCollection
}

func (e *JournalEvent) SetSession(session sitepages.WebSession) {
	LogSession(session)
	e.Session = session
}

func (e *JournalEvent) SetResult(result *mongo.UpdateResult) {
	e.UpdateResult = result
}

func (e *JournalEvent) SetBulkResult(result *mongo.BulkWriteResult) {
	e.BulkResult = result
}

func (e *JournalEvent) AddStat(statKey string, info string) {
	log.Printf("journal (%s) key: %s, msg: %s", e.ID.Hex(), statKey, info)
	e.Statistics = append(e.Statistics, EventStat{
		Key:  statKey,
		Info: info,
	})
}

func (e *JournalEvent) ToDocument() EventDocument {
	done := time.Now()

	return EventDocument{
		ID:           e.ID,
		SessionId:    e.Session.ID,
		TargetId:     e.TargetId,
		Statistics:   e.Statistics,
		EventAt:      e.EventAt,
		Duration:     done.Sub(e.EventAt),
		UpdateResult: e.UpdateResult,
		BulkResult:   e.BulkResult,
	}
}

func (e *JournalEvent) Done() {
	doc := e.ToDocument()
	log.Printf("Journal %s Event Done: %s, Session: %s, Target: %s, Start: %s, Duration micro sec: %d",
		e.EventCollectionClient.Name(), doc.ID.Hex(), doc.SessionId.Hex(), doc.TargetId.Hex(), doc.EventAt, doc.Duration.Microseconds())
	e.EventCollectionClient.InsertOne(context.Background(), doc)

}
