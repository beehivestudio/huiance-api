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

// PageController operations for Page
type PageController struct {
	CommonController
}

// URLMapping ...
func (c *PageController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Post
// @Description create Page
// @Param	body		body 	models.Page	true		"body for Page content"
// @Success 201 {int} models.Page
// @Failure 403 body is empty
// @router / [post]
// @author # Linlin.guo # 2020-06-17 15:00:05 #
func (c *PageController) Post() {
	ctx := GetAdminCntx()

	var v db_admin.Page
	var err error

	// 1.校验权限
	_, err = db_admin.ChargeAuthByUserId(ctx.Model.Mysql.O, c.userId)
	if nil != err {
		logs.Error("%s() ChargeAuthByUserId failed! body: %s errmsg:%s",
			utils.GetRunFuncName(), c.Ctx.Input.RequestBody, err.Error())
		c.ErrorMessage(comm.ERR_AUTH, err.Error())
		return
	}

	// 2.校验参数
	if err = jsoniter.Unmarshal(c.Ctx.Input.RequestBody, &v); nil != err {
		logs.Error("%s() Unmarshal parameter failed! body: %s errmsg:%s",
			utils.GetRunFuncName(), c.Ctx.Input.RequestBody, err.Error())
		c.ErrorMessage(comm.ERR_PARAM_INVALID, err.Error())
		return
	} else if 0 == len(v.Name) || 0 == len(v.Path) {
		logs.Error("%s() Parameter is invalid! body: %s",
			utils.GetRunFuncName(), c.Ctx.Input.RequestBody)
		c.ErrorMessage(comm.ERR_PARAM_MISS, "Parameter is invalid!")
		return
	}

	var id int64

	v.CreateUserId = c.userId

	if id, err = db_admin.AddPage(ctx.Model.Mysql.O, &v); err != nil {
		logs.Error("%s() Add page failed! errmsg:%s", utils.GetRunFuncName(), err.Error())
		code, e := database.MysqlFormatError(err)
		c.ErrorMessage(code, e)
		return
	}

	c.FormatAddResp(http.StatusOK, comm.OK, "Ok", id)
}

// Put ...
// @Title Put
// @Description update the Page
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Page	true		"body for Page content"
// @Success 200 {object} models.Page
// @Failure 403 :id is not int
// @router /:id([0-9]+) [put]
// @author # Linlin.guo # 2020-06-17 15:00:05 #
func (c *PageController) Put() {
	ctx := GetAdminCntx()

	// 校验权限
	_, err := db_admin.ChargeAuthByUserId(ctx.Model.Mysql.O, c.userId)
	if nil != err {
		logs.Error("%s() ChargeAuthByUserId failed! body: %s errmsg:%s",
			utils.GetRunFuncName(), c.Ctx.Input.RequestBody, err.Error())
		c.ErrorMessage(comm.ERR_AUTH, err.Error())
		return
	}

	// 校验参数
	idStr := c.Ctx.Input.Param(":id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if nil != err {
		logs.Error("%s() Get id failed! errmsg:%s", utils.GetRunFuncName(), err.Error())
		c.ErrorMessage(comm.ERR_PARAM_INVALID, err.Error())
		return
	}

	v := db_admin.Page{Id: id}

	if err := jsoniter.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
		logs.Error("%s() Unmarshal parameter failed! id:%d body: %s errmsg:%s",
			utils.GetRunFuncName(), id, c.Ctx.Input.RequestBody, err.Error())
		c.ErrorMessage(comm.ERR_PARAM_INVALID, err.Error())
		return
	} else if 0 == len(v.Name) || 0 == len(v.Path) {
		logs.Error("%s() Paramter is invalid! id:%d body: %s",
			utils.GetRunFuncName(), id, c.Ctx.Input.RequestBody)
		c.ErrorMessage(comm.ERR_PARAM_MISS, "Parameter is invalid!")
		return
	}

	v.UpdateUserId = c.userId

	if err := db_admin.UpdatePageById(ctx.Model.Mysql.O, &v); err != nil {
		if orm.ErrNoRows == err {
			logs.Error("%s() Page not exist! id:%d", utils.GetRunFuncName(), id)
			c.ErrorMessage(comm.ERR_PARAM_INVALID, "Page not exist!")
			return
		}
		logs.Error("%s() Update Page failed! id:%d errmsg:%s",
			utils.GetRunFuncName(), id, err.Error())
		code, e := database.MysqlFormatError(err)
		c.ErrorMessage(code, e)
		return
	}

	c.ErrorMessage(comm.OK, "Ok")
}

// GetOne ...
// @Title Get One
// @Description get Page by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Page
// @Failure 403 :id is empty
// @router /:id([0-9]+) [get]
// @author # Linlin.guo # 2020-06-17 15:00:05 #
func (c *PageController) GetOne() {
	ctx := GetAdminCntx()

	// 校验权限
	_, err := db_admin.ChargeAuthByUserId(ctx.Model.Mysql.O, c.userId)
	if nil != err {
		logs.Error("%s() ChargeAuthByUserId failed! body: %s errmsg:%s",
			utils.GetRunFuncName(), c.Ctx.Input.RequestBody, err.Error())
		c.ErrorMessage(comm.ERR_AUTH, err.Error())
		return
	}

	// 校验参数
	idStr := c.Ctx.Input.Param(":id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if nil != err {
		logs.Error("%s() Get id failed! errmsg:%s",
			utils.GetRunFuncName(), err.Error())
		c.ErrorMessage(comm.ERR_PARAM_INVALID, err.Error())
		return
	}

	v, err := db_admin.GetPageById(ctx.Model.Mysql.O, id)
	if nil != err {
		if orm.ErrNoRows == err {
			logs.Error("%s() Page not exist! id:%d", utils.GetRunFuncName(), id)
			c.ErrorMessage(comm.ERR_PARAM_INVALID, "Page not exist!")
			return
		}
		logs.Error("%s() Get page failed! id:%d errmsg:%s",
			utils.GetRunFuncName(), id, err.Error())
		c.ErrorMessage(comm.ERR_INTERNAL_SERVER_ERROR, err.Error())
		return
	}

	c.FormatInterfaceResp(http.StatusOK, comm.OK, "Ok", v)
}

// GetAll ...
// @Title Get All
// @Description get Page
// @Success 200 {object} models.Page
// @Failure 403
// @router /list [get]
// @author # Linlin.guo # 2020-06-17 15:00:05 #
func (c *PageController) GetAll() {
	ctx := GetAdminCntx()

	// 校验权限
	_, err := db_admin.ChargeAuthByUserId(ctx.Model.Mysql.O, c.userId)
	if nil != err {
		logs.Error("%s() ChargeAuthByUserId failed! body: %s errmsg:%s",
			utils.GetRunFuncName(), c.Ctx.Input.RequestBody, err.Error())
		c.ErrorMessage(comm.ERR_AUTH, err.Error())
		return
	}

	// 解析参数
	id, _ := c.GetInt64("id")
	name := c.GetString("name")

	l, err := db_admin.GetPageList(ctx.Model.Mysql.O, id, name)
	if err != nil {
		if orm.ErrNoRows == err {
			logs.Error("%s() Page not exist! id:%d", utils.GetRunFuncName(), id)
			c.ErrorMessage(comm.ERR_PARAM_INVALID, "Page not exist!")
			return
		}
		logs.Error("%s() Delete Page failed! id:%d errmsg:%s",
			utils.GetRunFuncName(), id, err.Error())
		code, e := database.MysqlFormatError(err)
		c.ErrorMessage(code, e)
		return
	}
	c.FormatInterfaceListResp(http.StatusOK, comm.OK, len(l), "Ok", l)
}

// Delete ...
// @Title Delete
// @Description delete the Page
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id([0-9]+) [delete]
// @author # Linlin.guo # 2020-06-17 15:00:05 #
func (c *PageController) Delete() {
	ctx := GetAdminCntx()

	// 校验权限
	_, err := db_admin.ChargeAuthByUserId(ctx.Model.Mysql.O, c.userId)
	if nil != err {
		logs.Error("%s() ChargeAuthByUserId failed! body: %s errmsg:%s",
			utils.GetRunFuncName(), c.Ctx.Input.RequestBody, err.Error())
		c.ErrorMessage(comm.ERR_AUTH, err.Error())
		return
	}

	// 校验参数
	idStr := c.Ctx.Input.Param(":id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if nil != err {
		logs.Error("%s() Parse id failed! id: %s errmsg:%s",
			utils.GetRunFuncName(), idStr, err.Error())
		c.ErrorMessage(comm.ERR_PARAM_INVALID, err.Error())
		return
	}

	if err := db_admin.DeletePage(ctx.Model.Mysql.O, id); err != nil {
		if orm.ErrNoRows == err {
			logs.Error("%s() Page not exist! id:%d", utils.GetRunFuncName(), id)
			c.ErrorMessage(comm.ERR_PARAM_INVALID, "Page not exist!")
			return
		}
		logs.Error("%s() Delete Page failed! id:%d errmsg:%s",
			utils.GetRunFuncName(), id, err.Error())
		code, e := database.MysqlFormatError(err)
		c.ErrorMessage(code, e)
		return

	}
	c.FormatResp(http.StatusOK, comm.OK, "Ok")
}
