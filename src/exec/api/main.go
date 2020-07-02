package main

import (
	"fmt"
	"github.com/astaxie/beego"
	"runtime"
	"upgrade-api/src/exec/api/controllers"
	"upgrade-api/src/exec/api/controllers/conf"
	_ "upgrade-api/src/exec/api/routers"
)

/* 初始化 */
func _init() *controllers.ApiContext {
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
func register(ctx *controllers.ApiContext) {
	controllers.ApiCntx = ctx
}

/* 启动服务 */
func launch(ctx *controllers.ApiContext) {
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

	runtime.GOMAXPROCS(runtime.NumCPU())

	/* > 注册回调 */
	register(ctx)

	/* > 启动服务 */
	launch(ctx)

}
