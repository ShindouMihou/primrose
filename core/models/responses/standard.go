package responses

import (
	"github.com/kataras/golog"
	"github.com/kataras/iris/v12/context"
)

func Reply(ctx *context.Context, response any) {
	if err := ctx.JSON(response); err != nil {
		golog.Error("[HTTP] An error occurred while trying to respond to ", ctx.Path(), ": ", err)
	}
}

func (err *ErrorResponse) Reply(ctx *context.Context) {
	ctx.StatusCode(err.Code)
	Reply(ctx, err)
}
