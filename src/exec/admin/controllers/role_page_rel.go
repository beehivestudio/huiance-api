package controllers

import (
	"net/http"
	"strconv"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/json-iterator/go"

	"upgrade-api/src/exec/admin/models"
	"upgrade-api/src/share/comm"
	"upgrade-api/src/share/database"
)

// RolePageRelController operations for RolePageRel
type RolePageRelController struct {
	CommonController
}

// URLMapping ...
func (c *RolePageRelController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetAll", c.GetAll)
}

// Post ...
// @Title Post
// @Description create RolePageRel
// @Param	id		path 	string	true		"The key for staticblock"
// @Param	body		body 	models.RolePageRel	true		"body for RolePageRel content"
// @Success 201 {int} models.RolePageRel
// @Failure 403 body is empty
// @router /:id/page [post]
// @author # Linlin.guo # 2020-06-17 15:00:05 #
func (c *RolePageRelController) Post() {
	ctx := GetAdminCntx()

	//解析参数&校验参数
	idStr := c.Ctx.Input.Param(":id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if nil != err {
		logs.Error("%s() Get id failed! errmsg:%s", utils.GetRunFuncName(), err.Error())
		c.ErrorMessage(comm.ERR_PARAM_INVALID, err.Error())
		return
	}

	req, code, err := c.getPostParam()
	if nil != err {
		logs.Error("%s() Get role page parameter format is invalid! errmsg:%s",
			utils.GetRunFuncName(), err.Error())
		c.ErrorMessage(code, err.Error())
		return
	}

	err = models.UpdateRolePage(ctx.Model.Mysql.O, id, c.userId, req)
	if nil != err {
		logs.Error("%s() Update role page failed! errmsg:%s", utils.GetRunFuncName(), err.Error())
		code, e := database.MysqlFormatError(err)
		c.ErrorMessage(code, e)
		return
	}

	c.FormatResp(http.StatusOK, comm.OK, "Ok")
}

//解析参数&校验参数
func (c *RolePageRelController) getPostParam() (rolePage models.RolePage, code int, err error) {

	if err := jsoniter.Unmarshal(c.Ctx.Input.RequestBody, &rolePage); nil != err {
		logs.Error("%s() Unmarshal parameter failed! body:%s errmsg:%s",
			utils.GetRunFuncName(), c.Ctx.Input.RequestBody, err.Error())
		return rolePage, comm.ERR_PARAM_INVALID, err
	}

	if len(rolePage.PageIds) == 0 {
		return rolePage, comm.ERR_PARAM_MISS, err
	}

	return rolePage, 0, nil
}

// GetAll ...
// @Title Get All
// @Description get RolePageRel
// @Success 200 {object} models.RolePageRel
// @Failure 403
// @router /page/list [get]
// @author # Linlin.guo # 2020-06-17 15:00:05 #
func (c *RolePageRelController) GetAll() {
	ctx := GetAdminCntx()

	// 获取参数
	id, _ := c.GetInt64("role_id")

	roleName := c.GetString("role_name")
	pagesName := c.GetString("pages_name")

	l, err := models.GetRolePageList(ctx.Model.Mysql.O, id, roleName, pagesName)
	if err != nil {
		if orm.ErrNoRows == err {
			logs.Error("%s() Page not exist! id:%d", utils.GetRunFuncName(), id)
			c.ErrorMessage(comm.ERR_PARAM_INVALID, "Page not exist!")
			return
		}
		logs.Error("%s() GetAll Page failed! id:%d errmsg:%s",
			utils.GetRunFuncName(), id, err.Error())
		code, e := database.MysqlFormatError(err)
		c.ErrorMessage(code, e)
		return
	}
	c.FormatInterfaceListResp(http.StatusOK, comm.OK, len(l), "Ok", l)

}
