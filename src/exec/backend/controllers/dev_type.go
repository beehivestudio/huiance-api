package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/json-iterator/go"

	"upgrade-api/src/share/comm"
	"upgrade-api/src/share/database"
	"upgrade-api/src/share/database/upgrade"
)

// DevTypeController operations for DevType
type DevTypeController struct {
	CommonController
}

// URLMapping ...
func (c *DevTypeController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

////////////////////////////////////////////////////////////////////////////////
// 新增设备种类

/* 新增设备种类请求 */
type CreateDevTypeReq struct {
	Name        string `json:"name"`        // 应用名
	Description string `json:"description"` // 描述信息
}

// Post ...
// @Title Post
// @Description create DevType
// @Param	body		body 	models.DevType	true		"body for DevType content"
// @Success 201 {int} models.DevType
// @Failure 403 body is empty
// @router / [post]
// @author # Linlin.guo # 2020-06-09 15:00:05 #
func (c *DevTypeController) Post() {
	ctx := GetBackendCntx()

	// 获取请求参数
	req := &CreateDevTypeReq{}
	var err error

	if err = jsoniter.Unmarshal(c.Ctx.Input.RequestBody, req); nil != err {
		logs.Error("Unmarshal request failed! errmsg:%s", err.Error())
		c.ErrorMessage(comm.ERR_PARAM_INVALID, err.Error())
		return
	} else if 0 == len(req.Name) { //校验参数
		logs.Error("Dev name is invalid!")
		c.ErrorMessage(comm.ERR_PARAM_MISS, "Dev name is invalid!")
		return
	}

	// 创建设备种类
	v := &upgrade.DevType{
		Name:        req.Name,        // 设备种类
		Description: req.Description, // 描述信息
		CreateTime:  time.Now(),      // 创建时间
		UpdateTime:  time.Now(),      // 更新时间
	}

	id := int64(0)
	if id, err = upgrade.AddDevType(ctx.Model.Mysql.O, v); nil != err {
		logs.Error("Create dev type failed! errmsg:%s", err.Error())
		code, e := database.MysqlFormatError(err)
		c.ErrorMessage(code, e)
		return
	}

	c.FormatAddResp(http.StatusOK, comm.OK, "Ok", id)
}

////////////////////////////////////////////////////////////////////////////////
// 更新设备种类

/* 更新设备种类请求 */
type UpdateDevTypeReq struct {
	Name        string `json:"name"`        // 应用名
	Description string `json:"description"` // 描述信息
}

// Put ...
// @Title Put
// @Description update the DevType
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.DevType	true		"body for DevType content"
// @Success 200 {object} models.DevType
// @Failure 403 :id is not int
// @router /:id([0-9]+) [put]
// @author # Linlin.guo # 2020-06-09 15:00:05 #
func (c *DevTypeController) Put() {
	ctx := GetBackendCntx()

	//校验参数
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	req := &UpdateDevTypeReq{}
	if err := jsoniter.Unmarshal(c.Ctx.Input.RequestBody, req); nil != err {
		logs.Error("Unmarshal request failed! errmsg:%s", err.Error())
		c.ErrorMessage(comm.ERR_PARAM_INVALID, err.Error())
		return
	} else if 0 == len(req.Name) { //校验参数
		logs.Error("Dev name is invalid!")
		c.ErrorMessage(comm.ERR_PARAM_MISS, "Dev name is invalid!")
		return
	}

	// 更新设备种类
	v := upgrade.DevType{
		Id:          id,              // 种类ID
		Name:        req.Name,        // 种类名称
		Description: req.Description, // 描述信息
		UpdateTime:  time.Now(),      // 更新时间
	}

	if err := upgrade.UpdateDevTypeById(ctx.Model.Mysql.O, &v); nil != err {
		if orm.ErrNoRows == err {
			logs.Error("DevType not exist! id:%d", id)
			c.ErrorMessage(comm.ERR_PARAM_INVALID, "DevType not exist!")
			return
		}
		logs.Error("Update dev type failed! id:%d errmsg:%s", id, err.Error())
		code, e := database.MysqlFormatError(err)
		c.ErrorMessage(code, e)
		return
	}

	c.ErrorMessage(comm.OK, "Ok")
}

////////////////////////////////////////////////////////////////////////////////

/* 设备种类信息 */
type DevTypeItem struct {
	Id          int64  `json:"id"`          // 设备类型ID
	Name        string `json:"name"`        // 设备类型名称
	Description string `json:"description"` // 描述信息
	CreateTime  string `json:"create_time"` // 创建时间
	UpdateTime  string `json:"update_time"` // 更新时间
}

// GetOne ...
// @Title Get One
// @Description get DevType by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.DevType
// @Failure 403 :id is empty
// @router /:id([0-9]+) [get]
// @author # Linlin.guo # 2020-06-09 15:00:05 #
func (c *DevTypeController) GetOne() {

	ctx := GetBackendCntx()

	// 获取参数信息
	idStr := c.Ctx.Input.Param(":id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if nil != err {
		logs.Error("Parse id failed! id:%s errmsg:%s", idStr, err.Error())
		c.ErrorMessage(comm.ERR_PARAM_INVALID, err.Error())
		return
	}

	// 查询设备种类
	v, err := upgrade.GetDevTypeById(ctx.Model.Mysql.O, id)
	if nil != err {
		if orm.ErrNoRows == err {
			logs.Error("DevType not exist! id:%d", id)
			c.ErrorMessage(comm.ERR_PARAM_INVALID, "DevType not exist!")
			return
		}
		logs.Error("Get dev type by id failed! id:%d errmsg:%s", id, err.Error())
		code, e := database.MysqlFormatError(err)
		c.ErrorMessage(code, e)
		return
	}

	// 返回查询结果
	devType := &DevTypeItem{
		Id:          v.Id,
		Name:        v.Name,
		Description: v.Description,
		CreateTime:  v.CreateTime.Format("2006-01-02 15:04:05"),
		UpdateTime:  v.CreateTime.Format("2006-01-02 15:04:05"),
	}

	c.FormatInterfaceResp(http.StatusOK, comm.OK, "Ok", devType)
}

// GetAll ...
// @Title Get All
// @Description get DevType
// @Success 200 {object} models.DevType
// @Failure 403
// @router /list [get]
// @author # Linlin.guo # 2020-06-09 15:00:05 #
func (c *DevTypeController) GetAll() {
	ctx := GetBackendCntx()

	// 获取设备种类列表
	list, err := upgrade.GetDevTypeList(ctx.Model.Mysql.O)
	if nil != err {
		logs.Error("Get dev type list failed! errmsg:%s", err.Error())
		code, e := database.MysqlFormatError(err)
		c.ErrorMessage(code, e)
		return
	}

	// 返回查询结果
	listData := make([]interface{}, len(list))
	for k, devType := range list {
		listData[k] = &DevTypeItem{
			Id:          devType.Id,
			Name:        devType.Name,
			Description: devType.Description,
			CreateTime:  devType.CreateTime.Format("2006-01-02 15:04:05"),
			UpdateTime:  devType.CreateTime.Format("2006-01-02 15:04:05"),
		}
	}

	c.FormatInterfaceListResp(http.StatusOK, comm.OK, len(list), "Ok", listData)
}

// Delete ...
// @Title Delete
// @Description delete the DevType
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id([0-9]+) [delete]
// @author # Linlin.guo # 2020-06-09 15:00:05 #
func (c *DevTypeController) Delete() {
	ctx := GetBackendCntx()

	// 校验参数合法性
	idStr := c.Ctx.Input.Param(":id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if nil != err {
		logs.Error("Parse id failed! id:%s errmsg:%s", idStr, err.Error())
		c.ErrorMessage(comm.ERR_PARAM_INVALID, err.Error())
		return
	}

	// 执行删除操作
	if err := upgrade.DeleteDevType(ctx.Model.Mysql.O, id); nil != err {
		if orm.ErrNoRows == err {
			logs.Error("DevType not exist! id:%d", id)
			c.ErrorMessage(comm.ERR_PARAM_INVALID, "DevType not exist!")
			return
		}
		logs.Error("Delete dev type failed! id:%d errmsg:%s", id, err.Error())
		code, e := database.MysqlFormatError(err)
		c.ErrorMessage(code, e)
		return
	}

	c.ErrorMessage(comm.OK, "Ok")
}
