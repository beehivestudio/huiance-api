package db_admin

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type ActionLog struct {
	Id         int64     `orm:"column(id);auto" description:"ID"`
	UserId     int64     `orm:"column(user_id)" description:"用户ID；外键：user.id"`
	Action     int64     `orm:"column(action)" description:"行为ID；用户行为ID取值：；0:未知行为；1：登录；2：退出；3：添加；4：修改；5：删除；6：查询（不用）"`
	Detail     string    `orm:"column(detail);size(1024);null" description:"行为详情"`
	CreateTime time.Time `orm:"column(create_time);type(timestamp);auto_now_add" description:"创建时间"`
}

func (t *ActionLog) TableName() string {
	return "action_log"
}

func init() {
	orm.RegisterModel(new(ActionLog))
}

// AddActionLog insert a new ActionLog into database and returns
// last inserted Id on success.
func AddActionLog(o orm.Ormer, m *ActionLog) (id int64, err error) {
	id, err = o.Insert(m)
	return
}

// GetActionLogById retrieves ActionLog by Id. Returns error if
// Id doesn't exist
func GetActionLogById(o orm.Ormer, id int64) (v *ActionLog, err error) {
	v = &ActionLog{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllActionLog retrieves all ActionLog matches certain condition. Returns empty list if
// no records exist
func GetAllActionLog(o orm.Ormer, query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	qs := o.QueryTable(new(ActionLog))
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

	var l []ActionLog
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

// UpdateActionLog updates ActionLog by Id and returns error if
// the record to be updated doesn't exist
func UpdateActionLogById(o orm.Ormer, m *ActionLog) (err error) {
	v := ActionLog{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteActionLog deletes ActionLog by Id and returns error if
// the record to be deleted doesn't exist
func DeleteActionLog(o orm.Ormer, id int64) (err error) {
	v := ActionLog{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&ActionLog{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
