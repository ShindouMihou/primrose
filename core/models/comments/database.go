package comments

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"primrose/utils"
)

var (
	lookupPipeline = bson.D{{"$lookup", bson.D{
		{"from", "users"},
		{"localField", "author"},
		{"foreignField", "_id"},
		{"as", "author"},
	}}}
	flattenPipeline = bson.D{{"$set", bson.D{
		{"author", bson.D{
			{"$first", "$author"},
		}},
	}}}
)

func List(limit uint8, post string, after *primitive.ObjectID) ([]JoinedComment, error) {
	var comments = make([]JoinedComment, limit)
	var query bson.M

	var apply = func(query bson.M) error {
		if after != nil {
			post, err := WithId(*after)
			if err != nil {
				return err
			}
			if post != nil {
				query["created_at"] = bson.M{"$gt": post}
			}
		}
		return nil
	}
	query = bson.M{"post": post}

	if err := apply(query); err != nil {
		return nil, err
	}

	res, err := GetCollection().Aggregate(context.TODO(), mongo.Pipeline{
		{{"$match", query}},
		{{"$sort", bson.D{
			{"created_at", 1}},
		}},
		{{"$limit", limit}},
		lookupPipeline,
		flattenPipeline,
	})
	if err != nil {
		return nil, err
	}
	if err = res.All(context.TODO(), &comments); err != nil {
		return nil, err
	}
	return comments, nil
}

func With(key string, value any) (*UnjoinedComment, error) {
	return utils.ReturningOne(func(t *UnjoinedComment) error {
		return GetCollection().FindOne(context.TODO(), bson.M{key: value}).Decode(t)
	})
}

func WithId(id primitive.ObjectID) (*UnjoinedComment, error) {
	return With("_id", id)
}

func Joined(id primitive.ObjectID) (*JoinedComment, error) {
	return utils.ReturningOne(func(t *JoinedComment) error {
		cursor, err := GetCollection().Aggregate(context.TODO(), mongo.Pipeline{
			{{"$match", bson.D{{"_id", id}}}},
			{{"$limit", 1}},
			lookupPipeline,
			flattenPipeline,
		})
		if err != nil {
			return err
		}
		for cursor.Next(context.TODO()) {
			if err := cursor.Decode(t); err != nil {
				return err
			}
			return nil
		}
		if err := cursor.Err(); err != nil {
			return err
		}
		return nil
	})
}
