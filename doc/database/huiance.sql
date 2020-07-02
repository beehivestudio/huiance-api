/*
 Navicat Premium Data Transfer

 Source Server         : tv-order-10.112.40.80-test
 Source Server Type    : MySQL
 Source Server Version : 100114
 Source Host           : 10.112.40.80:33107
 Source Schema         : upgrade

 Target Server Type    : MySQL
 Target Server Version : 100114
 File Encoding         : 65001

 Date: 17/06/2020 09:54:29
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for action_log
-- ----------------------------
DROP TABLE IF EXISTS `action_log`;
CREATE TABLE `action_log` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `user_id` bigint(20) NOT NULL COMMENT '用户ID；外键：user.id',
  `action` bigint(20) NOT NULL COMMENT '行为ID；用户行为ID取值：；0:未知行为；1：登录；2：退出；3：添加；4：修改；5：删除；6：查询（不用）',
  `detail` varchar(1024) DEFAULT NULL COMMENT '行为详情',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `action` (`action`) COMMENT '用户行为ID'
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='操作日志表；记录用户操作日志信息';

-- ----------------------------
-- Table structure for apk
-- ----------------------------
DROP TABLE IF EXISTS `apk`;
CREATE TABLE `apk` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `bcode` varchar(128) NOT NULL COMMENT '业务方唯一码',
  `app_id` bigint(20) NOT NULL COMMENT '应用ID；外键: app.id',
  `version_code` bigint(20) NOT NULL COMMENT '版本号；版本号数值越大，版本号越高；注：由程序从APK包中解析获取',
  `version_name` varchar(1024) NOT NULL COMMENT '版本名；注：由程序从APK包中解析获取',
  `url` varchar(1024) NOT NULL COMMENT '资源地址；Apk包的下载地址',
  `md5` varchar(64) NOT NULL COMMENT 'MD5值；APK包的MD5值（文件MD5值）',
  `size` bigint(20) NOT NULL COMMENT '包大小；APK包大小（单位：字节）',
  `eui_low_version` varchar(256) DEFAULT NULL COMMENT '依赖的EUI最低版本',
  `eui_high_version` varchar(256) DEFAULT NULL COMMENT '依赖的EUI最高版本',
  `status` int(11) NOT NULL COMMENT '处理状态：0：处理中1：正常2：版本验证失败3：秘钥验证失败4：其他异常5：超时（超过1分钟算超时）',
  `callback` varchar(2048) DEFAULT NULL COMMENT '业务回调',
  `description` varchar(2048) DEFAULT NULL COMMENT '描述信息',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `create_user` varchar(128) NOT NULL DEFAULT '0' COMMENT '页面创建者',
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  `update_user` varchar(128) NOT NULL DEFAULT '0' COMMENT '修改者账号',
  `callback_url` varchar(1024) NOT NULL DEFAULT '' COMMENT '资源地址；Apk回调地址',
  PRIMARY KEY (`id`),
  KEY `version` (`version_code`) COMMENT '版本号',
  KEY `version_name` (`version_name`(255)) COMMENT '版本名'
) ENGINE=InnoDB AUTO_INCREMENT=21 DEFAULT CHARSET=utf8 COMMENT='APK信息表；记录APK的基本信息';

-- ----------------------------
-- Table structure for apk_patch
-- ----------------------------
DROP TABLE IF EXISTS `apk_patch`;
CREATE TABLE `apk_patch` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `app_id` bigint(20) NOT NULL COMMENT '应用ID；外键: app.id',
  `high_version_code` bigint(20) NOT NULL COMMENT '高版本号',
  `low_version_code` bigint(20) NOT NULL COMMENT '低版本号',
  `patch_algo` bigint(20) NOT NULL COMMENT '差分算法',
  `status` int(20) NOT NULL COMMENT '处理状态：0：尚未进行差分处理；1：正在处理差分处理；2：上传cdn中；3：差分处理成功；4：差分包大于新版本全量包；5：差分过程超时；6：差分过程错误；7：差分包有问题',
  `url` varchar(1024) NOT NULL COMMENT '差分包的下载地址；注意：路径中应体现差分算法ID',
  `md5` varchar(64) NOT NULL COMMENT '差分包MD5值',
  `size` bigint(20) NOT NULL COMMENT '差分包大小；单位：字节',
  `description` varchar(2048) DEFAULT NULL COMMENT '描述信息',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  `local_path` varchar(255) NOT NULL COMMENT '差分包本地存储路径',
  PRIMARY KEY (`id`),
  UNIQUE KEY `app_id_patch_algo_high_version_low_version` (`app_id`,`patch_algo`,`high_version_code`,`low_version_code`) COMMENT '应用ID + 差分算法ID + 高版本ID + 低版本ID',
  KEY `high_version` (`high_version_code`) COMMENT '高版本号',
  KEY `low_version` (`low_version_code`) COMMENT '低版本号'
) ENGINE=InnoDB AUTO_INCREMENT=29 DEFAULT CHARSET=utf8 COMMENT='差分信息表；记录APK的差分基本信息';

-- ----------------------------
-- Table structure for apk_patch_algo
-- ----------------------------
DROP TABLE IF EXISTS `apk_patch_algo`;
CREATE TABLE `apk_patch_algo` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `name` varchar(128) NOT NULL COMMENT '算法名称',
  `enable` int(11) NOT NULL DEFAULT '0' COMMENT '是否启用',
  `description` varchar(1024) NOT NULL COMMENT '描述信息',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`) COMMENT '算法名称'
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8 COMMENT='差分算法表';

-- ----------------------------
-- Table structure for apk_upgrade_dev_group
-- ----------------------------
DROP TABLE IF EXISTS `apk_upgrade_dev_group`;
CREATE TABLE `apk_upgrade_dev_group` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `strategy_id` bigint(20) NOT NULL COMMENT 'APK升级策略ID；外键: apk_upgrade_strategy.id',
  `dev_plat_id` bigint(20) NOT NULL COMMENT '备平台ID. 外键：DevPlat.id',
  `dev_model_id` bigint(20) NOT NULL COMMENT '设备机型ID. 外键：DevModel.id',
  `dev_group_id` bigint(20) NOT NULL COMMENT '设备分组ID. 外键：DevGroup.id',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `strategy_id_dev_group_id` (`strategy_id`,`dev_group_id`) COMMENT '策略ID + 设备分组',
  KEY `dev_plat_id` (`dev_plat_id`) USING BTREE COMMENT '设备平台ID. 外键：DevPlat.id',
  KEY `dev_model_id` (`dev_model_id`) USING BTREE COMMENT '设备机型ID. 外键：DevModel.id',
  KEY `dev_group_id` (`dev_group_id`) USING BTREE COMMENT '设备分组ID. 外键：DevGroup.id'
) ENGINE=InnoDB AUTO_INCREMENT=90 DEFAULT CHARSET=utf8 COMMENT='升级策略-设备分组表；记录APK包的升级设备分组信息';

-- ----------------------------
-- Table structure for apk_upgrade_dev_list
-- ----------------------------
DROP TABLE IF EXISTS `apk_upgrade_dev_list`;
CREATE TABLE `apk_upgrade_dev_list` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `strategy_id` bigint(20) NOT NULL COMMENT 'APK升级策略ID；外键: apk_upgrade_strategy.id',
  `dev_id` varchar(255) NOT NULL COMMENT '设备ID（设备唯一标志）；1.电视：此值为电视MAC地址；2.手机：此值为手机IMEI号',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `strategy_id_dev_id` (`strategy_id`,`dev_id`) COMMENT '策略ID + 设备ID',
  KEY `dev_id` (`dev_id`) COMMENT '设备ID'
) ENGINE=InnoDB AUTO_INCREMENT=27 DEFAULT CHARSET=utf8 COMMENT='升级策略-设备清单表；记录APK包的升级设备ID列表';

-- ----------------------------
-- Table structure for apk_upgrade_flow_limit_strategy
-- ----------------------------
DROP TABLE IF EXISTS `apk_upgrade_flow_limit_strategy`;
CREATE TABLE `apk_upgrade_flow_limit_strategy` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `strategy_id` bigint(20) NOT NULL COMMENT 'APK升级策略ID；外键: apk_upgrade_strategy.id',
  `begin_time` int(11) NOT NULL DEFAULT '0' COMMENT '开始时间段；开始时间段[0, 23]（精确到时）',
  `end_time` int(11) NOT NULL DEFAULT '23' COMMENT '结束时间段；结束时间段[0, 23]（精确到时）',
  `dimension` int(11) NOT NULL DEFAULT '1' COMMENT '流控维度；流控维度：1：秒；2：分；3：时；4：天',
  `limit` bigint(20) NOT NULL COMMENT '频控限制；在流控维度上的次数',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=181 DEFAULT CHARSET=utf8 COMMENT='升级策略-流控策略表；配置APK包的升级流控策略';

-- ----------------------------
-- Table structure for apk_upgrade_stargazer_group
-- ----------------------------
DROP TABLE IF EXISTS `apk_upgrade_stargazer_group`;
CREATE TABLE `apk_upgrade_stargazer_group` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `strategy_id` bigint(20) NOT NULL COMMENT 'APK升级策略ID；外键: apk_upgrade_strategy.id',
  `stargazer_tid` bigint(255) NOT NULL COMMENT '观星用户画像分组ID',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `strategy_id_stargazer_tid` (`strategy_id`,`stargazer_tid`) COMMENT '策略ID + 观星用户画像分组ID'
) ENGINE=InnoDB AUTO_INCREMENT=12 DEFAULT CHARSET=utf8 COMMENT='升级策略-画像分组表；记录APK包的升级设备观星用户画像分组';

-- ----------------------------
-- Table structure for apk_upgrade_strategy
-- ----------------------------
DROP TABLE IF EXISTS `apk_upgrade_strategy`;
CREATE TABLE `apk_upgrade_strategy` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `apk_id` bigint(20) NOT NULL COMMENT 'APK_ID；外键: apk.id',
  `enable` int(11) NOT NULL COMMENT '是否生效；可通过此控制策略是否生效0：失效；1：生效',
  `begin_datetime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '升级开始时间',
  `end_datetime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '升级结束时间',
  `has_flow_limit` int(11) NOT NULL DEFAULT '0' COMMENT '有无流控；0：无流控；1：有流控',
  `upgrade_type` int(11) NOT NULL DEFAULT '0' COMMENT '升级方式：0：未知方式；1：可选升级（用户选择取消应用升级后，不再提示其是否升级）；2：推荐升级（用户选择取消应用升级后，下次依然提示是否升级）；3：强制升级（用户不可取消应用升级）；4：静默升级；5：强制安装（用户不知情的情况下安装应用）；6：强制卸载（用户不知情的情况下卸载应用）',
  `upgrade_dev_type` int(11) NOT NULL COMMENT '指定升级的设备范围：0：全部设备；1：指定设备ID（设备ID列表）；2：指定设备分组，1个设备只归属1个分组，1个分组只归属1个机型，1个机型只归属1个平台；3：观星用户画像组列表',
  `upgrade_dev_data` text NOT NULL COMMENT '升级设备数据：0.当为全部设备时，此字段为空；；1.当指定设备ID时，此字段存储设备ID列表；；2.当指定设备分组时，此字段存储选中的设备分组信息（JSON格式）；3.当指定观星用户画像组列表时，此字段存储用户画像ID列表；',
  `description` varchar(2048) DEFAULT NULL COMMENT '描述信息',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `create_user` varchar(128) NOT NULL DEFAULT '0' COMMENT '创建者账号',
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  `update_user` varchar(128) NOT NULL DEFAULT '0' COMMENT '修改者账号',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=92 DEFAULT CHARSET=utf8 COMMENT='升级策略表；配置APK包的升级策略';

-- ----------------------------
-- Table structure for apk_upload
-- ----------------------------
DROP TABLE IF EXISTS `apk_upload`;
CREATE TABLE `apk_upload` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `file_name` varchar(255) NOT NULL COMMENT '文件名',
  `md5` varchar(64) NOT NULL COMMENT '文件MD5值',
  `size` bigint(20) NOT NULL DEFAULT '0' COMMENT '文件大小	',
  `package_name` varchar(255) NOT NULL COMMENT '包名',
  `label` varchar(255) NOT NULL COMMENT '文件名',
  `version_code` bigint(20) NOT NULL DEFAULT '0' COMMENT '版本编号',
  `version_name` varchar(255) NOT NULL DEFAULT '0' COMMENT '版本名称',
  `public_key` varchar(2048) NOT NULL COMMENT '公钥信息',
  `local_path` varchar(255) NOT NULL COMMENT '本地相对存储路径',
  `cdn_url` varchar(128) NOT NULL COMMENT 'CDN存储路径',
  `status` int(11) NOT NULL DEFAULT '0' COMMENT '上传状态：0：正在上传；1：上传成功；2：上传失败；3：上传超时',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  `create_user` varchar(128) NOT NULL DEFAULT '0' COMMENT '创建者账号',
  `update_user` varchar(128) NOT NULL DEFAULT '0' COMMENT '修改者账号',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8 COMMENT='APK上传CDN的状态追踪表';

-- ----------------------------
-- Table structure for app
-- ----------------------------
DROP TABLE IF EXISTS `app`;
CREATE TABLE `app` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `name` varchar(128) NOT NULL COMMENT '应用名称',
  `package_name` varchar(128) NOT NULL COMMENT '应用包名（不唯一）；注：存在包名相同，但分属不同的平台',
  `business_id` bigint(20) NOT NULL COMMENT '业务ID；所属业务线',
  `cdn_plat_id` bigint(20) NOT NULL COMMENT 'CDN 平台ID',
  `cdn_splat_id` bigint(20) NOT NULL COMMENT 'CDN 子平台ID；SPLATID：子平台ID。用于业务计费，由CDN统一分配',
  `dev_type_id` bigint(20) NOT NULL COMMENT '设备种类ID；外键：devtype.id',
  `enable` int(11) NOT NULL DEFAULT '1' COMMENT '是否启用；可通过此开关控制应用上下线；0：禁用（下线）；1：启用（上线）',
  `app_public_key` varchar(512) DEFAULT NULL COMMENT '应用公钥：用于验证文件是否被篡改；注：1.创建应用时，该值为空；2.第一次创建APK时，将应用公钥信息录入此字段。',
  `has_dev_plat` int(11) NOT NULL COMMENT '是否关联设备平台；0：不关联。即：适用于所有设备平台；1：关联。即：只适用于指定设备平台。如果关联表中无设备平台，表示所有设备平台均不支持。',
  `dev_plats` varchar(512) NOT NULL DEFAULT '' COMMENT '关联设备平台，多个「,」分开，如 928,938',
  `description` varchar(2048) DEFAULT NULL COMMENT '描述信息',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `create_user` varchar(128) NOT NULL DEFAULT '0' COMMENT '创建者账号',
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  `update_user` varchar(128) NOT NULL DEFAULT '0' COMMENT '修改者账号',
  PRIMARY KEY (`id`),
  UNIQUE KEY `business_id_name` (`business_id`,`name`) COMMENT '业务ID + 应用名称',
  KEY `package_name` (`package_name`) COMMENT '应用包名',
  KEY `name` (`name`) COMMENT '应用名称',
  KEY `cdn_plat_id` (`cdn_plat_id`) COMMENT 'CDN PLATID',
  KEY `cdn_splat_id` (`cdn_splat_id`) COMMENT 'CDN SPLATID'
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8 COMMENT='应用信息表；记录应用的基本信息';

-- ----------------------------
-- Table structure for app_plat_rel
-- ----------------------------
DROP TABLE IF EXISTS `app_plat_rel`;
CREATE TABLE `app_plat_rel` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `app_id` bigint(20) NOT NULL COMMENT '应用ID；外键：app.id',
  `dev_plat_id` bigint(20) NOT NULL COMMENT '设备平台ID',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `app_id_dev_plat_id` (`app_id`,`dev_plat_id`) COMMENT '应用ID + 设备平台ID'
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8 COMMENT='应用平台关联表；记录应用和平台的关联信息。该表数据主要在配置APK升级策略时，用于筛选机型和设备分组。';

-- ----------------------------
-- Table structure for business
-- ----------------------------
DROP TABLE IF EXISTS `business`;
CREATE TABLE `business` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `name` varchar(128) NOT NULL COMMENT '业务名称',
  `key` varchar(64) NOT NULL COMMENT '业务秘钥；数据加解密秘钥',
  `enable` int(11) NOT NULL DEFAULT '0' COMMENT '是否启用',
  `has_flow_limit` int(11) NOT NULL DEFAULT '0' COMMENT '是否开启流控；有无流控；0：无流控；1：有流控',
  `manager` varchar(258) NOT NULL COMMENT '项目技术负责人',
  `description` varchar(255) DEFAULT NULL COMMENT '描述信息',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `create_user` varchar(128) NOT NULL DEFAULT '0' COMMENT '创建者账号',
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  `update_user` varchar(128) NOT NULL DEFAULT '0' COMMENT '更新者账号',
  PRIMARY KEY (`id`),
  KEY `name` (`name`) COMMENT '业务名称'
) ENGINE=InnoDB AUTO_INCREMENT=138 DEFAULT CHARSET=utf8 COMMENT='业务信息表';

-- ----------------------------
-- Table structure for business_flow_limit
-- ----------------------------
DROP TABLE IF EXISTS `business_flow_limit`;
CREATE TABLE `business_flow_limit` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `business_id` bigint(20) NOT NULL COMMENT '业务ID；外键：business.id',
  `begin_time` int(11) NOT NULL DEFAULT '0' COMMENT '开始时间；开始时间，取值范围[0, 23]',
  `end_time` int(11) NOT NULL DEFAULT '23' COMMENT '结束时间；结束时间，取值范围[0, 23]',
  `dimension` int(11) NOT NULL DEFAULT '1' COMMENT '流控维度；流控维度：1：秒；2：分；3：时；4：天',
  `limit` bigint(20) NOT NULL COMMENT '频控限制',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=26 DEFAULT CHARSET=utf8 COMMENT='业务流控表；记录业务的流控配置';

-- ----------------------------
-- Table structure for dev_group
-- ----------------------------
DROP TABLE IF EXISTS `dev_group`;
CREATE TABLE `dev_group` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `dev_model_id` bigint(20) NOT NULL COMMENT '设备机型ID；外键：dev_model.id',
  `dms_group_id` varchar(128) NOT NULL COMMENT 'DMS设备分组ID；从设备管理获取',
  `group_name` varchar(128) NOT NULL COMMENT '设备分组名称',
  `description` varchar(1024) DEFAULT NULL COMMENT '描述信息',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `dev_model_id_group_name` (`dev_model_id`,`group_name`) USING BTREE COMMENT '设备机型ID + 设备分组名称'
) ENGINE=InnoDB AUTO_INCREMENT=6545 DEFAULT CHARSET=utf8 COMMENT='设备分组表';

-- ----------------------------
-- Table structure for dev_model
-- ----------------------------
DROP TABLE IF EXISTS `dev_model`;
CREATE TABLE `dev_model` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `dev_plat_id` bigint(20) NOT NULL COMMENT '设备平台ID；外键：dev_plat.id',
  `model_name` varchar(128) NOT NULL COMMENT '设备机型名称',
  `description` varchar(1024) DEFAULT NULL COMMENT '描述信息',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `dev_plat_id_model_name` (`dev_plat_id`,`model_name`) USING BTREE COMMENT '设备平台ID + 设备机型名称'
) ENGINE=InnoDB AUTO_INCREMENT=11790 DEFAULT CHARSET=utf8 COMMENT='设备机型表';

-- ----------------------------
-- Table structure for dev_plat
-- ----------------------------
DROP TABLE IF EXISTS `dev_plat`;
CREATE TABLE `dev_plat` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `dev_type_id` bigint(20) NOT NULL COMMENT '设备种类；外键：dev_type.id',
  `name` varchar(128) NOT NULL COMMENT '设备平台名称',
  `description` varchar(1024) DEFAULT NULL COMMENT '描述信息',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `dev_type_id_name` (`dev_type_id`,`name`) USING BTREE COMMENT '设备种类ID + 设备平台名称'
) ENGINE=InnoDB AUTO_INCREMENT=52902 DEFAULT CHARSET=utf8 COMMENT='设备平台表';

-- ----------------------------
-- Table structure for dev_type
-- ----------------------------
DROP TABLE IF EXISTS `dev_type`;
CREATE TABLE `dev_type` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `name` varchar(128) NOT NULL COMMENT '设备名称；设备名称列表：1.乐视电视；2.乐视手机',
  `description` varchar(1024) DEFAULT NULL COMMENT '描述信息',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`) COMMENT '设备名称'
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8 COMMENT='设备种类表';

-- ----------------------------
-- Table structure for page
-- ----------------------------
DROP TABLE IF EXISTS `page`;
CREATE TABLE `page` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '页面ID',
  `name` varchar(128) NOT NULL COMMENT '页面名称',
  `path` varchar(255) NOT NULL COMMENT '资源路径',
  `priority` int(10) unsigned NOT NULL DEFAULT '1' COMMENT '优先级',
  `parent_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '父级页面ID',
  `description` varchar(1024) DEFAULT NULL COMMENT '描述信息',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `create_user_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '页面创建者',
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  `update_user_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '页面修改者',
  PRIMARY KEY (`id`),
  KEY `parent_id` (`parent_id`) COMMENT '父级页面ID'
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='页面信息表';

-- ----------------------------
-- Table structure for role
-- ----------------------------
DROP TABLE IF EXISTS `role`;
CREATE TABLE `role` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '角色ID',
  `name` varchar(128) NOT NULL COMMENT '角色名称',
  `description` varchar(1024) DEFAULT NULL COMMENT '描述信息',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `create_user_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '页面创建者',
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  `update_user_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '页面修改者',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='角色信息表';

-- ----------------------------
-- Table structure for role_page_rel
-- ----------------------------
DROP TABLE IF EXISTS `role_page_rel`;
CREATE TABLE `role_page_rel` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `role_id` bigint(20) NOT NULL COMMENT '角色ID；外键：role.id',
  `page_id` bigint(20) NOT NULL COMMENT '页面ID；外键：page.id',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `create_user_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '页面创建者',
  PRIMARY KEY (`id`),
  KEY `role_id` (`role_id`) COMMENT '角色ID',
  KEY `page_id` (`page_id`) COMMENT '页面ID'
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='角色页面关联表；记录角色对应的页面访问权限';

-- ----------------------------
-- Table structure for user
-- ----------------------------
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '用户ID',
  `name` varchar(128) NOT NULL COMMENT '显示姓名',
  `email` varchar(32) NOT NULL COMMENT '邮箱地址',
  `enable` int(11) NOT NULL DEFAULT '0' COMMENT '是否启用',
  `description` varchar(1024) DEFAULT NULL COMMENT '描述信息',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `create_user_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '页面创建者',
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  `update_user_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '页面修改者',
  PRIMARY KEY (`id`),
  KEY `email` (`email`) COMMENT '用户公司邮箱'
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='用户信息表';

-- ----------------------------
-- Table structure for user_role_rel
-- ----------------------------
DROP TABLE IF EXISTS `user_role_rel`;
CREATE TABLE `user_role_rel` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `user_id` bigint(20) NOT NULL COMMENT '用户ID；外键：user.id',
  `role_id` bigint(20) NOT NULL COMMENT '角色ID；外键：role.id',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `create_user_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '页面创建者',
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  `update_user_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '页面修改者',
  PRIMARY KEY (`id`),
  KEY `user_id` (`user_id`) COMMENT '用户ID；外键：user.id',
  KEY `role_id` (`role_id`) COMMENT '角色ID；外键：role.id',
  KEY `user_id_role_id` (`user_id`,`role_id`) COMMENT '用户ID + 角色ID'
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='用户角色关联表；记录用户所属角色';

-- ----------------------------
-- Table structure for action
-- ----------------------------
DROP TABLE IF EXISTS `action`;
CREATE TABLE `action` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '行为id',
  `action` varchar(255) NOT NULL COMMENT '行为名称 如：登录、退出、添加、修改、删除、查询等 ',
  `description` varchar(1024) DEFAULT NULL COMMENT '描述',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `action` (`action`) USING BTREE COMMENT '用户行为名称'
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8;

SET FOREIGN_KEY_CHECKS = 1;

