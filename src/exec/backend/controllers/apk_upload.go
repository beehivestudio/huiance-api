package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/gsp412/androidbinary/apk"

	"upgrade-api/src/share/cdn"
	"upgrade-api/src/share/comm"
	"upgrade-api/src/share/database"
	"upgrade-api/src/share/database/upgrade"
	"upgrade-api/src/share/lib/crypt"
)

// ApkUploadController operations for ApkUpload
type ApkUploadController struct {
	CommonController
}

// URLMapping ...
func (c *ApkUploadController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
}

// Post ...
// @Title Post
// @Description create ApkUpload
// @Param	body		body 	upgrade.ApkUpload	true		"body for ApkUpload content"
// @Success 201 {int} upgrade.ApkUpload
// @Failure 403 body is empty
// @router   / [post]
// @author # jiapeng # 2020-06-08 16:36:35 #
func (c *ApkUploadController) Post() {
	ctx := GetBackendCntx()

	file, information, err := c.GetFile("file")
	if nil != err {
		logs.Error("Get file failed! errmsg:%s", err.Error())
		c.ErrorMessage(comm.ERR_PARAM_INVALID, "Parameter is invalid!")
		return
	}
	defer file.Close()

	/* 保存到static文件目录 */
	downloadPath := GetDownloadPath(information.Filename)

	if err = c.SaveToFile("file", downloadPath); nil != err {
		logs.Error("Upload file failed! downloadPath:%s errmsg:%s",
			downloadPath, err.Error())
		c.ErrorMessage(comm.ERR_INTERNAL_SERVER_ERROR, err.Error())
		return
	}

	/* 计算MD5值 */
	md5Str, err := crypt.Md5ByFilePath(downloadPath)
	if nil != err {
		logs.Error("Get file md5 failed! errmsg:%s", err.Error())
		c.ErrorMessage(comm.ERR_INTERNAL_SERVER_ERROR, err.Error())
		return
	}

	/* 读取APK信息 */
	apk, err := apk.OpenFile(downloadPath)
	if nil != err {
		logs.Error("Read apk failed! downloadPath:%s errmsg:%s", downloadPath, err.Error())
		c.ErrorMessage(comm.ERR_INTERNAL_SERVER_ERROR, err.Error())
		return
	} else if apk.PublicKey() == "" { /* 校验APK文件是否有签名*/
		logs.Error("Apk signature is invalid! downloadPath:%s", downloadPath)
		c.ErrorMessage(comm.ERR_APK_SIGN_INVALID, "Apk signature is invalid!")
		return
	}

	label, err := apk.Label(comm.APK_CN_CONF)
	if nil != err {
		logs.Error("Read apk label failed! errmsg:%s", err.Error())
		c.ErrorMessage(comm.ERR_APK_DATA_INVALID, err.Error())
		return
	}

	versionCode, err := strconv.ParseInt(fmt.Sprintf("%v", apk.VersionCode()), 10, 32)
	if nil != err {
		logs.Error("Parse version code failed! VersionCode:%v errmsg:%s",
			apk.VersionCode(), err.Error())
		c.ErrorMessage(comm.ERR_APK_DATA_INVALID, err.Error())
		return
	}

	downloadUrl := GetDownloadUrl(information.Filename)
	//cdnFormat := fmt.Sprintf(comm.CDN_UPLOAD_ARGS, app.CdnPlatId, app.CdnSplatId)

	var v = &upgrade.ApkUpload{
		FileName:    information.Filename,
		Size:        information.Size,
		PackageName: apk.PackageName(),
		Label:       label,
		VersionCode: versionCode,
		VersionName: apk.VersionName(),
		PublicKey:   apk.PublicKey(),
		LocalPath:   downloadPath,
		CdnUrl:      "",
		Status:      upgrade.APKFILEUPLOAD_DOING,
		Md5:         md5Str,
		CreateTime:  time.Now(),
		CreateUser:  c.mail,
	}

	//判断数据库中是否已有上传的文件
	apkUpload, err := upgrade.GetApkUpload(ctx.Model.Mysql.O, md5Str, apk.PackageName(), versionCode)
	if nil != err {
		logs.Error("Get apk upload failed! errmsg :%s", err.Error())
		code, e := database.MysqlFormatError(err)
		c.ErrorMessage(code, e)
		return
	} else if apkUpload.Id != comm.ID_EMPTY {
		// APK文件状态为正在上传或上传成功，返回apk信息
		if apkUpload.Status == upgrade.APKFILEUPLOAD_SUCCESS ||
			apkUpload.Status == upgrade.APKFILEUPLOAD_DOING {
			v.Id = apkUpload.Id
			c.FormatInterfaceResp(http.StatusOK, comm.OK, "Ok", ApkUploadRspData(v))
			return
		}

		// APK文件上传失败,重新上传cdn
		err = ApkUploadCdn(int64(apkUpload.Id), downloadUrl, md5Str)
		if nil != err {
			logs.Error("Upload to cdn failed! downloadPath:%s errmsg:%v",
				downloadPath, err.Error())

			// 更改APK失败的状态
			err := upgrade.UpdateApkStatus(
				ctx.Model.Mysql.O, apkUpload.Id, upgrade.APKFILEUPLOAD_FAIL)
			if nil != err {
				logs.Error("Update apk status failed! downloadPath:%s errmsg:%v",
					downloadPath, err.Error())
				code, e := database.MysqlFormatError(err)
				c.ErrorMessage(code, e)
				return
			}
			c.ErrorMessage(comm.ERR_INTERNAL_SERVER_ERROR, err.Error())
			return
		}

		//更改文件状态
		err := upgrade.UpdateApkStatus(ctx.Model.Mysql.O, apkUpload.Id, upgrade.APKFILEUPLOAD_DOING)
		if nil != err {
			logs.Error("Update apk status failed! apkUpload.Id:%d errmsg:%s",
				apkUpload.Id, err.Error())
			code, e := database.MysqlFormatError(err)
			c.ErrorMessage(code, e)
			return
		}

		c.FormatInterfaceResp(http.StatusOK, comm.OK, "Ok", ApkUploadRspData(v))
		return
	}

	//数据库中无apk文件，增加apk文件
	id, err := upgrade.AddApkUpload(ctx.Model.Mysql.O, v)
	if nil != err {
		logs.Error("Add apk upload item failed! downloadPath:%s errmsg:%s",
			downloadPath, err.Error())
		code, e := database.MysqlFormatError(err)
		c.ErrorMessage(code, e)
		return
	}

	//cdn上传
	err = ApkUploadCdn(id, downloadUrl, md5Str)
	if nil != err {
		logs.Error("Upload to cdn failed  errmsg:%v", err.Error())

		//更改apk上传失败的状态
		err := upgrade.UpdateApkStatus(ctx.Model.Mysql.O, apkUpload.Id, upgrade.APKFILEUPLOAD_FAIL)
		if nil != err {
			code, e := database.MysqlFormatError(err)
			c.ErrorMessage(code, e)
			return
		}

		c.ErrorMessage(comm.ERR_INTERNAL_SERVER_ERROR, err.Error())
		return
	}

	c.FormatInterfaceResp(http.StatusOK, comm.OK, "Ok", ApkUploadRspData(v))
	return
}

//获取下载url
func GetDownloadUrl(fileName string) string {
	return fmt.Sprintf("%s%s%s%s",
		comm.DOMAIN_NAME, comm.STATIC_REQUEST_URL, "download/apk/", fileName)
}

//获取下载路径
func GetDownloadPath(fileName string) string {
	return fmt.Sprintf("%s%s%s",
		comm.STORAGE_PATH, comm.DOWNLOAD_APK_PATH, fileName)
}

//cdn上传
func ApkUploadCdn(apkUploadId int64, downloadUrl, md5Str string) error {

	outKey := cdn.GenApkOutKey(apkUploadId)

	err := cdn.Upload(downloadUrl, md5Str, outKey, "")
	if nil != err {

		//非重复上传，返回上传失败结果
		logs.Error("Upload to cdn failed  errmsg:%v", err.Error())
		if err.Error() != comm.CDN_UPLOAD_DUPLICATE {
			return err
		}

		//重复上传，重置任务状态，重新发回调
		logs.Info("Upload cdn duplicate task! md5 msg :%s", md5Str)
		err = cdn.GetCdnFileUrl(md5Str)
		if err != nil {
			return err
		}
	}

	return nil
}

/* apk上传返回参数 */
type ApkUploadRsp struct {
	Id          int    `json:"id"`
	FileName    string `json:"file_name"`
	Size        int64  `json:"size"`
	PackageName string `json:"package_name"`
	Label       string `json:"label"`
	VersionCode int64  `json:"version_code"`
	VersionName string `json:"version_name"`
	PublicKey   string `json:"public_key"`
}

//返回apk上传数据
func ApkUploadRspData(a *upgrade.ApkUpload) (rsp ApkUploadRsp) {
	rsp = ApkUploadRsp{
		Id:          a.Id,
		FileName:    a.FileName,
		Size:        a.Size,
		PackageName: a.PackageName,
		Label:       a.Label,
		VersionCode: a.VersionCode,
		VersionName: a.VersionName,
		PublicKey:   a.PublicKey,
	}
	return rsp
}

// GetOne ...
// @Title Get One
// @Description get ApkUpload by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success http.StatusOK {object} upgrade.ApkUpload
// @Failure 403 :id is empty
// @router /:id([0-9]+) [get]
// @author # jiapeng # 2020-06-08 16:36:35 #
func (c *ApkUploadController) GetOne() {
	ctx := GetBackendCntx()
	idStr := c.Ctx.Input.Param(":id")

	id, err := strconv.Atoi(idStr)
	if nil != err {
		logs.Error("Parse id failed! idStr:%s errmsg:%s",
			idStr, err.Error())
		c.ErrorMessage(comm.ERR_PARAM_INVALID, err.Error())
		return
	}

	v, err := upgrade.GetApkUploadById(ctx.Model.Mysql.O, id)
	if nil != err && err != orm.ErrNoRows {
		logs.Error("Get apk upload item failed! id:%d errmsg:%s",
			id, err.Error())
		code, e := database.MysqlFormatError(err)
		c.ErrorMessage(code, e)
		return
	} else if err == orm.ErrNoRows {
		logs.Error("Not found apk upload item! id:%d errmsg:%s",
			id, err.Error())
		c.ErrorMessage(comm.ERR_NOT_FOUND, err.Error())
		return
	}

	c.FormatInterfaceResp(http.StatusOK, comm.OK, "Ok", map[string]interface{}{
		"id":           v.Id,                                       // ID
		"file_name":    v.FileName,                                 // 文件名
		"status":       v.Status,                                   // 文件上传状态，1：处理中 2：已就绪 3：处理失败
		"size":         v.Size,                                     // 文件大小（单位字节）
		"url":          v.CdnUrl,                                   // CDN下载地址
		"package_name": v.PublicKey,                                // 应用报名
		"label":        v.Label,                                    // APK应用名称
		"version_code": v.VersionCode,                              // 版本code
		"version_name": v.VersionName,                              // 版本名
		"public_key":   v.PublicKey,                                // 公钥信息
		"create_time":  v.CreateTime.Format("2006-01-02 15:04:05"), // 创建时间
		"create_user":  v.CreateUser,                               // 创建者
		"update_time":  v.UpdateTime.Format("2006-01-02 15:04:05"), // 最后修改时间
		"update_user":  v.UpdateUser,                               // 修改者
	})
}

// GetAll ...
// @Title Get All
// @Description get ApkUpload
// @Success http.StatusOK {object} upgrade.ApkUpload
// @Failure 403
// @router /list [get]
// @author # jiapeng # 2020-06-08 16:36:35 #
func (c *ApkUploadController) GetAll() {
	ctx := GetBackendCntx()
	fields := []string{}
	values := []interface{}{}

	if id := c.GetString("id"); id != "" {
		fields = append(fields, "id = ?")
		values = append(values, id)
	}

	if fileName := c.GetString("file_name"); fileName != "" {
		fields = append(fields, "file_name like ?")
		values = append(values, fmt.Sprintf("%%%s%%", fileName))
	}

	if label := c.GetString("label"); label != "" {
		fields = append(fields, "label like ?")
		values = append(values, fmt.Sprintf("%%%s%%", label))
	}

	if packageName := c.GetString("package_name"); packageName != "" {
		fields = append(fields, "package_name like ?")
		values = append(values, fmt.Sprintf("%%%s%%", packageName))
	}

	if versionName := c.GetString("version_name"); versionName != "" {
		fields = append(fields, "version_name like ?")
		values = append(values, fmt.Sprintf("%%%s%%", versionName))
	}

	if status, err := c.GetInt("status"); err == nil {
		fields = append(fields, "status = ?")
		values = append(values, status)
	}

	beginTime := c.GetString("begin_time")
	endTime := c.GetString("end_time")

	if beginTime != "" && endTime != "" {
		fields = append(fields, fmt.Sprintf("create_time BETWEEN '%s' AND '%s' ", beginTime, endTime))
		//values = append(values, time)

	} else if beginTime == "" && endTime != "" {
		fields = append(fields, " create_time < ？")
		values = append(values, endTime)

	} else if beginTime != "" && endTime == "" {
		fields = append(fields, " create_time > ? ")
		values = append(values, beginTime)
	}

	pageSize, err := c.GetInt("page_size", comm.PAGE_SIZE)
	if nil != err {
		logs.Error("Get page size failed! errmsg:%s", err.Error())
		c.ErrorMessage(comm.ERR_PARAM_INVALID, err.Error())
		return
	}
	page, err := c.GetInt("page", comm.PAGE_START)
	if nil != err {
		logs.Error("Get page failed! errmsg:%s", err.Error())
		c.ErrorMessage(comm.ERR_PARAM_INVALID, err.Error())
		return
	}

	result, err := upgrade.GetAllApkUploadByFuzzy(ctx.Model.Mysql.O, values, fields, page, pageSize)
	if nil != err {
		logs.Error("Get page list failed! params:%v errmsg:%s", values, err.Error())
		code, e := database.MysqlFormatError(err)
		c.ErrorMessage(code, e)
		return
	}
	var data = []map[string]interface{}{}
	if ups, ok := result["list"]; ok {
		for _, res := range ups.([]upgrade.ApkUpload) {
			data = append(data, map[string]interface{}{
				"id":           res.Id,                                       // ID
				"file_name":    res.FileName,                                 // 文件名
				"status":       res.Status,                                   // 文件上传状态，1：处理中 2：已就绪 3：处理失败
				"size":         res.Size,                                     // 文件大小（单位字节）
				"url":          res.CdnUrl,                                   // CDN下载地址
				"package_name": res.PackageName,                              // 应用报名
				"label":        res.Label,                                    // APK应用名称
				"version_code": res.VersionCode,                              // 版本code
				"version_name": res.VersionName,                              // 版本名
				"public_key":   res.PublicKey,                                // 公钥信息
				"create_user":  res.CreateUser,                               // 创建者
				"create_time":  res.CreateTime.Format("2006-01-02 15:04:05"), // 创建时间
				"update_user":  res.UpdateUser,                               // 最后修改者
				"update_time":  res.UpdateTime.Format("2006-01-02 15:04:05"), // 最后修改时间
			})
		}
		result["list"] = data
	}

	c.FormatInterfaceResp(http.StatusOK, comm.OK, "Ok", result)
}
