package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["upgrade-api/src/exec/backend/controllers:ApkController"] = append(beego.GlobalControllerRouter["upgrade-api/src/exec/backend/controllers:ApkController"],
        beego.ControllerComments{
            Method: "Post",
            Router: `/`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["upgrade-api/src/exec/backend/controllers:ApkController"] = append(beego.GlobalControllerRouter["upgrade-api/src/exec/backend/controllers:ApkController"],
        beego.ControllerComments{
            Method: "Put",
            Router: `/:id([0-9]+)`,
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["upgrade-api/src/exec/backend/controllers:ApkController"] = append(beego.GlobalControllerRouter["upgrade-api/src/exec/backend/controllers:ApkController"],
        beego.ControllerComments{
            Method: "Get",
            Router: `/:id([0-9]+)`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["upgrade-api/src/exec/backend/controllers:ApkController"] = append(beego.GlobalControllerRouter["upgrade-api/src/exec/backend/controllers:ApkController"],
        beego.ControllerComments{
            Method: "GetList",
            Router: `/list`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["upgrade-api/src/exec/backend/controllers:ApkStrategyController"] = append(beego.GlobalControllerRouter["upgrade-api/src/exec/backend/controllers:ApkStrategyController"],
        beego.ControllerComments{
            Method: "Post",
            Router: `/`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["upgrade-api/src/exec/backend/controllers:ApkStrategyController"] = append(beego.GlobalControllerRouter["upgrade-api/src/exec/backend/controllers:ApkStrategyController"],
        beego.ControllerComments{
            Method: "Put",
            Router: `/:id([0-9]+)`,
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["upgrade-api/src/exec/backend/controllers:ApkStrategyController"] = append(beego.GlobalControllerRouter["upgrade-api/src/exec/backend/controllers:ApkStrategyController"],
        beego.ControllerComments{
            Method: "GetOne",
            Router: `/:id([0-9]+)`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["upgrade-api/src/exec/backend/controllers:ApkStrategyController"] = append(beego.GlobalControllerRouter["upgrade-api/src/exec/backend/controllers:ApkStrategyController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: `/list`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["upgrade-api/src/exec/backend/controllers:ApkUploadController"] = append(beego.GlobalControllerRouter["upgrade-api/src/exec/backend/controllers:ApkUploadController"],
        beego.ControllerComments{
            Method: "Post",
            Router: `/`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["upgrade-api/src/exec/backend/controllers:ApkUploadController"] = append(beego.GlobalControllerRouter["upgrade-api/src/exec/backend/controllers:ApkUploadController"],
        beego.ControllerComments{
            Method: "GetOne",
            Router: `/:id([0-9]+)`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["upgrade-api/src/exec/backend/controllers:ApkUploadController"] = append(beego.GlobalControllerRouter["upgrade-api/src/exec/backend/controllers:ApkUploadController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: `/list`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["upgrade-api/src/exec/backend/controllers:AppController"] = append(beego.GlobalControllerRouter["upgrade-api/src/exec/backend/controllers:AppController"],
        beego.ControllerComments{
            Method: "Post",
            Router: `/`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["upgrade-api/src/exec/backend/controllers:AppController"] = append(beego.GlobalControllerRouter["upgrade-api/src/exec/backend/controllers:AppController"],
        beego.ControllerComments{
            Method: "Put",
            Router: `/:id([0-9]+)`,
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["upgrade-api/src/exec/backend/controllers:AppController"] = append(beego.GlobalControllerRouter["upgrade-api/src/exec/backend/controllers:AppController"],
        beego.ControllerComments{
            Method: "Get",
            Router: `/:id([0-9]+)`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["upgrade-api/src/exec/backend/controllers:AppController"] = append(beego.GlobalControllerRouter["upgrade-api/src/exec/backend/controllers:AppController"],
        beego.ControllerComments{
            Method: "GetList",
            Router: `/list`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["upgrade-api/src/exec/backend/controllers:BusinessController"] = append(beego.GlobalControllerRouter["upgrade-api/src/exec/backend/controllers:BusinessController"],
        beego.ControllerComments{
            Method: "Post",
            Router: `/`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["upgrade-api/src/exec/backend/controllers:BusinessController"] = append(beego.GlobalControllerRouter["upgrade-api/src/exec/backend/controllers:BusinessController"],
        beego.ControllerComments{
            Method: "GetOne",
            Router: `/:id([0-9]+)`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["upgrade-api/src/exec/backend/controllers:BusinessController"] = append(beego.GlobalControllerRouter["upgrade-api/src/exec/backend/controllers:BusinessController"],
        beego.ControllerComments{
            Method: "Put",
            Router: `/:id([0-9]+)`,
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["upgrade-api/src/exec/backend/controllers:BusinessController"] = append(beego.GlobalControllerRouter["upgrade-api/src/exec/backend/controllers:BusinessController"],
        beego.ControllerComments{
            Method: "Delete",
            Router: `/:id([0-9]+)`,
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["upgrade-api/src/exec/backend/controllers:BusinessController"] = append(beego.GlobalControllerRouter["upgrade-api/src/exec/backend/controllers:BusinessController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: `/list`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["upgrade-api/src/exec/backend/controllers:CdnCallBackController"] = append(beego.GlobalControllerRouter["upgrade-api/src/exec/backend/controllers:CdnCallBackController"],
        beego.ControllerComments{
            Method: "Post",
            Router: `/`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["upgrade-api/src/exec/backend/controllers:DevGroupController"] = append(beego.GlobalControllerRouter["upgrade-api/src/exec/backend/controllers:DevGroupController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: `/list`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["upgrade-api/src/exec/backend/controllers:DevTypeController"] = append(beego.GlobalControllerRouter["upgrade-api/src/exec/backend/controllers:DevTypeController"],
        beego.ControllerComments{
            Method: "Post",
            Router: `/`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["upgrade-api/src/exec/backend/controllers:DevTypeController"] = append(beego.GlobalControllerRouter["upgrade-api/src/exec/backend/controllers:DevTypeController"],
        beego.ControllerComments{
            Method: "Put",
            Router: `/:id([0-9]+)`,
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["upgrade-api/src/exec/backend/controllers:DevTypeController"] = append(beego.GlobalControllerRouter["upgrade-api/src/exec/backend/controllers:DevTypeController"],
        beego.ControllerComments{
            Method: "GetOne",
            Router: `/:id([0-9]+)`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["upgrade-api/src/exec/backend/controllers:DevTypeController"] = append(beego.GlobalControllerRouter["upgrade-api/src/exec/backend/controllers:DevTypeController"],
        beego.ControllerComments{
            Method: "Delete",
            Router: `/:id([0-9]+)`,
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["upgrade-api/src/exec/backend/controllers:DevTypeController"] = append(beego.GlobalControllerRouter["upgrade-api/src/exec/backend/controllers:DevTypeController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: `/list`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["upgrade-api/src/exec/backend/controllers:PatchAlgoControleer"] = append(beego.GlobalControllerRouter["upgrade-api/src/exec/backend/controllers:PatchAlgoControleer"],
        beego.ControllerComments{
            Method: "Post",
            Router: `/`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["upgrade-api/src/exec/backend/controllers:PatchAlgoControleer"] = append(beego.GlobalControllerRouter["upgrade-api/src/exec/backend/controllers:PatchAlgoControleer"],
        beego.ControllerComments{
            Method: "Put",
            Router: `/:id([0-9]+)`,
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["upgrade-api/src/exec/backend/controllers:PatchAlgoControleer"] = append(beego.GlobalControllerRouter["upgrade-api/src/exec/backend/controllers:PatchAlgoControleer"],
        beego.ControllerComments{
            Method: "GetOne",
            Router: `/:id([0-9]+)`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["upgrade-api/src/exec/backend/controllers:PatchAlgoControleer"] = append(beego.GlobalControllerRouter["upgrade-api/src/exec/backend/controllers:PatchAlgoControleer"],
        beego.ControllerComments{
            Method: "Delete",
            Router: `/:id([0-9]+)`,
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["upgrade-api/src/exec/backend/controllers:PatchAlgoControleer"] = append(beego.GlobalControllerRouter["upgrade-api/src/exec/backend/controllers:PatchAlgoControleer"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: `/list`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
