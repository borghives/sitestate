package data

import (
	"context"
	"encoding/json"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Aggregation struct {
	pipeline mongo.Pipeline
}

func NewAggregation() *Aggregation {
	return &Aggregation{}
}

func (a *Aggregation) Match(filter bson.M) *Aggregation {
	a.pipeline = append(a.pipeline, bson.D{{Key: "$match", Value: filter}})
	return a
}

func (a *Aggregation) Group(group bson.M) *Aggregation {
	a.pipeline = append(a.pipeline, bson.D{{Key: "$group", Value: group}})
	return a
}

func (a *Aggregation) Lookup(lookup bson.M) *Aggregation {
	a.pipeline = append(a.pipeline, bson.D{{Key: "$lookup", Value: lookup}})
	return a
}

func (a *Aggregation) AddFields(field bson.M) *Aggregation {
	a.pipeline = append(a.pipeline, bson.D{{Key: "$addFields", Value: field}})
	return a
}

func (a *Aggregation) Project(fields bson.M) *Aggregation {
	a.pipeline = append(a.pipeline, bson.D{{Key: "$project", Value: fields}})
	return a
}

func (a *Aggregation) Sort(fields bson.D) *Aggregation {
	a.pipeline = append(a.pipeline, bson.D{{Key: "$sort", Value: fields}})
	return a
}

func (a *Aggregation) Limit(value int64) *Aggregation {
	a.pipeline = append(a.pipeline, bson.D{{Key: "$limit", Value: value}})
	return a
}

func (a *Aggregation) Search(fields bson.M) *Aggregation {
	a.pipeline = append(a.pipeline, bson.D{{Key: "$search", Value: fields}})
	return a
}

func (a *Aggregation) AppendFrom(agg *Aggregation) *Aggregation {
	a.pipeline = append(a.pipeline, agg.pipeline...)
	return a
}

func (a *Aggregation) AggregatePipeline(ctx context.Context, collection *mongo.Collection) (*mongo.Cursor, error) {
	return collection.Aggregate(ctx, a.Pipeline())
}

func (a *Aggregation) Pipeline() mongo.Pipeline {
	return a.pipeline
}

// mainly for debugging
func (a *Aggregation) JsonString() string {
	// Convert pipeline to bson.A
	bsonArray := bson.A{}
	for _, stage := range a.pipeline {
		bsonArray = append(bsonArray, stage)
	}

	// Marshal bson.A to JSON
	jsonString, err := json.Marshal(bsonArray)
	if err != nil {
		panic(err)
	}

	return string(jsonString)
}
