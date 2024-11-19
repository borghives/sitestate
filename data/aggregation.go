package data

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Aggregation struct {
	pipeline mongo.Pipeline
}

func NewAggregation() *Aggregation {
	return &Aggregation{}
}

func (a *Aggregation) Match(filter primitive.M) *Aggregation {
	a.pipeline = append(a.pipeline, primitive.D{{Key: "$match", Value: filter}})
	return a
}

func (a *Aggregation) Group(group primitive.M) *Aggregation {
	a.pipeline = append(a.pipeline, primitive.D{{Key: "$group", Value: group}})
	return a
}

func (a *Aggregation) Lookup(lookup primitive.M) *Aggregation {
	a.pipeline = append(a.pipeline, primitive.D{{Key: "$lookup", Value: lookup}})
	return a
}

func (a *Aggregation) AddFields(field primitive.M) *Aggregation {
	a.pipeline = append(a.pipeline, primitive.D{{Key: "$addFields", Value: field}})
	return a
}

func (a *Aggregation) Project(fields primitive.M) *Aggregation {
	a.pipeline = append(a.pipeline, primitive.D{{Key: "$project", Value: fields}})
	return a
}

func (a *Aggregation) Sort(fields primitive.D) *Aggregation {
	a.pipeline = append(a.pipeline, primitive.D{{Key: "$sort", Value: fields}})
	return a
}

func (a *Aggregation) Limit(value int64) *Aggregation {
	a.pipeline = append(a.pipeline, primitive.D{{Key: "$limit", Value: value}})
	return a
}

func (a *Aggregation) Search(fields primitive.D) *Aggregation {
	a.pipeline = append(a.pipeline, primitive.D{{Key: "$search", Value: fields}})
	return a
}

func (a *Aggregation) AppendFrom(agg *Aggregation) *Aggregation {
	a.pipeline = append(a.pipeline, agg.pipeline...)
	return a
}

func (a *Aggregation) Pipeline() mongo.Pipeline {
	return a.pipeline
}
