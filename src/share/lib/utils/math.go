package utils

import (
	"crypto/rand"
	"fmt"
	"math"
	"math/big"
	"strconv"
)

/* 浮点类型数据精度损失问题 */
func Round(f float64, n int) float64 {
	n10 := math.Pow10(n)

	return math.Trunc((f+0.5/n10)*n10) / n10
}

/* 浮点类型数据精度损失问题 */
func Round2(f float64, n int) float64 {
	floatStr := fmt.Sprintf("%."+strconv.Itoa(n)+"f", f)
	inst, _ := strconv.ParseFloat(floatStr, 64)
	return inst
}

/******************************************************************************
 **函数名称: Float64Sum
 **功    能: Float64 求和
 **输入参数:
 **      f1: 参数1
 **      f2: 参数2
 **      d: 位数
 **输出参数:
 **返    回:
 **实现描述:
 **注意事项:
 **作    者: # Shuangpeng.guo # 2018-12-11 13:15:42 #
 ******************************************************************************/
func Float64Sum(f1 float64, f2 float64, d int) float64 {
	return Round2(f1+f2, d)
}

/******************************************************************************
 **函数名称: RangeRand
 **功    能: 生成区间[-m, n]的安全随机数
 **输入参数:
 **      min: 最小值
 **      max: 最大值
 **输出参数:
 **返    回: 随机数
 **实现描述:
 **注意事项:
 **作    者: # taoshengbo # 2020-06-05 11:52:08 #
 ******************************************************************************/
func RangeRand(min, max int64) int64 {
	if max < min {
		min, max = max, min
	}
	if min < 0 {
		f64Min := math.Abs(float64(min))
		i64Min := int64(f64Min)
		result, _ := rand.Int(rand.Reader, big.NewInt(max+1+i64Min))
		return result.Int64() - i64Min
	} else {
		result, _ := rand.Int(rand.Reader, big.NewInt(max-min+1))
		return min + result.Int64()
	}
}
