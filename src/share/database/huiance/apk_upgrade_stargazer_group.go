package upgrade

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type ApkUpgradeStargazerGroup struct {
	Id           int64     `orm:"column(id);auto" description:"ID"`
	StrategyId   int64     `orm:"column(strategy_id)" description:"APK升级策略ID；外键: apk_upgrade_strategy.id"`
	StargazerTid int64     `orm:"column(stargazer_tid)" description:"观星用户画像分组ID"`
	CreateTime   time.Time `orm:"column(create_time);type(timestamp);auto_now_add" description:"创建时间"`
	UpdateTime   time.Time `orm:"column(update_time);type(timestamp);auto_now" description:"修改时间"`
}

func (t *ApkUpgradeStargazerGroup) TableName() string {
	return "apk_upgrade_stargazer_group"
}

func init() {
	orm.RegisterModel(new(ApkUpgradeStargazerGroup))
}

// AddApkUpgradeStargazerGroup insert a new ApkUpgradeStargazerGroup into database and returns
// last inserted Id on success.
func AddApkUpgradeStargazerGroup(o orm.Ormer, m *ApkUpgradeStargazerGroup) (id int64, err error) {
	id, err = o.Insert(m)
	return
}

// GetApkUpgradeStargazerGroupById retrieves ApkUpgradeStargazerGroup by Id. Returns error if
// Id doesn't exist
func GetApkUpgradeStargazerGroupById(o orm.Ormer, id int64) (v *ApkUpgradeStargazerGroup, err error) {
	v = &ApkUpgradeStargazerGroup{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllApkUpgradeStargazerGroup retrieves all ApkUpgradeStargazerGroup matches certain condition. Returns empty list if
// no records exist
func GetAllApkUpgradeStargazerGroup(o orm.Ormer, query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	qs := o.QueryTable(new(ApkUpgradeStargazerGroup))
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

	var l []ApkUpgradeStargazerGroup
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

// UpdateApkUpgradeStargazerGroup updates ApkUpgradeStargazerGroup by Id and returns error if
// the record to be updated doesn't exist
func UpdateApkUpgradeStargazerGroupById(o orm.Ormer, m *ApkUpgradeStargazerGroup) (err error) {
	v := ApkUpgradeStargazerGroup{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteApkUpgradeStargazerGroup deletes ApkUpgradeStargazerGroup by Id and returns error if
// the record to be deleted doesn't exist
func DeleteApkUpgradeStargazerGroup(o orm.Ormer, id int64) (err error) {
	v := ApkUpgradeStargazerGroup{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&ApkUpgradeStargazerGroup{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
