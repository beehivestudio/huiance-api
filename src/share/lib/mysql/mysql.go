package mysql

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

/* MYSQL错误码 */
const (
	ERR_SYNTAX             = 1064 // 语法错误
	ERR_SAFE_UPDATE        = 1175 // 不安全的修改,涉及全表删除修改操作
	ERR_DATA_TYPE          = 1292 // 数据类型错误
	ERR_ACCESS_DENIED      = 1045 // 用户访问受限
	ERR_MASTER_NODE        = 1236 // 主节点挂掉且关闭了实时同步
	ERR_TOO_MANY_OPEN_FILE = 24   // 超出系统文件打开限制
	ERR_DUPLICATE_ENTRY    = 1062 // 约束冲突
	ERR_CHARSET            = 1366 // 字符集不一致
	ERR_LIMIT              = 139  // 字段超出限制(数字过大或字符过长)
	ERR_FOREIGN_KEY        = 1215 // 外键约束错误
)

/* MYSQL配置信息 */
type Conf struct {
	Conn string // 数据库访问地址
}

/* 连接池对象 */
type Pool struct {
	AliasName string    // 注册数据库别名
	O         orm.Ormer // 注册数据库连接对象
}

/* 默认数据库名称 (保留) */
const (
	DR_MYSQL           = "mysql"
	DEFAULT_ALIAS_NAME = "default"
)

/******************************************************************************
 **函数名称: RegisterDb
 **功    能: 注册数据库，并返回通用ormer，在使用事务时需要额外注册
 **输入参数:
 **     addr: 数据库连接地址
 **     args: 当args不为空时, 则args[0]表示数据库别名.
 **输出参数: NONE
 **返    回: 连接池数组
 **实现描述:
 **注意事项:
 **作    者: # Shuangpeng.guo # 2018-11-30 10:55:24 #
 ******************************************************************************/
func RegisterDb(addr string, args ...string) {
	/* 注册数据库引擎 */
	// 参数1   driverName
	// 参数2   数据库类型
	// 这个用来设置 driverName 对应的数据库类型
	// mysql / sqlite3 / postgres 这三种是默认已经注册过的，所以可以无需设置
	//orm.RegisterDriver("mymysql", orm.DRMySQL)

	/* 注册数据库 */
	//orm.RegisterDataBase(dbName, "mysql", addr)

	/* 高级设置 设置最大空闲数、连接数 */
	maxIdle := 1000 // 设最大空闲连接
	maxConn := 2048 // 设置最大数据库连接 (go >= 1.2)

	aliasName := DEFAULT_ALIAS_NAME
	if len(args) > 0 {
		aliasName = args[0]
	}

	orm.RegisterDataBase(aliasName, DR_MYSQL, addr, maxIdle, maxConn)
	// orm.SetMaxIdleConns(dbName, maxIdle)  // 额外设置
	// orm.SetMaxOpenConns(dbName, maxConn)  // 额外设置

	/* 设置时区 */
	//orm.DefaultTimeLoc, _ = time.LoadLocation("Asia/Shanghai")
}

/******************************************************************************
 **函数名称: GetMysqlPool
 **功    能: 获取通用ormer，在使用事务时需要额外注册
 **输入参数:
 **     args: 当args不为空时, 则args[0]表示数据库别名.
 **输出参数: NONE
 **返    回: 连接池数组
 **实现描述:
 **注意事项:
 **作    者: # Shuangpeng.guo # 2018-11-30 10:55:24 #
 ******************************************************************************/
func GetMysqlPool(args ...string) *Pool {
	/* 数据库连接对象 */
	o := orm.NewOrm()

	aliasName := DEFAULT_ALIAS_NAME
	if len(args) > 0 {
		aliasName = args[0]
	}

	o.Using(aliasName)

	return &Pool{AliasName: aliasName, O: o}
}
