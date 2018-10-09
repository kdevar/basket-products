package errors

import (
	"errors"
	"github.com/gin-gonic/gin"
)

type WithError func(context *gin.Context) *ApiError

func (f WithError) Handle(context *gin.Context) {
	err := f(context)
	if err != nil {
		context.AbortWithStatusJSON(err.HttpCode, err)
	}
}

type ApiError struct {
	Err      error
	Message  string
	HttpCode int
}

func (a *ApiError) Error() string {
	return a.Message
}

func NotFoundError(message string) *ApiError {
	return &ApiError{
		errors.New(message),
		message,
		404,
	}
}

func ServerError(err error) *ApiError {

	return &ApiError{
		err,
		err.Error(),
		500,
	}

}
