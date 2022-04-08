package main

import (
	"log"

	"github.com/biges/mongo-db"
	"github.com/fatih/structs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Obj struct {
	ID    primitive.ObjectID `json:"id" bson:"_id,omitempty" structs:"_id,omitempty"`
	Field string             `json:"field,omitempty" bson:"field,omitempty" structs:"field,omitempty"`
}

func main() {
	obj := Obj{}
	objs := []Obj{}

	// create
	id, err := db.Mongo.Insert("mocks", &obj)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(id)

	// find
	err = db.Mongo.Find("mocks", bson.M{"field": "think"}, &objs, db.Mongo.NewPaginationParams())
	if err != nil {
		log.Fatal(err)
	}

	err = db.Mongo.FindOne("mocks", bson.M{"field": "think"}, &obj)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Mongo.FindOne("mocks", bson.M{"_id": id}, &obj)
	if err != nil {
		log.Fatal(err)
	}

	// update
	change := Obj{
		Field: "change",
	}

	fields := structs.Map(change)

	err = db.Mongo.Update("mocks", bson.M{"_id": id}, bson.M{"$set": fields})
	if err != nil {
		log.Fatal(err)
	}

	// count
	_, err = db.Mongo.Count("mocks", bson.M{"field": "cem"})
	if err != nil {
		log.Fatal(err)
	}

	// aggregate
	err = db.Mongo.Aggregate("mocks", bson.M{"field": "cem"}, objs)
	if err != nil {
		log.Fatal(err)
	}

	// index
	name := "new"
	err = db.Mongo.CreateIndex("mocks", bson.M{
		"created_at": 1,
	}, &options.IndexOptions{
		Name: &name,
	})
	if err != nil {
		log.Fatal(err)
	}
}
