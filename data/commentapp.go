package data

import "go.mongodb.org/mongo-driver/mongo"

var DB_COMMENT_COLLECTION_NAME = "comment"
var DB_USER_COMMENT_COLLECTION_NAME = "user_comment"

type Comment interface {
	Initialize()
	Comment() *mongo.Collection
	UserToCommentRelation() *mongo.Collection
}

type CommentApp struct{}

func (p *CommentApp) Initialize() {
}

func (p *CommentApp) Comment() *mongo.Collection {
	return GetDatabase(SEA_DATABASE).Collection(DB_COMMENT_COLLECTION_NAME)
}

func (p *CommentApp) UserToCommentRelation() *mongo.Collection {
	return GetDatabase(SEA_DATABASE).Collection(DB_USER_COMMENT_COLLECTION_NAME)
}

func CommentDataStore() Comment {
	return &CommentApp{}
}
