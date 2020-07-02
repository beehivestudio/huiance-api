package rds

import (
	"time"

	"github.com/garyburd/redigo/redis"
)

/* redis配置 */
type Conf struct {
	Conn    string // 地址
	Pwd     string // 密码
	MaxIdel int    // 最大空闲
}

/* 连接池对象 */
type Pool struct {
	session *redis.Pool /* 连接对象 */
}

/******************************************************************************
 **函数名称: CreatePool
 **功    能: 创建连接池
 **输入参数:
 **     addr: IP地址
 **     passwd: 登录密码
 **     max_idle: 最大空闲连接数
 **     max_active: 最大激活连接数
 **输出参数: NONE
 **返    回:
 **     pool: 连接池对象
 **实现描述:
 **注意事项:
 **     1. 如果max配置过小, 可能会出现连接池耗尽的情况.
 **     2. 如果idle配置过小, 可能会出现大量'TIMEWAIT'的TCP状态.
 **作    者: # Qifeng.zou # 2017.03.30 22:18:34 #
 ******************************************************************************/
func CreatePool(addr string, passwd string, max_idle int) *Pool {
	pool := &redis.Pool{
		MaxIdle: max_idle,
		Wait:    true,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", addr)
			if nil != err {
				//panic(err.Error())
				return nil, err
			}

			if 0 != len(passwd) {
				if _, err := c.Do("AUTH", passwd); nil != err {
					c.Close()
					//panic(err.Error())
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
	return &Pool{session: pool}
}

/******************************************************************************
 **函数名称: Get
 **功    能: 获取连接
 **输入参数: NONE
 **输出参数: NONE
 **返    回: 连接对象
 **实现描述:
 **注意事项:
 **作    者: # Zengyao.pang # 2017.06.12 23:28:27 #
 ******************************************************************************/
func (pool *Pool) Get() redis.Conn {
	return pool.session.Get()
}
