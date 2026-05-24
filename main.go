package main

import (
	"PersonalityTest/internal/cmd"
	_ "PersonalityTest/internal/logic"
	"PersonalityTest/internal/logic/initialData"
	"errors"
	"flag"
	"fmt"
	"os"
	"time"

	_ "github.com/gogf/gf/contrib/nosql/redis/v2" // Redis 驱动

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcfg"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gtime"
)

var (
	initFlag  = flag.Bool("init", false, "初始化测试数据")
	forceFlag = flag.Bool("force", false, "强制重新初始化数据")
)

func main() {
	flag.Parse()

	var err error

	// 全局设置i18n
	g.I18n().SetLanguage("zh-CN")

	// 设置全局CORS中间件
	setupGlobalCORS()

	// 添加调试信息
	fmt.Println("正在检查数据库连接...")

	// 检查数据库是否能连接
	err = connDb()
	if err != nil {
		fmt.Printf("数据库连接失败: %v\n", err)
		// 查看具体错误
		ctx := gctx.GetInitCtx()
		if dbErr := g.DB().GetCache().Close(ctx); dbErr != nil {
			fmt.Printf("数据库错误详情: %v\n", dbErr)
		}
		panic(err)
	}

	fmt.Println("数据库连接成功!")

	// 处理初始化参数
	if *initFlag {
		if *forceFlag {
			initialData.ForceImportData()
		} else {
			initialData.ImportData()
		}
		g.Log().Info(gctx.GetInitCtx(), "初始化完成，程序退出")
		return
	}

	// 显式添加配置文件搜索路径，防止找不到 config.yaml
	g.Cfg().GetAdapter().(*gcfg.AdapterFile).AddPath("manifest/config")

	// 正常启动服务
	fmt.Println("正在启动服务...")
	cmd.Main.Run(gctx.GetInitCtx())
}

func init() {
	// 终极修复：强制设置系统时区环境变量，确保底层驱动识别北京时间
	os.Setenv("TZ", "Asia/Shanghai")
	time.Local = time.FixedZone("CST", 8*3600)
	gtime.SetTimeZone("Asia/Shanghai")

	// 显式添加配置文件搜索路径，防止找不到 config.yaml
	g.Cfg().GetAdapter().(*gcfg.AdapterFile).AddPath("manifest/config")
}

// setupGlobalCORS 设置全局CORS中间件
func setupGlobalCORS() {
	s := g.Server()

	// 确保使用配置文件中的CORS设置
	corsConfig := g.Cfg().MustGet(gctx.GetInitCtx(), "server.cors").Map()
	if corsConfig != nil {
		// 动态设置CORS
		s.BindMiddlewareDefault(func(r *ghttp.Request) {
			r.Response.CORSDefault()
			r.Response.Header().Set("Access-Control-Allow-Origin", "*")
			r.Response.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
			r.Response.Header().Set("Access-Control-Allow-Headers", "Origin,Content-Type,Accept,Authorization,X-Requested-With,X-CSRF-Token")
			r.Middleware.Next()
		})
	} else {
		// 回退到简单CORS
		s.Use(ghttp.MiddlewareCORS)
	}

	fmt.Println("✅ CORS中间件已配置")
}

// connDb 检查数据库连接是否正常
func connDb() error {
	err := g.DB().PingMaster()
	if err != nil {
		return errors.New("连接到数据库失败: " + err.Error())
	}
	return nil
}
