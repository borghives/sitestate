package journal

import (
	"context"
	"log"

	"github.com/borghives/sitestate/data"
	"github.com/borghives/websession"
)

var webStore data.Web

func InitializeWebEvent(store data.Web) {
	webStore = store
}

func LogStartHost() {
	hostInfo := websession.GetHostInfo()
	log.Printf("START Instance@%s Build:%s Image:%s Run: %s", hostInfo.Id, hostInfo.BuildId, hostInfo.ImageId, hostInfo.AppCommand)
	webStore.HostInfo().InsertOne(context.Background(), hostInfo)
}

func LogStopHost() {
	hostInfo := websession.GetHostInfo()
	log.Printf("STOP Instance@%s Build:%s Image:%s ", hostInfo.Id, hostInfo.BuildId, hostInfo.ImageId)
	updateDirective := data.UpdateOperator{}
	updateDirective.CurrentDate("end_time")
	webStore.HostInfo().UpdateByID(context.Background(), hostInfo.Id, updateDirective.ToPrimitive())
}
