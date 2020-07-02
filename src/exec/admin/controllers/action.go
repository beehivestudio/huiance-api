package controllers

import (
	"net/http"
	"time"

	"github.com/astaxie/beego/logs"

	"upgrade-api/src/share/comm"
	"upgrade-api/src/share/database"
	db_admin "upgrade-api/src/share/database/admin"
	"upgrade-api/src/share/utils"
)

// ActionController operations for User
type ActionController struct {
	CommonController
}

// URLMapping ...
func (c *ActionController) URLMapping() {
	c.Mapping("GetAll", c.GetAll)
}

// 获取行为信息列表
// @Title 获取行为信息列表
// @Description 获取行为信息列表
// @Success 200 {object} comm.InterfaceResp
// @Failure 403
// @router /list [get]
// @author # Zhao.yang # 2020-06-17 14:01:08 #
func (c *ActionController) GetAll() {
	ctx := GetAdminCntx()

	//获取行为参数列表
	actions, err := db_admin.GetActionList(ctx.Model.Mysql.O)
	if nil != err {
		logs.Error("%s() Get action list failed! actions:%v errmsg:%s",
			utils.GetRunFuncName(), actions, err.Error())
		code, e := database.MysqlFormatError(err)
		c.ErrorMessage(code, e)
		return
	}

	var list []Action

	for _, action := range actions {
		l := Action{
			Id:          action.Id,
			Action:      action.Action,
			Description: action.Description,
			CreateTime:  action.CreateTime,
		}
		list = append(list, l)
	}

	c.FormatInterfaceListResp(http.StatusOK, comm.OK, len(actions), "Ok", list)
}

/* 获取行为列表返回数据 */
type Action struct {
	Id          int64     `json:"id" description:"行为ID"`
	Action      string    `json:"action"description:"角色ID"`
	Description string    `json:"description" description:"描述"`
	CreateTime  time.Time `json:"create_time" description:"创建时间"`
}
