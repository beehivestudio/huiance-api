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

type DevPlat struct {
	Id          int64     `orm:"column(id);auto" description:"ID"`
	DevTypeId   int64     `orm:"column(dev_type_id)" description:"设备种类；外键：dev_type.id"`
	Name        string    `orm:"column(name);size(128)" description:"设备平台名称"`
	Description string    `orm:"column(description);size(1024);null" description:"描述信息"`
	CreateTime  time.Time `orm:"column(create_time);type(timestamp);auto_now_add" description:"创建时间"`
	UpdateTime  time.Time `orm:"column(update_time);type(timestamp);auto_now" description:"修改时间"`
}

func (t *DevPlat) TableName() string {
	return "dev_plat"
}

func init() {
	orm.RegisterModel(new(DevPlat))
}

// AddDevPlat insert a new DevPlat into database and returns
// last inserted Id on success.
func AddDevPlat(o orm.Ormer, m *DevPlat) (id int64, err error) {

	id, err = o.Insert(m)
	return
}

// GetDevPlatById retrieves DevPlat by Id. Returns error if
// Id doesn't exist
func GetDevPlatById(o orm.Ormer, id int64) (v *DevPlat, err error) {

	v = &DevPlat{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllDevPlat retrieves all DevPlat matches certain condition. Returns empty list if
// no records exist
func GetAllDevPlat(o orm.Ormer, query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	qs := o.QueryTable(new(DevPlat))
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

	var l []DevPlat
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

// UpdateDevPlat updates DevPlat by Id and returns error if
// the record to be updated doesn't exist
func UpdateDevPlatById(o orm.Ormer, m *DevPlat) (err error) {
	v := DevPlat{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteDevPlat deletes DevPlat by Id and returns error if
// the record to be deleted doesn't exist
func DeleteDevPlat(o orm.Ormer, id int64) (err error) {
	v := DevPlat{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&DevPlat{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

//设备平台表数量
func GetDevPlatCount(o orm.Ormer) (count int64, err error) {
	return o.QueryTable(new(DevPlat)).Count()
}

//更新设备平台表数据
func UpdatePlat(o orm.Ormer, name, description string, devTypeId int64) error {

	logs.Info("Update or insert dev plat start.")

	devPlat := DevPlat{
		DevTypeId:   devTypeId,
		Name:        name,
		Description: description,
	}

	//如果数据表中存在，则更新
	p, err := GetPlatByPlatName(o, name, devTypeId)
	if err != nil {
		logs.Info("Get plat id failed! errmsg :%s", err.Error())
		return err
	}

	//数据存在，无须处理
	if p.Id != 0 {
		logs.Info("Plat already exists. msg:%v", p.Id)
		return nil
	}

	newId, err := o.Insert(&devPlat)
	if err != nil {
		logs.Error("Insert plat failed! errmsg :%s", err.Error())
		return err
	}

	logs.Info("Insert plat id :%v", newId)

	return nil

}

//查询设备平台信息
func GetPlatByPlatName(o orm.Ormer, platName string, devTypeId int64) (plat DevPlat, err error) {

	_, err = o.QueryTable(new(DevPlat)).Filter("dev_type_id", devTypeId).
		Filter("name", platName).All(&plat)

	if nil != err {
		logs.Error("Get plat failed! errmsg :%s", err.Error())
		return plat, err
	}

	return plat, nil
}

/******************************************************************************
 **函数名称: GetDevPlatByDevTypeId
 **功    能: 根据设备类型查询设备平台
 **输入参数:
 **      o: orm.Ormer
 **      devTypeId: 设备类型Id
 **输出参数: NONE
 **返    回:
 **      devPlatList: 设备平台列表
 **      err: 错误信息
 **实现描述:
 **注意事项:
 **作    者: # Linlin.guo # 2020-6-3 15:15:42 #
 ******************************************************************************/
func GetDevPlatByDevTypeId(o orm.Ormer, devTypeId int64) (devPlatList []DevPlat, err error) {

	logs.Info("GetDevPlatByDevTypeId request param. devTypeId: %d", devTypeId)

	qs := o.QueryTable(new(DevPlat)).Filter("dev_type_id", devTypeId).OrderBy("create_time")

	_, err = qs.All(&devPlatList)
	if err != nil {
		logs.Error("GetDevPlatByDevTypeId failed! err: %s", err.Error())
		return nil, err
	}

	return devPlatList, nil
}
