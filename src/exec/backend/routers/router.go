package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	"upgrade-api/src/share/comm"

	"upgrade-api/src/exec/backend/controllers"
)

func init() {

	ns := beego.NewNamespace("/upgrade/backend/v1",

		beego.NSNamespace("/apk/strategy",
			beego.NSInclude(
				&controllers.ApkStrategyController{},
			),
		),

		// 应用管理 @author gsp
		beego.NSNamespace("/app",
			beego.NSInclude(
				&controllers.AppController{},
			),
		),

		// APK管理 @author gsp
		beego.NSNamespace("/apk",
			beego.NSInclude(
				&controllers.ApkController{},
			),
		),

		//业务管理  @author Zhao.yang
		beego.NSNamespace("/business",
			beego.NSInclude(
				&controllers.BusinessController{},
			),
		),

		//设备分组管理 @author Linlin.guo
		beego.NSNamespace("/dev/group",
			beego.NSInclude(
				&controllers.DevGroupController{},
			),
		),

		//设备类型管理 @author Linlin.guo
		beego.NSNamespace("/dev/type",
			beego.NSInclude(
				&controllers.DevTypeController{},
			),
		),

		//差分算法管理 @author Linlin.guo
		beego.NSNamespace("/patch/algo",
			beego.NSInclude(
				&controllers.PatchAlgoControleer{},
			),
		),

		// TODO 应用管理 @author gsp
		beego.NSNamespace("/file",
			beego.NSInclude(
				&controllers.ApkUploadController{},
			),
		),

		//cdn回调 @author Zhao.yang
		beego.NSNamespace("/cdn/callback",
			beego.NSInclude(
				&controllers.CdnCallBackController{},
			),
		),
	)
	beego.AddNamespace(ns)
	//页面渲染 @author Linlin.guo
	beego.Router("/", &controllers.IndexPageController{})

	beego.SetStaticPath(comm.STATIC_REQUEST_URL, comm.STORAGE_PATH)
	logs.Info("%s", "routers init complete ")
}
