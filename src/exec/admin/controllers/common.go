package controllers

import (
	"github.com/astaxie/beego/logs"

	"upgrade-api/src/share/comm"
	"upgrade-api/src/share/database/admin"
)

// CommonController operations for Apk
type CommonController struct {
	comm.BaseController
	mail   string
	userId int64
}

func (this *CommonController) Prepare() {
	ctx := GetAdminCntx()
	m_tk := this.Ctx.Request.Header.Get("m-tk")
	if 0 == len(m_tk) {
		logs.Warn("%s() m-tk not allowed empty", utils.GetRunFuncName())
		this.ErrorMessage(comm.ERR_PARAM_MISS, "m-tk not allowed empty.")
		return
	}

	mail, code, err := ctx.Conf.InSso.InternalMTkIsValid(m_tk)
	if nil != err {
		logs.Warn("%s() m-tk check failed! errmsg:%s",
			utils.GetRunFuncName(), err.Error())
		this.ErrorMessage(code, err.Error())
		return
	}

	this.mail = mail
	user, err := db_admin.GetUserByMail(ctx.Model.Mysql.O, mail)
	if nil != err {
		logs.Warn("%s() Get user by mail failed! errmsg:%s",
			utils.GetRunFuncName(), err.Error())
		this.ErrorMessage(code, err.Error())
		return
	}
	this.userId = user.Id
}
