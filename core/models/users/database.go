package users

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"primrose/clients"
	"primrose/utils"
)

func WithEmail(email string) (*User, error) {
	return Get("email", email)
}

func WithId(id primitive.ObjectID) (*User, error) {
	return Get("_id", id)
}

func Get(key string, value any) (*User, error) {
	return utils.ReturningOne(func(t *User) error {
		return clients.Db.Collection("users").
			FindOne(context.TODO(), bson.M{key: value}).
			Decode(t)
	})
}
