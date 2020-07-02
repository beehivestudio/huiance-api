package comm

const PREFIX = "UPGRADE:"

/* Redis键定义 */
const (
	RDS_KEY_CDN_SPLAT_HASH  = PREFIX + "CDN:SPLAT:HASH"
	QUOTA_KEY_TYPE_BUSINESS = PREFIX + "QUOTA:BUSINESS:KEY:%d" // 业务频控
	QUOTA_KEY_TYPE_STRATEGY = PREFIX + "QUOTA:STRATEGY:KEY:%d" // 策略频控
	RDS_KEY_QUOTA_SYNC      = PREFIX + "QUOTA:SYNC"            // 同步key
	RDS_KEY_QUOTA_KEY_SEC   = PREFIX + "QUOTA:KEY:%s:SEC"      // 频控KEY(秒)
	RDS_KEY_QUOTA_KEY_MIN   = PREFIX + "QUOTA:KEY:%s:MIN"      // 频控KEY(分)
	RDS_KEY_QUOTA_KEY_HOUR  = PREFIX + "QUOTA:KEY:%s:HOUR"     // 频控KEY(时)
	RDS_KEY_QUOTA_KEY_DAY   = PREFIX + "QUOTA:KEY:%s:DAY"      // 频控KEY(天)
)

// 升级策略相关
const (
	RDS_KEY_BUSINESS            = PREFIX + "BUSINESS:ID:%v"                                                                           // BUSINESS
	RDS_KEY_APK_PATCH_TASK_LIST = PREFIX + "APK:PATCH:TASK:LIST"                                                                      // 差分包通知队列
	RDS_KEY_APP_APK             = PREFIX + "APP_APK:DEV_TYPE_ID:%v:BUS_ID:%v:PKG:%v:EUI_VER:%v:APPOINT_VER:%v:UP_TYPE:%v:VER_CODE:%v" // APP、APK 表缓存
	RDS_KEY_UPGRADE_STRATEGY    = PREFIX + "STRATEGY:APK_IDS:%v:UP_TYPE:%v"                                                           // 升级策略表缓存
	RDS_KEY_APP_PATCH           = PREFIX + "APK_PATCH:DEV_TYPE_ID:%v:HIGH_VER:%v:LOW_VER:%v:PATCH_ALGO:%v"                            // 差分包表缓存
	RDS_KEY_PKG_LIST            = PREFIX + "PKG_LIST:DEV_TYPE_ID:%v:BUS_ID:%v:PKG_NAME:%v"                                            // 安装包列表缓存
	RDS_KEY_MODEL_LIST          = PREFIX + "MODEL_LIST"                                                                               // 设备机型表缓存
	RDS_KEY_GROUP_LIST          = PREFIX + "GROUP_LIST"                                                                               // 设备分组表缓存
	RDS_KEY_PLAT_LIST           = PREFIX + "PLAT_LIST"                                                                                // 设备平台表缓存
	RDS_KEY_APP                 = PREFIX + "APP:APK_ID:%v"                                                                            // app 表缓存
	RDS_KEY_APP_PLAT_REL        = PREFIX + "APP_PLAT_REL:APP_ID:%v"                                                                   // app_plat_rel 表缓存
	RDS_KEY_LOW_APK             = PREFIX + "LOW_APK:APP_ID:%v:VER_CODE:%v:MD5:%v"                                                     // app_plat_rel 表缓存
)

//设备
const (
	RDS_KEY_DEVICE_GROUP_LIST = PREFIX + "DEVICE_GROUP_LIST" //设备分组列表
)
