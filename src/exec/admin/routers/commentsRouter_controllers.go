package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["upgrade-api/src/exec/admin/controllers:ActionController"] = append(beego.GlobalControllerRouter["upgrade-api/src/exec/admin/controllers:ActionController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: `/list`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["upgrade-api/src/exec/admin/controllers:ActionLogController"] = append(beego.GlobalControllerRouter["upgrade-api/src/exec/admin/controllers:ActionLogController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: `/`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["upgrade-api/src/exec/admin/controllers:PageController"] = append(beego.GlobalControllerRouter["upgrade-api/src/exec/admin/controllers:PageController"],
        beego.ControllerComments{
            Method: "Post",
            Router: `/`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["upgrade-api/src/exec/admin/controllers:PageController"] = append(beego.GlobalControllerRouter["upgrade-api/src/exec/admin/controllers:PageController"],
        beego.ControllerComments{
            Method: "Put",
            Router: `/:id([0-9]+)`,
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["upgrade-api/src/exec/admin/controllers:PageController"] = append(beego.GlobalControllerRouter["upgrade-api/src/exec/admin/controllers:PageController"],
        beego.ControllerComments{
            Method: "GetOne",
            Router: `/:id([0-9]+)`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["upgrade-api/src/exec/admin/controllers:PageController"] = append(beego.GlobalControllerRouter["upgrade-api/src/exec/admin/controllers:PageController"],
        beego.ControllerComments{
            Method: "Delete",
            Router: `/:id([0-9]+)`,
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["upgrade-api/src/exec/admin/controllers:PageController"] = append(beego.GlobalControllerRouter["upgrade-api/src/exec/admin/controllers:PageController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: `/list`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["upgrade-api/src/exec/admin/controllers:RoleController"] = append(beego.GlobalControllerRouter["upgrade-api/src/exec/admin/controllers:RoleController"],
        beego.ControllerComments{
            Method: "Post",
            Router: `/`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["upgrade-api/src/exec/admin/controllers:RoleController"] = append(beego.GlobalControllerRouter["upgrade-api/src/exec/admin/controllers:RoleController"],
        beego.ControllerComments{
            Method: "Put",
            Router: `/:id([0-9]+)`,
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["upgrade-api/src/exec/admin/controllers:RoleController"] = append(beego.GlobalControllerRouter["upgrade-api/src/exec/admin/controllers:RoleController"],
        beego.ControllerComments{
            Method: "GetOne",
            Router: `/:id([0-9]+)`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["upgrade-api/src/exec/admin/controllers:RoleController"] = append(beego.GlobalControllerRouter["upgrade-api/src/exec/admin/controllers:RoleController"],
        beego.ControllerComments{
            Method: "Delete",
            Router: `/:id([0-9]+)`,
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["upgrade-api/src/exec/admin/controllers:RoleController"] = append(beego.GlobalControllerRouter["upgrade-api/src/exec/admin/controllers:RoleController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: `/list`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["upgrade-api/src/exec/admin/controllers:RolePageRelController"] = append(beego.GlobalControllerRouter["upgrade-api/src/exec/admin/controllers:RolePageRelController"],
        beego.ControllerComments{
            Method: "Post",
            Router: `/:id/page`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["upgrade-api/src/exec/admin/controllers:RolePageRelController"] = append(beego.GlobalControllerRouter["upgrade-api/src/exec/admin/controllers:RolePageRelController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: `/page/list`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["upgrade-api/src/exec/admin/controllers:UserController"] = append(beego.GlobalControllerRouter["upgrade-api/src/exec/admin/controllers:UserController"],
        beego.ControllerComments{
            Method: "Post",
            Router: `/`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["upgrade-api/src/exec/admin/controllers:UserController"] = append(beego.GlobalControllerRouter["upgrade-api/src/exec/admin/controllers:UserController"],
        beego.ControllerComments{
            Method: "Put",
            Router: `/:id([0-9]+)`,
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["upgrade-api/src/exec/admin/controllers:UserController"] = append(beego.GlobalControllerRouter["upgrade-api/src/exec/admin/controllers:UserController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: `/list`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["upgrade-api/src/exec/admin/controllers:UserRoleRelController"] = append(beego.GlobalControllerRouter["upgrade-api/src/exec/admin/controllers:UserRoleRelController"],
        beego.ControllerComments{
            Method: "Post",
            Router: `/:id/role`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["upgrade-api/src/exec/admin/controllers:UserRoleRelController"] = append(beego.GlobalControllerRouter["upgrade-api/src/exec/admin/controllers:UserRoleRelController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: `/role/list`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
