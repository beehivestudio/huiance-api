package models

import (
	"upgrade-api/src/share/cdn"
	"upgrade-api/src/share/lib/mysql"
	"upgrade-api/src/share/lib/rds"
	"upgrade-api/src/share/lib/worker"
	"upgrade-api/src/share/quota"
)

type Models struct {
	Mysql    *mysql.Pool        // MYSQL连接池
	Redis    *rds.Pool          // redis连接池
	CdnSplat *cdn.Splat         // cdn splat句柄
	Worker   *worker.Worker     // 全局worker池
	Quota    *quota.QuotaWorker // 频控
}
