package data

import "go.mongodb.org/mongo-driver/mongo"

var DB_HOST_INFO_COLLECTION_NAME = "hostinfo"
var DB_SESSION_INFO_COLLECTION_NAME = "session_info"

type Web interface {
	Initialize()
	HostInfo() *mongo.Collection
	SessionInfo() *mongo.Collection
}

type WebApp struct{}

func (w *WebApp) Initialize() {
}

func (w *WebApp) HostInfo() *mongo.Collection {
	return GetDatabase(SEA_DATABASE).Collection(DB_HOST_INFO_COLLECTION_NAME)
}

func (w *WebApp) SessionInfo() *mongo.Collection {
	return GetDatabase(SEA_DATABASE).Collection(DB_SESSION_INFO_COLLECTION_NAME)
}

// webapp.go

func WebDataStore() Web {
	return &WebApp{}
}
