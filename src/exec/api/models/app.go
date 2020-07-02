package models

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"

	"upgrade-api/src/share/database/upgrade"
)

/******************************************************************************
 **函数名称: CreateApp
 **功    能: 创建应用
 **输入参数:
 **     app: 应用信息
 **		devPlats: 关联平台信息
 **输出参数: NONE
 **返    回: 创建Id
 **实现描述:
 **注意事项:
 **作    者: # Shuangpeng.guo # 2020-06-01 14:14:35 #
 ******************************************************************************/
func (mod *Models) CreateApp(app *upgrade.App, devPlatIds []int64) (int64, error) {
	o := orm.NewOrm()
	o.Using(mod.Mysql.AliasName)

	// 开始事物
	err := o.Begin()
	if nil != err {
		logs.Error("Begin transaction failed! errmsg:%s", err.Error())
		return 0, err
	}

	id, err := upgrade.AddApp(o, app)
	if nil != err {
		o.Rollback()
		logs.Error("App create failed, param:%v errmsg:%s", app, err.Error())
		return 0, err
	}
	if app.Enable == 0 || 0 == len(devPlatIds) {
		o.Commit()
		return id, nil
	}

	appPlatRels := make([]*upgrade.AppPlatRel, len(devPlatIds))

	for k, v := range devPlatIds {
		appPlatRels[k] = &upgrade.AppPlatRel{
			AppId:      id,
			DevPlatId:  v,
			CreateTime: time.Now(),
			UpdateTime: time.Now(),
		}
	}
	_, err = o.InsertMulti(len(devPlatIds), appPlatRels)
	if nil != err {
		o.Rollback()
		logs.Error("Insert app plats failed! param:%v errmsg:%s", appPlatRels, err.Error())
		return 0, err
	}

	o.Commit()

	return id, nil
}

/******************************************************************************
 **函数名称: UpdateApp
 **功    能: 修改应用
 **输入参数:
 **     app: 应用信息
 **		devPlats: 关联平台信息
 **输出参数: NONE
 **返    回:
 **实现描述:
 **注意事项:
 **作    者: # Shuangpeng.guo # 2020-06-02 10:14:35 #
 ******************************************************************************/
func (mod *Models) UpdateApp(app *upgrade.App, devPlatIds []int64) error {
	o := orm.NewOrm()
	o.Using(mod.Mysql.AliasName)

	// 开始事物
	err := o.Begin()
	if nil != err {
		logs.Error("Begin transaction failed! errmsg:%s", err.Error())
		return err
	}

	// 更新app信息
	err = upgrade.UpdateAppById(o, app)
	if nil != err {
		o.Rollback()
		logs.Error("App update failed! param:%v errmsg:%s", app, err.Error())
		return err
	}

	// 删除app平台关联信息
	deleteSql := fmt.Sprintf(
		`DELETE
		 FROM %s
         WHERE app_id=?`, upgrade.UPGRADE_TAB_APP_PLAT_REL)
	_, err = o.Raw(deleteSql, app.Id).Exec()
	if nil != err {
		o.Rollback()
		logs.Error("App plat rel delete failed! errmsg:%s", err.Error())
		return err
	}

	if app.HasDevPlat == 0 || 0 == len(devPlatIds) {
		o.Commit()
		return nil
	}

	appPlatRels := make([]*upgrade.AppPlatRel, len(devPlatIds))

	for k, v := range devPlatIds {
		appPlatRels[k] = &upgrade.AppPlatRel{
			AppId:      app.Id,
			DevPlatId:  v,
			CreateTime: time.Now(),
			UpdateTime: time.Now(),
		}
	}

	_, err = o.InsertMulti(len(devPlatIds), appPlatRels)
	if nil != err {
		o.Rollback()
		logs.Error("Insert app plats failed! param:%v errmsg:%s",
			appPlatRels, err.Error())
		return err
	}

	o.Commit()

	return nil
}

type App struct {
	Id           int64     `orm:"column(id)"`
	Name         string    `orm:"column(name)"`
	PackageName  string    `orm:"column(package_name)"`
	BusinessId   int64     `orm:"column(business_id)"`
	CdnPlatId    int64     `orm:"column(cdn_plat_id)"`
	CdnSplatId   int64     `orm:"column(cdn_splat_id)"`
	DevTypeId    int64     `orm:"column(dev_type_id)"`
	Enable       int       `orm:"column(enable)"`
	AppPublicKey string    `orm:"column(app_public_key)"`
	HasDevPlat   int       `orm:"column(has_dev_plat)"`
	DevPlatIds   string    `orm:"column(dev_plat_ids)"`
	Description  string    `orm:"column(description)"`
	CreateTime   time.Time `orm:"column(create_time)"`
	CreateUser   string    `orm:"column(create_user)"`
	UpdateTime   time.Time `orm:"column(update_time)"`
	UpdateUser   string    `orm:"column(update_user)"`
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
func (mod *Models) GetAppAndPlat(id int64) (*App, error) {
	app := new(App)
	appPlatRel := upgrade.AppPlatRel{}
	appTable := upgrade.App{}
	sql := fmt.Sprintf(`SELECT
		a.id,
		a.name,
		a.package_name,
		a.business_id,
		a.cdn_plat_id,
		a.cdn_splat_id,
		a.dev_type_id,
		a.enable,
		a.app_public_key,
		a.has_dev_plat,
		a.description,
		a.create_time,
		a.create_user,
		a.update_time,
		a.update_user,
		(SELECT group_concat(ap.dev_plat_id) FROM %s AS ap WHERE ap.app_id = a.id) AS dev_plat_ids
	FROM
		%s AS a
	WHERE
		a.id = ?
	`, appPlatRel.TableName(), appTable.TableName())

	err := mod.Mysql.O.Raw(sql, id).QueryRow(app)

	return app, err
}

/******************************************************************************
 **函数名称: GetAppAndPlatList
 **功    能: 获取应用及关联平台信息
 **输入参数:
 **     id: 应用ID
 **输出参数: NONE
 **返    回:
 **实现描述:
 **注意事项:
 **作    者: # Shuangpeng.guo # 2020-06-02 10:14:35 #
 ******************************************************************************/
func (mod *Models) GetAppAndPlatList(fields []string,
	values []interface{}, page, pageSize int) (int64, []*App, error) {

	apps := new([]*App)
	appPlatRel := upgrade.AppPlatRel{}
	appTable := upgrade.App{}

	sql := fmt.Sprintf(`SELECT
		a.id,
		a.name,
		a.package_name,
		a.business_id,
		a.cdn_plat_id,
		a.cdn_splat_id,
		a.dev_type_id,
		a.enable,
		a.app_public_key,
		a.has_dev_plat,
		a.description,
		a.create_time,
		a.create_user,
		a.update_time,
		a.update_user,
		(SELECT group_concat(ap.dev_plat_id) FROM %s AS ap WHERE ap.app_id = a.id) AS dev_plat_ids
	FROM
		%s AS a
	`, appPlatRel.TableName(), appTable.TableName())

	sqlCount := fmt.Sprintf(`SELECT
		COUNT(id) as total
	FROM
		%s
	`, appTable.TableName())

	// 如果where条件不为空追加 where 条件
	if 0 < len(fields) {
		where := strings.Join(fields, " AND ")
		sql = fmt.Sprintf("%s WHERE %s ", sql, where)
		sqlCount = fmt.Sprintf("%s WHERE %s ", sqlCount, where)
	}

	res := orm.ParamsList{}
	_, err := mod.Mysql.O.Raw(sqlCount, values).ValuesFlat(&res)
	if nil != err {
		logs.Error("Get app and plat list failed! sql:%s errmsg:%s", sql, err.Error())
		return 0, nil, err
	} else if len(res) == 0 {
		res = append(res, 0)
	}

	total := int64(0)
	totalStr, ok := res[0].(string)
	if ok {
		total, _ = strconv.ParseInt(totalStr, 10, 64)
	}

	// 添加排序&分页条件
	sql = fmt.Sprintf("%s ORDER BY a.id DESC LIMIT ?,? ", sql)
	values = append(values, (page-1)*pageSize, pageSize)

	fmt.Printf("\n%v---%v\n", sql, values)
	_, err = mod.Mysql.O.Raw(sql, values).QueryRows(apps)
	if nil != err {
		logs.Error("Get app and plat rows failed! sql:%s errmsg:%s", sql, err.Error())
		return 0, nil, err
	}

	return total, *apps, err
}
