package routes

import (
	"context"
	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	"github.com/kataras/iris/v12"
	inc "github.com/kataras/iris/v12/context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"primrose/clients"
	"primrose/env"
	"primrose/middlewares"
	"primrose/models/requests"
	"primrose/models/responses"
	"primrose/models/users"
	"time"
)

var _ = clients.Iris.Attach(func(app *iris.Application) {
	users := app.Party("/users")
	{
		users.Put("/", UserRoutes.Create)
		users.Post("/", UserRoutes.Login)
		me := users.Party("/@me")
		{
			me.Get("/", UserRoutes.Self)
		}
	}
})

type iUserCreate struct {
	Username string `json:"username" validate:"required,gte=0,lte=120"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type iUserLogin struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

var UserRoutes = struct {
	Create inc.Handler
	Login  inc.Handler
	Self   inc.Handler
}{
	Create: func(c *inc.Context) {
		ok, req := requests.Read(c, iUserCreate{})
		if !ok {
			return
		}
		hash, err := argon2id.CreateHash(req.Password, argon2id.DefaultParams)
		if err != nil {
			responses.Handle(c, err)
			return
		}
		user := users.User{
			Name:     req.Username,
			Email:    req.Email,
			Password: hash,
			Flags:    []string{},
		}
		res, err := clients.Db.
			Collection("users").
			InsertOne(context.TODO(), user)
		if err != nil {
			if mongo.IsDuplicateKeyError(err) {
				responses.EmailAlreadyUsed.Reply(c)
				return
			}
			responses.Handle(c, err)
			return
		}
		user.Id = res.InsertedID.(primitive.ObjectID).Hex()
		responses.Reply(c, user)
	},
	Login: func(c *inc.Context) {
		ok, req := requests.Read(c, iUserLogin{})
		if !ok {
			return
		}
		user, err := users.WithEmail(req.Email)
		if user == nil || err != nil {
			responses.InvalidEmailPassword.Reply(c)
			return
		}
		ok, err = argon2id.ComparePasswordAndHash(req.Password, user.Password)
		if err != nil {
			responses.VagueError.Reply(c)
			return
		}
		if !ok {
			responses.InvalidEmailPassword.Reply(c)
			return
		}
		key, err := user.GetOrCreateToken()
		if err != nil {
			responses.Handle(c, err)
			return
		}
		key += "/?=" + env.EnsureEnv("SIGNING_KEY")
		structuredToken := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.RegisteredClaims{
			Issuer:    "primrose",
			Subject:   user.Id,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * (24 * time.Hour))),
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        user.Id,
		})
		token, err := structuredToken.SignedString([]byte(key))
		if err != nil {
			responses.Handle(c, err)
			return
		}
		// TODO: Use cookie eventually since more secure.
		type TokenResponse struct {
			Token string `json:"token"`
		}
		responses.Reply(c, TokenResponse{Token: token})
	},
	Self: func(c *inc.Context) {
		if user, ok := middlewares.Secured(c); ok {
			responses.Reply(c, user)
		}
	},
}
