package cmd

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"

	"agigame/console/backend/internal/boot"
	"agigame/console/backend/internal/controller/auth"
	"agigame/console/backend/internal/controller/emu"
	"agigame/console/backend/internal/controller/system"
	_ "agigame/console/backend/internal/logic"
	"agigame/console/backend/internal/middleware"
)

// Main 服务启动入口
var Main = gcmd.Command{
	Name:  "main",
	Usage: "main",
	Brief: "GoBoy 控制台后端服务",
	Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
		if err = boot.InitDatabase(ctx); err != nil {
			return err
		}
		s := g.Server()
		s.Group("/api/v1", func(group *ghttp.RouterGroup) {
			group.Middleware(ghttp.MiddlewareCORS, middleware.Response, middleware.ErrorHandler)

			// 公开接口
			group.Bind(auth.NewLogin())

			// 需要登录的接口
			group.Group("/", func(protected *ghttp.RouterGroup) {
				protected.Middleware(middleware.Auth, middleware.OperLog)
				protected.Bind(
					auth.NewProfile(),
					system.NewUser(),
					system.NewRole(),
					system.NewMenu(),
					system.NewDept(),
					system.NewDict(),
					system.NewConfig(),
					system.NewLog(),
					emu.NewSession(),
				)
			})
		})

		// 模拟器实时流（WebSocket）：独立分组，仅挂 Auth，
		// 避开统一响应 Response 中间件对连接升级的干扰。
		s.Group("/ws", func(ws *ghttp.RouterGroup) {
			ws.Middleware(ghttp.MiddlewareCORS, middleware.Auth)
			ws.Bind(emu.NewStream())
		})

		s.Run()
		return nil
	},
}
