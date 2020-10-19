// filter project filter.go
package filter

import (
	"math"
)

/***********************************************
功能:实数除以一个复数
输入:
	r:float64:实数
	cmp:complex128:复数
输出:
	complex128;复数
说明:
编辑:wang_jp
时间:2020年10月16日
***********************************************/
func ComplexDivByReal(r float64, cmp complex128) complex128 {
	c_r := real(cmp)
	c_i := imag(cmp)
	return complex(c_r*r, -1*c_i*r) / (cmp * complex(c_r, -1*c_i))
}

/***********************************************
功能:实数除以一个复数数组中的每一个元素
输入:
	r:float64:实数
	cmps:[]complex128:复数
输出:
	[]complex128;复数
说明:
编辑:wang_jp
时间:2020年10月16日
***********************************************/
func ComplexArrDivByReal(r float64, cmps []complex128) []complex128 {
	var res []complex128
	for _, c := range cmps {
		res = append(res, ComplexDivByReal(r, c))
	}
	return res
}

/***********************************************
功能:实数与一个复数数组中的每一个元素相减
输入:
	r:float64:实数
	cmps:[]complex128:复数
输出:
	[]complex128;复数
说明:
编辑:wang_jp
时间:2020年10月16日
***********************************************/
func ComplexArrSubByReal(r float64, cmps []complex128) []complex128 {
	var res []complex128
	for _, c := range cmps {
		res = append(res, ComplexSubByReal(r, c))
	}
	return res
}

/***********************************************
功能:实数与一个复数数组中的每一个元素相加
输入:
	r:float64:实数
	cmps:[]complex128:复数
输出:
	[]complex128;复数
说明:
编辑:wang_jp
时间:2020年10月16日
***********************************************/
func ComplexArrAddReal(r float64, cmps []complex128) []complex128 {
	var res []complex128
	for _, c := range cmps {
		res = append(res, ComplexAddReal(r, c))
	}
	return res
}

/***********************************************
功能:两个复数数组中的每一个元素相除
输入:
	a,b:[]complex128:复数
输出:
	[]complex128;复数
说明:如果被除数组的实部和虚部都是0，则使结果元素等于被除数
编辑:wang_jp
时间:2020年10月16日
***********************************************/
func ComplexArrDiv(a, b []complex128) []complex128 {
	var res []complex128
	for i, c := range b {
		if real(c) == 0 && imag(c) == 0 {
			res = append(res, c)
		} else {
			res = append(res, a[i]/c)
		}
	}
	return res
}

/***********************************************
功能:复数加上一个实数
输入:
	r:float64:实数
	cmp:complex128:复数
输出:
	complex128;复数
说明:
编辑:wang_jp
时间:2020年10月16日
***********************************************/
func ComplexAddReal(r float64, cmp complex128) complex128 {
	return complex(real(cmp)+r, imag(cmp))
}

/***********************************************
功能:复数减去一个实数
输入:
	r:float64:实数
	cmp:complex128:复数
输出:
	complex128;复数
说明:
编辑:wang_jp
时间:2020年10月16日
***********************************************/
func ComplexSubReal(r float64, cmp complex128) complex128 {
	return complex(real(cmp)-r, imag(cmp))
}

/***********************************************
功能:实数减去一个复数
输入:
	r:float64:实数
	cmp:complex128:复数
输出:
	complex128;复数
说明:
编辑:wang_jp
时间:2020年10月16日
***********************************************/
func ComplexSubByReal(r float64, cmp complex128) complex128 {
	return complex(r-real(cmp), -imag(cmp))
}

/***********************************************
功能:实数乘以一个复数
输入:
	r:float64:实数
	cmp:complex128:复数
输出:
	complex128;复数
说明:
编辑:wang_jp
时间:2020年10月16日
***********************************************/
func ComplexMulReal(r float64, cmp complex128) complex128 {
	return complex(real(cmp)*r, imag(cmp)*r)
}

/***********************************************
功能:复数数组元素同时乘以一个实数
输入:
	cmps:[]complex128:复数数组切片
	r:float64;实数
输出:
	[]complex128;复数数组
说明:
编辑:wang_jp
时间:2020年10月16日
***********************************************/
func ComplexArrMulReal(cmps []complex128, r float64) []complex128 {
	res := make([]complex128, len(cmps))
	for i, cv := range cmps {
		res[i] = ComplexMulReal(r, cv)
	}
	return res
}

/***********************************************
功能:复数数组的乘积
输入:
	cmps:[]complex128:复数数组切片
输出:
	complex128;复数
说明:
编辑:wang_jp
时间:2020年10月16日
***********************************************/
func ComplexArrProd(cmps []complex128) complex128 {
	var res complex128 = 1 + 0i
	if len(cmps) < 1 {
		return res
	}
	for _, cv := range cmps {
		res *= cv
	}
	return res
}

/***********************************************
功能:比较两个复数数组是否相等
输入:
	a,b:[]complex128:复数数组切片
输出:
	如果两个数组长度相等且每个元素的实部和虚部都相等,返回true,否则返回false
说明:
编辑:wang_jp
时间:2020年10月16日
***********************************************/
func ComplexArrsIsEqual(a, b []complex128) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if math.Abs(real(v)-real(b[i])) > 1e-5 {
			return false
		}
		if math.Abs(imag(v)-imag(b[i])) > 1e-5 {
			return false
		}
	}
	return true
}

/***********************************************
功能:比较两个复数是否相等
输入:
	a,b:complex128:复数
输出:
	如果两个复数的实部和虚部都相等,返回true,否则返回false
说明:
编辑:wang_jp
时间:2020年10月16日
***********************************************/
func ComplexsIsEqual(a, b complex128) bool {
	if math.Abs(real(a)-real(b)) > 1e-5 {
		return false
	}
	if math.Abs(imag(a)-imag(b)) > 1e-5 {
		return false
	}
	return true
}

/***********************************************
功能:实数转换为复数
输入:
	a:float64 实数
输出:
	complex128
说明:
编辑:wang_jp
时间:2020年10月16日
***********************************************/
func Real2Complex(a float64) complex128 {
	return ComplexMulReal(a, 1+0i)
}

/***********************************************
功能:实数数组转换为复数数组
输入:
	arr:[]float64 实数
输出:
	[]complex128
说明:
编辑:wang_jp
时间:2020年10月16日
***********************************************/
func Reals2Complexs(arr []float64) []complex128 {
	var cmps []complex128
	for _, v := range arr {
		cmps = append(cmps, Real2Complex(v))
	}
	return cmps
}
