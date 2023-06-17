package mimi

import (
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type UpAction = func(ctx mongo.SessionContext, coll *mongo.Collection) error
type DownAction = func(ctx mongo.SessionContext, coll *mongo.Collection) error
type Migration struct {
	Key         string
	Collection  string
	MigrateUp   UpAction
	MigrateDown DownAction
}

type iMigration struct {
	Id        string    `bson:"_id" json:"key"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
}
