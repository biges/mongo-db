![Version](https://img.shields.io/badge/version-1.3.0-yellow.svg)

# db

Bu repository diğer servisleri etkilememek amacıyla Hybrone Sentinel Akbank projesi için açılmıştır.

DB is open a connection to database. DB could send data to  new relic.

Supported databases:

- Mongo

##  Set your env

Configuration is getting from your shell.

```bash
export GO_ENV=development
export MONGO_DB_URI=mongo://localhost:27017
export MONGO_DB_DEBUG_MODE=true
export MONGO_DB_NAME=my_awsome_db
export MONGO_DB_TIMEOUT=10
```

## Usage

```go
package main

import (
	"log"

	"github.com/biges/db"
	"github.com/fatih/structs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
```

---

## Rake Tasks

```bash
$ rake -T

rake default                    # show avaliable tasks (default task)
rake docker:build               # Build
rake release:check              # do release check
rake release:publish[revision]  # Publish project with revision: major,minor,patch, default: patch
rake serve_doc[port]            # run doc server at :port (default: 6060)
rake test[verbose]              # run tests
rake verify[tag]                # Verify package by tag
```

---

## Publish New Version

```bash
$ rake release:check

$ rake release:publish               # patch level -> 1.0.0 => 1.0.1
$ rake release:publish[minor]        # minor level -> 1.0.0 => 1.1.0
$ rake release:publish[major]        # major level -> 1.0.0 => 2.0.0

$ rake verify
```
