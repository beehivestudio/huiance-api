package models

import (
	"fmt"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"

	db_admin "upgrade-api/src/share/database/admin"
)

/* 修改用户角色请求数据 */
type UserRole struct {
	RoleIds []int `json:"role_ids" description:"角色ID集合"`
}

/******************************************************************************
 **函数名称: 更新用户角色
 **功    能: 更新用户角色
 **输入参数:
 **     userId: 用户id
 **     userRole: 用户请求数据
 **输出参数: NONE
 **返    回:
 **实现描述:
 **注意事项:
 **作    者: # Zhao.yang # 2020-06-15 18:12:31 #
 ******************************************************************************/
func UpdateUserRole(o orm.Ormer, userId int64, u *UserRole, id int64) error {

	// 开始事物
	err := o.Begin()
	if nil != err {
		logs.Error("%s() Begin transaction failed! errmsg:%s",
			utils.GetRunFuncName(), err.Error())
		return err
	}

	ids := u.RoleIds

	//查询用户的角色信息
	userRoles, err := db_admin.GetRoleByUserId(o, userId)
	if nil != err {
		logs.Error("%s() Get user role failed! errmsg:%s",
			utils.GetRunFuncName(), err.Error())
		o.Rollback()
		return err
	}

	//增加用户的角色
	if len(userRoles) == 0 {
		err := AddUserRole(o, ids, userId, id)
		if nil != err {
			o.Rollback()
			logs.Info("Add user role failed! errmsg:%s", err.Error())
			return err
		}
		o.Commit()
		return nil
	}

	//修改用户的角色
	err = db_admin.DelUserRoleRelByUserId(o, userId)
	if nil != err {
		logs.Info("Deleted user role failed! errmsg:%s", err.Error())
		o.Rollback()
		return err
	}

	err = AddUserRole(o, ids, userId, id)
	if nil != err {
		o.Rollback()
		logs.Info("Add user role failed! errmsg:%s", err.Error())
		return err
	}

	o.Commit()

	return nil
}

//增加用户角色
func AddUserRole(o orm.Ormer, roleIds []int, userId, id int64) error {

	for _, roleId := range roleIds {

		userRole := db_admin.UserRoleRel{
			UserId:       userId,
			RoleId:       int64(roleId),
			CreateTime:   time.Now(),
			CreateUserId: id,
		}

		id, err := db_admin.AddUserRoleRel(o, &userRole)
		if nil != err {
			logs.Error("%s() Add user role failed! errmsg:%s",
				utils.GetRunFuncName(), err.Error())
			return err
		}

		logs.Info("Insert user role id :%v", id)
	}

	return nil
}

/* 获取用户角色列表返回参数 */
type Roles struct {
	Id          int64  `json:"id"`          //角色Id
	Name        string `json:"name"`        //角色名称
	Description string `json:"description"` //角色备注信息
}

//获取用户角色列表
func GetUserRoles(o orm.Ormer, userId int64) (roles []Roles, err error) {

	sql := fmt.Sprintf(`select 
		r.id,r.name,
		r.description 
	from role r LEFT JOIN user_role_rel ur on ur.role_id= r.id and ur.user_id = ?`)

	_, err = o.Raw(sql, userId).QueryRows(&roles)
	if err != nil {
		logs.Error("%s() Get user roles list failed! errmsg:%s",
			utils.GetRunFuncName(), err.Error())
		return nil, err
	}

	return roles, nil
}
