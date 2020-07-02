package models

import (
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"

	"upgrade-api/src/share/database/admin"
)

/* 修改角色页面请求数据 */
type RolePage struct {
	PageIds []int `json:"page_ids" description:"页面ID集合"`
}

/******************************************************************************
 **函数名称: UpdateRolePage
 **功    能: 更新角色页面权限
 **输入参数:
 **      o: orm.Ormer
 **     roleId: 角色id
 **     rolePage: 页面集合
 **输出参数: NONE
 **返    回:
 **      err: 错误信息
 **实现描述:
 **注意事项:
 **作    者: # Linlin.guo # 2020-06-18 15:00:05 #
 ******************************************************************************/
func UpdateRolePage(o orm.Ormer, roleId, userId int64, rolePage RolePage) error {

	//开始事物
	err := o.Begin()
	if nil != err {
		logs.Error("%s() Begin transaction failed! errmsg:%s",
			utils.GetRunFuncName(), err.Error())
		return err
	}

	//如报错，回滚，返回错误信息
	defer func() {
		if nil != err {
			logs.Error("%s() Transaction happen rollback! errmsg:%s",
				utils.GetRunFuncName(), err.Error())
			_ = o.Rollback()
		}
	}()

	ids := rolePage.PageIds
	//查询角色的页面信息
	rolePages, err := db_admin.GetAllPageByRoleId(o, roleId)
	if nil != err {
		logs.Error("%s() Get role page failed! errmsg:%s", utils.GetRunFuncName(), err.Error())
		return err
	}

	//增加用户的角色
	if len(rolePages) > 0 {
		//删除用户已有角色
		err = db_admin.DelRolePageRelByRoleId(o, roleId)
		if nil != err {
			logs.Info("Deleted role page failed! errmsg:%s", err.Error())
			return err
		}
	}

	err = AddRolePage(o, ids, roleId, userId)
	if nil != err {
		logs.Info("Add role page failed! errmsg:%s", err.Error())
		return err
	}

	//提交事务
	err = o.Commit()
	if nil != err {
		logs.Error("%s() Commit transaction failed! errmsg:%s",
			utils.GetRunFuncName(), err.Error())
		return err
	}

	return nil
}

/******************************************************************************
 **函数名称: AddRolePage
 **功    能: 批量增加角色页面权限
 **输入参数:
 **      o: orm.Ormer
 **     pageIds: 页面集合
 **     roleId: 角色id
 **输出参数: NONE
 **返    回:
 **      err: 错误信息
 **实现描述:
 **注意事项:
 **作    者: # Linlin.guo # 2020-06-18 15:00:05 #
 ******************************************************************************/
func AddRolePage(o orm.Ormer, pageIds []int, roleId, userId int64) error {

	for _, pageId := range pageIds {

		rolePage := db_admin.RolePageRel{
			RoleId:       roleId,
			PageId:       int64(pageId),
			CreateTime:   time.Now(),
			CreateUserId: userId,
		}
		//todo 当单个页面不存在时，是忽略单个还是全部？
		_, err := db_admin.GetPageById(o, int64(pageId))
		if orm.ErrNoRows == err {
			logs.Info("AddRolePage page not found! pageId: %d", pageId)
			continue
		}

		_, err = db_admin.AddRolePageRel(o, &rolePage)
		if nil != err {
			logs.Error("%s() Add role page failed! errmsg:%s",
				utils.GetRunFuncName(), err.Error())
			return err
		}
	}

	return nil
}

type ReturnParam struct {
	db_admin.Role
	Pages []db_admin.Page `json:"pages"`
}

/******************************************************************************
 **函数名称: GetRolePageList
 **功    能: 模糊查询角色页面列表
 **输入参数:
 **      o: orm.Ormer
 **     pageIds: 页面集合
 **     roleId: 角色id
 **输出参数: NONE
 **返    回:
 **      err: 错误信息
 **实现描述:
 **注意事项:
 **作    者: # Linlin.guo # 2020-06-18 15:00:05 #
 ******************************************************************************/
func GetRolePageList(o orm.Ormer, roleId int64, roleName, pagesName string) (returnParam []ReturnParam, err error) {

	roleList, err := db_admin.GetRoleList(o, roleId, roleName)
	if err != nil {
		logs.Error("%s() GetRoleList failed! roleId:%d errmsg:%s",
			utils.GetRunFuncName(), roleId, err.Error())
		return returnParam, err
	}

	for _, role := range roleList {
		var tmpParam ReturnParam
		tmpParam.Role = role
		pageList, err := db_admin.GetRolePages(o, role.Id, pagesName)
		if err != nil {
			logs.Error("%s() GetRolePages failed! roleId:%d errmsg:%s",
				utils.GetRunFuncName(), roleId, err.Error())
			return returnParam, err
		}
		if len(pageList) > 0 {
			tmpParam.Pages = pageList
			returnParam = append(returnParam, tmpParam)
		}

	}
	return returnParam, err

}
