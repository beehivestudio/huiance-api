package db_admin

import (
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

type Action struct {
	Id          int64     `orm:"column(id);auto" description:"ID"`
	Action      string    `orm:"column(action)" description:"角色ID"`
	Description string    `orm:"column(description)" description:"描述"`
	CreateTime  time.Time `orm:"column(create_time);type(timestamp);" description:"创建时间"`
}

func (t *Action) TableName() string {
	return "action"
}

func init() {
	orm.RegisterModel(new(Action))
}

/******************************************************************************
 **函数名称: GetActionList
 **功    能: 获取行为列表
 **输入参数:
 **      o: orm.Ormer
 **输出参数: NONE
 **返    回:
 **      list: 行为列表
 **      err: 错误信息
 **实现描述:
 **注意事项:
 **作    者: # Zhao.yang # 2020-06-17 14:14:35 #
 ******************************************************************************/
func GetActionList(o orm.Ormer) (list []Action, err error) {

	_, err = o.QueryTable(new(Action)).All(&list)
	if nil != err {
		logs.Info("Get action list failed! errmsg :%s", err.Error())
		return list, err
	}

	return list, nil
}
