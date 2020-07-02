package utils

import "runtime"

/******************************************************************************
 **函数名称: GetRunFuncName
 **功    能: 获取正在运行的函数名
 **输入参数: NONE
 **输出参数: NONE
 **返    回: 函数名称
 **实现描述:
 **作    者: # Qifeng.zou # 2020-06-29 18:31:48 #
 ******************************************************************************/
func GetRunFuncName() string {
	pc := make([]uintptr, 1)

	runtime.Callers(2, pc)

	f := runtime.FuncForPC(pc[0])

	return f.Name()
}
