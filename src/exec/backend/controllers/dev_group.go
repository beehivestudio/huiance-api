package controllers

import (
	"net/http"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"

	"upgrade-api/src/share/comm"
	"upgrade-api/src/share/database"
	"upgrade-api/src/share/database/upgrade"
)

// DevGroupController operations for DevGroup
type DevGroupController struct {
	CommonController
}

// URLMapping ...
func (c *DevGroupController) URLMapping() {
	c.Mapping("GetAll", c.GetAll)
}

// Get ...
// @Title Get One
// @Description get GetDevGroupList by dev_type_id
// @Param	dev_type_id		query 	int 	true	"The dev_type_id for select"
// @Success 200 {object} models.DevGroup
// @Failure 403
// @router /list [get]
// @author # Linlin.guo # 2020-06-09 15:00:05 #
func (c *DevGroupController) GetAll() {
	var id int64
	var err error

	ctx := GetBackendCntx()

	if id, err = c.GetInt64("dev_type_id"); nil != err {
		logs.Error("Get dev type id failed! errmsg:%s", err.Error())
		c.ErrorMessage(comm.ERR_PARAM_INVALID, err.Error())
		return
	}

	list, err := upgrade.GetDevGroupList(ctx.Model.Mysql.O, ctx.Model.Redis.Get(), id)
	if nil != err {
		if orm.ErrNoRows == err {
			logs.Error("DevTypeId not exist! id:%d", id)
			c.ErrorMessage(comm.ERR_PARAM_INVALID, "DevTypeId not exist!")
			return
		}
		logs.Error("Get dev group list failed! errmsg:%s", err.Error())
		code, e := database.MysqlFormatError(err)
		c.ErrorMessage(code, e)
		return
	}

	c.FormatInterfaceResp(http.StatusOK, comm.OK, "Ok", list)
}
