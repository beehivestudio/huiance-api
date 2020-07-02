package comm

import (
	"time"

	"github.com/gsp412/androidbinary"
)

// 系统加载数据
const (
	DEFAULT_LANGUAGE = "zh-CN"
)

// 列表分页相关
const (
	PAGE_START = 1  // 默认起始页
	PAGE_SIZE  = 20 // 页面大小
)

// 时间相关
const (
	TM_MIN  = 60 * time.Second // 分
	TM_HOUR = 60 * TM_MIN      // 时
	TM_DAY  = 24 * TM_HOUR     // 天
	TM_WEEK = 7 * TM_DAY       // 周

	DEFAULT_APK_EXPIRE       = "10m" // APK验证超时时间
	DEFAULT_REQUEST_TIME_OUT = 3600  // 请求超时时间, 单位秒
)

const (
	SYSYTEM_USER = "SYSTEM" // 默认系统用户
)

// 接口相关
const (
	SSO_TK_CHECK       = "/api/checkTicket"     // 用户中心 sso_tk验证接口
	ACCESS_TK_CHECK    = "/oauthopen/userbasic" // 用户中心 access_tk验证接口
	SSO_GET_USER_BY_ID = "/api/getuserbyid"     // 用户中心 根据uid获取用户信息
)

// 定义APK解析语言
var APK_CN_CONF = &androidbinary.ResTableConfig{
	Language: [2]uint8{'z', 'h'},
	Country:  [2]uint8{'C', 'N'},
}

// worker 相关
const (
	WORKER_APK_CHANNEL_LEN_DEFAULT = 500
)

// 验证码过期时间
const (
	VERIF_CODE_EXPIRE_SEC = 120 // 验证码过期时间
)

const (
	DOMAIN_NAME         = "http://test.backend.upgrade.letv.com" //项目域名
	CDN_UPlOAD_URL      = "http://mcache.oss.letv.com/ext/addfile"
	CDN_UPLOAD_USER     = "scloud_apk_upgrade_test"
	CDN_UPLOAD_ARGS     = "?platid=%d&splatid=%d"
	CDN_UPLOAD_CDN_HOST = "http://g3.letv.cn/"
	CDN_FILE_STATUS     = "http://fid.oss.letv.com/ext/refresh?key="
	STORAGE_PATH        = "/letv/upgrade/storage/static/"        //文件存储的服务器目录
	STATIC_REQUEST_URL  = "/upgrade/backend/v1/static/file/url/" //静态文件请求路径
	DOWNLOAD_APK_PATH   = "download/apk/"                        //下载APK路径
	DOWNLOAD_PATCH_PATH = "download/patch/"                      //下载PATCH路径
	PATCH_PATH          = "patch/"                               //生成PATCH路径
	PATCH_FILE_SUFFIX   = ".patch"                               //差分包文件后缀
)

//cdn返回状态码
const (
	CDN_CODE_SUCC      = 200 //cdn返回成功状态
	CDN_CODE_DUPLICATE = 400 //cdn返回重复任务状态
)

//apk上传Cdn重复任务
const (
	CDN_UPLOAD_DUPLICATE = "400" //cdn返回重复任务状态
)

//设备管理相关
const (
	DEVICE_GROUPLS_RDS_EXPIRE_SEC = 10 * 60 // redis 缓存时间
)

// 升级策略相关
const (
	/* 通用开启相关 */
	NOT_ENABLE = 0
	ENABLE     = 1

	/* 升级方式 */
	UPGRADE_SILENCE   = 4 // 静默升级
	UPGRADE_UNINSTALL = 6 // 静默卸载

	/* 升级的设备范围 */
	DEV_TYPE_ALL       = 0 // 0：全部设备
	DEV_TYPE_ID        = 1 // 1：指定设备ID；设备ID列表
	DEV_TYPE_GROUP     = 2 // 2：指定设备分组，1个设备只归属1个分组，1个分组只归属1个机型，1个机型只归属1个平台
	DEV_TYPE_STARGAZER = 3 // 3：观星用户画像组列表

	/* 机型分组数据 */
	MODEL_TYPE_ALL   = 1 // 该机型下所有分组
	MODEL_TYPE_GROUP = 2 // 该机型下指定分组
	MODEL_All        = 0 // 0 表示该机型下所有分组，对应常量

	/* 操作行为 */
	ACTION_UPGRADE         = 1 // 1:只升级
	ACTION_INSTALL         = 2 // 2:只安装
	ACTION_UPGRADE_INSTALL = 3 // 3:升级+安装

	/* 是否人为操作 */
	HUMAN_NOT = 0 //0：非人为操作
	HUMAN_YES = 1 //1：人为操作

	/* redis 缓存时间，单位秒 */
	RDS_EXPIRE_TIME_SECOND = 20 //10 * 60
	RDS_EXPIER_TIME_RANDOM = 10 //60 // redis 缓存随机时间

	/* 流控时间范围 */
	LIMIT_TIME_MIN = 0
	LIMIT_TIME_MAX = 23

	/* 时间戳范围校验 */
	TIMESTAMP_VERIFY = 5 * 60 * 10 // 单位秒

	/* Redis 存储 nil 时对比字符串 */
	RDS_NULL = "null"

	/* 是否关联设备平台 */
	HAS_PLAT_NOT = 0 //0：不关联。即：适用于所有设备平台；
	HAS_PLAT_YES = 1 //1：关联。即：只适用于指定设备平台。如果关联表中无设备平台，表示所有设备平台均不支持。

	/* Redis 缓存数据为空时错误校验 */
	RDS_DATA_IS_NULL = "Data is null!"
)

// 接口状态码
const (
	HTTP_RESP_CODE = 10000
)

//业务区分关键字
const (
	CDN_CALLBACK_PATCH = "PATCH" //差分包
	CDN_CALLBACK_APK   = "APK"   //pak
)

//设备类型id
const (
	DEFAULT_DEV_TYPE_ID = 1 //设备类型
)

//id不存在
const (
	ID_EMPTY = 0 //id为0
)
