package controllers

import (
	"github.com/astaxie/beego/logs"

	"upgrade-api/src/share/comm"
)

// CommonController operations for Apk
type CommonController struct {
	comm.BaseController
	mail string
}

func (this *CommonController) Prepare() {
	ctx := GetBackendCntx()
	m_tk := this.Ctx.Request.Header.Get("m-tk")
	if 0 == len(m_tk) {
		logs.Warn("m-tk not allowed empty")
		this.ErrorMessage(comm.ERR_PARAM_MISS, "m-tk not allowed empty.")
		return
	}

	mail, code, err := ctx.Conf.InSso.InternalMTkIsValid(m_tk)
	if nil != err {
		logs.Warn("m-tk check failed, msg: %s", err.Error())
		this.ErrorMessage(code, err.Error())
		return
	}

	this.mail = mail

}
