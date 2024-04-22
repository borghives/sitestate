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
	Key   string   `bson:"key"`
	Infos []string `bson:"info, omitempty"`
}

type EventDocument struct {
	ID           primitive.ObjectID     `bson:"_id,omitempty"`
	HostInfo     primitive.ObjectID     `bson:"host_id"`
	SessionId    primitive.ObjectID     `bson:"session_id,omitempty"`
	TargetId     primitive.ObjectID     `bson:"target_id"`
	Statistics   []EventStat            `bson:"statistics,omitempty"`
	EventAt      time.Time              `bson:"event_at"`
	Duration     time.Duration          `bson:"duration"`
	UpdateResult *mongo.UpdateResult    `bson:"update_result,omitempty"`
	BulkResult   *mongo.BulkWriteResult `bson:"bulk_result,omitempty"`
}

type JournalEvent struct {
	ID               primitive.ObjectID
	HostInfo         sitepages.RutimeHostInfo
	Session          sitepages.WebSession
	TargetId         primitive.ObjectID
	JournalEventType data.EventType
	Statistics       []EventStat
	EventAt          time.Time
	DoneAt           time.Time
	UpdateResult     *mongo.UpdateResult
	BulkResult       *mongo.BulkWriteResult
}

var (
	eventStore data.Events
	webStore   data.WebApp
)

func InitializeJournal(store data.Events) {
	eventStore = store

}

func InitializeWebEvent(store data.WebApp) {
	webStore = store
}

func CreateJournalEvent() *JournalEvent {
	return &JournalEvent{
		ID:               primitive.NewObjectID(),
		HostInfo:         sitepages.GetHostInfo(),
		JournalEventType: data.EventDefault,
		EventAt:          time.Now(),
		Statistics:       []EventStat{},
	}
}

func (e *JournalEvent) SetTargetId(targetId primitive.ObjectID) {
	e.TargetId = targetId
}

func (e *JournalEvent) SetJournalEventType(eType data.EventType) {
	e.JournalEventType = eType
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

func (e *JournalEvent) AddStat(statKey string, infos ...string) {
	log.Printf("journal (%s:%s) key: %s, msg: %v", e.JournalEventType.String(), e.ID.Hex(), statKey, infos)
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
		e.JournalEventType.String(), doc.ID.Hex(), doc.SessionId.Hex(), doc.TargetId.Hex(), doc.EventAt, doc.Duration.Microseconds())

	eventStore.GetEventStoreByType(e.JournalEventType).InsertOne(context.Background(), doc)

}
