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

type JournalEvent struct {
	ID                    primitive.ObjectID
	Session               *sitepages.WebSession
	TargetId              primitive.ObjectID
	EventCollectionClient *mongo.Collection
	Messages              []string
	CreatedAt             time.Time
	DoneAt                time.Time
	UpdateResult          *mongo.UpdateResult
	BulkResult            *mongo.BulkWriteResult
}

func CreateJournalEvent() *JournalEvent {
	return &JournalEvent{
		ID:                    primitive.NewObjectID(),
		EventCollectionClient: data.GetDefaultEventsCollection(),
		CreatedAt:             time.Now(),
		Messages:              []string{},
	}
}

func newJournalEvent(targetId primitive.ObjectID, dbCollection *mongo.Collection) *JournalEvent {
	return &JournalEvent{
		ID:                    primitive.NewObjectID(),
		TargetId:              targetId,
		EventCollectionClient: dbCollection,
		CreatedAt:             time.Now(),
		Messages:              []string{},
	}
}

type EventDocument struct {
	ID           primitive.ObjectID     `bson:"_id,omitempty"`
	SessionId    primitive.ObjectID     `bson:"session_id"`
	TargetId     primitive.ObjectID     `bson:"target_id"`
	Messages     []string               `bson:"messages"`
	CreatedAt    time.Time              `bson:"created_at"`
	Duration     time.Duration          `bson:"duration"`
	UpdateResult *mongo.UpdateResult    `bson:"update_result"`
	BulkResult   *mongo.BulkWriteResult `bson:"bulk_result"`
}

func (e *JournalEvent) SetTargetId(targetId primitive.ObjectID) {
	e.TargetId = targetId
}

func (e *JournalEvent) SetEventCollection(dbCollection *mongo.Collection) {
	e.EventCollectionClient = dbCollection
}

func (e *JournalEvent) SetSession(session *sitepages.WebSession) {
	e.Session = session
}

func (e *JournalEvent) SetResult(result *mongo.UpdateResult) {
	e.UpdateResult = result
}

func (e *JournalEvent) SetBulkResult(result *mongo.BulkWriteResult) {
	e.BulkResult = result
}

func (e *JournalEvent) AddMessage(message string) {
	log.Printf("journal (%s) msg: %s", e.ID.Hex(), message)
	e.Messages = append(e.Messages, message)
}

func (e *JournalEvent) ToDocument() EventDocument {
	done := time.Now()
	sessionId := primitive.NilObjectID
	if e.Session != nil {
		sessionId = e.Session.ID
	}
	return EventDocument{
		ID:           e.ID,
		SessionId:    sessionId,
		TargetId:     e.TargetId,
		Messages:     e.Messages,
		CreatedAt:    e.CreatedAt,
		Duration:     done.Sub(e.CreatedAt),
		UpdateResult: e.UpdateResult,
		BulkResult:   e.BulkResult,
	}
}

func (e *JournalEvent) Done() {
	doc := e.ToDocument()
	log.Printf("Journal %s Event Done: %s, Session: %s, Target: %s, Start: %s, Duration micro sec: %d",
		e.EventCollectionClient.Name(), doc.ID.Hex(), doc.SessionId.Hex(), doc.TargetId.Hex(), doc.CreatedAt, doc.Duration.Microseconds())
	e.EventCollectionClient.InsertOne(context.Background(), doc)

}
