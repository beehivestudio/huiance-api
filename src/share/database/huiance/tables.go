package upgrade

/* mysql数据库表名定义 */
const (
	UPGRADE_TAB_APK_PATCH_ALGO                  = "ApkPatchAlgo"                // 差分算法表
	UPGRADE_TAB_CDN_PLAT                        = "CdnPlat"                     // CDN平台表
	UPGRADE_TAB_DEV_TYPE                        = "dev_type"                    // 设备种类表
	UPGRADE_TAB_DEV_PLAT                        = "dev_plat"                    // 设备平台表
	UPGRADE_TAB_DEV_MODEL                       = "DevModel"                    // 设备机型表
	UPGRADE_TAB_DEV_GROUP                       = "DevGroup"                    // 设备分组表
	UPGRADE_TAB_BUSINESS                        = "business"                    // 业务信息表
	UPGRADE_TAB_BUSINESS_FLOW_LIMIT             = "BusinessFlowLimit"           // 业务流控表
	UPGRADE_TAB_APP                             = "app"                         // 应用信息表
	UPGRADE_TAB_APP_PLAT_REL                    = "app_plat_rel"                // 应用与设备平台关联表
	UPGRADE_TAB_APK                             = "apk"                         // APK信息表
	UPGRADE_TAB_APK_UPLOAD                      = "apk_upload"                  // APK文件上传表
	UPGRADE_TAB_APK_PATCH                       = "apk_patch"                   // APK差分信息表
	UPGRADE_TAB_APK_UPGRADE_STRATEGY            = "ApkUpgradeStrategy"          // 升级策略表
	UPGRADE_TAB_APK_UPGRADE_FLOW_LIMIT_STRATEGY = "ApkUpgradeFlowLimitStrategy" // 升级策略-流控策略表
	UPGRADE_TAB_APK_UPGRADE_DEV_LIST            = "apk_upgrade_dev_list"        // 升级策略-设备清单表
	UPGRADE_TAB_APK_UPGRADE_DEV_GROUP           = "apk_upgrade_dev_group"       // 升级策略-设备分组表
	UPGRADE_TAB_APK_UPGRADE_STARGAZER_GROUP     = "apk_upgrade_stargazer_group" // 升级策略-观星分组表
	UPGRADE_TAB_BUSINESS_STATISTICS             = "BusinessStatistics"          // 业务统计表
	UPGRADE_TAB_APP_STATISTICS                  = "AppStatistics"               // 应用统计表
	UPGRADE_TAB_APK_STATISTICS                  = "ApkStatistics"               // APK统计表
)

type TableTotal struct {
	Total int64 `orm:"column(total)"`
}

//const (
//	APP_ENABLE_OFF        = 0 // 禁用（下线）
//	APP_ENABLE_ON         = 1 // 启用（上线）
//	APP_HASDEVPLAT_NO  = 0 // 不关联
//	APP_HASDEVPLAT_YES = 1 // 关联
//)
//
//////////////////////////////////////////////////////////////////////////////////
//// 应用信息表
//// @表名: App
//// @描述: 用于记录应用的基本信息.
//type App struct {
//	Id            uint64                            // 应用ID
//	Name          string                            // 应用名称
//	PackageName   string                            // 应用包名
//	BusinessId    uint64                            // 	业务ID
//	CdnPlatId     uint64                            // CDN 平台ID
//	CdnSplatId    uint64                            // CDN 子平台ID
//	DevTypeId  uint64                            // 设备种类ID (外键)
//	Enable        int                               // 是否启用 0：禁用（下线） 1：启用（上线）
//	AppPublicKey  string                            // 应用公钥
//	HasDevPlat string                            // 是否关联设备平台 0：不关联 1：关联
//	Description   string                            // 描述信息
//	CreateTime    time.Time                         // 创建时间
//	CreateUser    uint64                            //	创建者
//	UpdateTime    time.Time                         // 更新时间
//	UpdateUser    uint64                            // 修改者
//	AppPlatRel    *AppPlatRel `orm:"reverse(one)"`  // 反向关联应用平台
//	Apks          []*Apk      `orm:"reverse(many)"` // 反向关联apk
//	ApkPatchs      []*ApkPatch  `orn:"reverse(many)"` //反向关联拆分信息
//	//DevType    *DevType `orm:"rel(fk);cascade"`     // 设备种类ID (外键)
//}
//
//////////////////////////////////////////////////////////////////////////////////
//// 应用平台关联表
//// @表名: AppPlatRel
//// @描述: 记录应用和平台的关联信息。该表数据主要在配置APK升级策略时，用于筛选机型和设备分组。
//type AppPlatRel struct {
//	Id         uint64                       // ID
//	App        *App `orm:"rel(fk);cascade"` // 应用ID(外键)
//	Plat       string                       // 电视平台
//	CreateTime time.Time                    // 创建时间
//	UpdateTime time.Time                    // 更新时间
//}
//
const (
	APK_STATUS_PROCESSING   = 0 // 处理中
	APK_STATUS_NORMAL       = 1 // 正常
	APK_STATUS_VERSION_FAIL = 2 // 版本验证失败
	APK_STATUS_SK_FAIL      = 3 // 秘钥验证失败
	APK_STATUS_PKG_FAIL     = 4 // 包名验证失败
	APK_STATUS_OTHER_ERR    = 5 // 其他异常
	APK_STATUS_TIMEOUT      = 6 // 超时
)

//
//////////////////////////////////////////////////////////////////////////////////
//// Apk信息表
//// @表名: Apk
//// @描述: 记录APK的基本信息。
//type Apk struct {
//	Id                  uint64                                      // ID
//	App                 *App `orm:"rel(fk);cascade"`                // 应用ID(外键)
//	Version             uint64                                      // 版本号
//	VersionName         string                                      // 版本名
//	Url                 string                                      // 资源地址
//	Md5                 string                                      // MD5值
//	Size                uint64                                      // 包大小
//	EuiLowVersion       string                                      // 依赖的EUI最低版本
//	EuiHighVersion      string                                      // 依赖的EUI最高版本
//	Status              int                                         // 处理状态 0：处理中 1：正常 2：版本验证失败 3：秘钥验证失败 4：其他异常 5：超时（超过1分钟算超时）
//	Description         string                                      // 描述信息
//	CreateTime          time.Time                                   // 创建时间
//	CreateUser          uint64                                      //	创建者
//	UpdateTime          time.Time                                   // 更新时间
//	UpdateUser          uint64                                      // 修改者
//	ApkUpgradeStrategys []*ApkUpgradeStrategy `orm:"reverse(many)"` // 反向关联升级策略
//}
//
//const (
//	ApkPatch_STATUS_NOT_DONE        = 0 // 尚未进行差分处理
//	ApkPatch_STATUS_PROCESSING      = 1 // 正在处理差分处理
//	ApkPatch_STATUS_SUCCESS_        = 2 // 差分处理成功
//	ApkPatch_STATUS_BIGGER          = 3 // 差分包大于新版本全量包
//	ApkPatch_STATUS_TIMEOUT         = 4 // 差分过程超时
//	ApkPatch_STATUS_ERROR           = 5 // 差分过程错误
//	ApkPatch_STATUS_PACKAGE_ILLEGAL = 6 // 差分包有问题
//)
//
//////////////////////////////////////////////////////////////////////////////////
//// 差分信息表
//// @表名: ApkPatch
//// @描述: 记录差分信息表。
//type ApkPatch struct {
//	Id          uint64                       // ID
//	App         *App `orm:"rel(fk);cascade"` // 应用ID(外键)
//	HighVersion uint64                       // 高版本号
//	LowVersion  uint64                       // 低版本号
//	PatchAlgo    uint64                       // 拆分算法
//	Status      int                          // 处理状态 0：尚未进行差分处理 1：正在处理差分处理 2：差分处理成功 3：差分包大于新版本全量包 4：差分过程超时 5：差分过程错误 6：差分包有问题
//	Url         string                       // 差分包地址
//	Md5         string                       // 差分包MD5值
//	Size        uint64                       // 拆分包大小
//	Description string                       // 描述信息
//	CreateTime  time.Time                    // 创建时间
//	UpdateTime  time.Time                    // 更新时间
//}
//
//const (
//	APKUPGRADESTRATEGY_ENABLE_OFF                        = 0 // 失效
//	APKUPGRADESTRATEGY_ENABLE_ON                         = 1 // 生效
//	APKUPGRADESTRATEGY_HASFLOWLIMIT_NO                   = 0 // 无流控
//	APKUPGRADESTRATEGY_HASFLOWLIMIT_YES                  = 1 // 有流控
//	APKUPGRADESTRATEGY_UPGRADE_TYPE_UNKNOWN              = 0 // 未知方式
//	APKUPGRADESTRATEGY_UPGRADE_TYPE_CHOOSEABLE_UPGRADE   = 1 // 可选升级
//	APKUPGRADESTRATEGY_UPGRADE_TYPE_RECOMMEND_UPGRADE    = 2 // 推荐升级
//	APKUPGRADESTRATEGY_UPGRADE_TYPE_FORCE_UPGRADE        = 3 // 强制升级
//	APKUPGRADESTRATEGY_UPGRADE_TYPE_QUIET_UPGRADE        = 4 // 静默升级
//	APKUPGRADESTRATEGY_UPGRADE_TYPE_FORCE_INSTALL        = 5 // 强制安装
//	APKUPGRADESTRATEGY_UPGRADE_TYPE_FORCE_UNINSTALL      = 6 // 强制卸载
//	APKUPGRADESTRATEGY_UPGRADE_DEV_TYPE_ALL              = 0 // 全部设备
//	APKUPGRADESTRATEGY_UPGRADE_DEV_TYPE_ASSIGN_DEV_ID    = 1 // 指定设备ID
//	APKUPGRADESTRATEGY_UPGRADE_DEV_TYPE_ASSIGN_DEV_GROUP = 2 // 指定设备分组
//	APKUPGRADESTRATEGY_UPGRADE_DEV_TYPE_STARGAZER        = 3 // 观星用户画像组列表
//)
//
//////////////////////////////////////////////////////////////////////////////////
//// 升级策略表
//// @表名: ApkUpgradeStrategy
//// @描述: 配置APK包的升级策略。
//type ApkUpgradeStrategy struct {
//	Id                          uint64                                            // ID
//	Apk                         *Apk `orm:"rel(fk)";cascade`                      // ApkID(外键)
//	Enable                      int                                               // 可通过此控制策略是否生效 0：失效 1：生效
//	Begin_Datetime              time.Time                                         // 开始时间
//	End_DataTime                time.Time                                         // 结束时间
//	HasFlowLimit                int                                               // 有无流控 0：无流控 1：有流控
//	Upgrade_Type                int                                               // 升级方式 0：未知方式 1：可选升级 2：推荐升级 3：强制升级 4：静默升级 5：强制安装 6：强制卸载
//	Upgrade_Dev_Type            int                                               // 指定升级的设备范围 0：全部设备 1：指定设备ID 2：指定设备分组 3：观星用户画像组列表
//	Upgrade_Dev_Data            string                                            // 升级设备数据：0.当为全部设备时，此字段为空；1.当指定设备ID时，此字段存储设备ID列表；2.当指定设备分组时，此字段存储选中的设备分组信息（JSON格式）3.当指定观星用户画像组列表时，此字段存储用户画像ID列表；
//	Description                 string                                            // 描述信息
//	CreateTime                  time.Time                                         // 创建时间
//	CreateUser                  uint64                                            //	创建者
//	UpdateTime                  time.Time                                         // 更新时间
//	UpdateUser                  uint64                                            // 修改者
//	ApkUpgradeFlowLimitStrategy []*ApkUpgradeFlowLimitStrategy `orm:"reverse(one)"` // 反向关联流控策略
//	ApkUpgradeDevList           *ApkUpgradeDevList           `orm:"reverse(one)"` // 反向关联设备清单
//	ApkUpgradeDevGroup          *ApkUpgradeDevGroup          `orm:"reverse(one)"` //反向关联设备分组
//	ApkUpgradeStargazerGroup    *ApkUpgradeStargazerGroup    `orm:"reverse(one)"` //反向代理观星用户画像分组
//}
//
//const (
//	APKUPGRADEFLOWLIMITSTRATEGY_DIMENSION_S = 1 // 秒
//	APKUPGRADEFLOWLIMITSTRATEGY_DIMENSION_M = 2 // 分
//	APKUPGRADEFLOWLIMITSTRATEGY_DIMENSION_H = 3 // 时
//	APKUPGRADEFLOWLIMITSTRATEGY_DIMENSION_D = 4 // 天
//)
//
//////////////////////////////////////////////////////////////////////////////////
//// 升级策略-流控策略表
//// @表名: ApkUpgradeFlowLimitStrategy
//// @描述: 配置APK包的升级流控策略。
//type ApkUpgradeFlowLimitStrategy struct {
//	Id                 uint64                                                           // ID
//	ApkUpgradeStrategy *ApkUpgradeStrategy `orm:"column(strategy_id);rel(fk);cascade" ` // APK升级策略ID （外键）
//	BeginTime          int                                                              // 开始时间段[0, 24）（精确到时）
//	EndTime            int                                                              // 结束时间段[0, 24）（精确到时）
//	Dimension          int                                                              // 流控维度 1：秒 2：分 3：时 4：天
//	Limit              int                                                              // 频控限制
//	CreateTime         time.Time                                                        // 创建时间
//	UpdateTime         time.Time                                                        // 更新时间
//}
//
//////////////////////////////////////////////////////////////////////////////////
//// 升级策略-设备清单表
//// @表名: ApkUpgradeDevList
//// @描述: 记录APK包的升级设备ID列表
//type ApkUpgradeDevList struct {
//	Id                 uint64                                                           // ID
//	ApkUpgradeStrategy *ApkUpgradeStrategy `orm:"column(strategy_id);rel(fk);cascade" ` // APK升级策略ID （外键）
//	DevId              string                                                           //设备ID 1.电视：此值为电视MAC地址 2.手机：此值为手机IMEI号
//	CreateTime         time.Time                                                        // 创建时间
//	UpdateTime         time.Time                                                        // 更新时间
//}
//
//////////////////////////////////////////////////////////////////////////////////
//// 升级策略-设备分组表
//// @表名: ApkUpgradeDevGroup
//// @描述: 记录APK包的升级设备分组信息
//type ApkUpgradeDevGroup struct {
//	Id                 uint64                                                   // ID
//	ApkUpgradeStrategy *ApkUpgradeStrategy `orm:"column(strategy_id);rel(fk);cascade" ` // APK升级策略ID （外键）
//	TvPlatId           uint64                                                   // TV设备平台ID
//	TvModel            uint64                                                   // TV机型
//	TvGroupId          uint64                                                   // TV设备分组ID
//	CreateTime         time.Time                                                // 创建时间
//	UpdateTime         time.Time                                                // 更新时间
//}
//
//////////////////////////////////////////////////////////////////////////////////
//// 升级策略-画像分组表
//// @表名: ApkUpgradeStargazerGroup
//// @描述: 记录APK包的升级设备观星用户画像分组
//type ApkUpgradeStargazerGroup struct {
//	Id                 uint64                                                   // ID
//	ApkUpgradeStrategy *ApkUpgradeStrategy `orm:"column(strategy_id);rel(fk);cascade"` // APK升级策略ID （外键）
//	StargazerTid       string                                                   // 观星用户画像分组ID
//	CreateTime         time.Time                                                // 创建时间
//	UpdateTime         time.Time                                                // 更新时间
//}
