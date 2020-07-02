package controllers

import (
	"strings"
	"upgrade-api/src/share/comm"
)

type IndexPageController struct {
	comm.BaseController
}

/******************************************************************************
 **函数名称: IndexPageController Get
 **功    能: 首页渲染、加载昵称到首页
 **输入参数:
 **输出参数: NONE
 **返    回:
 **实现描述:
 **注意事项:
 **作    者: # Linlin.guo # 2020-6-29 15:15:42 #
 ******************************************************************************/
func (c *IndexPageController) Get() {
	ctx := GetBackendCntx()

	var mail string
	data := strings.Split(c.Ctx.Request.RequestURI, "?")
	var m_tk string
	for _, v := range data {
		items := strings.Split(v, "=")
		if "m_tk" == items[0] {
			m_tk = items[1]
			break
		}
	}
	if 0 != len(m_tk) {
		mail, _, _ = ctx.Conf.InSso.InternalMTkIsValid(m_tk)
	}

	c.Data["email"] = mail

	data = strings.Split(mail, "@")
	if data != nil {
		c.Data["nickName"] = data[0]
	}

	c.TplName = "index.html"
}
