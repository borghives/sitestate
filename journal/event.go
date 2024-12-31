package journal

import (
	"context"
	"log"
	"time"

	"github.com/borghives/sitestate/data"
	"github.com/borghives/websession"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type EventStat struct {
	Key   string   `bson:"key"`
	Infos []string `bson:"info, omitempty"`
}

type EventDocument struct {
	ID           primitive.ObjectID     `bson:"_id,omitempty"`
	HostInfo     primitive.ObjectID     `bson:"host_id"`
	SessionId    primitive.ObjectID     `bson:"session_id,omitempty"`
	TargetId     primitive.ObjectID     `bson:"target_id"`
	Topic        string                 `bson:"topic"`
	Statistics   []EventStat            `bson:"statistics,omitempty"`
	EventAt      time.Time              `bson:"event_at"`
	Duration     time.Duration          `bson:"duration"`
	UpdateResult *mongo.UpdateResult    `bson:"update_result,omitempty"`
	BulkResult   *mongo.BulkWriteResult `bson:"bulk_result,omitempty"`
}

type JournalEvent struct {
	ID           primitive.ObjectID
	HostInfo     websession.RutimeHostInfo
	Session      websession.Session
	TargetId     primitive.ObjectID
	Topic        string
	Statistics   []EventStat
	EventAt      time.Time
	DoneAt       time.Time
	UpdateResult *mongo.UpdateResult
	BulkResult   *mongo.BulkWriteResult
}

var (
	eventStore data.Events
)

func InitializeJournal(event data.Events) {
	eventStore = event
}

func CreateJournalEvent() *JournalEvent {
	return &JournalEvent{
		ID:         primitive.NewObjectID(),
		HostInfo:   websession.GetHostInfo(),
		Topic:      data.JOURNAL_DEFAULT_EVENTS_NAME,
		EventAt:    time.Now(),
		Statistics: []EventStat{},
	}
}

func (e *JournalEvent) SetTargetId(targetId primitive.ObjectID) {
	e.TargetId = targetId
}

func (e *JournalEvent) SetJournalTopic(topic string) {
	e.Topic = topic
}

func (e *JournalEvent) SetSession(session websession.Session) {
	LogSession(session)
	e.Session = session
}

func (e *JournalEvent) SetResult(result *mongo.UpdateResult) {
	e.UpdateResult = result
}

func (e *JournalEvent) SetBulkResult(result *mongo.BulkWriteResult) {
	e.BulkResult = result
}

func (e *JournalEvent) AddStat(statKey string, infos ...string) {
	log.Printf("journal (%s:%s) key: %s, msg: %v", e.Topic, e.ID.Hex(), statKey, infos)
	e.Statistics = append(e.Statistics, EventStat{
		Key:   statKey,
		Infos: infos,
	})
}

func (e *JournalEvent) ToDocument() EventDocument {
	done := time.Now()

	return EventDocument{
		ID:           e.ID,
		HostInfo:     e.HostInfo.Id,
		SessionId:    e.Session.ID,
		TargetId:     e.TargetId,
		Topic:        e.Topic,
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
		e.Topic, doc.ID.Hex(), doc.SessionId.Hex(), doc.TargetId.Hex(), doc.EventAt, doc.Duration.Microseconds())

	eventStore.GetStoreByTopic(e.Topic).InsertOne(context.Background(), doc)

}
