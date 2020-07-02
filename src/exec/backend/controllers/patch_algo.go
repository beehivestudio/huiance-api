package controllers

import (
	"net/http"
	"strconv"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/json-iterator/go"

	"upgrade-api/src/share/comm"
	"upgrade-api/src/share/database"
	"upgrade-api/src/share/database/upgrade"
)

// PatchAlgoControleer operations for ApkPatchAlgo
type PatchAlgoControleer struct {
	CommonController
}

// URLMapping ...
func (c *PatchAlgoControleer) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Post
// @Description create ApkPatchAlgo
// @Param	body		body 	models.ApkPatchAlgo	true		"body for ApkPatchAlgo content"
// @Success 201 {int} models.ApkPatchAlgo
// @Failure 403 body is empty
// @router / [post]
// @author # Linlin.guo # 2020-06-09 15:00:05 #
func (c *PatchAlgoControleer) Post() {
	ctx := GetBackendCntx()

	// 校验参数
	var v upgrade.ApkPatchAlgo
	var err error
	if err = jsoniter.Unmarshal(c.Ctx.Input.RequestBody, &v); nil != err {
		logs.Error("Unmarshal parameter failed! body:%s errmsg:%s",
			c.Ctx.Input.RequestBody, err.Error())
		c.ErrorMessage(comm.ERR_PARAM_INVALID, err.Error())
		return
	} else if 0 == len(v.Name) || 0 == len(v.Description) {
		logs.Error("Parameter is invalid! body:%s", c.Ctx.Input.RequestBody)
		c.ErrorMessage(comm.ERR_PARAM_MISS, "Parameter is invalid!")
		return
	}

	var id int64
	if id, err = upgrade.AddApkPatchAlgo(ctx.Model.Mysql.O, &v); nil != err {
		logs.Error("Add patch algo failed! errmsg:%s", err.Error())
		code, e := database.MysqlFormatError(err)
		c.ErrorMessage(code, e)
		return
	}

	c.FormatAddResp(http.StatusOK, comm.OK, "Ok", id)
}

// Put ...
// @Title Put
// @Description update the ApkPatchAlgo
// @Param	id		path	string				true	"The id you want to update"
// @Param	body	body	models.ApkPatchAlgo	true	"body for ApkPatchAlgo content"
// @Success 200 {object} models.ApkPatchAlgo
// @Failure 403 :id is not int
// @router /:id([0-9]+) [put]
// @author # Linlin.guo # 2020-06-09 15:00:05 #
func (c *PatchAlgoControleer) Put() {
	ctx := GetBackendCntx()

	// 校验参数
	idStr := c.Ctx.Input.Param(":id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if nil != err {
		logs.Error("Get id failed! errmsg:%s", err.Error())
		c.ErrorMessage(comm.ERR_PARAM_INVALID, err.Error())
		return
	}

	v := upgrade.ApkPatchAlgo{Id: id}

	if err := jsoniter.Unmarshal(c.Ctx.Input.RequestBody, &v); nil != err {
		logs.Error("Unmarshal parameter failed! id:%d body:%s errmsg:%s",
			id, c.Ctx.Input.RequestBody, err.Error())
		c.ErrorMessage(comm.ERR_PARAM_INVALID, err.Error())
		return
	} else if 0 == len(v.Name) || 0 == len(v.Description) {
		logs.Error("Paramter is invalid! id:%d body:%s", id, c.Ctx.Input.RequestBody)
		c.ErrorMessage(comm.ERR_PARAM_MISS, "Parameter is invalid!")
		return
	}

	if err := upgrade.UpdateApkPatchAlgoById(ctx.Model.Mysql.O, &v); nil != err {
		if orm.ErrNoRows == err {
			logs.Error("Patch not exist! id:%d", id)
			c.ErrorMessage(comm.ERR_PARAM_INVALID, "Patch not exist!")
			return
		}
		logs.Error("Update patch algo failed! id:%d errmsg:%s", id, err.Error())
		code, e := database.MysqlFormatError(err)
		c.ErrorMessage(code, e)
		return
	}

	c.ErrorMessage(comm.OK, "Ok")
}

////////////////////////////////////////////////////////////////////////////////

/* ApkPatchAlog信息 */
type ApkPatchAlog struct {
	Id          int64  `json:"id"`          // 算法ID
	Name        string `json:"name"`        // 算法名称
	Enable      int    `json:"enable"`      // 是否启用
	Description string `json:"description"` // 描述信息
	CreateTime  string `json:"create_time"` // 创建时间
	UpdateTime  string `json:"update_time"` // 更新时间
}

// GetOne ...
// @Title Get One
// @Description get ApkPatchAlgo by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.ApkPatchAlgo
// @Failure 403 :id is empty
// @router /:id([0-9]+) [get]
// @author # Linlin.guo # 2020-06-09 15:00:05 #
func (c *PatchAlgoControleer) GetOne() {
	ctx := GetBackendCntx()

	// 校验参数
	idStr := c.Ctx.Input.Param(":id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if nil != err {
		logs.Error("Parse id failed! id:%s errmsg:%s", idStr, err.Error())
		c.ErrorMessage(comm.ERR_PARAM_INVALID, err.Error())
		return
	}

	// 获取数据
	v, err := upgrade.GetApkPatchAlgoById(ctx.Model.Mysql.O, id)
	if nil != err {
		if orm.ErrNoRows == err {
			logs.Error("Patch not exist! id:%d", id)
			c.ErrorMessage(comm.ERR_PARAM_INVALID, "Patch not exist!")
			return
		}
		logs.Error("Get patch algo failed! id:%d errmsg:%s", id, err.Error())
		code, e := database.MysqlFormatError(err)
		c.ErrorMessage(code, e)
		return
	}

	// 返回结果
	apkPatchAlog := &ApkPatchAlog{
		Id:          v.Id,
		Name:        v.Name,
		Enable:      v.Enable,
		Description: v.Description,
		CreateTime:  v.CreateTime.Format("2006-01-02 15:04:05"),
		UpdateTime:  v.CreateTime.Format("2006-01-02 15:04:05"),
	}

	c.FormatInterfaceResp(http.StatusOK, comm.OK, "Ok", apkPatchAlog)
}

// GetAll ...
// @Title Get All
// @Description get ApkPatchAlgo
// @Success 200 {object} models.ApkPatchAlgo
// @Failure 403
// @router /list [get]
// @author # Linlin.guo # 2020-06-09 15:00:05 #
func (c *PatchAlgoControleer) GetAll() {
	ctx := GetBackendCntx()

	// 获取数据
	patchAlgoList, err := upgrade.GetPatchAlgoList(ctx.Model.Mysql.O, 0)
	if nil != err {
		logs.Error("Get patch algo list failed! errmsg:%s", err.Error())
		code, e := database.MysqlFormatError(err)
		c.ErrorMessage(code, e)
		return
	}

	// 返回结果
	listData := make([]interface{}, len(patchAlgoList))
	for k, patchAlgo := range patchAlgoList {
		listData[k] = &ApkPatchAlog{
			Id:          patchAlgo.Id,
			Name:        patchAlgo.Name,
			Enable:      patchAlgo.Enable,
			Description: patchAlgo.Description,
			CreateTime:  patchAlgo.CreateTime.Format("2006-01-02 15:04:05"),
			UpdateTime:  patchAlgo.CreateTime.Format("2006-01-02 15:04:05"),
		}
	}

	c.FormatInterfaceListResp(http.StatusOK, comm.OK, len(patchAlgoList), "Ok", listData)
}

// Delete ...
// @Title Delete
// @Description delete the ApkPatchAlgo
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id([0-9]+) [delete]
// @author # Linlin.guo # 2020-06-09 15:00:05 #
func (c *PatchAlgoControleer) Delete() {
	ctx := GetBackendCntx()

	// 校验参数
	idStr := c.Ctx.Input.Param(":id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if nil != err {
		logs.Error("Parse id failed! id:%s errmsg:%s", idStr, err.Error())
		c.ErrorMessage(comm.ERR_PARAM_INVALID, err.Error())
		return
	}

	// 执行操作
	if err := upgrade.DeleteApkPatchAlgo(ctx.Model.Mysql.O, id); nil != err {
		if orm.ErrNoRows == err {
			logs.Error("Patch not exist! id:%d", id)
			c.ErrorMessage(comm.ERR_PARAM_INVALID, "Patch not exist!")
			return
		}
		logs.Error("Delete patch algo failed! id:%d errmsg:%s", id, err.Error())
		code, e := database.MysqlFormatError(err)
		c.ErrorMessage(code, e)
		return
	}

	c.FormatResp(http.StatusOK, comm.OK, "Ok")
}
