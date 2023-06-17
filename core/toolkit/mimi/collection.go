package mimi

import "go.mongodb.org/mongo-driver/mongo"

type CollectionAction = func(database *mongo.Database) error
type Collection struct {
	Key    string
	Action CollectionAction
}
