package upgrade

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/garyburd/redigo/redis"

	"upgrade-api/src/share/comm"
	"upgrade-api/src/share/lib/rds"
)

type DevType struct {
	Id          int64     `orm:"column(id);auto" description:"ID"`
	Name        string    `orm:"column(name);size(128)" description:"设备名称；设备名称列表：1.乐视电视；2.乐视手机"`
	Description string    `orm:"column(description);size(1024);null" description:"描述信息"`
	CreateTime  time.Time `orm:"column(create_time);type(timestamp);auto_now_add" description:"创建时间"`
	UpdateTime  time.Time `orm:"column(update_time);type(timestamp);auto_now" description:"修改时间"`
}

func (t *DevType) TableName() string {
	return "dev_type"
}

func init() {
	orm.RegisterModel(new(DevType))
}

// AddDevType insert a new DevType into database and returns
// last inserted Id on success.
func AddDevType(o orm.Ormer, m *DevType) (id int64, err error) {
	id, err = o.Insert(m)
	return
}

// GetDevTypeById retrieves DevType by Id. Returns error if
// Id doesn't exist
func GetDevTypeById(o orm.Ormer, id int64) (v *DevType, err error) {
	v = &DevType{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllDevType retrieves all DevType matches certain condition. Returns empty list if
// no records exist
func GetAllDevType(o orm.Ormer, query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	qs := o.QueryTable(new(DevType))
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

	var l []DevType
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

// UpdateDevType updates DevType by Id and returns error if
// the record to be updated doesn't exist
func UpdateDevTypeById(o orm.Ormer, m *DevType) (err error) {
	v := DevType{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteDevType deletes DevType by Id and returns error if
// the record to be deleted doesn't exist
func DeleteDevType(o orm.Ormer, id int64) (err error) {
	v := DevType{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&DevType{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

/******************************************************************************
 **函数名称: GetDevTypeList
 **功    能: 查询所有设备类型列表
 **输入参数:
 **      o: orm.Ormer
 **输出参数: NONE
 **返    回:
 **      devTypeList: 设备类型列表
 **      err: 错误信息
 **实现描述:
 **注意事项:
 **作    者: # Linlin.guo # 2020-5-25 15:15:42 #
 ******************************************************************************/
func GetDevTypeList(o orm.Ormer) (devTypeList []DevType, err error) {

	qs := o.QueryTable(new(DevType)).OrderBy("create_time")

	_, err = qs.All(&devTypeList)
	if err != nil {
		logs.Error("GetDevTypeList failed! err: %s", err.Error())
		return nil, err
	}

	return devTypeList, nil
}

type DevReturn struct {
	DevTypeId int64           `json:"dev_type_id"` //设备类型ID
	Len       int64           `json:"len"`         //平台集合长度
	Plats     []DevPlatReturn `json:"plats"`       //平台
}

type DevPlatReturn struct {
	Id       int64            `json:"id"`
	PlatName string           `json:"plat_name"` //平台名称
	Len      int64            `json:"len"`       //型号集合长度
	Models   []DevModelReturn `json:"models"`    //型号
}

type DevModelReturn struct {
	Id        int64            `json:"id"`
	ModelName string           `json:"model_name"` //型号名称
	Len       int64            `json:"len"`        //分组集合长度
	Groups    []DevGroupReturn `json:"groups"`     //分组
}

type DevGroupReturn struct {
	Id         int64  `json:"id"`           //分组ID
	DmsGroupId string `json:"dms_group_id"` //DMS设备分组ID
	GroupName  string `json:"group_name"`   //分组名称
}

/******************************************************************************
 **函数名称: GetDevGroupList
 **功    能: 根据设备类型查询分组列表
 **输入参数:
 **      o: orm.Ormer
 **      conn: redis.Conn
 **      devTypeId: 设备类型Id
 **输出参数: NONE
 **返    回:
 **      devReturn: 分组返回列表
 **      err: 错误信息
 **实现描述:
 **注意事项:
 **作    者: # Linlin.guo # 2020-6-4 15:15:42 #
 ******************************************************************************/
func GetDevGroupList(o orm.Ormer, conn redis.Conn, devTypeId int64) (devReturn DevReturn, err error) {

	logs.Info("GetDevGroupList request param. devTypeId: %d", devTypeId)

	//1.验证devTypeI的有效性
	_, err = GetDevTypeById(o, devTypeId)
	if err != nil {
		logs.Error("GetDevTypeById failed! err: %s", err.Error())
		return devReturn, err
	}

	//2.数据从缓存中获取
	err = rds.RedisGetJsonData(conn, comm.RDS_KEY_DEVICE_GROUP_LIST, &devReturn)
	if err == nil {
		return devReturn, err
	}

	//3.获取设备分组映射关系列表
	devReturn, err = HandleDevGroupList(o, devTypeId)
	if err != nil {
		logs.Error("GetDevGroupList HandleDevGroupList failed! err: %s", err.Error())
		return devReturn, err
	}

	//4.数据放入缓存中
	err = rds.RedisSaveJsonData(conn, comm.RDS_KEY_DEVICE_GROUP_LIST, comm.DEVICE_GROUPLS_RDS_EXPIRE_SEC, devReturn)
	if err != nil && devReturn.Len != 0 {
		logs.Error("GetDevGroupList RedisGetJsonData failed! err: %s", err.Error())
		return devReturn, err
	}

	return devReturn, err
}

/******************************************************************************
 **函数名称: HandleDevGroupList
 **功    能: 查询并处理设备映射关系
 **输入参数:
 **      o: orm.Ormer
 **      devTypeId: 设备类型Id
 **输出参数: NONE
 **返    回:
 **      devReturn: 分组返回列表
 **      err: 错误信息
 **实现描述:
 **注意事项:
 **作    者: # Linlin.guo # 2020-6-4 15:15:42 #
 ******************************************************************************/
func HandleDevGroupList(o orm.Ormer, devTypeId int64) (devReturn DevReturn, err error) {
	devReturn.DevTypeId = devTypeId

	//1.根据设备类型查设备平台集合
	devPlatList, err := GetDevPlatByDevTypeId(o, devTypeId)
	if err != nil {
		logs.Error("GetDevPlatByDevTypeId failed! err: %s", err.Error())
		return devReturn, err
	}

	devReturn.Len = int64(len(devPlatList))

	//循环设备平台集合
	for _, devPlat := range devPlatList {
		var devPlatReturn DevPlatReturn

		//2.根据设备平台查设备型号集合
		devModelList, err := GetDevModelByPlatId(o, devPlat.Id)
		if err != nil {
			logs.Error("GetDevModelByPlatId failed! err: %s", err.Error())
			return devReturn, err
		}

		devPlatReturn.Len = int64(len(devModelList))
		devPlatReturn.Id = devPlat.Id
		devPlatReturn.PlatName = devPlat.Name

		//循环设备型号集合
		for _, devModel := range devModelList {
			var devModelReturn DevModelReturn

			//3.根据设备型号查设备分组集合
			devGroupList, err := GetDevGroupByModelId(o, devModel.Id)
			if err != nil {
				logs.Error("GetDevGroupByModelId failed! err: %s", err.Error())
				return devReturn, err
			}

			devModelReturn.Len = int64(len(devGroupList))
			devModelReturn.Id = devModel.Id
			devModelReturn.ModelName = devModel.ModelName

			//循环设备分组集合
			for _, devGroup := range devGroupList {
				var devGroupReturn DevGroupReturn

				devGroupReturn.Id = devGroup.Id
				devGroupReturn.DmsGroupId = devGroup.DmsGroupId
				devGroupReturn.GroupName = devGroup.GroupName

				devModelReturn.Groups = append(devModelReturn.Groups, devGroupReturn)
			}
			devPlatReturn.Models = append(devPlatReturn.Models, devModelReturn)
		}
		devReturn.Plats = append(devReturn.Plats, devPlatReturn)
	}

	return devReturn, nil
}
