package mimi

import "go.mongodb.org/mongo-driver/mongo"

type IndexAction = func(coll *mongo.Collection) error
type Index struct {
	Key        string
	Collection string
	Action     IndexAction
}
