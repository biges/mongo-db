package cnf

import (
	"log"

	"github.com/codingconcepts/env"
)

var (
	// Cnf is a configuration struct.
	Cnf Config
)

// Mongo is configuration about mongodb.
type Mongo struct {
	MongoDBURI       string `env:"MONGO_DB_URI" required:"true"`
	MongoDBDebugMode bool   `env:"MONGO_DB_DEBUG_MODE" default:"false"`
	MongoDBName      string `env:"MONGO_DB_NAME" default:"cyclops"`
	MongoDBTimeout   int64  `env:"MONGO_DB_TIMEOUT" default:"30"`
}

// environment is a type for GO_ENV.
type environment string

const (
	// Production is live system.
	Production environment = "production"
	// Staging is test system.
	Staging environment = "staging"
	// Development is your environment.
	Development environment = "development"
)

// Config represents configuration parameters.
type Config struct {
	Env                environment `env:"GO_ENV" default:"development"`
	ReleaseVersionInfo string      `env:"VERSION_INFO" default:"unknown"`
	Mongo              Mongo
}

func init() {
	if err := env.Set(&Cnf); err != nil {
		log.Fatal(err)
	}

	// package does not support sub structs. you should set one by one.
	if err := env.Set(&Cnf.Mongo); err != nil {
		log.Fatal(err)
	}
}
