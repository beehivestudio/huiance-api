package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/json-iterator/go"

	"upgrade-api/src/exec/backend/models"
	"upgrade-api/src/share/comm"
	"upgrade-api/src/share/database"
)

// BusinessController operations for Business
type BusinessController struct {
	CommonController
}

// URLMapping ...
func (c *BusinessController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post 新增业务
// @Title 新增业务
// @Description 新增业务
// @Param	body		body 	models.Business	true		"body for Business content"
// @Success 201 {int} models.Business
// @Failure 403 body is empty
// @router / [post]
// @author # yangzhao # 2020-06-08 18:00:05 #
func (c *BusinessController) Post() {
	ctx := GetBackendCntx()

	// 校验参数和获取请求数据
	req, code, err := c.getBusinessPost()
	if nil != err {
		logs.Error("Business parameter is invalid! errmsg:%s", err.Error())
		c.ErrorMessage(code, err.Error())
		return
	}

	logs.Info("Create business parameter. req:%v user:%s", req, c.mail)

	//创建业务和业务流控数据
	id, err := ctx.Model.CreateBusiness(ctx.Model.Quota, req, c.mail)
	if nil != err {
		logs.Error("Create business or business flow limit failed! errmsg:%s", err.Error())
		code, e := database.MysqlFormatError(err)
		c.ErrorMessage(code, e)
		return
	}

	c.FormatAddResp(http.StatusCreated, comm.OK, "Ok", id)
}

// 校验业务参数和获取请求数据
func (c *BusinessController) getBusinessPost() (*models.BusinessReq, int, error) {
	req := &models.BusinessReq{}

	if err := jsoniter.Unmarshal(c.Ctx.Input.RequestBody, req); nil != err {
		logs.Error("Unmarshal parameter failed! body:%s errmsg:%s",
			c.Ctx.Input.RequestBody, err.Error())
		return nil, comm.ERR_PARAM_INVALID, err
	} else if 0 == len(req.Name) {
		logs.Error("Business name is empty!")
		return nil, comm.ERR_PARAM_MISS, errors.New("Business name is invalid!")
	} else if 0 == len(req.Manager) {
		return nil, comm.ERR_PARAM_MISS, errors.New("Business manager is invalid!")
	}

	return req, 0, nil

}

// 查询单条业务
// @Title 查询单条业务
// @Description 根据id查询单条业务
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Business
// @Failure 403 :id is empty
// @router /:id([0-9]+) [get]
// @author # yangzhao # 2020-06-08 18:00:05 #
func (c *BusinessController) GetOne() {
	ctx := GetBackendCntx()

	// 获取参数
	idStr := c.Ctx.Input.Param(":id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if nil != err {
		logs.Error("Parse id failed! id:%s errmsg:%s", idStr, err.Error())
		c.ErrorMessage(comm.ERR_PARAM_INVALID, err.Error())
		return
	}

	// 查询业务信息和流控信息
	business, err := ctx.Model.GetBusinessAndPlat(id)
	if nil != err && orm.ErrNoRows != err {
		code, e := database.MysqlFormatError(err)
		logs.Error("Get business and plat failed! id:%d errmsg:%s", id, err.Error())
		c.ErrorMessage(code, e)
		return
	} else if orm.ErrNoRows == err {
		logs.Error("Not found business and plat! id:%d errmsg:%s", id, err.Error())
		c.ErrorMessage(comm.ERR_NOT_FOUND, "Business not exist!")
		return
	}

	c.FormatInterfaceResp(http.StatusOK, comm.OK, "Ok", business)
}

// 获取业务列表
// @Title 获取业务列表
// @Description 获取业务列表
// @Success 200 {object} models.Business
// @Failure 403
// @router /list [get]
// @author # yangzhao # 2020-06-08 18:00:05 #
func (c *BusinessController) GetAll() {
	ctx := GetBackendCntx()

	list, err := ctx.Model.GetBusinessList()
	if nil != err {
		logs.Error("Business list query failed! errmsg:%s", err.Error())
		code, e := database.MysqlFormatError(err)
		c.ErrorMessage(code, e)
		return
	}

	c.FormatInterfaceListResp(http.StatusOK, comm.OK, len(list), "Ok", list)
}

// 修改业务
// @Title 修改业务
// @Description 更加业务id修改业务
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body	body 	models.Business	true		"body for Business content"
// @Success 200 {object} models.Business
// @Failure 403 :id is not int
// @router /:id([0-9]+) [put]
// @author # yangzhao # 2020-06-08 18:00:05 #
func (c *BusinessController) Put() {
	ctx := GetBackendCntx()

	id, err := strconv.ParseInt(c.Ctx.Input.Param(":id"), 10, 64)
	if nil != err {
		logs.Error("Parse id failed! id:%s errmsg:%s", id, err.Error())
		c.ErrorMessage(comm.ERR_PARAM_INVALID, err.Error())
		return
	}

	// 解析请求数据和校验参数
	req, code, err := c.getBusinessPost()
	if nil != err {
		logs.Error("Business parameter is invalid! errmsg:%s", err.Error())
		c.ErrorMessage(code, err.Error())
		return
	}

	//更改业务或业务流控数据
	err = ctx.Model.UpdateBusiness(*ctx.Model.Quota, req, id, c.mail)
	if nil != err {
		logs.Error("Update business failed! id:%d req:%v errmsg:%s",
			id, req, err.Error())
		code, e := database.MysqlFormatError(err)
		c.ErrorMessage(code, e)
		return
	}

	c.FormatResp(http.StatusCreated, comm.OK, "Ok")
}

// 删除业务
// @Title 删除业务
// @Description 根据业务id删除业务和业务流控数据
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id([0-9]+) [delete]
// @author # yangzhao # 2020-06-08 18:00:05 #
func (c *BusinessController) Delete() {
	ctx := GetBackendCntx()

	id, err := strconv.ParseInt(c.Ctx.Input.Param(":id"), 10, 64)
	if nil != err {
		logs.Error("Parse id failed! id:%s errmsg:%s", id, err.Error())
		c.ErrorMessage(comm.ERR_PARAM_INVALID, err.Error())
		return
	}

	err = ctx.Model.DelBusinessById(ctx.Model.Quota, id)
	if nil != err {
		logs.Error("Del business failed! id:%d errmsg:%s", id, err.Error())
		code, e := database.MysqlFormatError(err)
		c.ErrorMessage(code, e)
		return
	}

	c.FormatResp(http.StatusCreated, comm.OK, "Ok")
}
