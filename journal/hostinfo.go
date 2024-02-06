package journal

import (
	"context"
	"log"

	"github.com/borghives/sitepages"
	"github.com/borghives/sitestate/data"
)

func GetHostInfo() sitepages.RutimeHostInfo {
	return sitepages.NewHostInstanceInfo()
}

func LogStartHost() {
	hostInfo := GetHostInfo()
	log.Printf("START Instance@%s Build:%s Image:%s Run: %s", hostInfo.Id, hostInfo.BuildId, hostInfo.ImageId, hostInfo.AppCommand)
	data.GetHostInfoCollection().InsertOne(context.Background(), hostInfo)
}

func LogStopHost() {
	hostInfo := GetHostInfo()
	log.Printf("STOP Instance@%s Build:%s Image:%s ", hostInfo.Id, hostInfo.BuildId, hostInfo.ImageId)
	updateDirective := data.UpdateOperator{}
	updateDirective.CurrentDate("end_time")
	data.GetHostInfoCollection().UpdateByID(context.Background(), hostInfo.Id, updateDirective.ToPrimitive())
}
