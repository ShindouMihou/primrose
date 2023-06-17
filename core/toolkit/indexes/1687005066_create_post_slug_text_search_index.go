package indexes

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"primrose/toolkit/mimi"
)

var CreatePostSlugTextSearchIndexMigration = mimi.AddIndex(mimi.Index{
	Key:        "create_post_slug_text_search_index",
	Collection: "posts",
	Action: func(coll *mongo.Collection) error {
		if _, err := coll.Indexes().CreateOne(context.TODO(), mongo.IndexModel{
			Keys:    bson.M{"slug": "text"},
			Options: options.Index().SetUnique(true),
		}); err != nil {
			return err
		}
		return nil
	},
})
