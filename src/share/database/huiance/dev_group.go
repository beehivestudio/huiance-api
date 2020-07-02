package upgrade

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"

	"upgrade-api/src/share/dms"
)

type DevGroup struct {
	Id          int64     `orm:"column(id);auto" description:"ID"`
	DevModelId  int64     `orm:"column(dev_model_id)" description:"设备机型ID；外键：dev_model.id"`
	DmsGroupId  string    `orm:"column(dms_group_id);size(128)" description:"DMS设备分组ID；从设备管理获取"`
	GroupName   string    `orm:"column(group_name);size(128)" description:"设备分组名称"`
	Description string    `orm:"column(description);size(1024);null" description:"描述信息"`
	CreateTime  time.Time `orm:"column(create_time);type(timestamp);auto_now_add" description:"创建时间"`
	UpdateTime  time.Time `orm:"column(update_time);type(timestamp);auto_now" description:"修改时间"`
}

func (t *DevGroup) TableName() string {
	return "dev_group"
}

func init() {
	orm.RegisterModel(new(DevGroup))
}

// AddDevGroup insert a new DevGroup into database and returns
// last inserted Id on success.
func AddDevGroup(o orm.Ormer, m *DevGroup) (id int64, err error) {
	id, err = o.Insert(m)
	return
}

// GetDevGroupById retrieves DevGroup by Id. Returns error if
// Id doesn't exist
func GetDevGroupById(o orm.Ormer, id int64) (v *DevGroup, err error) {
	v = &DevGroup{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllDevGroup retrieves all DevGroup matches certain condition. Returns empty list if
// no records exist
func GetAllDevGroup(o orm.Ormer, query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	qs := o.QueryTable(new(DevGroup))
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

	var l []DevGroup
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

// UpdateDevGroup updates DevGroup by Id and returns error if
// the record to be updated doesn't exist
func UpdateDevGroupById(o orm.Ormer, m *DevGroup) (err error) {
	v := DevGroup{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteDevGroup deletes DevGroup by Id and returns error if
// the record to be deleted doesn't exist
func DeleteDevGroup(o orm.Ormer, id int64) (err error) {
	v := DevGroup{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&DevGroup{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

//更新设备分组表
func UpdateOrInsertDevGroup(o orm.Ormer, groupId, name, description string, devModelid int64) error {

	if devModelid == 0 {
		return errors.New("Model id is empty ")
	}

	logs.Info("Update or insert dev group start.")

	devGroup := DevGroup{
		DevModelId:  devModelid,
		DmsGroupId:  groupId,
		GroupName:   name,
		Description: description,
	}

	//如果数据表中存在，则更新
	group, err := GetGroupByGroupName(o, groupId)
	if err != nil {
		logs.Info("Get id failed! errmsg :%s", err.Error())
		return err
	}

	//数据存在，无须处理
	if group.Id == 0 {

		newId, err := o.Insert(&devGroup)
		if err != nil {
			logs.Error("Insert group failed! errmsg :%s", err.Error())
			return err
		}

		logs.Info("Insert group id :%v", newId)

		return nil
	}

	//更新名称和描述
	if group.GroupName != name || group.Description != description {

		err := UpdateGroupById(o, groupId, name, description)
		if err != nil {
			logs.Error("Update group info failed! errmsg :%s", err.Error())
			return err
		}
	}

	logs.Info("Group already exists. msg :%v", group.Id)

	return nil
}

//更新设备分组信息
func UpdateGroupById(o orm.Ormer, id, name, description string) error {

	sql := "update dev_group set group_name = ?,description = ? WHERE dms_group_id= ?"

	_, err := o.Raw(sql, name, description, id).Exec()
	if err != nil {
		logs.Error("Update group id failed! errmsg :%s", err.Error())
		return err
	}
	return nil

}

//查询设备分组信息
func GetGroupByGroupName(o orm.Ormer, groupId string) (group DevGroup, err error) {

	_, err = o.QueryTable(new(DevGroup)).Filter("dms_group_id", groupId).All(&group)

	if nil != err {
		logs.Error("Get group failed! errmsg :%s", err.Error())
		return group, err
	}

	return group, nil
}

//根据分组id获取机型id
func GetModelIdByGroupId(o orm.Ormer, groupId string) (modelId int64, err error) {

	modelName, err := dms.GetTvModelByGroupId(groupId)
	if err != nil {
		logs.Warn("Get model is failed errmsg :%s", err.Error())
		return 0, err
	}

	return GetModelByModelId(o, modelName)
}

//通过设备机型名称获取机型id
func GetModelByModelId(o orm.Ormer, modelName string) (modelId int64, err error) {

	var devModel DevModel

	_, err = o.QueryTable(new(DevModel)).Filter("model_name", modelName).All(&devModel)
	if err != nil {
		logs.Warn("Query modelId is failed ,errmsg :%s", err.Error())
		return 0, err
	}

	return devModel.Id, nil
}

/******************************************************************************
 **函数名称: GetDevGroupByModelId
 **功    能: 根据设备型号Id查询设备分组
 **输入参数:
 **      o: orm.Ormer
 **      modelId: 设备型号Id
 **输出参数: NONE
 **返    回:
 **      devGroupList: 设备分组列表
 **      err: 错误信息
 **实现描述:
 **注意事项:
 **作    者: # Linlin.guo # 2020-6-3 15:15:42 #
 ******************************************************************************/
func GetDevGroupByModelId(o orm.Ormer, modelId int64) (devGroupList []DevGroup, err error) {

	logs.Info("GetDevGroupByModelId request param. modelId: %d", modelId)

	qs := o.QueryTable(new(DevGroup)).Filter("dev_model_id", modelId).OrderBy("create_time")

	_, err = qs.All(&devGroupList)
	if err != nil {
		logs.Error("GetDevGroupByModelId failed! err: %s", err.Error())
		return nil, err
	}

	return devGroupList, nil
}
