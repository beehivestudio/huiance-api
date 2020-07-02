package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	"upgrade-api/src/exec/admin/controllers"
)

func init() {
	ns := beego.NewNamespace("/upgrade/auth/v1",

		beego.NSNamespace("/action/log",
			beego.NSInclude(
				&controllers.ActionLogController{},
			),
		),

		beego.NSNamespace("/page",
			beego.NSInclude(
				&controllers.PageController{},
			),
		),

		beego.NSNamespace("/role",
			beego.NSInclude(
				&controllers.RoleController{},
			),
		),

		beego.NSNamespace("/role",
			beego.NSInclude(
				&controllers.RolePageRelController{},
			),
		),

		beego.NSNamespace("/user",
			beego.NSInclude(
				&controllers.UserController{},
			),
		),

		beego.NSNamespace("/user",
			beego.NSInclude(
				&controllers.UserRoleRelController{},
			),
		),

		beego.NSNamespace("/action",
			beego.NSInclude(
				&controllers.ActionController{},
			),
		),
	)
	beego.AddNamespace(ns)
	logs.Info("%s", "routers init complete")
}
