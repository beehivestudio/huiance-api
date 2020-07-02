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

type DevModel struct {
	Id          int64     `orm:"column(id);auto" description:"ID"`
	DevPlatId   int64     `orm:"column(dev_plat_id)" description:"设备平台ID；外键：dev_plat.id"`
	ModelName   string    `orm:"column(model_name);size(128)" description:"设备机型名称"`
	Description string    `orm:"column(description);size(1024);null" description:"描述信息"`
	CreateTime  time.Time `orm:"column(create_time);type(timestamp);auto_now_add" description:"创建时间"`
	UpdateTime  time.Time `orm:"column(update_time);type(timestamp);auto_now" description:"修改时间"`
}

func (t *DevModel) TableName() string {
	return "dev_model"
}

func init() {
	orm.RegisterModel(new(DevModel))
}

// AddDevModel insert a new DevModel into database and returns
// last inserted Id on success.
func AddDevModel(o orm.Ormer, m *DevModel) (id int64, err error) {
	id, err = o.Insert(m)
	return
}

// GetDevModelById retrieves DevModel by Id. Returns error if
// Id doesn't exist
func GetDevModelById(o orm.Ormer, id int64) (v *DevModel, err error) {
	v = &DevModel{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllDevModel retrieves all DevModel matches certain condition. Returns empty list if
// no records exist
func GetAllDevModel(o orm.Ormer, query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	qs := o.QueryTable(new(DevModel))
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

	var l []DevModel
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

// UpdateDevModel updates DevModel by Id and returns error if
// the record to be updated doesn't exist
func UpdateDevModelById(o orm.Ormer, m *DevModel) (err error) {
	v := DevModel{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteDevModel deletes DevModel by Id and returns error if
// the record to be deleted doesn't exist
func DeleteDevModel(o orm.Ormer, id int64) (err error) {
	v := DevModel{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&DevModel{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

type GroupsInfor struct {
	Group       string `json:"group"`
	Description string `json:"description"`
}

//查询设备机型表数量
func GetDevModelCount(o orm.Ormer) (count int64, err error) {
	return o.QueryTable(new(DevModel)).Count()
}

//更新设备机型表
func UpdateModel(o orm.Ormer, model, description string, devPlatId int64) error {

	logs.Info("Update or insert dev model start.")

	devModel := DevModel{
		DevPlatId:   devPlatId,
		ModelName:   model,
		Description: description,
	}

	//如果数据表中存在，则更新
	m, err := GetModelByModelName(o, model, devPlatId)
	if err != nil {
		logs.Info("Get id failed! errmsg :%s", err.Error())
		return err
	}

	//数据存在，无须处理
	if m.Id != 0 {
		logs.Info("Model already exists. mes :%v", m.Id)
		return nil
	}

	newId, err := o.Insert(&devModel)
	if err != nil {
		logs.Error("Insert model failed! errmsg :%s", err.Error())
		return err
	}

	logs.Info("Insert model id :%v", newId)

	return nil

}

//查询设备机型信息
func GetModelByModelName(o orm.Ormer, modelName string, devPlatId int64) (model DevModel, err error) {

	_, err = o.QueryTable(new(DevModel)).Filter("dev_plat_id", devPlatId).
		Filter("model_name", modelName).All(&model)

	if nil != err {
		logs.Error("Get Model failed! errmsg :%s", err.Error())
		return model, err
	}

	return model, nil
}

//通过设备机型获取平台id
func GetPlatIdByModel(o orm.Ormer, model string) (platId int64, err error) {

	//获取平台名称
	plat, err := dms.GetTvPlatByModel(model)

	if err != nil {
		return 0, errors.New("get plat err")
	} else if plat == "" {
		logs.Warn("Get plat is empty. msg :%s", model)
		return 0, nil
	}

	return GetPlatIdlByPlat(o, plat)
}

//通过平台获取平台id
func GetPlatIdlByPlat(o orm.Ormer, plat string) (platId int64, err error) {

	var devplat DevPlat

	_, err = o.QueryTable(new(DevPlat)).Filter("name", plat).All(&devplat)
	if err != nil {
		logs.Warn("Query platId is failed ,errmsg :%s", err.Error())
		return 0, err
	}

	return devplat.Id, nil
}

/******************************************************************************
 **函数名称: GetDevModelByPlatId
 **功    能: 根据设备平台Id查询设备型号
 **输入参数:
 **      o: orm.Ormer
 **      platId: 设备平台Id
 **输出参数: NONE
 **返    回:
 **      devModelList: 设备平台列表
 **      err: 错误信息
 **实现描述:
 **注意事项:
 **作    者: # Linlin.guo # 2020-6-3 15:15:42 #
 ******************************************************************************/
func GetDevModelByPlatId(o orm.Ormer, platId int64) (devModelList []DevModel, err error) {

	logs.Info("GetDevModelByPlatId request param. platId: %d", platId)

	qs := o.QueryTable(new(DevModel)).Filter("dev_plat_id", platId).OrderBy("create_time")

	_, err = qs.All(&devModelList)
	if err != nil {
		logs.Error("GetDevModelByPlatId failed! err: %s", err.Error())
		return nil, err
	}

	return devModelList, nil
}
