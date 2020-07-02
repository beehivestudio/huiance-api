package main

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"testing"
	"time"
	"upgrade-api/src/exec/api/controllers"
	"upgrade-api/src/share/comm"
	"upgrade-api/src/share/database/upgrade"
	"upgrade-api/src/share/lib/rds"
	"upgrade-api/src/share/strategy"
)

var ctx = &controllers.ApiContext{}

func init() {
	ctx = _init()
	if nil == ctx {
		fmt.Printf("Initialize context failed!\n")
		return
	}
	register(ctx)

	/* > quota */
	//quota.Worker.Load(&quota.Component{
	//	Mysql: ctx.Mysql,
	//	Redis: ctx.Redis,
	//})
	//defer quota.Worker.ShutDown()
}

func Test_strategy_create(t *testing.T) {

	var limitStrategy []*upgrade.ApkUpgradeFlowLimitStrategy
	limitStrategy = append(limitStrategy, &upgrade.ApkUpgradeFlowLimitStrategy{BeginTime: 1, EndTime: 5, Limit: 10})
	limitStrategy = append(limitStrategy, &upgrade.ApkUpgradeFlowLimitStrategy{BeginTime: 6, EndTime: 10, Limit: 5})

	upgradeStrategy := upgrade.ApkUpgradeStrategy{
		ApkId:         1,
		Enable:        1,
		BeginDatetime: time.Now(),
		EndDatetime:   time.Now().AddDate(0, 0, 7),
		HasFlowLimit:  1,
		UpgradeType:   1,
		//UpgradeDevType:1,
		//UpgradeDevData:"  aaaaa,  bbbb,  cccc  ",
		//		UpgradeDevType:2,
		//		UpgradeDevData:`{
		//	"plats": [{
		//		"id": 938,
		//		"models": [{
		//			"model": "Y55C",
		//			"group_type": 2,
		//			"group_ids": ["aaa", "bbb"]
		//		}]
		//	}, {
		//		"id": 928,
		//		"models": [{
		//			"model": "Y32",
		//			"group_type": 1,
		//			"group_ids": []
		//		}]
		//	}]
		//}`,
		UpgradeDevType:              3,
		UpgradeDevData:              " 1111,  2222,  3333  ",
		Description:                 "空",
		ApkUpgradeFlowLimitStrategy: limitStrategy,
	}

	id, err := strategy.Create(ctx.Model.Mysql, ctx.Model.Redis, ctx.Model.Quota, &upgradeStrategy)

	logs.Info("id=", id, " err=", err)

}

func Test_strategy_update(t *testing.T) {

	var limitStrategy []*upgrade.ApkUpgradeFlowLimitStrategy
	limitStrategy = append(limitStrategy, &upgrade.ApkUpgradeFlowLimitStrategy{BeginTime: 1, EndTime: 5, Limit: 10})
	limitStrategy = append(limitStrategy, &upgrade.ApkUpgradeFlowLimitStrategy{BeginTime: 6, EndTime: 10, Limit: 5})

	upgradeStrategy := upgrade.ApkUpgradeStrategy{
		Id:            31,
		ApkId:         12,
		Enable:        1,
		BeginDatetime: time.Now(),
		EndDatetime:   time.Now().AddDate(0, 0, 7),
		HasFlowLimit:  1,
		UpgradeType:   1,
		//UpgradeDevType:1,
		//UpgradeDevData:"  aaaaa,  bbbb,  cccc  ",
		UpgradeDevType: 2,
		UpgradeDevData: `{
	"plats": [{
		"id": 11938111,
		"models": [{
			"model": "Y55C",
			"group_type": 2,
			"group_ids": ["aaa", "bbb"]
		}]
	}, {
		"id": 928,
		"models": [{
			"model": "Y32",
			"group_type": 1,
			"group_ids": []
		}]
	}]
}`,
		Description:                 "空",
		ApkUpgradeFlowLimitStrategy: limitStrategy,
	}

	err := strategy.Update(ctx.Model.Mysql, ctx.Model.Redis, ctx.Model.Quota, &upgradeStrategy)

	logs.Info(" err=", err)
}

func Test_GetStrategyById(t *testing.T) {
	v, err := strategy.GetStrategyById(ctx.Model.Mysql.O, 34)
	logs.Info(v, err)
	for _, _v := range v.ApkUpgradeFlowLimitStrategy {
		logs.Info(_v)
	}
}

func Test_GetStrategyListByApkId(t *testing.T) {
	v, err := strategy.GetStrategyListByApkId(ctx.Model.Mysql.O, 10)
	logs.Info(v, err)
	for _, _v := range v {
		logs.Info(_v)
		for _, __v := range _v.ApkUpgradeFlowLimitStrategy {
			logs.Info(__v)
		}
	}
}

func Test_tvGroup(t *testing.T) {

	//mac := "0013a0702601"
	////mac = "aa"
	//
	//model, groupId, err := strategy.GetTvModelAndGroupId(mac)
	//logs.Info(model, groupId, err)

}

func Test_row(t *testing.T) {
	sql_dev_list := `select id from apk_upgrade_dev_list where strategy_id = ? and dev_id = ? `

	var ids orm.ParamsList

	result, err := ctx.Model.Mysql.O.Raw(sql_dev_list, 32, "aaa").ValuesFlat(&ids)
	logs.Info("result, err = ", result, err)
	logs.Info("rows, err = ", ids, len(ids), err)
}

func Test_CreateApkPatch(t *testing.T) {
	//apk_patch, err := strategy.CreateApkPatch(ctx.Mysql.O, 1, 100, 50, 1, "", "")
	//logs.Info("apk_patch, err:", apk_patch, err)
}

func Test_Redis(t *testing.T) {
	conn := ctx.Model.Redis.Get()
	defer conn.Close()

	key := "test"

	err := rds.RedisLPUSH(conn, key, 1)
	logs.Info("err=", err)
	err = rds.RedisLPUSH(conn, key, 2)
	logs.Info("err=", err)
	err = rds.RedisLPUSH(conn, key, 3)
	logs.Info("err=", err)

	res, err := rds.RedisRPOP(conn, key)
	logs.Info("res, err ", res, err)
	//res, err = strategy.RedisBRPOP(conn, key)
	//logs.Info("res, err ", res, err)
}

func Test_RedisSaveJsonData(t *testing.T) {
	conn := ctx.Model.Redis.Get()
	defer conn.Close()
	key := "test1"

	upgradeInfo := &strategy.UpgradeInfo{
		BusinessId: 123,
		PatchAlgo:  0,
		DevTypeId:  1,
		DevId:      "0013a0702601",
		DevModel:   "Y32",
		DevPlat:    "918",
		EuiVersion: "100",
		//AppointVersion:  103,
		PackageName:    "com.my.test",
		ApkVersionCode: 0,
		ApkMd5:         "",
	}

	err := rds.RedisSaveJsonData(conn, key, 1000, upgradeInfo)
	logs.Info("err= ", err)

}

func Test_RedisGetJsonData(t *testing.T) {
	conn := ctx.Model.Redis.Get()
	defer conn.Close()
	key := "test1"

	upgradeInfo := &strategy.UpgradeInfo{}

	err := rds.RedisGetJsonData(conn, key, upgradeInfo)
	logs.Info("err=", err, upgradeInfo)

}

func Test_GetBusinessById(t *testing.T) {
	business, err := upgrade.GetBusinessById(ctx.Model.Mysql.O, 11111)
	logs.Info(business, err)
}

func Test_ApkUpgradeOne(t *testing.T) {

	upgradeInfo := strategy.UpgradeInfo{
		BusinessId: 123,
		PatchAlgo:  1,
		DevTypeId:  1,
		DevId:      "0013a0702601",
		DevModel:   "Y32",
		DevPlat:    "918",
		EuiVersion: "100",
		//AppointVersion:  103,
		PackageName:    "com.my.test",
		ApkVersionCode: 0,
		ApkMd5:         "",
	}

	resp, err := strategy.ApkUpgradeOne(ctx.Model.Mysql.O, ctx.Model.Redis, ctx.Model.Quota, comm.UPGRADE_SILENCE, 1, &upgradeInfo)
	logs.Info("resp, err :", resp, err)
}

func Test_ApkUpgradeBatch(t *testing.T) {

	//upgradeInfo := &strategy.UpgradeInfo{
	//	BusinessId:   123,
	//	PatchAlgo:    0,
	//	DevTypeId: 1,
	//	DevId:     "0013a0702601",
	//	DevModel:  "Y32",
	//	DevPlat:   "918",
	//	EuiVersion:   "100",
	//	//AppointVersion:  103,
	//	PackageName: "com.my.test",
	//	ApkVersionCode:  0,
	//	ApkMd5:      "",
	//}
	//
	//var upgradeInfos []*strategy.UpgradeInfo
	//
	//upgradeInfos = append(upgradeInfos, upgradeInfo)
	//
	//resp, err := strategy.ApkUpgradeBatch(ctx.Model.Mysql.O, ctx.Model.Redis, ctx.Model.Quota, comm.UPGRADE_SILENCE,1,1, upgradeInfos)
	//logs.Info("resp, err :", resp, err)
	//for i, v := range resp {
	//	logs.Info("----: i=", i, " v=", v)
	//}
}

func Test_CacheModel(t *testing.T) {
	err := strategy.CacheModel(ctx.Model.Mysql.O, ctx.Model.Redis)
	logs.Info("--:", err)
}

func Test_GetModelFromCacheById(t *testing.T) {
	model, err := strategy.GetModelFromCacheById(ctx.Model.Redis, 11660)
	logs.Info("model, err", model, err)
}

func Test_CacheGroup(t *testing.T) {
	err := strategy.CacheGroup(ctx.Model.Mysql.O, ctx.Model.Redis)
	logs.Info("-----:", err)
}

func Test_GetGroupFromCacheById(t *testing.T) {
	group, err := strategy.GetDmsGroupIdFromCacheById(ctx.Model.Redis, 6496)
	logs.Info("-----:", group, err)
}

func Test_CachePlat(t *testing.T) {
	err := strategy.CachePlat(ctx.Model.Mysql.O, ctx.Model.Redis)
	logs.Info("----:", err)
}

func Test_GetPlatCacheById(t *testing.T) {
	plat, err := strategy.GetPlatCacheById(ctx.Model.Redis, 52888)
	logs.Info("---:", plat, err)
}

func Test_DealEui(t *testing.T) {
	fmt.Println(strategy.DealEui("5.8"))
	fmt.Println(strategy.DealEui("5.8.020"))
	fmt.Println(strategy.DealEui("V2406RCN02C060256B08011S"))
}

func Test_DelCache(t *testing.T) {
	err := rds.DelCache(ctx.Model.Redis.Get(), "test")
	fmt.Println("err:", err)
}

func Test_A(t *testing.T) {
	a := new(map[string]string)
	fmt.Println("---", a, &a, *a)
}

func Test_GetAppApkPlatsInfo(t *testing.T) {
	upgradeInfo := strategy.UpgradeInfo{
		BusinessId: 123,
		PatchAlgo:  1,
		DevTypeId:  1,
		DevId:      "0013a0702601",
		DevModel:   "Y32",
		DevPlat:    "918",
		EuiVersion: "100",
		//AppointVersion:  103,
		PackageName:    "com.my.test",
		ApkVersionCode: 0,
		ApkMd5:         "",
	}
	apps, err := strategy.GetAppApkPlatsInfo(ctx.Model.Mysql.O, ctx.Model.Redis, &upgradeInfo, 1)
	fmt.Println("-------:", apps, err)
	for k, v := range apps[0].Apks {
		fmt.Println("-----:", k, v)
	}
	for k, v := range apps[0].DevPlats {
		fmt.Println("-----:", k, v)
	}
}
