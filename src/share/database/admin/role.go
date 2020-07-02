package db_admin

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

type Role struct {
	Id           int64     `orm:"column(id);auto" description:"角色ID"`
	Name         string    `orm:"column(name);size(128)" description:"角色名称"`
	Description  string    `orm:"column(description);size(1024);null" description:"描述信息"`
	CreateTime   time.Time `orm:"column(create_time);type(timestamp);auto_now_add" description:"创建时间"`
	CreateUserId int64     `orm:"column(create_user_id)" description:"页面创建者"`
	UpdateTime   time.Time `orm:"column(update_time);type(timestamp);auto_now" description:"修改时间"`
	UpdateUserId int64     `orm:"column(update_user_id)" description:"页面修改者"`
}

func (t *Role) TableName() string {
	return "role"
}

func init() {
	orm.RegisterModel(new(Role))
}

// AddRole insert a new Role into database and returns
// last inserted Id on success.
func AddRole(o orm.Ormer, m *Role) (id int64, err error) {
	id, err = o.Insert(m)
	return
}

// GetRoleById retrieves Role by Id. Returns error if
// Id doesn't exist
func GetRoleById(o orm.Ormer, id int64) (v *Role, err error) {
	v = &Role{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllRole retrieves all Role matches certain condition. Returns empty list if
// no records exist
func GetAllRole(o orm.Ormer, query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	qs := o.QueryTable(new(Role))
	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		if strings.Contains(k, "isnull") {
			qs = qs.Filter(k, (v == "true" || v == "1"))
		} else {
			qs = qs.Filter(k, v)
		}
	}
	// order by:
	var sortFields []string
	if len(sortby) != 0 {
		if len(sortby) == len(order) {
			// 1) for each sort field, there is an associated order
			for i, v := range sortby {
				orderby := ""
				if order[i] == "desc" {
					orderby = "-" + v
				} else if order[i] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
			qs = qs.OrderBy(sortFields...)
		} else if len(sortby) != len(order) && len(order) == 1 {
			// 2) there is exactly one order, all the sorted fields will be sorted by this order
			for _, v := range sortby {
				orderby := ""
				if order[0] == "desc" {
					orderby = "-" + v
				} else if order[0] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return nil, errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return nil, errors.New("Error: unused 'order' fields")
		}
	}

	var l []Role
	qs = qs.OrderBy(sortFields...)
	if _, err = qs.Limit(limit, offset).All(&l, fields...); err == nil {
		if len(fields) == 0 {
			for _, v := range l {
				ml = append(ml, v)
			}
		} else {
			// trim unused fields
			for _, v := range l {
				m := make(map[string]interface{})
				val := reflect.ValueOf(v)
				for _, fname := range fields {
					m[fname] = val.FieldByName(fname).Interface()
				}
				ml = append(ml, m)
			}
		}
		return ml, nil
	}
	return nil, err
}

// UpdateRole updates Role by Id and returns error if
// the record to be updated doesn't exist
func UpdateRoleById(o orm.Ormer, m *Role) (err error) {
	v := Role{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteRole deletes Role by Id and returns error if
// the record to be deleted doesn't exist
func DeleteRole(o orm.Ormer, id int64) (err error) {
	v := Role{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Role{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

/******************************************************************************
 **函数名称: GetRoleList
 **功    能: 查询角色列表（支持角色名称模糊查询）
 **输入参数:
 **      o: orm.Ormer
 **      id: 角色Id
 **      name: 角色名称
 **输出参数: NONE
 **返    回:
 **      roleList: 角色列表
 **      err: 错误信息
 **实现描述:
 **注意事项:
 **作    者: # Linlin.guo # 2020-6-17 15:15:42 #
 ******************************************************************************/
func GetRoleList(o orm.Ormer, id int64, name string) (roleList []Role, err error) {

	qs := o.QueryTable(new(Role))
	if id != 0 {
		qs = qs.Filter("id", id)
	}
	if len(name) > 0 {
		qs = qs.Filter("name__contains", name)
	}

	_, err = qs.All(&roleList)
	if err != nil {
		logs.Error("GetRoleList err: %s", err.Error())
		return nil, err
	}

	return roleList, nil
}
