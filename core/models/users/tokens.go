package users

import (
	"context"
	"github.com/dchest/uniuri"
	"go.mongodb.org/mongo-driver/bson"
	"primrose/clients"
)

func (user *User) UpdateToken() (string, error) {
	token := uniuri.NewLen(256)
	err := clients.Db.Collection("users").
		FindOneAndUpdate(context.TODO(), bson.M{"email": user.Email}, bson.M{"$set": bson.M{"token": token}}).
		Err()
	return token, err
}

func (user *User) GetOrCreateToken() (string, error) {
	if user.Token == "" {
		return user.UpdateToken()
	}
	return user.Token, nil
}
