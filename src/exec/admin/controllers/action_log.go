package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/astaxie/beego/logs"

	"upgrade-api/src/exec/admin/models"
	"upgrade-api/src/share/comm"
	"upgrade-api/src/share/database"
)

// ActionLogController operations for ActionLog
type ActionLogController struct {
	CommonController
}

// URLMapping ...
func (c *ActionLogController) URLMapping() {
	c.Mapping("GetAll", c.GetAll)
}

// 获取用户行为
// @Title 获取用户行为
// @Description 获取用户行为
// @Param	user_id			query	int64	true	"用户id"
// @Param	action			query	string	true	"用户行为id取值"
// @Param	user_email		query	string	true	"用户邮箱"
// @Success 200 {object} comm.InterfaceResp
// @Failure 403
// @router / [get]
// @author # Zhao.yang # 2020-06-18 10:01:08 #
func (c *ActionLogController) GetAll() {

	ctx := GetAdminCntx()

	//获取参数&校验参数
	fields, values, page, pageSize, err := c.getListParam()
	if nil != err {
		logs.Error("%s() Get action log parameter failed! errmsg:%s",
			utils.GetRunFuncName(), err.Error())
		c.ErrorMessage(comm.ERR_PARAM_INVALID, err.Error())
		return
	}

	//获取日志列表
	total, logList, err := models.GetActionLogList(ctx.Model.Mysql.O, fields, values, page, pageSize)
	if nil != err {
		logs.Error("%s() Get action log list failed! fields:%v errmsg:%s",
			utils.GetRunFuncName(), fields, err.Error())
		code, e := database.MysqlFormatError(err)
		c.ErrorMessage(code, e)
		return
	}

	l := len(logList)

	list := &comm.ListData{
		Page:     page,
		PageSize: pageSize,
		Total:    int(total),
		Len:      l,
		List:     make([]interface{}, l),
	}

	for u, log := range logList {
		list.List[u] = &LogList{
			Id:         log.Id,
			UserName:   log.Name,
			UserEmail:  log.Email,
			ActionId:   log.ActionId,
			Action:     log.Action,
			Detail:     log.Detail,
			CreateTime: log.CreateTime,
		}
	}

	logs.Info("%s() Get action log list num msg :%v", utils.GetRunFuncName(), len(list.List))

	c.FormatInterfaceResp(http.StatusOK, comm.OK, "Ok", list)
}

/* 返回用户行为数据 */
type LogList struct {
	Id         int64     `json:"id" description:"日志id"`
	UserName   string    `json:"user_name" description:"操作人名称"`
	UserEmail  string    `json:"user_email" description:"操作人邮箱"`
	ActionId   int64     `json:"action_id" description:"用户行为id"`
	Action     string    `json:"action" description:"用户行为描述"`
	Detail     string    `json:"detail" description:"日志详细信息"`
	CreateTime time.Time `json:"create_time" description:"创建时间"`
}

//获取日志参数列表
func (c *ActionLogController) getListParam() (
	fields []string, values []interface{}, page, pageSize int, err error) {

	page = comm.PAGE_START
	pageSize = comm.PAGE_SIZE

	idStr := c.GetString("user_id")
	if "" != idStr {
		uid, err := strconv.ParseInt(idStr, 10, 64)
		if nil != err {
			logs.Error("%s() Get user id failed! errmsg:%s",
				utils.GetRunFuncName(), err.Error())
			return nil, nil, 0, 0, err
		}
		fields = append(fields, "actionlog.user_id = ?")
		values = append(values, uid)
	}

	action := c.GetString("action_id")
	if "" != action {
		a, err := strconv.ParseInt(action, 10, 64)
		if nil != err {
			logs.Error("%s() Get action failed! errmsg:%s",
				utils.GetRunFuncName(), err.Error())
			return nil, nil, 0, 0, err
		}
		fields = append(fields, "actionlog.action_id = ?")
		values = append(values, a)
	}

	email := c.GetString("user_email")
	if "" != email {
		fields = append(fields, "userinfor.email LIKE ?")
		values = append(values, fmt.Sprintf("%%%s%%", email))
	}

	return fields, values, page, pageSize, nil

}
