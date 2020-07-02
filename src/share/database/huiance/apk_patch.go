package upgrade

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

/* 差分处理的状态 */
const (
	APK_PATCH_STATUS_WAIT       = 0 // 差分处理尚未开始(待处理)
	APK_PATCH_STATUS_DEALING    = 1 // 差分处理中
	APK_PATCH_STATUS_UPLOAD_CDN = 2 // 差分包上传cdn中
	APK_PATCH_STATUS_SUCC       = 3 // 差分处理成功
	APK_PATCH_STATUS_OVERSIZE   = 4 // 差分包大于新版本全量包
	APK_PATCH_STATUS_TIMEOUT    = 5 // 差分过程超时
	APK_PATCH_STATUS_ERROR      = 6 // 差分过程错误
	APK_PATCH_STATUS_PROBLEM    = 7 // 差分包有问题
)

type ApkPatch struct {
	Id              int64     `orm:"column(id);auto" description:"ID"`
	AppId           int64     `orm:"column(app_id)" description:"应用ID；外键: app.id"`
	HighVersionCode int64     `orm:"column(high_version_code)" description:"高版本号"`
	LowVersionCode  int64     `orm:"column(low_version_code)" description:"低版本号"`
	PatchAlgo       int64     `orm:"column(patch_algo)" description:"差分算法"`
	Status          int       `orm:"column(status)" description:"处理状态"`
	Url             string    `orm:"column(url);size(1024)" description:"差分包的下载地址；注意：路径中应体现差分算法ID"`
	Md5             string    `orm:"column(md5);size(64)" description:"差分包MD5值"`
	Size            int64     `orm:"column(size)" description:"差分包大小；单位：字节"`
	Description     string    `orm:"column(description);size(2048);null" description:"描述信息"`
	CreateTime      time.Time `orm:"column(create_time);type(timestamp);auto_now_add" description:"创建时间"`
	UpdateTime      time.Time `orm:"column(update_time);type(timestamp);auto_now" description:"修改时间"`
}

func (t *ApkPatch) TableName() string {
	return "apk_patch"
}

func init() {
	orm.RegisterModel(new(ApkPatch))
}

// AddApkPatch insert a new ApkPatch into database and returns
// last inserted Id on success.
func AddApkPatch(o orm.Ormer, m *ApkPatch) (id int64, err error) {
	id, err = o.Insert(m)
	return
}

// GetApkPatchById retrieves ApkPatch by Id. Returns error if
// Id doesn't exist
func GetApkPatchById(o orm.Ormer, id int64) (v *ApkPatch, err error) {
	v = &ApkPatch{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllApkPatch retrieves all ApkPatch matches certain condition. Returns empty list if
// no records exist
func GetAllApkPatch(o orm.Ormer, query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	qs := o.QueryTable(new(ApkPatch))
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

	var l []ApkPatch
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

// UpdateApkPatch updates ApkPatch by Id and returns error if
// the record to be updated doesn't exist
func UpdateApkPatchById(o orm.Ormer, m *ApkPatch) (err error) {
	v := ApkPatch{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteApkPatch deletes ApkPatch by Id and returns error if
// the record to be deleted doesn't exist
func DeleteApkPatch(o orm.Ormer, id int64) (err error) {
	v := ApkPatch{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&ApkPatch{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

/******************************************************************************
 **函数名称: GetApkPatchByStatus
 **功    能: 查询apkPatch 符合处理要求的信息
 **输入参数:
 **      o: orm.Ormer
 **      status: 差分信息状态
 **输出参数:
 **返    回:
 **      err: 错误信息
 **实现描述:
 **注意事项:
 **作    者: # Linlin.guo # 2020-5-25 15:15:42 #
 ******************************************************************************/
func GetApkPatchByStatus(o orm.Ormer, status int64) (apkPatch []ApkPatch, err error) {
	sql := fmt.Sprintf(`
		SELECT
			ap.id,
			apa.enable,
			ap.update_time
		FROM apk_patch ap
		LEFT JOIN apk_patch_algo apa ON ap.patch_algo=apa.id 
		WHERE apa.enable=1 AND status=%d
		ORDER BY ap.update_time`, status)

	_, err = o.Raw(sql).QueryRows(&apkPatch)
	if nil != err {
		logs.Error("Get apk patch by status failed! status:%d errmsg:%s",
			status, err.Error())
		return apkPatch, err
	}

	return apkPatch, nil
}

/******************************************************************************
 **函数名称: ApkPatchUpdateByMap
 **功    能: apkPatch根据map更新部分字段
 **输入参数:
 **      o: orm.Ormer
 **      mapParam: map
 **      apkPatchId: id
 **输出参数: NONE
 **返    回:
 **      err: 错误信息
 **实现描述:
 **注意事项:
 **作    者: # Linlin.guo # 2020-5-25 15:15:42 #
 ******************************************************************************/
func ApkPatchUpdateByMap(o orm.Ormer, mapParam map[string]interface{}, apkPatchId int64) error {
	apkPatchTemp := ""
	for k, v := range mapParam {
		switch v := v.(type) {
		case int:
			apkPatchTemp += k + "= " + strconv.Itoa(v) + ","
		case int64:
			apkPatchTemp += k + "= " + strconv.FormatInt(v, 10) + ","
		case string:
			apkPatchTemp += k + "='" + v + "',"
		}
	}

	apkPatchTemp = apkPatchTemp[0 : len(apkPatchTemp)-1]
	sql := fmt.Sprintf(`update %s set %s WHERE id= %d`, "apk_patch", apkPatchTemp, apkPatchId)

	_, err := o.Raw(sql).Exec()
	if nil != err {
		logs.Error("ApkPatchUpdateByMap: err=", err)
		return err
	}
	return nil
}

/******************************************************************************
 **函数名称: DeletePatch
 **功    能:删除关联patch包
 **输入参数:
 ** 	appId: 应用ID
 ** 	verCode: 关联包版本
 **输出参数:
 **返    回:
 **实现描述:
 **注意事项:
 **作    者: # Shuangpeng.guo # 2020-06-17 13:42:44 #
 ******************************************************************************/
func DeletePatch(o orm.Ormer, appId, verCode int64) error {
	sql := fmt.Sprintf(`
	DELETE
	FROM
		%s
	WHERE
		app_id = ?
	AND
		high_version_code = ?
	OR
		low_version_code = ?
	`, UPGRADE_TAB_APK_PATCH)

	if _, err := o.Raw(sql, appId, verCode, verCode).Exec(); nil != err {
		logs.Error("delete apk patch by appId and verCode failed, appId: %d, verCode:%d",
			appId, verCode)
		return err
	}
	return nil
}
