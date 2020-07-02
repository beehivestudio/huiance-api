package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	"upgrade-api/src/exec/api/controllers"
)

func init() {
	ns := beego.NewNamespace("/upgrade/api/v1/",

		/* 策略相关 */
		beego.NSNamespace("/apk/strategy",
			beego.NSInclude(
				&controllers.ApkStrategyController{},
			),
		),

		/* 升级相关 */
		beego.NSNamespace("/apk/upgrade",
			beego.NSInclude(
				&controllers.ApkUpgradeController{},
			),
		),

		/* 卸载相关 */
		beego.NSNamespace("/apk/uninstall",
			beego.NSInclude(
				&controllers.ApkUninstallController{},
			),
		),

		beego.NSNamespace("/apk",
			beego.NSInclude(
				&controllers.ApkController{},
			),
		),

		beego.NSNamespace("/app",
			beego.NSInclude(
				&controllers.AppController{},
			),
		),
	)
	beego.AddNamespace(ns)

	logs.Info("%s", "routers init complete")
}
