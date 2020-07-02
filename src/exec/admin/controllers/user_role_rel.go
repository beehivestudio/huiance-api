package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/astaxie/beego/logs"
	jsoniter "github.com/json-iterator/go"

	"upgrade-api/src/exec/admin/models"
	"upgrade-api/src/share/comm"
	"upgrade-api/src/share/database"
)

// UserRoleRelController operations for UserRoleRel
type UserRoleRelController struct {
	CommonController
}

// URLMapping ...
func (c *UserRoleRelController) URLMapping() {
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Post", c.Post)

}

// 用户角色修改
// @Title 用户角色修改
// @Description 用户角色修改
// @Param	body		body 	models.UserRoleRel	true		"body for UserRoleRel content"
// @Success 201 {int} models.UserRoleRel
// @Failure 403 body is empty
// @router /:id/role [post]
// @author # Zhao.yang # 2020-06-17 15:01:08 #
func (c *UserRoleRelController) Post() {
	ctx := GetAdminCntx()

	id, err := strconv.ParseInt(c.Ctx.Input.Param(":id"), 10, 64)
	if nil != err {
		logs.Error("%s() Parse id failed! errmsg:%s", utils.GetRunFuncName(), err.Error())
		c.ErrorMessage(comm.ERR_PARAM_INVALID, err.Error())
		return
	}

	//解析参数&校验参数
	req, code, err := c.getPostParam()
	if nil != err {
		logs.Error("%s() User role parameter format is invalid! errmsg:%s",
			utils.GetRunFuncName(), err.Error())
		c.ErrorMessage(code, err.Error())
		return
	}

	err = models.UpdateUserRole(ctx.Model.Mysql.O, id, &req, c.userId)
	if nil != err {
		logs.Error("%s() Update user role failed! errmsg:%s",
			utils.GetRunFuncName(), err.Error())
		code, e := database.MysqlFormatError(err)
		c.ErrorMessage(code, e)
		return
	}

	c.FormatResp(http.StatusOK, comm.OK, "Ok")
}

//解析参数&校验参数
func (c *UserRoleRelController) getPostParam() (u models.UserRole, code int, err error) {

	if err := jsoniter.Unmarshal(c.Ctx.Input.RequestBody, &u); nil != err {
		logs.Error("%s() Unmarshal parameter failed! body:%s errmsg:%s",
			utils.GetRunFuncName(), c.Ctx.Input.RequestBody, err.Error())
		return u, comm.ERR_PARAM_INVALID, err
	}

	if len(u.RoleIds) == 0 {
		return u, comm.ERR_PARAM_MISS, err
	}

	return u, 0, nil
}

// 获取用户角色列表
// @Title 获取用户角色列表
// @Description 获取用户角色列表
// @Success 200 {object} comm.InterfaceResp
// @Failure 403
// @router /role/list [get]
// @Author Zhao.yang 2020-06-17 13:38:02
func (c *UserRoleRelController) GetAll() {
	ctx := GetAdminCntx()

	fields, page, pageSize, err := c.getListParam()
	if nil != err {
		logs.Error("%s() Get user role list parameter failed! errmsg:%s",
			utils.GetRunFuncName(), err.Error())
		c.ErrorMessage(comm.ERR_PARAM_INVALID, err.Error())
		return
	}

	//获取用户列表
	total, users, err := models.GetUserList(ctx.Model.Mysql.O, fields, page, pageSize)
	if nil != err {
		logs.Error("%s() Get user list failed! fields:%v errmsg:%s",
			utils.GetRunFuncName(), fields, err.Error())
		code, e := database.MysqlFormatError(err)
		c.ErrorMessage(code, e)
		return
	}

	l := len(users)

	list := &comm.ListData{
		Page:     page,
		PageSize: pageSize,
		Total:    int(total),
		Len:      l,
		List:     make([]interface{}, l),
	}

	for u, user := range users {

		//获取用户角色
		roles, err := models.GetUserRoles(ctx.Model.Mysql.O, user.Id)
		if nil != err {
			logs.Error("%s() Get user roles list failed! errmsg:%s",
				utils.GetRunFuncName(), err.Error())
			code, e := database.MysqlFormatError(err)
			c.ErrorMessage(code, e)
			return
		}

		list.List[u] = &UserRoleList{
			Id:           user.Id,
			Name:         user.Name,
			Email:        user.Email,
			Description:  user.Description,
			CreateTime:   user.CreateTime,
			UpdateTime:   user.UpdateTime,
			CreateUserId: user.CreateUserId,
			UpdateUserId: user.UpdateUserId,
			Roles:        roles,
		}
	}

	logs.Info("Get user role list num msg :%v", len(list.List))

	c.FormatInterfaceResp(http.StatusOK, comm.OK, "Ok", list)

}

/* 获取用户角色列表返回参数 */
type UserRoleList struct {
	Id           int64          `json:"id"`
	Name         string         `json:"name"`
	Email        string         `json:"email"`
	Enable       int            `json:"enable"`
	Description  string         `json:"description"`
	CreateTime   time.Time      `json:"create_time"`
	CreateUserId int64          `json:"create_user_id"`
	UpdateTime   time.Time      `json:"update_time"`
	UpdateUserId int64          `json:"update_user_id"`
	Roles        []models.Roles `json:"roles"`
}

//获取用户参数列表
func (c *UserRoleRelController) getListParam() (fields map[string]interface{},
	page, pageSize int, err error) {

	page = comm.PAGE_START
	pageSize = comm.PAGE_SIZE

	fields = make(map[string]interface{})

	userId := c.GetString("user_id")
	if "" != userId {
		fields["id"] = userId
	}

	return fields, page, pageSize, nil
}
