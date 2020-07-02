package models

import (
	"github.com/astaxie/beego/orm"

	db_admin "upgrade-api/src/share/database/admin"
)

/******************************************************************************
 **函数名称: GetUserList
 **功    能: 获取用户列表
 **输入参数:
 **     fields: 过滤条件
 **输出参数: NONE
 **返    回:
 **实现描述:
 **注意事项:
 **作    者: # Zhao.yang # 2020-06-17 14:14:35 #
 ******************************************************************************/
func GetUserList(o orm.Ormer, fields map[string]interface{}, page, pageSize int) (int64, []*db_admin.User, error) {

	user := new([]*db_admin.User)

	cond := orm.NewCondition()

	for k, v := range fields {
		cond = cond.And(k, v)
	}

	qs := o.QueryTable(new(db_admin.User)).SetCond(cond)

	// 查询总量
	total, err := qs.Count()
	if nil != err {
		return 0, nil, err
	}

	_, err = qs.OrderBy("-id").Limit(pageSize, (page-1)*pageSize).All(user)
	if nil != err {
		return 0, nil, err
	}

	return total, *user, err
}
