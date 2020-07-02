package controllers

import (
	"fmt"
	"strconv"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/garyburd/redigo/redis"

	"upgrade-api/src/share/comm"
	"upgrade-api/src/share/database/upgrade"
	"upgrade-api/src/share/lib/rds"
	"upgrade-api/src/share/strategy"
)

type ApiBaseController struct {
	comm.BaseController
}

/* 请求预处理 */
// @author # taoshengbo # 2020-06-08 16:36:35 #
func (c *ApiBaseController) Prepare() {

	/* 获取参数信息 */
	_timestamp := c.Ctx.Request.Header.Get("timestamp")
	_businessId := c.Ctx.Request.Header.Get("business-id")
	sign := c.Ctx.Request.Header.Get("sign")

	timestamp, _ := strconv.ParseInt(_timestamp, 10, 64)
	businessId, _ := strconv.ParseInt(_businessId, 10, 64)
	if businessId == 0 || timestamp == 0 || sign == "" {
		logs.Warn("Header is invalid! businessId:%d timestamp:%d sign:%s",
			businessId, timestamp, sign)
		c.ErrorMessage(comm.ERR_AUTH, `Header is invalid!`)
		return
	}

	now := time.Now().Unix()
	logs.Info("timestamp=%v, now=%v, businessId=%v, sign=%v", timestamp, now, businessId, sign)

	if now < timestamp - comm.TIMESTAMP_VERIFY || now > timestamp + comm.TIMESTAMP_VERIFY {
		c.ErrorMessage(comm.ERR_AUTH, "Timestamp overtime")
		return
	}

	business, err := GetBusinessCache(businessId)
	if nil != err {
		logs.Error("Get business cache failed! businessId:%d errmsg:%s",
			businessId, err.Error())
		c.ErrorMessage(comm.ERR_AUTH, err.Error())
		return
	} else if 1 != business.Enable {
		logs.Warn("Business is disable! businessId:%d", businessId)
		c.ErrorMessage(comm.ERR_AUTH, "Business is disable!")
		return
	}

	signBody := map[string]interface{}{}
	signBody["body"] = string(c.Ctx.Input.RequestBody)
	signBody["timestamp"] = timestamp
	signBody["url"] = c.Ctx.Request.RequestURI

	logs.Info("Get sign body msg :%v, business key msg :%s", signBody, business.Key)

	sign2, err := comm.SignMapByMd5(signBody, business.Key, true)
	if nil != err {
		logs.Error("Sign md5 failed! body:%s errmsg:%s",
			c.Ctx.Input.RequestBody, err.Error())
		c.ErrorMessage(comm.ERR_AUTH, err.Error())
		return
	} else if sign != sign2 {
		logs.Error("Sign is invalid! body:%s sign:%s sign2:%s",
			c.Ctx.Input.RequestBody, sign, sign2)
		c.ErrorMessage(comm.ERR_AUTH, "Sign is invalid!")
		return
	}
}

func GetBusinessCache(businessId int64) (business *upgrade.Business, err error) {

	/**********************************************缓存***********************************************/
	key := fmt.Sprintf(comm.RDS_KEY_BUSINESS, businessId)

	conn := ApiCntx.Model.Redis.Get()
	defer conn.Close()
	err = rds.RedisGetJsonData(conn, key, &business)
	if nil != err && err != redis.ErrNil {
		return nil, err
	}
	if business != nil {
		return business, nil
	}

	defer func() {
		_ = rds.RedisSaveJsonData(conn, key, strategy.GetRedisExpireTime(), &business)
	}()

	/**********************************************缓存***********************************************/

	business, err = upgrade.GetBusinessById(ApiCntx.Model.Mysql.O, businessId)
	if nil != err {
		logs.Error("Get business by id failed! businessId:%d errmsg:%s",
			businessId, err.Error())
		return nil, err
	}

	return business, nil
}

