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

type ApkUpload struct {
	Id          int       `orm:"column(id);auto" description:"ID"`
	FileName    string    `orm:"column(file_name);size(255)" description:"文件名"`
	Md5         string    `orm:"column(md5);size(64)" description:"文件MD5值"`
	Size        int64     `orm:"column(size)" description:"文件大小	"`
	PackageName string    `orm:"column(package_name);size(255)" description:"包名"`
	Label       string    `orm:"column(label);size(255)" description:"文件名"`
	VersionCode int64     `orm:"column(version_code)" description:"版本编号"`
	VersionName string    `orm:"column(version_name);size(255)" description:"版本名称"`
	PublicKey   string    `orm:"column(public_key);size(2048)" description:"公钥信息"`
	LocalPath   string    `orm:"column(local_path);size(255)" description:"本地相对存储路径"`
	CdnUrl      string    `orm:"column(cdn_url);size(128)" description:"CDN存储路径"`
	Status      int       `orm:"column(status)" description:"上传状态：0：正在上传；1：上传成功；2：上传失败；3：上传超时"`
	CreateTime  time.Time `orm:"column(create_time);type(timestamp);auto_now_add" description:"创建时间"`
	UpdateTime  time.Time `orm:"column(update_time);type(timestamp);auto_now" description:"修改时间"`
	CreateUser  string    `orm:"column(create_user)" description:"页面创建者"`
	UpdateUser  string    `orm:"column(update_user)" description:"页面修改者"`
}

func (t *ApkUpload) TableName() string {
	return "apk_upload"
}

func init() {
	orm.RegisterModel(new(ApkUpload))
}

const (
	APKFILEUPLOAD_DOING   = 0
	APKFILEUPLOAD_SUCCESS = 1
	APKFILEUPLOAD_FAIL    = 2
	APKFILEUPLOAD_TIMEOUT = 3
)

// AddApkUpload insert a new ApkUpload into database and returns
// last inserted Id on success.
func AddApkUpload(o orm.Ormer, m *ApkUpload) (id int64, err error) {
	id, err = o.Insert(m)
	return id, err
}

// GetApkUploadById retrieves ApkUpload by Id. Returns error if
// Id doesn't exist
func GetApkUploadById(o orm.Ormer, id int) (v *ApkUpload, err error) {
	v = &ApkUpload{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllApkUpload retrieves all ApkUpload matches certain condition. Returns empty list if
// no records exist
func GetAllApkUpload(o orm.Ormer, query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	qs := o.QueryTable(new(ApkUpload))
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

	var l []ApkUpload
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

// UpdateApkUpload updates ApkUpload by Id and returns error if
// the record to be updated doesn't exist
func UpdateApkUploadById(o orm.Ormer, m *ApkUpload) (err error) {
	v := ApkUpload{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return err
}

// DeleteApkUpload deletes ApkUpload by Id and returns error if
// the record to be deleted doesn't exist
func DeleteApkUpload(o orm.Ormer, id int) (err error) {
	v := ApkUpload{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&ApkUpload{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return err
}

func GetAllApkUploadByFuzzy(o orm.Ormer, values []interface{}, fields []string, page, pageSize int) (res map[string]interface{},
	err error) {
	var aus []ApkUpload
	var apkUpload = ApkUpload{}
	sql := fmt.Sprintf("select * from %s ", apkUpload.TableName())

	if 0 < len(fields) {
		where := strings.Join(fields, " AND ")
		sql = fmt.Sprintf("%s where %s ", sql, where)
	}

	// 添加排序&分页条件
	sql = fmt.Sprintf("%s ORDER BY create_time DESC LIMIT ? OFFSET ? ", sql)
	values = append(values, pageSize, (page-1)*pageSize)

	if _, err := o.Raw(sql, values...).QueryRows(&aus); nil != err {
		return nil, err
	}

	res = map[string]interface{}{
		"page":      page,
		"page_size": pageSize,
		"total":     len(aus),
		"len":       len(aus),
		"list":      aus,
	}
	return res, nil
}

/******************************************************************************
 **函数名称: UpdateApkUpload
 **功    能: 更新apk上传url和状态
 **输入参数:
 **     out: 异常返回
 **输出参数:
 **     error: 异常返回
 **返    回:
 **实现描述:
 **注意事项:
 **作    者: # peng.jia # 2020-06-02 10:14:35 #
 ******************************************************************************/
func UpdateApkUpload(o orm.Ormer, url string, outkey string) error {
	id, err := strconv.ParseInt(outkey[4:len(outkey)], 10, 64)
	if err != nil {
		return err
	}

	apkUpload := &ApkUpload{Id: int(id)}
	err = o.Read(apkUpload)
	if err != nil {
		errMsg := fmt.Sprintf("apkUpload not exist, errMsg:%v", id)
		return errors.New(errMsg)
	}

	app := &App{}
	sql := fmt.Sprintf("select * from %s where package_name = ? limit 1", app.TableName())
	if err = o.Raw(sql,
		apkUpload.PackageName).QueryRow(app); err != nil {
		errMsg := fmt.Sprintf("app not exist, package name:%s errMsg:%v", apkUpload.PackageName, err)
		return errors.New(errMsg)
	}
	//url = fmt.Sprintf("%s?platid=%d&splatid=%d", url, app.CdnPlatId, app.CdnSplatId)
	_, err = o.Raw("update apk_upload set cdn_url= ? , status = ? where id = ? ",
		url, APKFILEUPLOAD_SUCCESS, id).Exec()
	if err != nil {
		return err
	}
	return nil
}

/******************************************************************************
 **函数名称: GetApkUploadByMd5
 **功    能: 根据apk的md5查询对应的apk包的存放地址
 **输入参数:
 **      o: orm.Ormer
 **      md5: md5
 **      packageName: 包名
 **      versionCode: 版本号
 **输出参数:
 **返    回:
 **      apkUploadList: apk下载包
 **      err: 错误信息
 **实现描述:
 **注意事项:
 **作    者: # Linlin.guo # 2020-6-1 13:15:42 #
 ******************************************************************************/
func GetApkUploadByMd5(o orm.Ormer, md5, packageName string, versionCode int64) (apkUpload ApkUpload, err error) {

	logs.Info("GetApkUploadByMd5 request param md5: %s", md5)

	qs := o.QueryTable(new(ApkUpload)).Filter("md5", md5).Filter("package_name", packageName).
		Filter("version_code", versionCode).Filter("status", APKFILEUPLOAD_SUCCESS)

	err = qs.One(&apkUpload)
	if err != nil {
		logs.Error("GetApkUploadByMd5 err: %s", err.Error())
		return apkUpload, err
	}

	return apkUpload, nil
}

//根据apk的md5,包名,版本编号,查询对应的apk上传信息
func GetApkUpload(o orm.Ormer, md5, packageName string, versionCode int64) (apkUpload ApkUpload, err error) {

	qs := o.QueryTable(new(ApkUpload)).Filter("md5", md5).Filter("package_name", packageName).Filter("version_code",
		versionCode)

	_, err = qs.All(&apkUpload)
	if err != nil {
		logs.Error("Get apk upload errmsg: %s", err.Error())
		return apkUpload, err
	}

	return apkUpload, nil
}

//更改文件状态
func UpdateApkStatus(o orm.Ormer, id, status int) (err error) {
	_, err = o.Raw("update apk_upload set status = ? where id = ? ", status, id).Exec()
	if err != nil {
		logs.Error("Update apk status failed! errmsg :%s", err.Error())
		return err
	}
	return nil
}
