package users

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/kataras/iris/v12/context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"primrose/env"
	"strings"
)

func From(c *context.Context) (*User, error) {
	token := c.GetHeader("Authorization")
	if token == "" {
		return nil, nil
	}
	token, _ = strings.CutPrefix(token, "Bearer ")
	var user *User
	claims, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		subject, err := token.Claims.GetSubject()
		if err != nil {
			return nil, err
		}
		id, err := primitive.ObjectIDFromHex(subject)
		if err != nil {
			return nil, err
		}
		user, err = WithId(id)
		if err != nil {
			return nil, err
		}
		key, err := user.GetOrCreateToken()
		if err != nil {
			return nil, err
		}
		return []byte(key + "/?=" + env.EnsureEnv("SIGNING_KEY")), nil
	})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		if errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet) || errors.Is(err, jwt.ErrTokenMalformed) || errors.Is(err, jwt.ErrTokenSignatureInvalid) {
			return nil, nil
		}
		return nil, err
	}
	if !claims.Valid {
		return nil, nil
	}
	if user.Flags == nil {
		user.Flags = []string{}
	}
	_ = c.SetUser(user)
	return user, nil
}
