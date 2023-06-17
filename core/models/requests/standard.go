package requests

import (
	"github.com/kataras/iris/v12/context"
	"primrose/models/responses"
)

func Read[T any](ctx *context.Context, payload T) (bool, *T) {
	if err := ctx.ReadJSON(&payload); err != nil {
		responses.InvalidPayload.Reply(ctx)
		return false, nil
	}
	return true, &payload
}
