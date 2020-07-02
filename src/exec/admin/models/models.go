package models

import (
	"upgrade-api/src/share/lib/mysql"
)

type Models struct {
	Mysql *mysql.Pool // mysql 连接池
}
