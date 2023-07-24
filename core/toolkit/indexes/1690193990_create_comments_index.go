package indexes

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"primrose/toolkit/mimi"
)

var CreateCommentsIndexMigration = mimi.AddIndex(mimi.Index{
	Key:        "create_comments_index",
	Collection: "comments",
	Action: func(coll *mongo.Collection) error {
		if _, err := coll.Indexes().CreateOne(context.TODO(), mongo.IndexModel{
			Keys: bson.M{"post": -1},
		}); err != nil {
			return err
		}
		return nil
	},
})
