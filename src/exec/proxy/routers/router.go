package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	"upgrade-api/src/exec/proxy/controllers"
)

func init() {

	ns := beego.NewNamespace("/upgrade/proxy/v1",

		beego.NSNamespace("/apk/upgrade",
			beego.NSInclude(
				&controllers.ApkUpgradeController{},
			),
		),

		beego.NSNamespace("/apk/uninstall",
			beego.NSInclude(
				&controllers.ApkUninstallController{},
			),
		),
	)

	beego.AddNamespace(ns)
	logs.Info("%s", "routers init complete ")
}
