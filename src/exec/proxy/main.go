package main

import (
	"fmt"
	"runtime"

	"github.com/astaxie/beego"

	"upgrade-api/src/exec/proxy/controllers"
	"upgrade-api/src/exec/proxy/controllers/conf"
	_ "upgrade-api/src/exec/proxy/routers"
)

/* 初始化 */
func _init() *controllers.ProxyContext {
	runtime.GOMAXPROCS(runtime.NumCPU())

	/* > 加载配置 */
	c, err := conf.Load()
	if nil != err {
		fmt.Printf("Load configuration failed! errmsg:%s\n", err.Error())
		return nil
	}

	/* > 初始化环境 */
	ctx, err := controllers.Init(c)
	if nil != err {
		fmt.Printf("Initialize context failed! errmsg:%s\n", err.Error())
		return nil
	}

	return ctx
}

/* 注册回调 */
func register(ctx *controllers.ProxyContext) {
	controllers.ProxyCntx = ctx
}

/* 启动服务 */
func launch(ctx *controllers.ProxyContext) {
	controllers.Launch(ctx)

	beego.Run()
}

func main() {

	/* > 初始化 */
	ctx := _init()
	if nil == ctx {
		fmt.Printf("Initialize context failed!\n")
		return
	}

	/* > 注册回调 */
	register(ctx)

	/* > 启动服务 */
	launch(ctx)

}
