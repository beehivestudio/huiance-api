package db_admin

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

type RolePageRel struct {
	Id           int64     `orm:"column(id);auto" description:"ID"`
	RoleId       int64     `orm:"column(role_id)" description:"角色ID；外键：role.id"`
	PageId       int64     `orm:"column(page_id)" description:"页面ID；外键：page.id"`
	CreateTime   time.Time `orm:"column(create_time);type(timestamp);auto_now_add" description:"创建时间"`
	CreateUserId int64     `orm:"column(create_user_id)" description:"页面创建者"`
}

func (t *RolePageRel) TableName() string {
	return "role_page_rel"
}

func init() {
	orm.RegisterModel(new(RolePageRel))
}

// AddRolePageRel insert a new RolePageRel into database and returns
// last inserted Id on success.
func AddRolePageRel(o orm.Ormer, m *RolePageRel) (id int64, err error) {
	id, err = o.Insert(m)
	return
}

// GetRolePageRelById retrieves RolePageRel by Id. Returns error if
// Id doesn't exist
func GetRolePageRelById(o orm.Ormer, id int64) (v *RolePageRel, err error) {
	v = &RolePageRel{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllRolePageRel retrieves all RolePageRel matches certain condition. Returns empty list if
// no records exist
func GetAllRolePageRel(o orm.Ormer, query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	qs := o.QueryTable(new(RolePageRel))
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

	var l []RolePageRel
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

// UpdateRolePageRel updates RolePageRel by Id and returns error if
// the record to be updated doesn't exist
func UpdateRolePageRelById(o orm.Ormer, m *RolePageRel) (err error) {
	v := RolePageRel{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteRolePageRel deletes RolePageRel by Id and returns error if
// the record to be deleted doesn't exist
func DeleteRolePageRel(o orm.Ormer, id int64) (err error) {
	v := RolePageRel{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&RolePageRel{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

/******************************************************************************
 **函数名称: GetAllPageByRoleId
 **功    能: 根据角色Id查询页面权限列表
 **输入参数:
 **      o: orm.Ormer
 **      roleId: 角色Id
 **输出参数: NONE
 **返    回:
 **      rolePageRel: 角色页面权限列表
 **      err: 错误信息
 **实现描述:
 **注意事项:
 **作    者: # Linlin.guo # 2020-6-18 15:15:42 #
 ******************************************************************************/
func GetAllPageByRoleId(o orm.Ormer, roleId int64) (rolePageRel []RolePageRel, err error) {

	if _, err := o.QueryTable(new(RolePageRel)).Filter("role_id", roleId).All(&rolePageRel); nil != err {

		logs.Error("GetDevPlatByDevTypeId failed! err: %s", err.Error())
		return nil, err
	}

	return rolePageRel, nil
}

/******************************************************************************
 **函数名称: DelRolePageRelByRoleId
 **功    能: 根据roleId删除对应页面列表
 **输入参数:
 **      o: orm.Ormer
 **      roleId: 角色Id
 **输出参数: NONE
 **返    回:
 **      err: 错误信息
 **实现描述:
 **注意事项:
 **作    者: # Linlin.guo # 2020-6-18 15:15:42 #
 ******************************************************************************/
func DelRolePageRelByRoleId(o orm.Ormer, roleId int64) (err error) {

	if _, err := o.QueryTable(new(RolePageRel)).Filter("role_id", roleId).Delete(); nil != err {

		logs.Error("GetDevPlatByDevTypeId failed! err: %s", err.Error())
		return err
	}
	return nil
}

/******************************************************************************
 **函数名称: DelRolePageRelByRoleId
 **功    能: 根据roleId删除对应页面列表
 **输入参数:
 **      o: orm.Ormer
 **      roleId: 角色Id
 **输出参数: NONE
 **返    回:
 **      err: 错误信息
 **实现描述:
 **注意事项:
 **作    者: # Linlin.guo # 2020-6-18 15:15:42 #
 ******************************************************************************/
func GetRolePages(o orm.Ormer, roleId int64, pagesName string) (pages []Page, err error) {

	sql := fmt.Sprintf(`SELECT
	p.id,
	p.name,
	p.path,
	p.description 
    FROM
	role_page_rel rp
	LEFT JOIN page p ON rp.page_id = p.id 
    WHERE rp.role_id = ? `)
	if len(pagesName) > 0 {
		sql += "AND p.name like '%" + pagesName + "%'"
	}

	_, err = o.Raw(sql, roleId).QueryRows(&pages)
	if err != nil {
		logs.Error("Get role pages list failed! errmsg: %s", err.Error())
		return nil, err
	}

	return pages, nil
}
