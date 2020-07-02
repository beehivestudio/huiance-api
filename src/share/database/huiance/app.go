package upgrade

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"

	"upgrade-api/src/share/lib/utils"
)

type App struct {
	Id           int64         `orm:"column(id);auto" description:"ID"`
	Name         string        `orm:"column(name);size(128)" description:"应用名称"`
	PackageName  string        `orm:"column(package_name);size(128)" description:"应用包名（不唯一）；注：存在包名相同，但分属不同的平台"`
	BusinessId   int64         `orm:"column(business_id)" description:"业务ID；所属业务线"`
	CdnPlatId    int64         `orm:"column(cdn_plat_id)" description:"CDN 平台ID"`
	CdnSplatId   int64         `orm:"column(cdn_splat_id)" description:"CDN 子平台ID；SPLATID：子平台ID。用于业务计费，由CDN统一分配"`
	DevTypeId    int64         `orm:"column(dev_type_id)" description:"设备种类ID；外键：DevType.id"`
	Enable       int           `orm:"column(enable)" description:"是否启用；可通过此开关控制应用上下线；0：禁用（下线）；1：启用（上线）"`
	AppPublicKey string        `orm:"column(app_public_key);size(512);null" description:"应用公钥：用于验证文件是否被篡改；注：1.创建应用时，该值为空；2.第一次创建APK时，将应用公钥信息录入此字段。"`
	HasDevPlat   int           `orm:"column(has_dev_plat)" description:"是否关联设备平台；0：不关联。即：适用于所有设备平台；1：关联。即：只适用于指定设备平台。如果关联表中无设备平台，表示所有设备平台均不支持。"`
	Description  string        `orm:"column(description);size(2048);null" description:"描述信息"`
	CreateTime   time.Time     `orm:"column(create_time);type(timestamp);auto_now_add" description:"创建时间"`
	CreateUser   string        `orm:"column(create_user)" description:"页面创建者"`
	UpdateTime   time.Time     `orm:"column(update_time);type(timestamp);auto_now" description:"修改时间"`
	UpdateUser   string        `orm:"column(update_user)" description:"页面修改者"`
	Apks         []*Apk        `orm:"-"`
	DevPlats     []*DevPlatOut `orm:"-"`
}

func (t *App) TableName() string {
	return "app"
}

func init() {
	orm.RegisterModel(new(App))
}

// AddApp insert a new App into database and returns
// last inserted Id on success.
func AddApp(o orm.Ormer, m *App) (id int64, err error) {
	id, err = o.Insert(m)
	return
}

// GetAppById retrieves App by Id. Returns error if
// Id doesn't exist
func GetAppById(o orm.Ormer, id int64) (v *App, err error) {
	v = &App{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

func GetAppByPackageNameLimitOne(o orm.Ormer, pkgName string) (v *App, err error) {
	v = &App{}
	sql := fmt.Sprintf("select * from %s where package_name = ? limit 1", v.TableName())
	if err = o.Raw(sql, pkgName).QueryRow(v); err == nil {
		return v, nil
	}
	return nil, err
}

func GetAppByPkgName(o orm.Ormer, pkgName string) ([]*App, error) {
	var res []*App
	_, err := o.QueryTable(UPGRADE_TAB_APK).Filter("package_name", pkgName).All(&res)
	return res, err
}

// GetAllApp retrieves all App matches certain condition. Returns empty list if
// no records exist
func GetAllApp(o orm.Ormer, query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	qs := o.QueryTable(new(App))
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

	var l []App
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

// UpdateApp updates App by Id and returns error if
// the record to be updated doesn't exist
func UpdateAppById(o orm.Ormer, m *App) (err error) {
	v := App{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteApp deletes App by Id and returns error if
// the record to be deleted doesn't exist
func DeleteApp(o orm.Ormer, id int64) (err error) {
	v := App{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&App{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

type AppRelSql struct {
	Id           int64     `gsp:"Id, pk"`           // appid
	Name         string    `gsp:"Name"`             // app名称
	PackageName  string    `gsp:"PackageName"`      // 应用包名
	BusinessId   int64     `gsp:"BusinessId"`       // 业务线ID
	BusinessName string    `gsp:"BusinessName"`     // 业务线名称
	CdnPlatId    int64     `gsp:"CdnPlatId"`        // cdn平台ID
	CdnSplatId   int64     `gsp:"CdnSplatId"`       // cdn平台子ID
	DevTypeId    int64     `gsp:"DevTypeId"`        // 设备类型ID
	DevTypeName  string    `gsp:"DevTypeName"`      // 设备类型名
	Enable       int       `gsp:"Enable"`           // 是否启用
	AppPublicKey string    `gsp:"AppPublicKey"`     // APK签名公钥
	HasDevPlat   int       `gsp:"HasDevPlat"`       // 是否关联设备平台
	Description  string    `gsp:"Description"`      // 描述信息
	CreateTime   time.Time `gsp:"CreateTime"`       // 创建时间
	CreateUser   string    `gsp:"CreateUser"`       // 创建用户
	UpdateTime   time.Time `gsp:"UpdateTime"`       // 修改时间
	UpdateUser   string    `gsp:"UpdateUser"`       // 修改用户
	DevPlatId    int64     `gsp:"DevPlats__Id, pk"` // 关联 设备平台ID
	DevPlatName  string    `gsp:"DevPlats__Name"`   // 关联设备平台名
}

type AppRel struct {
	Id           int64         `orm:"column(id)"`             // appid
	Name         string        `orm:"column(name)"`           // app名称
	PackageName  string        `orm:"column(package_name)"`   // 应用包名
	BusinessId   int64         `orm:"column(business_id)"`    // 业务线ID
	BusinessName string        `orm:"column(business_name)"`  // 业务线名称
	CdnPlatId    int64         `orm:"column(cdn_plat_id)"`    // cdn平台ID
	CdnSplatId   int64         `orm:"column(cdn_splat_id)"`   // cdn平台子ID
	DevTypeId    int64         `orm:"column(dev_type_id)"`    // 设备类型ID
	DevTypeName  string        `orm:"column(dev_type_name)"`  // 设备类型名
	Enable       int           `orm:"column(enable)"`         // 是否启用
	AppPublicKey string        `orm:"column(app_public_key)"` // APK签名公钥
	HasDevPlat   int           `orm:"column(has_dev_plat)"`   // 是否关联设备平台
	Description  string        `orm:"column(description)"`    // 描述信息
	CreateTime   time.Time     `orm:"column(create_time)"`    // 创建时间
	CreateUser   string        `orm:"column(create_user)"`    // 修改时间
	UpdateTime   time.Time     `orm:"column(update_time)"`    // 创建用户
	UpdateUser   string        `orm:"column(update_user)"`    // 修改用户
	DevPlats     []*DevPlatOut `orm:"-"`                      // 关联设备平台信息
	Apks         []*Apk        `orm:"-"`                      // 关联平台信息
}

type DevPlatOut struct {
	Id   int64  `json:"id"`   // 设备平台ID
	Name string `json:"name"` // 设备平台名
}

/******************************************************************************
 **函数名称: GetAppAndPlatByPkg
 **功    能: 根据包名获取应用及关联平台信息
 **输入参数:
 **     pkgName: 包名
 **输出参数: NONE
 **返    回:
 **实现描述:
 **注意事项:
 **作    者: # Shuangpeng.guo # 2020-06-02 10:14:35 #
 ******************************************************************************/
func GetAppAndPlatByPkg(o orm.Ormer, pkgName string) ([]*AppRel, error) {

	var appSqls []*AppRelSql
	var apps []*AppRel

	sql := fmt.Sprintf(`SELECT
		app.id,
		app.name,
		app.package_name,
		app.business_id,
		app.cdn_plat_id,
		app.cdn_splat_id,
		app.dev_type_id,
		app.enable,
		app.app_public_key,
		app.has_dev_plat,
		app.description,
		app.create_time,
		app.create_user,
		app.update_time,
		app.update_user,
		dev_plat.id as dev_plat_id,
		dev_plat.name as dev_plat_name
	FROM
		%s AS app
	LEFT JOIN %s AS app_plat_rel
		ON app.id = app_plat_rel.app_id
	LEFT JOIN %s AS dev_plat
		ON app_plat_rel.dev_plat_id = dev_plat.id
	WHERE
		app.package_name = ?
	`, UPGRADE_TAB_APP, UPGRADE_TAB_APP_PLAT_REL, UPGRADE_TAB_DEV_PLAT)

	if _, err := o.Raw(sql, pkgName).QueryRows(&appSqls); nil != err {
		return nil, err
	} else if len(appSqls) == 0 {
		return apps, nil
	}

	if err := utils.ReflectStructs(&appSqls, &apps); nil != err {
		return nil, err
	}

	return apps, nil
}

/******************************************************************************
 **函数名称: GetAppAndPlat
 **功    能: 获取应用及关联平台信息
 **输入参数:
 **     id: 应用ID
 **输出参数: NONE
 **返    回:
 **实现描述:
 **注意事项:
 **作    者: # Shuangpeng.guo # 2020-06-02 10:14:35 #
 ******************************************************************************/
func GetAppAndPlatById(o orm.Ormer, id int64) (*AppRel, error) {

	var appSqls []*AppRelSql
	var apps []*AppRel

	sql := fmt.Sprintf(`SELECT
		app.id,
		app.name,
		app.package_name,
		app.business_id,
		business.name as business_name,
		app.cdn_plat_id,
		app.cdn_splat_id,
		app.dev_type_id,
		dev_type.name as dev_type_name,
		app.enable,
		app.app_public_key,
		app.has_dev_plat,
		app.description,
		app.create_time,
		app.create_user,
		app.update_time,
		app.update_user,
		dev_plat.id as dev_plat_id,
		dev_plat.name as dev_plat_name
	FROM
		%s AS app
	LEFT JOIN %s AS business
		ON app.business_id = business.id
	LEFT JOIN %s AS dev_type
		ON app.dev_type_id = dev_type.id
	LEFT JOIN %s AS app_plat_rel
		ON app.id = app_plat_rel.app_id
	LEFT JOIN %s AS dev_plat
		ON app_plat_rel.dev_plat_id = dev_plat.id
	WHERE
		app.id = ?
	`, UPGRADE_TAB_APP, UPGRADE_TAB_BUSINESS, UPGRADE_TAB_DEV_TYPE,
		UPGRADE_TAB_APP_PLAT_REL, UPGRADE_TAB_DEV_PLAT)

	if _, err := o.Raw(sql, id).QueryRows(&appSqls); nil != err {
		return nil, err
	} else if len(appSqls) == 0 {
		return nil, orm.ErrNoRows
	}

	if err := utils.ReflectStructs(&appSqls, &apps); nil != err {
		return nil, err
	}

	return apps[0], nil
}

/******************************************************************************
 **函数名称: GetAppAndPlatList
 **功    能: 获取应用及关联平台信息
 **输入参数:
 **     fields: 查询字段
 **     values: 对应值
 **输出参数: NONE
 **返    回:
 **实现描述:
 **注意事项:
 **作    者: # Shuangpeng.guo # 2020-06-02 10:14:35 #
 ******************************************************************************/
func GetAppAndPlatList(o orm.Ormer, fields []string,
	values []interface{}, page, pageSize int) (int64, []*AppRel, error) {

	var appSqls []*AppRelSql
	var apps []*AppRel

	// 优化分页查询，先根据条件查询限定行，再查询具体数据
	sql := fmt.Sprintf(`SELECT
		app.id,
		app.name,
		app.package_name,
		app.business_id,
		business.name as business_name,
		app.cdn_plat_id,
		app.cdn_splat_id,
		app.dev_type_id,
		dev_type.name as dev_type_name,
		app.enable,
		app.app_public_key,
		app.has_dev_plat,
		app.description,
		app.create_time,
		app.create_user,
		app.update_time,
		app.update_user,
		dev_plat.id as dev_plat_id,
		dev_plat.name as dev_plat_name
	FROM
		%s AS app
	LEFT JOIN %s AS business
		ON app.business_id = business.id
	LEFT JOIN %s AS dev_type
		ON app.dev_type_id = dev_type.id
	LEFT JOIN %s AS app_plat_rel
		ON app.id = app_plat_rel.app_id
	LEFT JOIN %s AS dev_plat
		ON app_plat_rel.dev_plat_id = dev_plat.id
	`, UPGRADE_TAB_APP, UPGRADE_TAB_BUSINESS, UPGRADE_TAB_DEV_TYPE,
		UPGRADE_TAB_APP_PLAT_REL, UPGRADE_TAB_DEV_PLAT)

	sqlCount := fmt.Sprintf(`SELECT
		COUNT(id) as total
	FROM
		%s AS app
	`, UPGRADE_TAB_APP)

	sqlInt := fmt.Sprintf(`SELECT
		id
	FROM
		%s AS app
	`, UPGRADE_TAB_APP)

	// 如果where条件不为空追加 where 条件
	if 0 < len(fields) {
		where := strings.Join(fields, " AND ")
		sqlInt = fmt.Sprintf("%s WHERE %s ", sqlInt, where)
		sqlCount = fmt.Sprintf("%s WHERE %s ", sqlCount, where)
	}

	// 查询总数
	total := &TableTotal{}
	if err := o.Raw(sqlCount, values).QueryRow(total); nil != err {
		return 0, nil, err
	} else if 0 == total.Total {
		return 0, apps, nil
	}

	// 查询索引ID
	ids := orm.ParamsList{}
	sqlInt = fmt.Sprintf("%s ORDER BY app.id DESC LIMIT ? OFFSET ? ", sqlInt)
	valueInt := append(values, pageSize, (page-1)*pageSize)
	if num, err := o.Raw(sqlInt, valueInt).ValuesFlat(&ids); nil != err {
		return 0, nil, err
	} else if 0 == num {
		return total.Total, apps, nil
	}

	// 查询实际数据
	fields = append(fields, fmt.Sprintf("app.id in (%s)",
		strings.Replace(strings.Trim(fmt.Sprint(ids), "[]"), " ", ",", -1)))

	// 添加排序&分页条件
	sql = fmt.Sprintf("%s WHERE %s ORDER BY app.id DESC", sql, strings.Join(fields, " AND "))

	if _num, err := o.Raw(sql, values...).QueryRows(&appSqls); nil != err {
		return 0, nil, err
	} else if 0 == _num {
		return total.Total, apps, nil
	}

	if err := utils.ReflectStructs(&appSqls, &apps); nil != err {
		return 0, nil, err
	}

	return total.Total, apps, nil
}
