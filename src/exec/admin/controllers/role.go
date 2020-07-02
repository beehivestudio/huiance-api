package controllers

import (
	"net/http"
	"strconv"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/json-iterator/go"

	"upgrade-api/src/share/comm"
	"upgrade-api/src/share/database"
	"upgrade-api/src/share/database/admin"
)

// RoleController operations for Role
type RoleController struct {
	CommonController
}

// URLMapping ...
func (c *RoleController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Post
// @Description create Role
// @Param	body		body 	models.Role	true		"body for Role content"
// @Success 201 {int} models.Role
// @Failure 403 body is empty
// @router / [post]
// @author # Linlin.guo # 2020-06-17 15:00:05 #
func (c *RoleController) Post() {
	ctx := GetAdminCntx()

	// 校验参数
	var v db_admin.Role
	var err error

	if err = jsoniter.Unmarshal(c.Ctx.Input.RequestBody, &v); nil != err {
		logs.Error("%s() Unmarshal parameter failed! body:%s errmsg:%s",
			utils.GetRunFuncName(), c.Ctx.Input.RequestBody, err.Error())
		c.ErrorMessage(comm.ERR_PARAM_INVALID, err.Error())
		return
	} else if 0 == len(v.Name) {
		logs.Error("%s() Parameter is invalid! body:%s",
			utils.GetRunFuncName(), c.Ctx.Input.RequestBody)
		c.ErrorMessage(comm.ERR_PARAM_MISS, "Parameter is invalid!")
		return
	}

	var id int64
	v.CreateUserId = c.userId

	if id, err = db_admin.AddRole(ctx.Model.Mysql.O, &v); err != nil {
		logs.Error("%s() Add role failed! errmsg:%s", utils.GetRunFuncName(), err.Error())
		code, e := database.MysqlFormatError(err)
		c.ErrorMessage(code, e)
		return
	}
	c.FormatAddResp(http.StatusOK, comm.OK, "Ok", id)
}

// Put ...
// @Title Put
// @Description update the Role
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Role	true		"body for Role content"
// @Success 200 {object} models.Role
// @Failure 403 :id is not int
// @router /:id([0-9]+) [put]
// @author # Linlin.guo # 2020-06-17 15:00:05 #
func (c *RoleController) Put() {
	ctx := GetAdminCntx()

	// 校验参数
	idStr := c.Ctx.Input.Param(":id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if nil != err {
		logs.Error("%s() Get id failed! errmsg:%s", utils.GetRunFuncName(), err.Error())
		c.ErrorMessage(comm.ERR_PARAM_INVALID, err.Error())
		return
	}

	v := db_admin.Role{Id: id}

	if err := jsoniter.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
		logs.Error("%s() Unmarshal parameter failed! id:%d body:%s errmsg:%s",
			utils.GetRunFuncName(), id, c.Ctx.Input.RequestBody, err.Error())
		c.ErrorMessage(comm.ERR_PARAM_INVALID, err.Error())
		return
	} else if 0 == len(v.Name) {
		logs.Error("%s() Paramter is invalid! id:%d body:%s",
			utils.GetRunFuncName(), id, c.Ctx.Input.RequestBody)
		c.ErrorMessage(comm.ERR_PARAM_MISS, "Parameter is invalid!")
		return
	}

	v.UpdateUserId = c.userId
	if err := db_admin.UpdateRoleById(ctx.Model.Mysql.O, &v); err != nil {
		if orm.ErrNoRows == err {
			logs.Error("%s() Role not exist! id:%d", utils.GetRunFuncName(), id)
			c.ErrorMessage(comm.ERR_PARAM_INVALID, "Role not exist!")
			return
		}
		logs.Error("%s() Update Role failed! id:%d errmsg:%s",
			utils.GetRunFuncName(), id, err.Error())
		code, e := database.MysqlFormatError(err)
		c.ErrorMessage(code, e)
		return
	}

	c.ErrorMessage(comm.OK, "Ok")
}

// GetOne ...
// @Title Get One
// @Description get Role by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Role
// @Failure 403 :id is empty
// @router /:id([0-9]+) [get]
// @author # Linlin.guo # 2020-06-17 15:00:05 #
func (c *RoleController) GetOne() {
	ctx := GetAdminCntx()

	// 校验参数
	idStr := c.Ctx.Input.Param(":id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if nil != err {
		logs.Error("%s() Parse id failed! id:%s errmsg:%s",
			utils.GetRunFuncName(), idStr, err.Error())
		c.ErrorMessage(comm.ERR_PARAM_INVALID, err.Error())
		return
	}

	v, err := db_admin.GetRoleById(ctx.Model.Mysql.O, id)
	if err != nil {
		if orm.ErrNoRows == err {
			logs.Error("%s() Role not exist! id:%d", utils.GetRunFuncName(), id)
			c.ErrorMessage(comm.ERR_PARAM_INVALID, "Role not exist!")
			return
		}
		logs.Error("%s() Get Role failed! id:%d errmsg:%s",
			utils.GetRunFuncName(), id, err.Error())
		code, e := database.MysqlFormatError(err)
		c.ErrorMessage(code, e)
		return
	}

	c.FormatInterfaceResp(http.StatusOK, comm.OK, "Ok", v)
}

// GetAll ...
// @Title Get All
// @Description get Role
// @Success 200 {object} models.Role
// @Failure 403
// @router /list [get]
// @author # Linlin.guo # 2020-06-17 15:00:05 #
func (c *RoleController) GetAll() {
	ctx := GetAdminCntx()

	// 校验参数
	id, _ := c.GetInt64("id")
	name := c.GetString("name")

	l, err := db_admin.GetRoleList(ctx.Model.Mysql.O, id, name)
	if err != nil {
		if orm.ErrNoRows == err {
			logs.Error("%s() Role not exist! id:%d", utils.GetRunFuncName(), id)
			c.ErrorMessage(comm.ERR_PARAM_INVALID, "Role not exist!")
			return
		}
		logs.Error("%s() GetRoleList failed! id:%d errmsg:%s",
			utils.GetRunFuncName(), id, err.Error())
		code, e := database.MysqlFormatError(err)
		c.ErrorMessage(code, e)
		return
	}

	c.FormatInterfaceListResp(http.StatusOK, comm.OK, len(l), "Ok", l)
}

// Delete ...
// @Title Delete
// @Description delete the Role
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id([0-9]+) [delete]
// @author # Linlin.guo # 2020-06-17 15:00:05 #
func (c *RoleController) Delete() {
	ctx := GetAdminCntx()

	// 校验参数
	idStr := c.Ctx.Input.Param(":id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if nil != err {
		logs.Error("%s() Parse id failed! id:%s errmsg:%s",
			utils.GetRunFuncName(), idStr, err.Error())
		c.ErrorMessage(comm.ERR_PARAM_INVALID, err.Error())
		return
	}

	if err := db_admin.DeleteRole(ctx.Model.Mysql.O, id); err != nil {
		if orm.ErrNoRows == err {
			logs.Error("%s() Role not exist! id:%d", utils.GetRunFuncName(), id)
			c.ErrorMessage(comm.ERR_PARAM_INVALID, "Role not exist!")
			return
		}
		logs.Error("%s() Delete Role failed! id:%d errmsg:%s",
			utils.GetRunFuncName(), id, err.Error())
		code, e := database.MysqlFormatError(err)
		c.ErrorMessage(code, e)
		return

	}
	c.FormatResp(http.StatusOK, comm.OK, "Ok")
}
