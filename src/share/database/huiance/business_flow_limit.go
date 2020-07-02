package upgrade

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type BusinessFlowLimit struct {
	Id         int64     `orm:"column(id);auto" description:"ID"`
	BusinessId int64     `orm:"column(business_id)" description:"业务ID；外键：business.id"`
	BeginTime  int       `orm:"column(begin_time)" description:"开始时间；开始时间，取值范围[0, 23]"`
	EndTime    int       `orm:"column(end_time)" description:"结束时间；结束时间，取值范围[0, 23]"`
	Dimension  int       `orm:"column(dimension)" description:"流控维度；流控维度：1：秒；2：分；3：时；4：天"`
	Limit      int64     `orm:"column(limit)" description:"频控限制"`
	CreateTime time.Time `orm:"column(create_time);type(timestamp);auto_now_add" description:"创建时间"`
	UpdateTime time.Time `orm:"column(update_time);type(timestamp);auto_now" description:"修改时间"`
}

func (t *BusinessFlowLimit) TableName() string {
	return "business_flow_limit"
}

func init() {
	orm.RegisterModel(new(BusinessFlowLimit))
}

// AddBusinessFlowLimit insert a new BusinessFlowLimit into database and returns
// last inserted Id on success.
func AddBusinessFlowLimit(o orm.Ormer, m *BusinessFlowLimit) (id int64, err error) {
	id, err = o.Insert(m)
	return
}

// GetBusinessFlowLimitById retrieves BusinessFlowLimit by Id. Returns error if
// Id doesn't exist
func GetBusinessFlowLimitById(o orm.Ormer, id int64) (v *BusinessFlowLimit, err error) {
	v = &BusinessFlowLimit{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllBusinessFlowLimit retrieves all BusinessFlowLimit matches certain condition. Returns empty list if
// no records exist
func GetAllBusinessFlowLimit(o orm.Ormer, query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	qs := o.QueryTable(new(BusinessFlowLimit))
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

	var l []BusinessFlowLimit
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

// UpdateBusinessFlowLimit updates BusinessFlowLimit by Id and returns error if
// the record to be updated doesn't exist
func UpdateBusinessFlowLimitById(o orm.Ormer, m *BusinessFlowLimit) (err error) {
	v := BusinessFlowLimit{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteBusinessFlowLimit deletes BusinessFlowLimit by Id and returns error if
// the record to be deleted doesn't exist
func DeleteBusinessFlowLimit(o orm.Ormer, id int64) (err error) {
	v := BusinessFlowLimit{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&BusinessFlowLimit{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

func GetBusinessFlowLimitAll(o orm.Ormer) (v []BusinessFlowLimit, err error) {
	if _, err = o.QueryTable(new(BusinessFlowLimit)).All(&v); err == nil {
		return v, nil
	}
	return nil, err
}

func GetBusinessFlowLimitAllWithLimit(o orm.Ormer, limit, offset int) (v []BusinessFlowLimit, err error) {
	if _, err = o.QueryTable(new(BusinessFlowLimit)).Limit(limit, offset).All(&v); err == nil {
		return v, nil
	}
	return nil, err
}

//根据业务id获取流控数据
func GetFlowLimitByBusinessId(o orm.Ormer, businessId int64) (v []BusinessFlowLimit, err error) {
	if _, err = o.QueryTable(new(BusinessFlowLimit)).Filter("business_id", businessId).All(&v); err == nil {
		return v, nil
	}
	return nil, err
}

//根据业务id删除流控数据
func DelFlowLimitByBusinessId(o orm.Ormer, id int64) (num int64, err error) {
	if num, err = o.QueryTable(new(BusinessFlowLimit)).Filter("business_id", id).Delete(); err == nil {
		return num, nil
	}
	return
}
