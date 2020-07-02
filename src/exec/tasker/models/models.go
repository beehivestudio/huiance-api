package models

import (
	"upgrade-api/src/share/lib/mysql"
	"upgrade-api/src/share/lib/rds"
)

type Models struct {
	Mysql *mysql.Pool // MySQL连接池
	Redis *rds.Pool   // Redis连接池
}
