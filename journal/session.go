package journal

import (
	"context"
	"log"
	"time"

	"github.com/borghives/sitepages"
	"github.com/borghives/sitestate/data"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB_SESSION_INFO_COLLECTION_NAME = "session_info"

func LogSession(session sitepages.WebSession) {
	opt := options.Update().SetUpsert(true)
	update := data.NewUpdate().SetDoc(session).SetOnInsert("created_time", time.Now()).CurrentDate("last_seen")

	result, err := data.GetSessionInfoCollection().UpdateByID(context.Background(), session.ID, update.ToPrimitive(), opt)
	if err != nil {
		log.Printf("Error logging session: %v", err)
		return
	}

	if result.UpsertedCount == 0 {
		return
	}
}
