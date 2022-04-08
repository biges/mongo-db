package db

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/biges/mongo-db/cnf"
	"github.com/biges/logger"
	nr "github.com/biges/newrelic"
	"github.com/newrelic/go-agent/v3/integrations/nrmongo"
	"github.com/newrelic/go-agent/v3/newrelic"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoStorageOfficial holds Session and dial info of MongoDB connection
type MongoStorageOfficial struct {
	options                 *options.ClientOptions
	Client                  *mongo.Client
	Session                 *mongo.Database
	newRelicApp             *newrelic.Application
	DefaultPaginationParams *PaginationParams
}

var (
	// Mongo is a constant definition for mongo.
	Mongo *MongoStorageOfficial
)

// NewMongoStorageOfficial returns a new MongoStorage with an active Session
func NewMongoStorageOfficial(uri string) (*MongoStorageOfficial, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cnf.Cnf.Mongo.MongoDBTimeout)*time.Second)
	defer cancel()

	// set monitor
	monitor := &event.CommandMonitor{}
	var err error

	// if cnf.Cnf.Mongo.MongoDBDebugMode {
	// 	monitor = mongonitor.NewMongonitor()
	// } else {
	// 	monitor = nil
	// }

	nrMon := nrmongo.NewCommandMonitor(monitor)
	clientOptions := options.Client().ApplyURI(uri).SetMonitor(nrMon).
		SetCompressors([]string{"zstd"}).SetZstdLevel(10)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("client can't connect to MongoDB: %v", err)
	}

	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("mongodb connection error: %v", err)
	}

	// new database
	database := client.Database(cnf.Cnf.Mongo.MongoDBName)

	return &MongoStorageOfficial{
		Session:     database,
		Client:      client,
		options:     clientOptions,
		newRelicApp: nr.App,
		DefaultPaginationParams: &PaginationParams{
			Limit:  50,
			SortBy: "_id",
			Page:   1,
		},
	}, nil
}

func init() {
	var dbErr error
	Mongo, dbErr = NewMongoStorageOfficial(cnf.Cnf.Mongo.MongoDBURI)
	if dbErr != nil {
		logger.Zap.Fatal(dbErr)
	}
}

// Find returns all matching documents with filter and pagination params.
// page starts with 1.
func (s *MongoStorageOfficial) Find(collectionName string, filter interface{}, result interface{}, pagination *PaginationParams) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cnf.Cnf.Mongo.MongoDBTimeout)*time.Second)
	defer cancel()

	txn := new(newrelic.Transaction)
	if s.newRelicApp != nil {
		txn = s.newRelicApp.StartTransaction("mongo-find")
		ctx = newrelic.NewContext(ctx, txn)
	}

	//filter options
	filterOptions := options.Find()
	if pagination == nil {
		pagination = s.DefaultPaginationParams
	}

	// fix 0 pagination
	if pagination.Page == 0 {
		pagination.Page = 1
	}

	skipVal := int64((pagination.Page - 1) * pagination.Limit)
	limitVal := int64(pagination.Limit)

	if pagination.SortBy == "" {
		pagination.SortBy = "_id"
	}

	sortBy := bson.D{}
	for _, sortOpt := range strings.Split(pagination.SortBy, ",") {
		if string(sortOpt[0]) == "-" {
			sortBy = append(sortBy, bson.E{
				Key:   string(sortOpt[1:]),
				Value: -1,
			})
		} else {
			sortBy = append(sortBy, bson.E{
				Key:   sortOpt,
				Value: 1,
			})
		}
	}

	filterOptions.SetSort(sortBy)
	filterOptions.Skip = &skipVal
	filterOptions.Limit = &limitVal

	collection := s.Session.Collection(collectionName)
	cur, err := collection.Find(ctx, filter, filterOptions)
	if err != nil {
		return err
	}

	txn.End()

	if err := cur.Err(); err != nil {
		return err
	}

	if err := cur.All(ctx, result); err != nil {
		return err
	}

	return nil
}

// FindOne returns matching document
func (s *MongoStorageOfficial) FindOne(collectionName string, filter interface{}, result interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cnf.Cnf.Mongo.MongoDBTimeout)*time.Second)
	defer cancel()

	txn := new(newrelic.Transaction)
	if s.newRelicApp != nil {
		txn = s.newRelicApp.StartTransaction("mongo-findone")
		ctx = newrelic.NewContext(ctx, txn)
	}

	collection := s.Session.Collection(collectionName)
	err := collection.FindOne(ctx, filter).Decode(result)
	if err != nil {
		return err
	}

	txn.End()

	return nil
}

// FindWithSkip returns all matching documents with filter and skip params.
// page starts with 1.
func (s *MongoStorageOfficial) FindWithSkip(collectionName string, filter interface{}, result interface{}, skip, limit int64, sortBy string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cnf.Cnf.Mongo.MongoDBTimeout)*time.Second)
	defer cancel()

	txn := new(newrelic.Transaction)
	if s.newRelicApp != nil {
		txn = s.newRelicApp.StartTransaction("mongo-find")
		ctx = newrelic.NewContext(ctx, txn)
	}

	//filter options
	filterOptions := options.Find()

	if sortBy == "" {
		sortBy = "_id"
	}

	sortByD := bson.D{}
	for _, sortOpt := range strings.Split(sortBy, ",") {
		if string(sortOpt[0]) == "-" {
			sortByD = append(sortByD, bson.E{
				Key:   string(sortOpt[1:]),
				Value: -1,
			})
		} else {
			sortByD = append(sortByD, bson.E{
				Key:   sortOpt,
				Value: 1,
			})
		}
	}

	filterOptions.SetSort(sortByD)
	filterOptions.Skip = &skip
	filterOptions.Limit = &limit

	collection := s.Session.Collection(collectionName)
	cur, err := collection.Find(ctx, filter, filterOptions)
	if err != nil {
		return err
	}

	txn.End()

	if err := cur.Err(); err != nil {
		return err
	}

	if err := cur.All(ctx, result); err != nil {
		return err
	}

	return nil
}

// Insert add given object to store.
func (s *MongoStorageOfficial) Insert(collectionName string, object interface{}) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cnf.Cnf.Mongo.MongoDBTimeout)*time.Second)
	defer cancel()
	txn := new(newrelic.Transaction)
	if s.newRelicApp != nil {
		txn = s.newRelicApp.StartTransaction("mongo-insert")
		ctx = newrelic.NewContext(ctx, txn)
	}

	collection := s.Session.Collection(collectionName)
	res, err := collection.InsertOne(ctx, object)
	if err != nil {
		return primitive.ObjectID{}, err
	}

	txn.End()

	return res.InsertedID.(primitive.ObjectID), nil
}

// InsertMany inserts given list of object to store
func (s *MongoStorageOfficial) InsertMany(collectionName string,
	objects []interface{}) ([]primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cnf.Cnf.Mongo.MongoDBTimeout)*time.Second)
	defer cancel()

	txn := new(newrelic.Transaction)
	if s.newRelicApp != nil {
		txn = s.newRelicApp.StartTransaction("mongo-insertmany")
		ctx = newrelic.NewContext(ctx, txn)
	}

	collection := s.Session.Collection(collectionName)
	res, err := collection.InsertMany(ctx, objects)
	if err != nil {
		return nil, err
	}
	txn.End()

	ids := []primitive.ObjectID{}
	for _, insertedID := range res.InsertedIDs {
		ids = append(ids, insertedID.(primitive.ObjectID))
	}

	return ids, nil
}

// Update updates record with given object
func (s *MongoStorageOfficial) Update(collectionName string, filter interface{},
	change interface{}, opts ...*options.UpdateOptions) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cnf.Cnf.Mongo.MongoDBTimeout)*time.Second)
	defer cancel()

	txn := new(newrelic.Transaction)
	if s.newRelicApp != nil {
		txn = s.newRelicApp.StartTransaction("mongo-update")
		ctx = newrelic.NewContext(ctx, txn)
	}

	collection := s.Session.Collection(collectionName)
	_, err := collection.UpdateOne(ctx, filter, change)
	if err != nil {
		return err
	}

	txn.End()

	return nil
}

// UpdateMany updates record with given lis of object object
func (s *MongoStorageOfficial) UpdateMany(collectionName string, filter interface{}, change interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cnf.Cnf.Mongo.MongoDBTimeout)*time.Second)
	defer cancel()

	txn := new(newrelic.Transaction)
	if s.newRelicApp != nil {
		txn = s.newRelicApp.StartTransaction("mongo-updatemany")
		ctx = newrelic.NewContext(ctx, txn)
	}

	collection := s.Session.Collection(collectionName)
	_, err := collection.UpdateMany(ctx, filter, change)
	if err != nil {
		return err
	}

	txn.End()

	return nil
}

// Delete remove object with given id from store
func (s *MongoStorageOfficial) Delete(collectionName string, filter interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cnf.Cnf.Mongo.MongoDBTimeout)*time.Second)
	defer cancel()

	txn := new(newrelic.Transaction)
	if s.newRelicApp != nil {
		txn = s.newRelicApp.StartTransaction("mongo-delete")
		ctx = newrelic.NewContext(ctx, txn)
	}

	collection := s.Session.Collection(collectionName)
	_, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	txn.End()

	return nil
}

// DeleteMany remove object with given list of ids from store
func (s *MongoStorageOfficial) DeleteMany(collectionName string, filter interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cnf.Cnf.Mongo.MongoDBTimeout)*time.Second)
	defer cancel()

	txn := new(newrelic.Transaction)
	if s.newRelicApp != nil {
		txn = s.newRelicApp.StartTransaction("mongo-deletemany")
		ctx = newrelic.NewContext(ctx, txn)
	}

	collection := s.Session.Collection(collectionName)
	_, err := collection.DeleteMany(ctx, filter)
	if err != nil {
		return err
	}

	txn.End()

	return nil
}

// Count retrieves object count directly from dbms
func (s *MongoStorageOfficial) Count(collectionName string, filter interface{}) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cnf.Cnf.Mongo.MongoDBTimeout)*time.Second)
	defer cancel()

	txn := new(newrelic.Transaction)
	if s.newRelicApp != nil {
		txn = s.newRelicApp.StartTransaction("mongo-count")
		ctx = newrelic.NewContext(ctx, txn)
	}

	collection := s.Session.Collection(collectionName)
	docCount, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}

	txn.End()

	return int(docCount), nil
}

// Aggregate aggregate object(s) directly from dbms
func (s *MongoStorageOfficial) Aggregate(collectionName string, filter interface{}, result interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cnf.Cnf.Mongo.MongoDBTimeout)*time.Second)
	defer cancel()

	txn := new(newrelic.Transaction)
	if s.newRelicApp != nil {
		txn = s.newRelicApp.StartTransaction("mongo-aggregate")
		ctx = newrelic.NewContext(ctx, txn)
	}

	collection := s.Session.Collection(collectionName)
	cur, err := collection.Aggregate(ctx, filter)
	if err != nil {
		return err
	}

	txn.End()

	var results []bson.M
	if err := cur.All(ctx, &results); err != nil {
		return err
	}

	if err := cur.Err(); err != nil {
		return err
	}

	jsonString, _ := json.Marshal(results)
	err = json.Unmarshal(jsonString, result)
	if err != nil {
		return err
	}

	return nil
}

// CreateIndex creates index to colelction.
func (s *MongoStorageOfficial) CreateIndex(collectionName string, keys bson.M,
	opts *options.IndexOptions) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cnf.Cnf.Mongo.MongoDBTimeout)*time.Second)
	defer cancel()

	collection := s.Session.Collection(collectionName)

	// Declare an index model object to pass to CreateOne()
	// db.members.createIndex( { "SOME_FIELD": 1 }, { unique: true } )
	mod := mongo.IndexModel{
		Keys:    keys,
		Options: opts,
	}

	_, err := collection.Indexes().CreateOne(ctx, mod)
	if err != nil {
		return err
	}

	return nil
}

// Close connection
func (s *MongoStorageOfficial) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cnf.Cnf.Mongo.MongoDBTimeout)*time.Second)
	defer cancel()

	return s.Client.Disconnect(ctx)
}

// NewPaginationParams returns default pagination params
func (s *MongoStorageOfficial) NewPaginationParams() *PaginationParams {
	return &PaginationParams{
		SortBy: "_id",
		Page:   0,
		Limit:  50,
	}
}
