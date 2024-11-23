package data

import "go.mongodb.org/mongo-driver/mongo"

var DB_HOST_INFO_COLLECTION_NAME = "hostinfo"
var DB_SESSION_INFO_COLLECTION_NAME = "session_info"

type Web interface {
	Initialize()
	HostInfo() *mongo.Collection
	SessionInfo() *mongo.Collection
}

type WebApp struct {
	database_name string
}

func (w *WebApp) Initialize() {
}

func (w *WebApp) HostInfo() *mongo.Collection {
	return GetDatabase(w.database_name).Collection(DB_HOST_INFO_COLLECTION_NAME)
}

func (w *WebApp) SessionInfo() *mongo.Collection {
	return GetDatabase(w.database_name).Collection(DB_SESSION_INFO_COLLECTION_NAME)
}

// webapp.go

func WebDataStore(database_name string) Web {
	return &WebApp{
		database_name: database_name,
	}
}
