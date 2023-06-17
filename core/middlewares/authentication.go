package middlewares

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/kataras/iris/v12/context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"primrose/env"
	"primrose/models/responses"
	"primrose/models/users"
	"strings"
)

func Secured(c *context.Context) (*users.User, bool) {
	token := c.GetHeader("Authorization")
	if token == "" {
		responses.AuthenticationRequired.Reply(c)
		return nil, false
	}
	token, _ = strings.CutPrefix(token, "Bearer ")
	var user *users.User
	claims, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		subject, err := token.Claims.GetSubject()
		if err != nil {
			return nil, err
		}
		id, err := primitive.ObjectIDFromHex(subject)
		if err != nil {
			return nil, err
		}
		user, err = users.WithId(id)
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
			responses.AuthenticationRequired.Reply(c)
			return nil, false
		}
		responses.Handle(c, err)
		return nil, false
	}
	if !claims.Valid {
		responses.AuthenticationRequired.Reply(c)
		return nil, false
	}
	if user.Flags == nil {
		user.Flags = []string{}
	}
	_ = c.SetUser(user)
	return user, true
}
