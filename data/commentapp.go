package data

import "go.mongodb.org/mongo-driver/mongo"

var DB_COMMENT_COLLECTION_NAME = "comment"
var DB_USER_COMMENT_COLLECTION_NAME = "user_comment"

type MsgFeed interface {
	Initialize()
	Comment() *mongo.Collection
	UserToCommentRelation() *mongo.Collection
}

type MsgFeedApp struct{}

func (p *MsgFeedApp) Initialize() {
}

func (p *MsgFeedApp) Comment() *mongo.Collection {
	return GetDatabase(SEA_DATABASE).Collection(DB_COMMENT_COLLECTION_NAME)
}

func (p *MsgFeedApp) UserToCommentRelation() *mongo.Collection {
	return GetDatabase(SEA_DATABASE).Collection(DB_USER_COMMENT_COLLECTION_NAME)
}

func CommentDataStore() MsgFeed {
	return &MsgFeedApp{}
}
