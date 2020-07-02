package upgrade

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

type ApkPatchAlgo struct {
	Id          int64     `orm:"column(id);auto" description:"ID"`
	Name        string    `orm:"column(name);size(128)" description:"算法名称"`
	Enable      int       `orm:"column(enable)" description:"是否启用"`
	Description string    `orm:"column(description);size(1024)" description:"描述信息"`
	CreateTime  time.Time `orm:"column(create_time);type(timestamp);auto_now_add" description:"创建时间"`
	UpdateTime  time.Time `orm:"column(update_time);type(timestamp);auto_now" description:"修改时间"`
}

func (t *ApkPatchAlgo) TableName() string {
	return "apk_patch_algo"
}

func init() {
	orm.RegisterModel(new(ApkPatchAlgo))
}

// AddApkPatchAlgo insert a new ApkPatchAlgo into database and returns
// last inserted Id on success.
func AddApkPatchAlgo(o orm.Ormer, m *ApkPatchAlgo) (id int64, err error) {

	id, err = o.Insert(m)
	return
}

// GetApkPatchAlgoById retrieves ApkPatchAlgo by Id. Returns error if
// Id doesn't exist
func GetApkPatchAlgoById(o orm.Ormer, id int64) (v *ApkPatchAlgo, err error) {
	v = &ApkPatchAlgo{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllApkPatchAlgo retrieves all ApkPatchAlgo matches certain condition. Returns empty list if
// no records exist
func GetAllApkPatchAlgo(o orm.Ormer, query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	qs := o.QueryTable(new(ApkPatchAlgo))
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

	var l []ApkPatchAlgo
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

// UpdateApkPatchAlgo updates ApkPatchAlgo by Id and returns error if
// the record to be updated doesn't exist
func UpdateApkPatchAlgoById(o orm.Ormer, m *ApkPatchAlgo) (err error) {
	v := ApkPatchAlgo{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteApkPatchAlgo deletes ApkPatchAlgo by Id and returns error if
// the record to be deleted doesn't exist
func DeleteApkPatchAlgo(o orm.Ormer, id int64) (err error) {
	v := ApkPatchAlgo{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&ApkPatchAlgo{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

/******************************************************************************
 **函数名称: GetPatchAlgoList
 **功    能: 查询所有差分算法列表
 **输入参数:
 **      o: orm.Ormer
 **      id: 差分算法Id
 **输出参数: NONE
 **返    回:
 **      patchAlgoList: 差分算法列表
 **      err: 错误信息
 **实现描述:
 **注意事项:
 **作    者: # Linlin.guo # 2020-6-2 15:15:42 #
 ******************************************************************************/
func GetPatchAlgoList(o orm.Ormer, id int64) (patchAlgoList []ApkPatchAlgo, err error) {

	qs := o.QueryTable(new(ApkPatchAlgo)).Filter("enable", 1)
	if id != 0 {
		qs = qs.Filter("id", id)
	}
	qs = qs.OrderBy("create_time")

	_, err = qs.All(&patchAlgoList)
	if err != nil {
		logs.Error("GetPatchAlgoList err: %s", err.Error())
		return nil, err
	}

	return patchAlgoList, nil
}
