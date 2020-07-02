package upgrade

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type Business struct {
	Id           int64     `orm:"column(id);auto" description:"ID"`
	Name         string    `orm:"column(name);size(128)" description:"业务名称"`
	Key          string    `orm:"column(key);size(64)" description:"业务秘钥；数据加解密秘钥"`
	Enable       int       `orm:"column(enable)" description:"是否启用"`
	HasFlowLimit int       `orm:"column(has_flow_limit)" description:"是否开启流控；有无流控；0：无流控；1：有流控"`
	Manager      string    `orm:"column(manager)" description:"项目技术负责人"`
	Description  string    `orm:"column(description)" description:"描述信息"`
	CreateTime   time.Time `orm:"column(create_time);type(timestamp);auto_now_add" description:"创建时间"`
	CreateUser   string    `orm:"column(create_user)" description:"页面创建者"`
	UpdateTime   time.Time `orm:"column(update_time);type(timestamp);auto_now" description:"修改时间"`
	UpdateUser   string    `orm:"column(update_user)" description:"页面修改者"`
}

func (t *Business) TableName() string {
	return "business"
}

func init() {
	orm.RegisterModel(new(Business))
}

// AddBusiness insert a new Business into database and returns
// last inserted Id on success.
func AddBusiness(o orm.Ormer, m *Business) (id int64, err error) {
	id, err = o.Insert(m)
	return
}

// GetBusinessById retrieves Business by Id. Returns error if
// Id doesn't exist
func GetBusinessById(o orm.Ormer, id int64) (v *Business, err error) {
	v = &Business{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllBusiness retrieves all Business matches certain condition. Returns empty list if
// no records exist
func GetAllBusiness(o orm.Ormer, query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	qs := o.QueryTable(new(Business))
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

	var l []Business
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

// UpdateBusiness updates Business by Id and returns error if
// the record to be updated doesn't exist
func UpdateBusinessById(o orm.Ormer, m *Business) (err error) {
	v := Business{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteBusiness deletes Business by Id and returns error if
// the record to be deleted doesn't exist
func DeleteBusiness(o orm.Ormer, id int64) (err error) {
	v := Business{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Business{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

//获取业务列表
func GetBusinessList(o orm.Ormer) (business []Business, err error) {
	if _, err = o.QueryTable(new(Business)).All(&business); err == nil {
		return business, nil
	}
	return nil, err
}
