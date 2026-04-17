package journal

import (
	"log"

	"github.com/borghives/websession"
)

func LogStartHost() {
	hostInfo := websession.GetHostInfo()
	log.Printf("START Instance@%s Build:%s Image:%s Run: %s", hostInfo.ID, hostInfo.BuildId, hostInfo.ImageId, hostInfo.AppCommand)
}

func LogStopHost() {
	hostInfo := websession.GetHostInfo()
	log.Printf("STOP Instance@%s Build:%s Image:%s ", hostInfo.ID, hostInfo.BuildId, hostInfo.ImageId)
}
