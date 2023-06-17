package routes

import (
	"github.com/kataras/iris/v12"
	inc "github.com/kataras/iris/v12/context"
	"primrose/clients"
	"primrose/middlewares"
	"primrose/models/posts"
	"primrose/models/requests"
	"primrose/models/responses"
	"primrose/models/users"
	"primrose/utils"
	"time"
)

var _ = clients.Iris.Attach(func(app *iris.Application) {
	posts := app.Party("/posts")
	{
		posts.Put("/save", PostRoutes.Save)
		posts.Delete("/delete", PostRoutes.Del)
		posts.Get("/view", PostRoutes.View)
		posts.Get("/list", PostRoutes.List)
	}
})

var PostRoutes = struct {
	Save inc.Handler
	View inc.Handler
	List inc.Handler
	Del  inc.Handler
}{
	Save: func(ctx *inc.Context) {
		if user, ok := middlewares.Secured(ctx); ok {
			if !user.IsAdmin() {
				responses.Unauthorized.Reply(ctx)
				return
			}
			ok, request := requests.Read(ctx, posts.Post{})
			if !ok {
				return
			}
			var post *posts.Post
			post = request
			var reset = func() {
				post.Id = ""
				post.CreatedAt = time.Now()
			}
			if request.Id != "" {
				t, err := posts.WithId(request.Id)
				if err != nil {
					responses.Handle(ctx, err)
					return
				}
				if t == nil {
					reset()
				} else {
					post = t
					request.Transfer(post)
				}
			} else {
				reset()
			}
			if err := post.Save(); err != nil {
				responses.Handle(ctx, err)
				return
			}
			responses.Reply(ctx, post)
		}
	},
	View: func(ctx *inc.Context) {
		if !ctx.URLParamExists("key") {
			responses.InvalidPayload.Reply(ctx)
			return
		}
		key := ctx.URLParam("key")

		isSlug := ctx.URLParamBoolDefault("isSlug", false)
		findClosest := ctx.URLParamBoolDefault("findClosest", false)

		user, err := users.From(ctx)
		if err != nil {
			responses.Handle(ctx, err)
			return
		}

		var flush = func(post *posts.Post) {
			if post == nil {
				responses.NotFound.Reply(ctx)
				return
			}
			if !post.Published && (user == nil || !user.IsAdmin()) {
				responses.Unauthorized.Reply(ctx)
				return
			}
			responses.Reply(ctx, post)
		}

		if isSlug {
			post, err := posts.WithSlug(key)
			if err != nil {
				responses.Handle(ctx, err)
				return
			}
			if post == nil {
				if findClosest {
					published := ctx.URLParamBoolDefault("published", true)
					if !published && (user == nil || !user.IsAdmin()) {
						published = false
					}
					psts, err := posts.List(1, published, &key, nil)
					if err != nil {
						responses.Handle(ctx, err)
						return
					}
					if len(psts) == 0 {
						responses.NotFound.Reply(ctx)
						return
					}
					responses.Reply(ctx, psts[0])
					return
				}
				responses.NotFound.Reply(ctx)
				return
			}
			flush(post)
			return
		}
		post, err := posts.WithId(key)
		if err != nil {
			responses.Handle(ctx, err)
			return
		}
		flush(post)
	},
	List: func(ctx *inc.Context) {
		var search *string
		if ctx.URLParamExists("search") {
			search = utils.Ptr(ctx.URLParam("search"))
		}

		var after *string
		if ctx.URLParamExists("after") {
			after = utils.Ptr(ctx.URLParam("after"))
		}

		limit := ctx.URLParamInt32Default("limit", 100)
		if limit > 100 {
			limit = 100
		}

		published := ctx.URLParamBoolDefault("published", true)
		if !published {
			user, err := users.From(ctx)
			if err != nil {
				responses.Handle(ctx, err)
				return
			}
			if user == nil || !user.IsAdmin() {
				responses.Unauthorized.Reply(ctx)
				return
			}
		}

		posts, err := posts.List(uint8(limit), published, search, after)
		if err != nil {
			responses.Handle(ctx, err)
			return
		}
		responses.Reply(ctx, responses.Arrayed{Data: posts})
	},
	Del: func(ctx *inc.Context) {
		if user, ok := middlewares.Secured(ctx); ok {
			if !user.IsAdmin() {
				responses.Unauthorized.Reply(ctx)
				return
			}
			if !ctx.URLParamExists("key") {
				responses.InvalidPayload.Reply(ctx)
				return
			}
			key := ctx.URLParam("key")
			post, err := posts.WithId(key)
			if err != nil {
				responses.Handle(ctx, err)
				return
			}
			if post == nil {
				responses.NotFound.Reply(ctx)
				return
			}
			err = post.Delete()
			if err != nil {
				responses.Handle(ctx, err)
				return
			}
			responses.Reply(ctx, responses.Acknowledged)
		}
	},
}
