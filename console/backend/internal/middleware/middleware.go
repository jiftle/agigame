package middleware

import (
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"

	"agigame/console/backend/internal/consts"
	"agigame/console/backend/internal/model"
)

// Response 统一响应中间件
func Response(r *ghttp.Request) {
	r.Middleware.Next()

	// 已有自定义输出（如文件下载）则直接返回
	if r.Response.BufferLength() > 0 {
		return
	}

	var (
		err = r.GetError()
		res = r.GetHandlerResponse()
	)
	if err != nil {
		code := consts.CodeError
		if c := gerror.Code(err); c != nil && c != gcode.CodeNil {
			code = c.Code()
		}
		r.Response.WriteJson(model.Response{
			Code:    code,
			Message: err.Error(),
		})
		return
	}

	r.Response.WriteJson(model.Response{
		Code:    consts.CodeSuccess,
		Message: "success",
		Data:    res,
	})
}

// ErrorHandler 异常恢复中间件
func ErrorHandler(r *ghttp.Request) {
	defer func() {
		if e := recover(); e != nil {
			g.Log().Errorf(r.Context(), "panic recovered: %v", e)
			r.Response.ClearBuffer()
			r.SetError(gerror.NewCodef(
				gcode.New(consts.CodeError, "系统内部错误", nil), "%v", e,
			))
		}
	}()
	r.Middleware.Next()
}
