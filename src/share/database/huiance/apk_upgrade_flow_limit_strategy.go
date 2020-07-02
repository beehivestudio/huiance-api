package upgrade

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type ApkUpgradeFlowLimitStrategy struct {
	Id         int64     `json:"id" orm:"column(id);auto" description:"ID"`
	StrategyId int64     `json:"strategy_id" orm:"column(strategy_id)" description:"APK升级策略ID；外键: apk_upgrade_strategy.id"`
	BeginTime  int       `json:"begin_time" orm:"column(begin_time)" description:"开始时间段；开始时间段[0, 23]（精确到时）"`
	EndTime    int       `json:"end_time" orm:"column(end_time)" description:"结束时间段；结束时间段[0, 23]（精确到时）"`
	Dimension  int       `json:"dimension" orm:"column(dimension)" description:"流控维度；流控维度：1：秒；2：分；3：时；4：天"`
	Limit      int64     `json:"limit" orm:"column(limit)" description:"频控限制；在流控维度上的次数"`
	CreateTime time.Time `json:"create_time" orm:"column(create_time);type(timestamp);auto_now_add" description:"创建时间"`
	UpdateTime time.Time `json:"update_time" orm:"column(update_time);type(timestamp);auto_now" description:"修改时间"`
}

func (t *ApkUpgradeFlowLimitStrategy) TableName() string {
	return "apk_upgrade_flow_limit_strategy"
}

func init() {
	orm.RegisterModel(new(ApkUpgradeFlowLimitStrategy))
}

// AddApkUpgradeFlowLimitStrategy insert a new ApkUpgradeFlowLimitStrategy into database and returns
// last inserted Id on success.
func AddApkUpgradeFlowLimitStrategy(o orm.Ormer, m *ApkUpgradeFlowLimitStrategy) (id int64, err error) {
	id, err = o.Insert(m)
	return
}

// GetApkUpgradeFlowLimitStrategyById retrieves ApkUpgradeFlowLimitStrategy by Id. Returns error if
// Id doesn't exist
func GetApkUpgradeFlowLimitStrategyById(o orm.Ormer, id int64) (v *ApkUpgradeFlowLimitStrategy, err error) {
	v = &ApkUpgradeFlowLimitStrategy{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllApkUpgradeFlowLimitStrategy retrieves all ApkUpgradeFlowLimitStrategy matches certain condition. Returns empty list if
// no records exist
func GetAllApkUpgradeFlowLimitStrategy(o orm.Ormer, query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	qs := o.QueryTable(new(ApkUpgradeFlowLimitStrategy))
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

	var l []ApkUpgradeFlowLimitStrategy
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

// UpdateApkUpgradeFlowLimitStrategy updates ApkUpgradeFlowLimitStrategy by Id and returns error if
// the record to be updated doesn't exist
func UpdateApkUpgradeFlowLimitStrategyById(o orm.Ormer, m *ApkUpgradeFlowLimitStrategy) (err error) {
	v := ApkUpgradeFlowLimitStrategy{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteApkUpgradeFlowLimitStrategy deletes ApkUpgradeFlowLimitStrategy by Id and returns error if
// the record to be deleted doesn't exist
func DeleteApkUpgradeFlowLimitStrategy(o orm.Ormer, id int64) (err error) {
	v := ApkUpgradeFlowLimitStrategy{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&ApkUpgradeFlowLimitStrategy{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

func GetApkUpgradeFlowLimitStrategyAll(o orm.Ormer) (v []ApkUpgradeFlowLimitStrategy, err error) {
	if _, err = o.QueryTable(new(ApkUpgradeFlowLimitStrategy)).All(&v); err == nil {
		return v, nil
	}
	return nil, err
}

func GetApkUpgradeFlowLimitStrategyAllWithLimit(o orm.Ormer, limit, offset int) (v []ApkUpgradeFlowLimitStrategy, err error) {
	if _, err = o.QueryTable(new(ApkUpgradeFlowLimitStrategy)).Limit(limit, offset).All(&v); err == nil {
		return v, nil
	}
	return nil, err
}
