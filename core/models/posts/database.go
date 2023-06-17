package posts

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"primrose/utils"
)

func List(limit uint8, published bool, search *string, after *string) ([]Post, error) {
	var posts = make([]Post, limit)
	var query bson.M

	var apply = func(query bson.M) error {
		if after != nil {
			post, err := WithId(*after)
			if err != nil {
				return err
			}
			if post != nil {
				query["created_at"] = bson.M{"$gt": post.CreatedAt}
			}
		}
		if published {
			query["published"] = published
		}
		return nil
	}

	opts := options.Find().SetLimit(int64(limit))
	if search != nil {
		query = bson.M{"$text": bson.M{"$search": *search, "$caseSensitive": false}}
	} else {
		query = bson.M{}
	}

	if err := apply(query); err != nil {
		return nil, err
	}

	res, err := GetCollection().Find(context.TODO(), query, opts.SetSort(bson.M{"created_at": -1}))
	if err != nil {
		return nil, err
	}
	if err = res.All(context.TODO(), &posts); err != nil {
		return nil, err
	}
	return posts, nil
}

func With(key string, value string) (*Post, error) {
	return utils.ReturningOne(func(t *Post) error {
		return GetCollection().FindOne(context.TODO(), bson.M{key: value}).Decode(t)
	})
}

func WithSlug(slug string) (*Post, error) {
	return With("slug", slug)
}

func WithId(id string) (*Post, error) {
	return With("_id", id)
}
