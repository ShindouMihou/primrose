package responses

import (
	"github.com/kataras/golog"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"strings"
)

type ErrorResponse struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}

func (errorResponse ErrorResponse) Format(args ...string) ErrorResponse {
	text := errorResponse.Error
	for _, arg := range args {
		text = strings.Replace(text, "{$PLACEHOLDER}", arg, 1)
	}
	return ErrorResponse{Code: errorResponse.Code, Error: text}
}

func Handle(ctx *context.Context, err error) {
	VagueError.Reply(ctx)
	golog.Error("[HTTP] An error occurred while trying to execute ", ctx.Path(), ": ", err)
}

var InvalidPayload = ErrorResponse{Code: iris.StatusBadRequest, Error: "Invalid payload."}
var VagueError = ErrorResponse{Code: iris.StatusBadRequest, Error: "An error occurred while trying to execute this task."}
var EmailAlreadyUsed = ErrorResponse{Code: iris.StatusBadRequest, Error: "Another user is already using the provided e-mail address."}
var InvalidEmailPassword = ErrorResponse{Code: iris.StatusBadRequest, Error: "Invalid email or password."}
var AuthenticationRequired = ErrorResponse{Code: iris.StatusForbidden, Error: "You cannot perform this task, or access this resource."}
var Unauthorized = ErrorResponse{Code: iris.StatusUnauthorized, Error: "You do not have the privilege to perform this task or access this resource."}
var NotFound = ErrorResponse{Code: iris.StatusNotFound, Error: "We cannot find any resource that matches."}
