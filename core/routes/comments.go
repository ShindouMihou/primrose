package routes

import (
	"github.com/kataras/iris/v12"
	inc "github.com/kataras/iris/v12/context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"primrose/clients"
	"primrose/middlewares"
	"primrose/models/comments"
	"primrose/models/posts"
	"primrose/models/requests"
	"primrose/models/responses"
	"time"
)

var _ = clients.Iris.Attach(func(app *iris.Application) {
	comments := app.Party("/comments/")
	{
		post := comments.Party("/{post}")
		{
			post.Put("/", CommentRoutes.Save)
			post.Get("/", CommentRoutes.List)
		}
		comments.Patch("/edit/{comment}", CommentRoutes.Edit)
		comments.Delete("/delete/{comment}", CommentRoutes.Del)
		comments.Get("/view/{comment}", CommentRoutes.View)
	}
})

func getCommentId(ctx *inc.Context) (*primitive.ObjectID, bool) {
	cid := ctx.Params().Get("comment")
	if cid == "" {
		responses.InvalidPayload.Reply(ctx)
		return nil, false
	}
	commentId, err := primitive.ObjectIDFromHex(cid)
	if err != nil {
		responses.NotFound.Reply(ctx)
		return nil, false
	}
	return &commentId, true
}

func getJoinedComment(ctx *inc.Context) (*comments.JoinedComment, bool) {
	commentId, ok := getCommentId(ctx)
	if !ok {
		return nil, false
	}
	comment, err := comments.Joined(*commentId)
	if err != nil {
		responses.Handle(ctx, err)
		return nil, false
	}
	if comment == nil {
		responses.NotFound.Reply(ctx)
		return nil, false
	}
	return comment, true
}

func getUnjoinedComment(ctx *inc.Context) (*comments.UnjoinedComment, bool) {
	commentId, ok := getCommentId(ctx)
	if !ok {
		return nil, false
	}
	comment, err := comments.WithId(*commentId)
	if err != nil {
		responses.Handle(ctx, err)
		return nil, false
	}
	if comment == nil {
		responses.NotFound.Reply(ctx)
		return nil, false
	}
	return comment, true
}

var CommentRoutes = struct {
	Save inc.Handler
	Edit inc.Handler
	View inc.Handler
	List inc.Handler
	Del  inc.Handler
}{
	Save: func(ctx *inc.Context) {
		user, ok := middlewares.Secured(ctx)
		if !ok {
			return
		}
		postId := ctx.Params().Get("post")
		if postId == "" {
			responses.InvalidPayload.Reply(ctx)
			return
		}
		post, err := posts.WithId(postId)
		if err != nil {
			responses.Handle(ctx, err)
			return
		}
		if post == nil {
			responses.NotFound.Reply(ctx)
			return
		}
		rid := ctx.URLParam("reply_to")
		var replyTo *primitive.ObjectID
		if rid != "" {
			rid, err := primitive.ObjectIDFromHex(rid)
			if err != nil {
				responses.InvalidReplyToComment.Reply(ctx)
				return
			}
			replyTo = &rid
		}
		ok, comment := requests.Read(ctx, comments.Comment{})
		if comment.Content == "" {
			responses.InvalidPayload.Reply(ctx)
			return
		}
		if len(comment.Content) > comments.MaximumContentLength {
			responses.ExceededMaximumCommentLength.Reply(ctx)
			return
		}
		comment.UpdatedAt = time.Now()
		comment.CreatedAt = time.Now()
		comment.Post = post.Id

		unjoined := comments.UnjoinedComment{Comment: *comment, ReplyTo: replyTo, Author: user.ObjectId()}
		if err := unjoined.Save(); err != nil {
			if err == comments.CannotReplyToUnknownCommentErr {
				responses.InvalidReplyToComment.Reply(ctx)
				return
			}
			responses.Handle(ctx, err)
			return
		}
		responses.Reply(ctx, unjoined)
	},
	Edit: func(ctx *inc.Context) {
		user, ok := middlewares.Secured(ctx)
		if !ok {
			return
		}
		comment, ok := getUnjoinedComment(ctx)
		if !ok {
			return
		}
		if comment.Author != user.ObjectId() {
			responses.Unauthorized.Reply(ctx)
			return
		}
		type EditCommentRequest struct {
			Content string `json:"content" validate:"required"`
		}
		ok, request := requests.Read(ctx, EditCommentRequest{})
		if !ok {
			return
		}
		if request.Content == "" {
			responses.InvalidPayload.Reply(ctx)
			return
		}
		if len(request.Content) > comments.MaximumContentLength {
			responses.ExceededMaximumCommentLength.Reply(ctx)
			return
		}
		comment.Content = request.Content
		if err := comment.Save(); err != nil {
			if err == comments.CannotReplyToUnknownCommentErr {
				responses.InvalidReplyToComment.Reply(ctx)
				return
			}
			responses.Handle(ctx, err)
			return
		}
		responses.Reply(ctx, comment)
	},
	View: func(ctx *inc.Context) {
		if comment, ok := getJoinedComment(ctx); ok {
			responses.Reply(ctx, comment)
		}
	},
	Del: func(ctx *inc.Context) {
		user, ok := middlewares.Secured(ctx)
		if !ok {
			return
		}
		comment, ok := getUnjoinedComment(ctx)
		if !ok {
			return
		}
		if comment.Author != user.ObjectId() {
			responses.Unauthorized.Reply(ctx)
			return
		}
		err := comment.Delete()
		if err != nil {
			responses.Handle(ctx, err)
			return
		}
		responses.Reply(ctx, responses.Acknowledged)
	},
	List: func(ctx *inc.Context) {
		postId := ctx.Params().Get("post")
		if postId == "" {
			responses.InvalidPayload.Reply(ctx)
			return
		}
		var after *primitive.ObjectID
		if ctx.URLParamExists("after") {
			id, err := primitive.ObjectIDFromHex(ctx.URLParam("after"))
			if err != nil {
				responses.NotFound.Reply(ctx)
				return
			}
			after = &id
		}
		limit := ctx.URLParamInt32Default("limit", 100)
		if limit > 100 {
			limit = 100
		}
		comments, err := comments.List(uint8(limit), postId, after)
		if err != nil {
			responses.Handle(ctx, err)
			return
		}
		responses.Reply(ctx, responses.Arrayed{Data: comments})
	},
}
