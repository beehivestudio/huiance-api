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

type UserRoleRel struct {
	Id           int64     `orm:"column(id);auto" description:"ID"`
	UserId       int64     `orm:"column(user_id)" description:"用户ID；外键：user.id"`
	RoleId       int64     `orm:"column(role_id)" description:"角色ID；外键：role.id"`
	CreateTime   time.Time `orm:"column(create_time);type(timestamp);auto_now_add" description:"创建时间"`
	CreateUserId int64     `orm:"column(create_user_id)" description:"页面创建者"`
	UpdateTime   time.Time `orm:"column(update_time);type(timestamp);auto_now" description:"修改时间"`
	UpdateUserId int64     `orm:"column(update_user_id)" description:"页面修改者"`
}

func (t *UserRoleRel) TableName() string {
	return "user_role_rel"
}

func init() {
	orm.RegisterModel(new(UserRoleRel))
}

// AddUserRoleRel insert a new UserRoleRel into database and returns
// last inserted Id on success.
func AddUserRoleRel(o orm.Ormer, m *UserRoleRel) (id int64, err error) {
	id, err = o.Insert(m)
	return
}

// GetUserRoleRelById retrieves UserRoleRel by Id. Returns error if
// Id doesn't exist
func GetUserRoleRelById(o orm.Ormer, id int64) (v *UserRoleRel, err error) {
	v = &UserRoleRel{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllUserRoleRel retrieves all UserRoleRel matches certain condition. Returns empty list if
// no records exist
func GetAllUserRoleRel(o orm.Ormer, query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	qs := o.QueryTable(new(UserRoleRel))
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

	var l []UserRoleRel
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

// UpdateUserRoleRel updates UserRoleRel by Id and returns error if
// the record to be updated doesn't exist
func UpdateUserRoleRelById(o orm.Ormer, m *UserRoleRel) (err error) {
	v := UserRoleRel{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteUserRoleRel deletes UserRoleRel by Id and returns error if
// the record to be deleted doesn't exist
func DeleteUserRoleRel(o orm.Ormer, id int64) (err error) {
	v := UserRoleRel{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&UserRoleRel{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

//根据用户id查询角色信息
func GetRoleByUserId(o orm.Ormer, userId int64) (user []UserRoleRel, err error) {
	if _, err := o.QueryTable(new(UserRoleRel)).Filter("user_id", userId).All(&user); nil == err {
		return user, nil
	}
	return
}

//删除用户角色信息
func DelUserRoleRelByUserId(o orm.Ormer, userId int64) (err error) {
	if _, err := o.QueryTable(new(UserRoleRel)).Filter("user_id", userId).Delete(); nil == err {
		return nil
	}
	return
}

//根据用户Id判断是否是超级管理员
func ChargeAuthByUserId(o orm.Ormer, userId int64) (role Role, err error) {

	sql := fmt.Sprintf(`SELECT
	r.id,
	r.NAME,
	r.description 
     FROM
	role r
	LEFT JOIN user_role_rel ur ON ur.role_id = r.id 
    WHERE
	r.NAME = "超级管理员" 
	AND ur.user_id = ?`)

	err = o.Raw(sql, userId).QueryRow(&role)
	if err != nil {
		logs.Error("ChargeAuthByUserId failed! errmsg:%s", err.Error())
		return role, err
	}
	return role, err
}
