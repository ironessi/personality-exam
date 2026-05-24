package cmd

import (
	"PersonalityTest/internal/controller/auth"
	"PersonalityTest/internal/controller/category"
	"PersonalityTest/internal/controller/favorites"
	"PersonalityTest/internal/controller/initialize"
	"PersonalityTest/internal/controller/payment"
	"PersonalityTest/internal/controller/sms"
	"PersonalityTest/internal/controller/test_type"
	"PersonalityTest/internal/controller/user"
	"PersonalityTest/internal/controller/user_test"
	"PersonalityTest/internal/logic/middleware"
	"PersonalityTest/internal/logic/pay"
	"context"

	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcfg"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/gcron"
	"github.com/gogf/gf/v2/os/gfile"
	"os"
	"time"
)

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start http server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			// 显式添加配置文件搜索路径，防止找不到 config.yaml
			g.Cfg().GetAdapter().(*gcfg.AdapterFile).AddPath("manifest/config")

			s := g.Server()
			s.EnableAdmin() // 开启平滑重启支持
			s.SetGraceful(true)

			// 注册每日凌晨2点执行的对账定时任务
			_, err = gcron.AddSingleton(ctx, "0 0 2 * * *", pay.New().DailyReconciliation, "daily_reconciliation")
			if err != nil {
				g.Log().Warningf(ctx, "Failed to register daily reconciliation task: %v", err)
			}
			// 周期性处理待退款任务，避免请求内 goroutine 因重启丢失。
			_, err = gcron.AddSingleton(ctx, "*/10 * * * * *", func(ctx context.Context) {
				pay.NewRefunds().ProcessPendingRefunds(ctx, 10)
			}, "pending_refunds")
			if err != nil {
				g.Log().Warningf(ctx, "Failed to register refund worker task: %v", err)
			}

			// 全局中间件
			s.Use(
				ghttp.MiddlewareHandlerResponse,
			)

			s.Group("/", func(group *ghttp.RouterGroup) {
				// 用户认证路由
				group.Group("/api/v1", func(group *ghttp.RouterGroup) {
					group.Bind(
						auth.NewV1(),
					)
					// 需要登录的接口
					group.Group("/", func(group *ghttp.RouterGroup) {
						group.Middleware(middleware.Auth)
						group.Bind(
							user.NewV1(),
						)
					})
				})

				// 初始化路由
				group.Group("/api", func(group *ghttp.RouterGroup) {
					group.Bind(
						initialize.New(),
					)
				})

				// 短信API (严格匹配文档 /api/v1/sms)
				group.Group("/api/v1/sms", func(group *ghttp.RouterGroup) {
					smsCtrl := sms.NewV1()
					group.Bind(smsCtrl)
					// 添加 Status/:mobile 支持
					// Bind 默认生成 /status，此处手动添加带参数的路由映射到同一个 Status 方法
					group.GET("/status/:mobile", smsCtrl.Status)
				})

				// 用户测试相关接口
				group.Group("/api/user/tests", func(group *ghttp.RouterGroup) {
					group.Middleware(middleware.Auth)
					group.Bind(
						user_test.New(),
					)
				})

				group.Group("/api/user/favorites", func(group *ghttp.RouterGroup) {
					group.Middleware(middleware.Auth)
					group.Bind(
						favorites.New(),
					)
				})

				// 测试管理相关接口
				group.Group("/api/tests", func(group *ghttp.RouterGroup) {
					group.Bind(
						test_type.NewV1(),
					)
				})

				// 分类相关接口
				group.Group("/api/categories", func(group *ghttp.RouterGroup) {
					group.Bind(
						category.NewV1(),
					)
				})

				// 支付相关接口
				group.Group("/api/v1", func(group *ghttp.RouterGroup) {
					// 支付宝/微信回调不走 JWT（第三方服务器直接调用）
					group.POST("/payment/notify/alipay", payment.New().AlipayNotify)
					group.POST("/payment/notify/wechat", payment.New().WechatNotify)
					// 需要登录的支付接口
					group.Group("/", func(group *ghttp.RouterGroup) {
						group.Middleware(middleware.Auth)
						payment.RegisterRoutes(group)
					})
				})

				group.GET("/api/health", func(r *ghttp.Request) {
					r.Response.Write("OK")
				})
			})

			// 自动监控二进制文件更新并触发平滑重启 (针对无法操作终端的部署环境)
			go func() {
				selfPath := gfile.SelfPath()
				stat, err := os.Stat(selfPath)
				if err != nil {
					return
				}
				lastModTime := stat.ModTime()
				for {
					time.Sleep(5 * time.Second)
					newStat, err := os.Stat(selfPath)
					if err != nil {
						continue
					}
					if newStat.ModTime().After(lastModTime) {
						g.Log().Infof(ctx, "检测到可执行文件已更新 [%s]，正在触发平滑重启...", selfPath)
						_ = ghttp.RestartAllServer(ctx, "")
						return
					}
				}
			}()

			s.Run()
			return nil
		},
	}
)
