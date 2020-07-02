package controllers

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/astaxie/beego/logs"
	jsoniter "github.com/json-iterator/go"

	"upgrade-api/src/exec/admin/models"
	"upgrade-api/src/share/comm"
	"upgrade-api/src/share/database"
	db_admin "upgrade-api/src/share/database/admin"
)

// UserController operations for User
type UserController struct {
	CommonController
}

// URLMapping ...
func (c *UserController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("Put", c.Put)
	c.Mapping("GetAll", c.GetAll)
}

/* 增加用户请求数据 */
type User struct {
	Name        string `json:"name" description:"用户名"`
	Email       string `json:"email" description:"邮箱地址"`
	Description string `json:"description" description:"描述"`
}

// 增加用户
// @Title 增加用户
// @Description 增加用户
// @Param	body		body 	models.User	true		"body for User content"
// @Success 201 {int} models.User
// @Failure 403 body is empty
// @router / [post]
// @author # Zhao.yang # 2020-06-17 14:01:08 #
func (c *UserController) Post() {
	ctx := GetAdminCntx()

	// 获取请求数据&检查参数
	req, code, err := c.getUserPost()
	if nil != err {
		logs.Error("%s() User add parameter format is invalid! errmsg:%s",
			utils.GetRunFuncName(), err.Error())
		c.ErrorMessage(code, err.Error())
		return
	}

	user := db_admin.User{
		Name:         req.Name,
		Email:        req.Email,
		Description:  req.Description,
		CreateTime:   time.Now(),
		CreateUserId: c.userId,
	}

	//增加用户
	id, err := db_admin.AddUser(ctx.Model.Mysql.O, &user)
	if nil != err {
		logs.Error("%s() User user is failed! user:%s errmsg:%s",
			utils.GetRunFuncName(), user, err.Error())
		code, e := database.MysqlFormatError(err)
		c.ErrorMessage(code, e)
		return
	}

	c.FormatAddResp(http.StatusCreated, comm.OK, "Ok", id)
}

//解析参数&验证参数
func (c *UserController) getUserPost() (u User, code int, err error) {

	if err := jsoniter.Unmarshal(c.Ctx.Input.RequestBody, &u); nil != err {
		logs.Error("%s() User parameter format is invalid! body:%s errmsg:%s",
			utils.GetRunFuncName(), c.Ctx.Input.RequestBody, err.Error())
		return u, comm.ERR_PARAM_INVALID, err
	}

	if 0 == len(u.Name) {
		return u, comm.ERR_PARAM_MISS, errors.New("user name is empty")
	}

	if 0 == len(u.Email) {
		return u, comm.ERR_PARAM_MISS, errors.New("user email is empty")
	}

	if 0 == len(u.Description) {
		return u, comm.ERR_PARAM_MISS, errors.New("user description is empty")
	}

	return u, 0, nil
}

/* 修改用户请求数据 */
type UserInfo struct {
	Name        string `json:"name" description:"用户名"`
	Email       string `json:"email" description:"邮箱地址"`
	Enable      int    `json:"enable" description:"邮箱地址"`
	Description string `json:"description" description:"描述"`
}

// 修改用户
// @Title 修改用户
// @Description 修改用户
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.User	true		"body for User content"
// @Success 200 {object} models.User
// @Failure 403 :id is not int
// @router /:id([0-9]+) [put]
// @author # Zhao.yang # 2020-06-17 14:01:08 #
func (c *UserController) Put() {
	ctx := GetAdminCntx()

	id, err := strconv.ParseInt(c.Ctx.Input.Param(":id"), 10, 64)
	if nil != err {
		logs.Error("%s() Get user id failed! errmsg:%s", utils.GetRunFuncName(), err.Error())
		c.ErrorMessage(comm.ERR_PARAM_INVALID, err.Error())
		return
	}

	//验证参数
	req, code, err := c.getUserPut()
	if nil != err {
		logs.Error("%s() User parameter format is invalid! errmsg:%s",
			utils.GetRunFuncName(), err.Error())
		c.ErrorMessage(code, err.Error())
		return
	}

	user := db_admin.User{
		Id:           id,
		Enable:       req.Enable,
		Description:  req.Description,
		UpdateTime:   time.Now(),
		UpdateUserId: c.userId,
	}

	logs.Info("Update user body msg :%v", user)

	//修改用户
	err = db_admin.UpdateUserById(ctx.Model.Mysql.O, &user)
	if nil != err {
		logs.Error("%s() Update user is failed! errmsg:%s",
			utils.GetRunFuncName(), err.Error())
		c.ErrorMessage(code, err.Error())
		return
	}

	c.FormatResp(http.StatusCreated, comm.OK, "Ok")
}

//获取修改用户参数&参数校验
func (c *UserController) getUserPut() (req *UserInfo, code int, err error) {

	if err := jsoniter.Unmarshal(c.Ctx.Input.RequestBody, &req); nil != err {
		logs.Error("%s() User parameter format is invalid! body:%s errmsg:%s",
			utils.GetRunFuncName(), c.Ctx.Input.RequestBody, err.Error())
		return nil, comm.ERR_PARAM_INVALID, err
	}

	return req, 0, nil

}

// 获取用户列表
// @Title 获取用户列表
// @Description 获取用户列表
// @Param	user_id			query	int64	true	"用户id"
// @Param	name			query	string	true	"用户名"
// @Param	email			query	string	true	"用户邮箱"
// @Success 200 {object} comm.InterfaceResp
// @Failure 403
// @router /list [get]
// @author # Zhao.yang # 2020-06-17 14:01:08 #
func (c *UserController) GetAll() {
	ctx := GetAdminCntx()

	fields, page, pageSize, err := c.getListParam()
	if nil != err {
		logs.Error("%s() Get parameter failed! errmsg:%s",
			utils.GetRunFuncName(), err.Error())
		c.ErrorMessage(comm.ERR_PARAM_INVALID, err.Error())
		return
	}

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
		list.List[u] = &db_admin.User{
			Id:           user.Id,
			Name:         user.Name,
			Email:        user.Email,
			Description:  user.Description,
			CreateTime:   user.CreateTime,
			UpdateTime:   user.UpdateTime,
			CreateUserId: user.CreateUserId,
			UpdateUserId: user.UpdateUserId,
		}
	}

	logs.Info("Get user list num msg :%v", len(list.List))

	c.FormatInterfaceResp(http.StatusOK, comm.OK, "Ok", list)
}

//获取用户参数列表
func (c *UserController) getListParam() (fields map[string]interface{},
	page, pageSize int, err error) {

	fields = make(map[string]interface{})
	page = comm.PAGE_START
	pageSize = comm.PAGE_SIZE

	id := c.GetString("user_id")
	if "" != id {
		fields["id"] = id
	}

	name := c.GetString("name")
	if "" != name {
		fields["name"] = name
	}

	email := c.GetString("email")
	if "" != email {
		fields["email"] = email
	}

	return fields, page, pageSize, nil
}
