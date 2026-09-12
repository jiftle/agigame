package errcode

import (
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"agigame/console/backend/internal/consts"
)

// New 创建带业务码的错误
func New(code int, msg string) error {
	return gerror.NewCode(gcode.New(code, msg, nil), msg)
}

// BadRequest 400
func BadRequest(msg string) error {
	return New(consts.CodeBadRequest, msg)
}

// Unauthorized 401
func Unauthorized(msg string) error {
	return New(consts.CodeUnauthorized, msg)
}

// TokenExpired 402
func TokenExpired(msg string) error {
	return New(consts.CodeTokenExpired, msg)
}

// Forbidden 403
func Forbidden(msg string) error {
	return New(consts.CodeForbidden, msg)
}

// NotFound 404
func NotFound(msg string) error {
	return New(consts.CodeNotFound, msg)
}

// Business 通用业务错误
func Business(msg string) error {
	return New(consts.CodeError, msg)
}
