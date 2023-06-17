package clients

import (
	"context"
	"github.com/kataras/golog"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"primrose/env"
)

var MongoClient *mongo.Client
var Db *mongo.Database

type IMongo struct {
	Client
	Logger *golog.Logger
}

var M = IMongo{Logger: golog.New().SetPrefix("[DB] ")}

func (m *IMongo) Init() {
	uri := env.EnsureEnv("MONGO_URI")
	client, err := mongo.Connect(context.TODO(), options.Client().
		SetAppName("Primrose").
		SetReadPreference(readpref.Primary()).
		SetMonitor(&event.CommandMonitor{
			Started: func(ctx context.Context, startedEvent *event.CommandStartedEvent) {
				golog.Debug(startedEvent.Command)
			},
			Failed: func(ctx context.Context, failedEvent *event.CommandFailedEvent) {
				golog.Error(failedEvent.Failure)
			},
		}).
		ApplyURI(uri))
	if err != nil {
		m.Logger.Fatal("Failed to connect to Mongo: ", err)
	}
	err = client.Ping(context.TODO(), readpref.Primary())
	if err != nil {
		m.Logger.Fatal("Failed to ping to cluster: ", err)
	}
	MongoClient = client
	Db = client.Database("primrose")
}
