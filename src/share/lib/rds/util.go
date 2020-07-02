package rds

import (
	"errors"
	"fmt"

	"github.com/astaxie/beego/logs"
	"github.com/garyburd/redigo/redis"
	"github.com/json-iterator/go"

	"upgrade-api/src/share/comm"
)

/******************************************************************************
 **函数名称: RedisLPUSH
 **功    能: RedisLPUSH
 **输入参数: conn: redis.Conn
 **       : key: redis_key
 **       : value: value
 **输出参数: NONE
 **返    回:
 **实现描述:
 **注意事项:
 **作    者: # taoshengbo # 2020-05-29 20:06:06 #
 ******************************************************************************/
func RedisLPUSH(conn redis.Conn, key string, value interface{}) (err error) {
	_, err = conn.Do("lpush", key, value)
	if nil != err {
		return err
	}

	return nil
}

/******************************************************************************
 **函数名称: RedisLPUSHJsonData
 **功    能: 任意结构体存储到 Redis LIST
 **输入参数: conn: redis.Conn
 **       : key: redis_key
 **       : value: 任意结构体
 **输出参数: NONE
 **返    回:
 **实现描述:
 **注意事项:
 **作    者: # taoshengbo # 2020-05-29 20:06:06 #
 ******************************************************************************/
func RedisLPUSHJsonData(conn redis.Conn, key string, value interface{}) (err error) {
	jsonByteData, err := jsoniter.Marshal(&value)
	if nil != err {
		return err
	}
	err = RedisLPUSH(conn, key, string(jsonByteData))
	if nil != err {
		return err
	}

	return nil
}

/******************************************************************************
 **函数名称: RedisRPOP
 **功    能: RedisRPOP
 **输入参数: conn: redis.Conn
 **       : key: redis_key
 **输出参数: NONE
 **返    回: res: 取出字符串
 **实现描述:
 **注意事项:
 **作    者: # taoshengbo # 2020-05-29 20:06:06 #
 ******************************************************************************/
func RedisRPOP(conn redis.Conn, key string) (res string, err error) {
	res, err = redis.String(conn.Do("rpop", key))
	if nil != err {
		return "", err
	}

	return res, nil
}

/******************************************************************************
 **函数名称: RedisSaveJsonData
 **功    能: Redis 存储结构体
 **输入参数: conn: redis.Conn
 **       : key: redis_key
 **       : second: 过期秒数
 **       : structData: 结构体
 **输出参数: NONE
 **返    回:
 **实现描述:
 **注意事项:
 **作    者: # taoshengbo # 2020-05-29 20:08:43 #
 ******************************************************************************/
func RedisSaveJsonData(conn redis.Conn, key string, second int, structData interface{}) (err error) {
	defer func() {
		logs.Info("RedisSaveJsonData key=", key, " second=", second, " minute=", second/60, " structData=", structData)
		if nil != err {
			logs.Error("RedisSaveJsonData err=", err, " key=", key)
		}
	}()

	jsonByteData, err := jsoniter.Marshal(&structData)
	if nil != err {
		return err
	}

	_, err = conn.Do("setex", key, second, string(jsonByteData))
	if nil != err {
		return err
	}

	return nil
}

/******************************************************************************
 **函数名称: RedisGetJsonData
 **功    能: Redis 取出结构体
 **输入参数: conn: redis.Conn
 **       : key: redis_key
 **       : structData: 结构体
 **输出参数: structData: 结构体
 **返    回:
 **实现描述:
 **注意事项:
 **作    者: # taoshengbo # 2020-05-29 20:09:12 #
 ******************************************************************************/
func RedisGetJsonData(conn redis.Conn, key string, structData interface{}) (err error) {
	defer func() {
		logs.Info("RedisGetJsonData key=", key, " structData=", structData)
		if nil != err && err != redis.ErrNil {
			logs.Error("RedisGetJsonData err=", err, " key=", key)
		}
	}()

	var data string
	data, err = redis.String(conn.Do("get", key))
	if nil != err {
		return err
	}

	if data == comm.RDS_NULL {
		err = errors.New(fmt.Sprintf(comm.RDS_DATA_IS_NULL+", key=%v", key))
		return err
	}

	err = jsoniter.Unmarshal([]byte(data), &structData)
	if nil != err {
		return err
	}

	return nil
}

/******************************************************************************
 **函数名称: DelCache
 **功    能: 删除缓存
 **输入参数:
 **       : key: redis_key
 **输出参数:
 **返    回:
 **实现描述:
 **注意事项:
 **作    者: # taoshengbo # 2020-06-19 19:18:33 #
 ******************************************************************************/
func DelCache(conn redis.Conn, key string) (err error) {
	_, err = conn.Do("del", key)
	if err != nil {
		return err
	}
	return nil
}
