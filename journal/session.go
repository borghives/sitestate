package journal

import (
	"context"
	"log"
	"time"

	"github.com/borghives/sitestate/data"
	"github.com/borghives/websession"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB_SESSION_INFO_COLLECTION_NAME = "session_info"

func LogSession(session websession.Session) {
	opt := options.Update().SetUpsert(true)
	update := data.NewUpdate().SetDoc(session).SetOnInsert("event_at", time.Now()).CurrentDate("last_seen")

	result, err := webStore.SessionInfo().UpdateByID(context.Background(), session.ID, update.ToPrimitive(), opt)
	if err != nil {
		log.Printf("Error logging session: %v", err)
		return
	}

	if result.UpsertedCount == 0 {
		return
	}
}
