package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["upgrade-api/src/exec/api/controllers:ApkController"] = append(beego.GlobalControllerRouter["upgrade-api/src/exec/api/controllers:ApkController"],
        beego.ControllerComments{
            Method: "Post",
            Router: `/`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["upgrade-api/src/exec/api/controllers:ApkController"] = append(beego.GlobalControllerRouter["upgrade-api/src/exec/api/controllers:ApkController"],
        beego.ControllerComments{
            Method: "GetOne",
            Router: `/:id([0-9]+)`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["upgrade-api/src/exec/api/controllers:ApkController"] = append(beego.GlobalControllerRouter["upgrade-api/src/exec/api/controllers:ApkController"],
        beego.ControllerComments{
            Method: "Put",
            Router: `/:id([0-9]+)`,
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["upgrade-api/src/exec/api/controllers:ApkController"] = append(beego.GlobalControllerRouter["upgrade-api/src/exec/api/controllers:ApkController"],
        beego.ControllerComments{
            Method: "GetOneByBcode",
            Router: `/bcode/:bcode`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["upgrade-api/src/exec/api/controllers:ApkController"] = append(beego.GlobalControllerRouter["upgrade-api/src/exec/api/controllers:ApkController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: `/list`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["upgrade-api/src/exec/api/controllers:ApkStrategyController"] = append(beego.GlobalControllerRouter["upgrade-api/src/exec/api/controllers:ApkStrategyController"],
        beego.ControllerComments{
            Method: "Post",
            Router: `/`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["upgrade-api/src/exec/api/controllers:ApkStrategyController"] = append(beego.GlobalControllerRouter["upgrade-api/src/exec/api/controllers:ApkStrategyController"],
        beego.ControllerComments{
            Method: "Put",
            Router: `/:id([0-9]+)`,
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["upgrade-api/src/exec/api/controllers:ApkStrategyController"] = append(beego.GlobalControllerRouter["upgrade-api/src/exec/api/controllers:ApkStrategyController"],
        beego.ControllerComments{
            Method: "GetOne",
            Router: `/:id([0-9]+)`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["upgrade-api/src/exec/api/controllers:ApkStrategyController"] = append(beego.GlobalControllerRouter["upgrade-api/src/exec/api/controllers:ApkStrategyController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: `/list`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["upgrade-api/src/exec/api/controllers:ApkUninstallController"] = append(beego.GlobalControllerRouter["upgrade-api/src/exec/api/controllers:ApkUninstallController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: `/list`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["upgrade-api/src/exec/api/controllers:ApkUpgradeController"] = append(beego.GlobalControllerRouter["upgrade-api/src/exec/api/controllers:ApkUpgradeController"],
        beego.ControllerComments{
            Method: "GetOne",
            Router: `/`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["upgrade-api/src/exec/api/controllers:ApkUpgradeController"] = append(beego.GlobalControllerRouter["upgrade-api/src/exec/api/controllers:ApkUpgradeController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: `/list`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["upgrade-api/src/exec/api/controllers:AppController"] = append(beego.GlobalControllerRouter["upgrade-api/src/exec/api/controllers:AppController"],
        beego.ControllerComments{
            Method: "Post",
            Router: `/`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["upgrade-api/src/exec/api/controllers:AppController"] = append(beego.GlobalControllerRouter["upgrade-api/src/exec/api/controllers:AppController"],
        beego.ControllerComments{
            Method: "GetOne",
            Router: `/:id([0-9]+)`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["upgrade-api/src/exec/api/controllers:AppController"] = append(beego.GlobalControllerRouter["upgrade-api/src/exec/api/controllers:AppController"],
        beego.ControllerComments{
            Method: "Put",
            Router: `/:id([0-9]+)`,
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["upgrade-api/src/exec/api/controllers:AppController"] = append(beego.GlobalControllerRouter["upgrade-api/src/exec/api/controllers:AppController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: `/list`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
