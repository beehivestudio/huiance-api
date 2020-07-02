package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["upgrade-api/src/exec/proxy/controllers:ApkUninstallController"] = append(beego.GlobalControllerRouter["upgrade-api/src/exec/proxy/controllers:ApkUninstallController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: `/list`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["upgrade-api/src/exec/proxy/controllers:ApkUpgradeController"] = append(beego.GlobalControllerRouter["upgrade-api/src/exec/proxy/controllers:ApkUpgradeController"],
        beego.ControllerComments{
            Method: "GetOne",
            Router: `/`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["upgrade-api/src/exec/proxy/controllers:ApkUpgradeController"] = append(beego.GlobalControllerRouter["upgrade-api/src/exec/proxy/controllers:ApkUpgradeController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: `/list`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
