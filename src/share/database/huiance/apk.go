package upgrade

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"
	"upgrade-api/src/share/comm"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

type Apk struct {
	Id                int64     `orm:"column(id);auto" description:"ID"`
	Bcode             string    `orm:"column(bcode);size(128)" description:"业务方唯一码"`
	Enable            int       `orm:"column(enable)" description:"是否启用；可通过此开关控制APK上下线；0：禁用（下线）；1：启用（上线）"`
	AppId             int64     `orm:"column(app_id)" description:"应用ID；外键: app.id"`
	VersionCode       int64     `orm:"column(version_code)" description:"版本号；版本号数值越大，版本号越高；注：由程序从APK包中解析获取"`
	VersionName       string    `orm:"column(version_name);size(1024)" description:"版本名；注：由程序从APK包中解析获取"`
	Url               string    `orm:"column(url);size(1024)" description:"资源地址；Apk包的下载地址"`
	Md5               string    `orm:"column(md5);size(64)" description:"MD5值；APK包的MD5值（文件MD5值）"`
	Size              int64     `orm:"column(size)" description:"包大小；APK包大小（单位：字节）"`
	EuiLowVersion     string    `orm:"column(eui_low_version);size(256);null" description:"依赖的EUI最低版本"`
	EuiLowVersionInt  int64     `orm:"column(eui_low_version_int)" description:"依赖的EUI最低版本（数字）"`
	EuiHighVersion    string    `orm:"column(eui_high_version);size(256);null" description:"依赖的EUI最高版本"`
	EuiHighVersionInt int64     `orm:"column(eui_high_version_int)" description:"依赖的EUI最高版本"`
	Status            int       `orm:"column(status)" description:"处理状态：0：处理中1：正常2：版本验证失败3：秘钥验证失败4：其他异常5：超时（超过1分钟算超时）"`
	Description       string    `orm:"column(description);size(2048);null" description:"描述信息"`
	Memo              string    `orm:"column(memo);size(2048);null" description:"记录内部的沟通结果"`
	CallbackUrl       string    `orm:"column(callback_url)" description:"回调url"`
	CreateTime        time.Time `orm:"column(create_time);type(timestamp);auto_now_add" description:"创建时间"`
	CreateUser        string    `orm:"column(create_user)" description:"页面创建者"`
	UpdateTime        time.Time `orm:"column(update_time);type(timestamp);auto_now" description:"修改时间"`
	UpdateUser        string    `orm:"column(update_user)" description:"页面修改者"`
}

func (t *Apk) TableName() string {
	return "apk"
}

func init() {
	orm.RegisterModel(new(Apk))
}

// AddApk insert a new Apk into database and returns
// last inserted Id on success.
func AddApk(o orm.Ormer, m *Apk) (id int64, err error) {
	id, err = o.Insert(m)
	return
}

// GetApkById retrieves Apk by Id. Returns error if
// Id doesn't exist
func GetApkById(o orm.Ormer, id int64) (v *Apk, err error) {

	v = &Apk{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

func GetApkByBcode(o orm.Ormer, bcode string) (v *Apk, err error) {

	v = &Apk{Bcode: bcode}
	if err = o.Read(v, "bcode"); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllApk retrieves all Apk matches certain condition. Returns empty list if
// no records exist
func GetAllApk(o orm.Ormer, query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {

	qs := o.QueryTable(new(Apk))
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

	var l []Apk
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

// UpdateApk updates Apk by Id and returns error if
// the record to be updated doesn't exist
func UpdateApkById(o orm.Ormer, m *Apk) (err error) {

	v := Apk{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteApk deletes Apk by Id and returns error if
// the record to be deleted doesn't exist
func DeleteApk(o orm.Ormer, id int64) (err error) {

	v := Apk{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Apk{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

/******************************************************************************
 **函数名称: GetApkByAppIdAndVersion
 **功    能: 根据应用ID和应用版本查询对应的apk
 **输入参数:
 **      o: orm.Ormer
 **      appId: 应用ID
 **      version: 应用版本号
 **输出参数:
 **返    回:
 **      apk: apk信息
 **      err: 错误信息
 **实现描述:
 **注意事项:
 **作    者: # Linlin.guo # 2020-6-1 13:15:42 #
 ******************************************************************************/
func GetApkByAppIdAndVersion(o orm.Ormer, appId, version int64) (apk Apk, err error) {

	logs.Info("GetApkByAppIdAndVersion request param appId: %d,version: %d", appId, version)

	qs := o.QueryTable(new(Apk)).Filter("appId", appId).Filter("version_code",
		version).Filter("status", APK_STATUS_NORMAL).Limit(1)

	_, err = qs.All(&apk)
	if err != nil {
		logs.Error("GetApkPatchByStatus err: %s", err.Error())
		return apk, err
	}
	if len(apk.Url) == 0 {
		return apk, errors.New("Error:ApkUrl is null")
	}

	return apk, err
}

/******************************************************************************
 **函数名称: GetApkMaxVersionCodeByAppId
 **功    能: 根据应用ID查询apk版本号最大
 **输入参数:
 **      o: orm.Ormer
 **      appId: 应用ID
 **输出参数:
 **返    回:
 **      versionCode: apk版本号
 **      err: 错误信息
 **实现描述:
 **注意事项:
 **作    者: # Shuangpeng.guo # 2020-6-1 13:15:42 #
 ******************************************************************************/
func GetApkMaxVersionCodeByAppId(o orm.Ormer, appId int64) (*Apk, error) {
	apk := new(Apk)
	sql := fmt.Sprintf(`SELECT 
		id,
		version_code
	FROM 
		%s 
	WHERE
		enable = ?
	AND
		app_id = ?
	ORDER BY 
		version_code DESC 
	LIMIT 1
	`, apk.TableName())

	if err := o.Raw(sql, comm.ENABLE, appId).QueryRow(apk); err != nil {
		return nil, err
	}

	return apk, nil
}
