package upgrade

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type ApkUpgradeStrategy struct {
	Id                          int64                          `json:"id"orm:"column(id);auto" description:"ID"`
	ApkId                       int64                          `json:"apk_id"orm:"column(apk_id)" description:"APK_ID；外键: apk.id"`
	Enable                      int                            `json:"enable"orm:"column(enable)" description:"是否生效；可通过此控制策略是否生效0：失效；1：生效"`
	BeginDatetime               time.Time                      `json:"-"orm:"column(begin_datetime);type(timestamp)" description:"升级开始时间"`
	EndDatetime                 time.Time                      `json:"-"orm:"column(end_datetime);type(timestamp)" description:"升级结束时间"`
	ReqBeginDatetime            int64                          `json:"begin_datetime"orm:"-"`
	ReqEndDatetime              int64                          `json:"end_datetime"orm:"-"`
	HasFlowLimit                int                            `json:"has_flow_limit"orm:"column(has_flow_limit)" description:"有无流控；0：无流控；1：有流控"`
	UpgradeType                 int                            `json:"upgrade_type"orm:"column(upgrade_type)" description:"升级方式：0：未知方式；1：可选升级（用户选择取消应用升级后，不再提示其是否升级）；2：推荐升级（用户选择取消应用升级后，下次依然提示是否升级）；3：强制升级（用户不可取消应用升级）；4：静默升级；5：强制安装（用户不知情的情况下安装应用）；6：强制卸载（用户不知情的情况下卸载应用）"`
	UpgradeDevType              int                            `json:"upgrade_dev_type"orm:"column(upgrade_dev_type)" description:"指定升级的设备范围：0：全部设备；1：指定设备ID（设备ID列表）；2：指定设备分组，1个设备只归属1个分组，1个分组只归属1个机型，1个机型只归属1个平台；3：观星用户画像组列表"`
	UpgradeDevData              string                         `json:"upgrade_dev_data"orm:"column(upgrade_dev_data)" description:"升级设备数据：0.当为全部设备时，此字段为空；；1.当指定设备ID时，此字段存储设备ID列表；；2.当指定设备分组时，此字段存储选中的设备分组信息（JSON格式）；3.当指定观星用户画像组列表时，此字段存储用户画像ID列表；"`
	Description                 string                         `json:"description"orm:"column(description);size(2048);null" description:"描述信息"`
	CreateTime                  time.Time                      `json:"create_time"orm:"column(create_time);type(timestamp);auto_now_add" description:"创建时间"`
	CreateUser                  string                         `json:"create_user"orm:"column(create_user)" description:"页面创建者"`
	UpdateTime                  time.Time                      `json:"update_time"orm:"column(update_time);type(timestamp);auto_now" description:"修改时间"`
	UpdateUser                  string                         `json:"update_user"orm:"column(update_user)" description:"页面修改者"`
	ApkUpgradeFlowLimitStrategy []*ApkUpgradeFlowLimitStrategy `json:"flow_limit" orm:"-"` // 反向关联流控策略
}

func (t *ApkUpgradeStrategy) TableName() string {
	return "apk_upgrade_strategy"
}

func init() {
	orm.RegisterModel(new(ApkUpgradeStrategy))
}

// AddApkUpgradeStrategy insert a new ApkUpgradeStrategy into database and returns
// last inserted Id on success.
func AddApkUpgradeStrategy(o orm.Ormer, m *ApkUpgradeStrategy) (id int64, err error) {
	id, err = o.Insert(m)
	return
}

// GetApkUpgradeStrategyById retrieves ApkUpgradeStrategy by Id. Returns error if
// Id doesn't exist
func GetApkUpgradeStrategyById(o orm.Ormer, id int64) (v *ApkUpgradeStrategy, err error) {
	v = &ApkUpgradeStrategy{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllApkUpgradeStrategy retrieves all ApkUpgradeStrategy matches certain condition. Returns empty list if
// no records exist
func GetAllApkUpgradeStrategy(o orm.Ormer, query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	qs := o.QueryTable(new(ApkUpgradeStrategy))
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

	var l []ApkUpgradeStrategy
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

// UpdateApkUpgradeStrategy updates ApkUpgradeStrategy by Id and returns error if
// the record to be updated doesn't exist
func UpdateApkUpgradeStrategyById(o orm.Ormer, m *ApkUpgradeStrategy) (err error) {
	v := ApkUpgradeStrategy{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteApkUpgradeStrategy deletes ApkUpgradeStrategy by Id and returns error if
// the record to be deleted doesn't exist
func DeleteApkUpgradeStrategy(o orm.Ormer, id int64) (err error) {
	v := ApkUpgradeStrategy{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&ApkUpgradeStrategy{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
