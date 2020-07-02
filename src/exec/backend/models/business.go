package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"

	"upgrade-api/src/share/comm"
	"upgrade-api/src/share/database/upgrade"
	"upgrade-api/src/share/lib/crypt"
	"upgrade-api/src/share/quota"
)

/******************************************************************************
 **函数名称: DelBusinessById
 **功    能: 删除业务
 **输入参数:
 **     id: 业务id
 **输出参数: NONE
 **返    回: NONE
 **实现描述:
 **注意事项:
 **作    者: # zhao.yang # 2020-06-02 14:12:31 #
 ******************************************************************************/
func (ctx *Models) DelBusinessById(q *quota.QuotaWorker, businessId int64) error {

	o := ctx.Mysql.O

	err := o.Begin()
	if nil != err {
		logs.Error("Begin transaction failed! errmsg:%s", err.Error())
		return err
	}

	//删除业务数据
	err = upgrade.DeleteBusiness(o, businessId)
	if err != nil {
		o.Rollback()
		logs.Error("Del business failed! errmsg:%s id:%d", err.Error(), businessId)
		return err
	}

	delFlowlimitIds, err := upgrade.GetFlowLimitByBusinessId(o, businessId)
	if err != nil && err != orm.ErrNoRows {
		logs.Error("Query flow limit failed! errmsg:%s businessId:%d", err.Error(), businessId)
		return err
	}

	//删除流控数据
	num, err := upgrade.DelFlowLimitByBusinessId(o, businessId)
	logs.Info("Del flow limit num :%d", num)
	if err != nil && orm.ErrNoRows != err {
		o.Rollback()
		logs.Error("Del flow limit failed! errmsg:%s businessId:%d",
			err.Error(), businessId)
		return err
	} else if orm.ErrNoRows == err {
		logs.Error("Del flow limit is empty, businessId :%d", businessId)
		return errors.New("flow limit is empty")
	}

	err = o.Commit()
	if nil != err {
		logs.Error("Close transaction failed! errmsg:%s", err.Error())
		return err
	}

	//同步流控
	for _, flowlimit := range delFlowlimitIds {
		err = q.Sync(getQuotaKey(comm.QUOTA_KEY_TYPE_STRATEGY, businessId), quota.QUOTA_TYPE_STRATEGY, flowlimit.Id)
		if err != nil {
			logs.Error("Business  sync flow limit failed! id:%v errmsg:%s", flowlimit.Id, err.Error())
			return err
		}
	}

	return nil
}

/******************************************************************************
 **函数名称: UpdateBusiness
 **功    能: 修改业务
 **输入参数:
 **     Business: 业务model
 **输出参数: NONE
 **返    回: NONE
 **实现描述:
 **注意事项:
 **作    者: # zhao.yang # 2020-06-02 14:12:31 #
 ******************************************************************************/
func (ctx *Models) UpdateBusiness(q quota.QuotaWorker, req *BusinessReq, id int64, user string) (err error) {

	business := upgrade.Business{
		Id:           id,
		Name:         req.Name,
		Key:          crypt.Md5Sum(fmt.Sprintf("%d", time.Now().UnixNano())),
		Enable:       req.Enable,
		HasFlowLimit: req.HasFlowLimit,
		Description:  req.Description,
		Manager:      req.Manager,
		UpdateUser:   user,
		UpdateTime:   time.Now(),
	}

	o := ctx.Mysql.O

	// 开始事物
	err = o.Begin()
	if nil != err {
		logs.Error("Begin transaction failed! errmsg:%s", err.Error())
		return err
	}

	//修改业务
	err = upgrade.UpdateBusinessById(o, &business)
	if nil != err {
		o.Rollback()
		logs.Warn("Business update failed! param:%v errmsg:%s", business, err.Error())
		return err
	}

	delFlowlimitIds, err := upgrade.GetFlowLimitByBusinessId(o, id)
	if err != nil && err != orm.ErrNoRows {
		logs.Error("Query flow limit failed errmsg:%s businessId:%d", err.Error(), id)
		return err
	}

	//删除流控数据
	num, err := upgrade.DelFlowLimitByBusinessId(o, id)
	logs.Info("Del flow limit num :%d", num)
	if err != nil {
		logs.Error("Del flow limit failed errmsg:%s businessId:%d", err.Error(), id)
		return err
	}

	newFlowlimitIds := []int64{}

	//添加业务流控数据
	for _, businessFlowLimit := range req.FlowLimit {

		businessFlowLimit := upgrade.BusinessFlowLimit{
			BusinessId: id,
			BeginTime:  businessFlowLimit.BeginTime,
			EndTime:    businessFlowLimit.EndTime,
			Dimension:  businessFlowLimit.Dimension,
			Limit:      int64(businessFlowLimit.Limit),
		}

		flowLimitId, err := upgrade.AddBusinessFlowLimit(o, &businessFlowLimit)
		logs.Info("Create business flow limit id :%s", flowLimitId)
		if nil != err {
			o.Rollback()
			logs.Warn("Business flow limitId create failed! param:%v errmsg:%s", businessFlowLimit, err.Error())
			return err
		}
		newFlowlimitIds = append(newFlowlimitIds, flowLimitId)
	}

	err = ctx.Mysql.O.Commit()
	if nil != err {
		logs.Error("Close transaction failed! errmsg:%s", err.Error())
		return err
	}

	//同步流控
	for _, flowlimit := range delFlowlimitIds {
		err = q.Sync(getQuotaKey(comm.QUOTA_KEY_TYPE_STRATEGY, id), quota.QUOTA_TYPE_STRATEGY, flowlimit.Id)
		if err != nil {
			logs.Error("Business  sync flow limit failed!id:%v errmsg:%s", flowlimit.Id, err.Error())
			return err
		}
	}
	for _, flowlimitId := range newFlowlimitIds {
		err = q.Sync(getQuotaKey(comm.QUOTA_KEY_TYPE_STRATEGY, id), quota.QUOTA_TYPE_STRATEGY, flowlimitId)
		if err != nil {
			logs.Error("Business  sync flow limit failed! id:%v errmsg:%s", flowlimitId, err.Error())
			return err
		}
	}
	return
}

/******************************************************************************
 **函数名称: GetBusinessList
 **功    能: 获取业务信息列表
 **输入参数:NONE
 **输出参数: NONE
 **返    回:
 **     BusinessRet 业务列表信息
 **     err 错误信息
 **实现描述:
 **注意事项:
 **作    者: # zhao.yang # 2020-06-03 13:12:15 #
 ******************************************************************************/
func (ctx *Models) GetBusinessList() (rets []BusinessResp, err error) {
	//获取业务列表
	list, err := upgrade.GetBusinessList(ctx.Mysql.O)
	if err != nil {
		logs.Error("Query business list failed! errmsg:%s businessList:%d",
			err.Error(), len(list))
		return nil, err
	}

	for _, business := range list {

		//查询业务下对应的流控数据
		businessFlowLimits, err := upgrade.GetFlowLimitByBusinessId(ctx.Mysql.O, business.Id)
		if nil != err && orm.ErrNoRows != err {
			logs.Error("Get business flow limits failed! businessId:%v errmsg:%s",
				business.Id, err.Error())
			return nil, err
		} else if orm.ErrNoRows == err {
			logs.Error("Get business flow limits is empty ,businessId:%v", business.Id)
		}

		ret := BusinessResp{
			Id:           int(business.Id),
			Name:         business.Name,
			Enable:       business.Enable,
			HasFlowLimit: business.HasFlowLimit,
			Manager:      business.Manager,
			Key:          business.Key,
			FlowLimit:    GetFlowLimit(businessFlowLimits),
			Description:  business.Description,
			CreateUser:   business.CreateUser,
			CreateTime:   business.CreateTime,
			UpdateUser:   business.UpdateUser,
			UpdateTime:   business.UpdateTime,
		}

		rets = append(rets, ret)
		logs.Info("Query business list len :%s", len(rets))
	}

	return rets, nil
}

/* 查询单个业务返回参数 */
type BusinessResp struct {
	Id           int         `json:"id" description:"业务id"`
	Name         string      `json:"name" description:"业务名称"`
	Enable       int         `json:"enable" description:"是否启用"`
	HasFlowLimit int         `json:"has_flow_limit" description:"是否开启流控；有无流控；0：无流控；1：有流控"`
	Manager      string      `json:"manager" description:"项目技术负责人"`
	Key          string      `json:"key" description:"业务key"`
	FlowLimit    []FlowLimit `json:"flow_limit" description:"流控策略"`
	Description  string      `json:"description" description:"描述信息"`
	CreateUser   string      `json:"create_user"`
	CreateTime   time.Time   `json:"create_time"`
	UpdateUser   string      `json:"update_user"`
	UpdateTime   time.Time   `json:"update_time"`
}

/******************************************************************************
 **函数名称: GetBusinessAndPlat
 **功    能: 获取单个业务信息
 **输入参数:
 **     id: 业务ID
 **输出参数: 业务model
 **返    回:
 **实现描述:
 **注意事项:
 **作    者: # zhao.yang # 2020-06-02 15:12:15 #
 ******************************************************************************/
func (ctx *Models) GetBusinessAndPlat(id int64) (*BusinessResp, error) {

	//查询业务数据
	business, err := upgrade.GetBusinessById(ctx.Mysql.O, id)
	if err != nil {
		logs.Error("Get business by id failed! id:%v errmsg:%s", id, err.Error())
		return nil, err
	}

	//查询流控数据
	businessFlowLimits, err := upgrade.GetFlowLimitByBusinessId(ctx.Mysql.O, id)
	if err != nil {
		logs.Error("Get business flow limits failed! businessId:%v errmsg:%s", id, err.Error())
	}

	businessRet := &BusinessResp{
		Id:           int(business.Id),
		Name:         business.Name,
		Enable:       business.Enable,
		HasFlowLimit: business.HasFlowLimit,
		Manager:      business.Manager,
		Key:          business.Key,
		FlowLimit:    GetFlowLimit(businessFlowLimits),
		Description:  business.Description,
		CreateUser:   business.CreateUser,
		CreateTime:   business.CreateTime,
		UpdateUser:   business.UpdateUser,
		UpdateTime:   business.UpdateTime,
	}

	return businessRet, nil
}

/* 查询流控数据 */
func GetFlowLimit(data []upgrade.BusinessFlowLimit) (flowLimits []FlowLimit) {

	var flowLimit FlowLimit
	for _, v := range data {
		flowLimit = FlowLimit{
			BeginTime: v.BeginTime,
			EndTime:   v.EndTime,
			Dimension: v.Dimension,
			Limit:     int(v.Limit),
		}
		flowLimits = append(flowLimits, flowLimit)
	}
	return flowLimits
}

/* 查询业务返回参数 */
type BusinessRet struct {
	Id           int         `json:"id" description:"业务id"`
	Name         string      `json:"name" description:"业务名称"`
	Enable       int         `json:"enable" description:"是否启用"`
	HasFlowLimit int         `json:"has_flow_limit" description:"是否开启流控；有无流控；0：无流控；1：有流控"`
	FlowLimit    []FlowLimit `json:"flow_limit" description:"流控策略"`
	Description  string      `json:"description" description:"描述信息"`
}

/******************************************************************************
 **函数名称: CreateBusiness
 **功    能: 创建业务
 **输入参数:
 **     business: 业务model
 **     user: 创建用户
 **输出参数: NONE
 **返    回: NONE
 **实现描述:
 **注意事项:
 **作    者: # zhao.yang # 2020-06-02 14:12:31 #
 ******************************************************************************/
func (ctx *Models) CreateBusiness(q *quota.QuotaWorker, req *BusinessReq, user string) (int64, error) {

	//生成业务key
	key := crypt.Md5Sum(fmt.Sprintf("%d", time.Now().UnixNano()))

	business := upgrade.Business{
		Name:         req.Name,
		Key:          key,
		Enable:       req.Enable,
		HasFlowLimit: req.HasFlowLimit,
		Manager:      req.Manager,
		Description:  req.Description,
		CreateUser:   user,
		CreateTime:   time.Now(),
	}

	o := ctx.Mysql.O

	// 开始事物
	err := o.Begin()
	if nil != err {
		logs.Error("Begin transaction failed! errmsg:%s", err.Error())
		return 0, err
	}

	//创建业务
	id, err := upgrade.AddBusiness(o, &business)
	if nil != err {
		o.Rollback()
		logs.Warn("Business create failed! param:%v errmsg:%s", business, err.Error())
		return 0, err
	}

	logs.Info("Create business id :%d", id)

	flowlimitIds := []int64{}

	for _, flowLimit := range req.FlowLimit {

		bf := upgrade.BusinessFlowLimit{
			BusinessId: id,
			BeginTime:  flowLimit.BeginTime,
			Dimension:  flowLimit.Dimension,
			EndTime:    flowLimit.EndTime,
			Limit:      int64(flowLimit.Limit),
		}

		//新增业务流控数据
		flowLimitId, err := upgrade.AddBusinessFlowLimit(o, &bf)
		logs.Info("Create business flow limit id :%s", flowLimitId)
		if nil != err {
			o.Rollback()
			logs.Warn("Business flow limitId create failed! param:%v errmsg:%s", bf, err.Error())
			return 0, err
		}

		flowlimitIds = append(flowlimitIds, flowLimitId)

	}

	err = o.Commit()
	if nil != err {
		logs.Error("Close transaction failed! errmsg:%s", err.Error())
		return 0, err
	}

	//同步流控
	for _, flowlimitId := range flowlimitIds {
		err = q.Sync(getQuotaKey(comm.QUOTA_KEY_TYPE_STRATEGY, id), quota.QUOTA_TYPE_STRATEGY, flowlimitId)
		if err != nil {
			logs.Error("Business  sync flow limit failed! id:%v errmsg:%s", flowlimitId, err.Error())
			return 0, err
		}
	}

	return id, nil
}

/******************************************************************************
 **函数名称: getQuotaKey
 **功    能: 获取流控 reids key
 **输入参数: pre: redis key 前缀
 **       : id: strategy_id
 **输出参数: NONE
 **返    回: model: 拼接好的 key
 **实现描述:
 **注意事项:
 **作    者: # taoshengbo # 2020-05-28 19:40:31 #
 ******************************************************************************/
func getQuotaKey(pre string, id int64) string {

	logs.Info("getQuotaKey:", fmt.Sprintf(pre, id))

	return fmt.Sprintf(pre, id)
}

/* 返回业务结参数 */
type BusinessRe struct {
	Name         string    `json:"name" description:"业务名称"`
	Enable       int       `json:"enable" description:"是否启用"`
	HasFlowLimit int       `json:"has_flow_limit" description:"是否开启流控；有无流控；0：无流控；1：有流控"`
	FlowLimit    FlowLimit `json:"flow_limit" description:"流控策略"`
	Description  string    `json:"description" description:"描述信息"`
}

/* 创建业务请求参数 */
type BusinessReq struct {
	Name         string      `json:"name" description:"业务名称"`
	Enable       int         `json:"enable" description:"是否启用"`
	HasFlowLimit int         `json:"has_flow_limit" description:"是否开启流控；有无流控；0：无流控；1：有流控"`
	Manager      string      `json:"manager" description:"项目技术负责人"`
	FlowLimit    []FlowLimit `json:"flow_limit" description:"流控策略"`
	Description  string      `json:"description" description:"描述信息"`
}

/* 流控策略 */
type FlowLimit struct {
	BeginTime int `json:"begin_time" description:"开始时间[0, 23]（时)"`
	EndTime   int `json:"end_time" description:"结束时间[0, 23]（时)"`
	Dimension int `json:"dimension" description:"流控维度（1：秒 2：分 3：时 4：天）"`
	Limit     int `json:"limit" description:"流控值"`
}
