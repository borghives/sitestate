package data

import (
	"bytes"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestUpdate(t *testing.T) {
	// ...

	updateNative := bson.M{
		"$set": bson.M{
			"name": "John",
		},
	}

	blobNative, err := bson.Marshal(updateNative)
	if err != nil {
		t.Fatal(err)
	}

	if len(blobNative) != 31 {
		t.Errorf("expected 31 bytes, got %d", len(blobNative))
	}

	var update UpdateOperator
	update.Set("name", "John")
	updatePrimitive := update.ToPrimitive()
	blob, err := bson.Marshal(updatePrimitive)
	if err != nil {
		t.Fatal(err)
	}

	if len(blob) != len(blobNative) {
		t.Errorf("expected %d bytes, got %d", len(blobNative), len(blob))
	}

	updateCheck := bson.M{}
	err = bson.Unmarshal(blob, &updateCheck)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(blobNative, blob) {
		t.Errorf("expected %v, got %v", updateNative, updateCheck)
	}

	// ...
}

type TestDocObject struct {
	Name string `bson:"name"`
	Age  int    `bson:"age"`
}

func TestUpdateDoc(t *testing.T) {
	doc := TestDocObject{
		Name: "John",
		Age:  20,
	}
	blob, err := bson.Marshal(doc)
	if err != nil {
		t.Fatal(err)
	}

	var doc2 primitive.D
	err = bson.Unmarshal(blob, &doc2)
	if err != nil {
		t.Fatal(err)
	}

	if len(doc2) != 2 {
		t.Errorf("expected 2 elements, got %d", len(doc2))
	}

	blob2, err := bson.Marshal(doc2)
	if err != nil {
		t.Fatal(err)
	}

	var doc3 TestDocObject
	err = bson.Unmarshal(blob2, &doc3)
	if err != nil {
		t.Fatal(err)
	}

	if doc3 != doc {
		t.Errorf("expected %v, got %v", doc, doc3)
	}
}
