// filter project filter.go
package filter

import (
	"fmt"
	"math"
	"math/cmplx"
	"strings"

	"github.com/bkzy-wangjp/MicEngine/MicScript/numgo"
)

const (
	pi = 3.141592653589793
)

/***********************************************
功能:返回n阶Butterworth 滤波器初始参数
输入:
	n: 滤波阶数,n不可小于0,如果输入为负数自动转换为正数使用
输出:
	z:[]float64;零点参数
	p:[]complex128;极点参数
	k:float64;增益值
说明:
编辑:wang_jp
时间:2020年10月16日
***********************************************/
func buttap(n int) (numgo.Array, []complex128, float64) {
	if n < 0 {
		n *= -1
	}

	var z, m numgo.Array
	for i := (-1*n + 1); i < n; i = i + 2 {
		m = append(m, float64(i))
	}
	var p []complex128
	for i := 0; i < n; i++ {
		tmp := pi * m[i] / (2 * float64(n))
		p = append(p, cmplx.Exp(complex(0, tmp))*-1)
	}
	k := 1.0
	return z, p, k
}

/***********************************************
功能:从零点和极点返回传递函数的相对阶数
输入:
	z:[]float64;零点参数
	p:[]complex128;极点参数
输出:
	degree:int;相对阶数
说明:Return relative degree of transfer function from zeros and poles
编辑:wang_jp
时间:2020年10月16日
***********************************************/
func relativeDegree(z interface{}, p []complex128) (int, error) {
	zlen := 0
	zav, ok := z.(numgo.Array)
	if ok {
		zlen = len(zav)
	} else {
		zcv, ok := z.([]complex128)
		if ok {
			zlen = len(zcv)
		}
	}

	degree := len(p) - zlen
	if degree < 0 {
		return degree, fmt.Errorf("Improper transfer function. Must have at least as many poles as zeros[传递函数不正确,极点的阶数不能小于零点的阶数]")
	}
	return degree, nil
}

/***********************************************
功能:将低通滤波器原型转换为不同的频率
输入:
	z:[]float64;零点参数,运行完毕后直接改变
	p:[]complex128;极点参数,运行完毕后直接改变
	k:float64;增益值,运行完毕后直接改变
	wo:float64 截止频率
输出:错误信息
说明:从具有单位截止频率的模拟低通滤波器原型返回截止频率为“wo”的模拟低通滤波器
	 使用零、极点和增益（'zpk'）表示
编辑:wang_jp
时间:2020年10月16日
***********************************************/
func lp2lpZpk(z numgo.Array, p []complex128, k *float64, wo float64) error {
	degree, err := relativeDegree(z, p)
	if err != nil {
		return err
	}
	//Scale all points radially from origin to shift cutoff frequency
	z.MulScalar(wo)
	for i, pv := range p {
		p[i] = ComplexMulReal(wo, pv)
	}
	//Each shifted pole decreases gain by wo, each shifted zero increases it.
	//Cancel out the net change to keep overall gain the same.
	*k *= math.Pow(wo, float64(degree))
	return nil
}

/***********************************************
功能:将低通滤波器原型转换为不同的频率
输入:
	z:[]float64;零点参数,运行完毕后直接改变
	p:[]complex128;极点参数,运行完毕后直接改变
	k:float64;增益值,运行完毕后直接改变
	wo:float64 截止频率
输出:错误信息
说明:从具有单位截止频率的模拟高通滤波器原型返回截止频率为“wo”的模拟低通滤波器
	 使用零、极点和增益（'zpk'）表示
编辑:wang_jp
时间:2020年10月16日
***********************************************/
func lp2hpZpk(z numgo.Array, p []complex128, k *float64, wo float64) error {
	degree, err := relativeDegree(z, p)
	if err != nil {
		return err
	}
	fz := make(numgo.Array, len(z))
	copy(fz, z)
	fz.MulScalar(-1.0)
	fp := ComplexArrMulReal(p, -1.0)
	// Invert positions radially about unit circle to convert LPF to HPF
	// Scale all points radially from origin to shift cutoff frequency
	z.DivByScalar(wo)
	for i, pv := range p {
		p[i] = ComplexDivByReal(wo, pv)
	}
	// If lowpass had zeros at infinity, inverting moves them to origin.
	z = append(z, make(numgo.Array, degree)...)

	// Cancel out gain change caused by inversion
	*k *= real(ComplexDivByReal(fz.Product(), ComplexArrProd(fp)))
	return nil
}

/***********************************************
功能:基带移位
输入:
	cmps:[]complex128;
	wo:float64 中心频率,默认为1
输出:[]complex128;
说明:np.concatenate((z_lp + np.sqrt(z_lp**2 - wo**2),
                  z_lp - np.sqrt(z_lp**2 - wo**2)))
编辑:wang_jp
时间:2020年10月16日
***********************************************/
func shiftband(cmps []complex128, wo float64) []complex128 {
	wo = wo * wo
	var a, b []complex128
	for _, v := range cmps {
		tmp := v * v
		wc := Real2Complex(wo)
		tmp -= wc
		sqt := cmplx.Sqrt(tmp)
		a = append(a, v+sqt)
		b = append(b, v-sqt)
	}
	return append(a, b...)
}

/***********************************************
功能:将低通滤波器原型转换为不同的频率
输入:
	z:[]complex128;零点参数,运行完毕后直接改变
	p:[]complex128;极点参数,运行完毕后直接改变
	k:float64;增益值,运行完毕后直接改变
	wo:float64 中心频率,默认为1
	bw:float64 带宽,默认为1
输出:错误信息
说明:从具有单位截止频率的模拟低通滤波器原型返回中心频率为“wo”且带宽为“bw”的模拟带通滤波器器
	 使用零、极点和增益（'zpk'）表示
编辑:wang_jp
时间:2020年10月16日
***********************************************/
func lp2bpZpk(z, p []complex128, k, wo, bw float64) ([]complex128, []complex128, float64, error) {
	degree, err := relativeDegree(z, p)
	if err != nil {
		return nil, nil, 0, err
	}
	//缩放零点和极点参数到所需的带宽
	z_lp := ComplexArrMulReal(z, bw/2)
	p_lp := ComplexArrMulReal(p, bw/2)
	// 复制极点和零，从基带移到+wo和-wo
	z_bp := shiftband(z_lp, wo)
	p_bp := shiftband(p_lp, wo)
	// Move degree zeros to origin, leaving degree zeros at infinity for BPF.
	z_bp = append(z_bp, make([]complex128, degree)...)
	// Cancel out gain change caused by inversion
	k_bp := k * math.Pow(bw, float64(degree))
	return z_bp, p_bp, k_bp, nil
}

/***********************************************
功能:将低通滤波器原型转换为不同的频率
输入:
	z:[]complex128;零点参数,运行完毕后直接改变
	p:[]complex128;极点参数,运行完毕后直接改变
	k:float64;增益值,运行完毕后直接改变
	wo:float64 中心频率,默认为1
	bw:float64 带宽,默认为1
输出:错误信息
说明:从具有单位截止频率的模拟低通滤波器原型返回中心频率为“wo”且阻带宽度为“bw”的模拟带阻滤波器
	 使用零、极点和增益（'zpk'）表示
编辑:wang_jp
时间:2020年10月16日
***********************************************/
func lp2bsZpk(z, p []complex128, k, wo, bw float64) ([]complex128, []complex128, float64, error) {
	degree, err := relativeDegree(z, p)
	if err != nil {
		return nil, nil, 0, err
	}

	fz := ComplexArrMulReal(z, -1.0)
	fp := ComplexArrMulReal(p, -1.0)
	//缩放零点和极点参数到所需的带宽

	z_hp := ComplexArrDivByReal(bw/2, z)
	p_hp := ComplexArrDivByReal(bw/2, p)
	// 复制极点和零，从基带移到+wo和-wo
	z_bs := shiftband(z_hp, wo)
	p_bs := shiftband(p_hp, wo)
	// Move degree zeros to origin, leaving degree zeros at infinity for BPF.
	tmp := make([]complex128, degree)
	for i, _ := range tmp {
		tmp[i] = ComplexMulReal(wo, +1i)
	}
	z_bs = append(z_bs, tmp...)
	for i, _ := range tmp {
		tmp[i] = ComplexMulReal(wo, -1i)
	}
	z_bs = append(z_bs, tmp...)

	// Cancel out gain change caused by inversion
	k_bs := k * real(ComplexArrProd(fz)/ComplexArrProd(fp))

	return z_bs, p_bs, k_bs, nil
}

/***********************************************
功能:使用双线性变换从模拟滤波器返回数字IIR滤波器
输入:
	z:[]complex128/[]float64;零点参数,运行完毕后直接改变
	p:[]complex128;极点参数,运行完毕后直接改变
	k:float64;增益值,运行完毕后直接改变
	fs:float64
输出:z,p,k
说明:
编辑:wang_jp
时间:2020年10月18日
***********************************************/
func bilinearZpk(z interface{}, p []complex128, k, fs float64) (interface{}, []complex128, float64, error) {
	degree, err := relativeDegree(z, p)
	if err != nil {
		return nil, nil, 0, err
	}
	fs2 := 2.0 * fs

	p_z := ComplexArrAddReal(fs2, p)     //fs2+p
	p_tmp := ComplexArrSubByReal(fs2, p) //fs2-p
	p_z = ComplexArrDiv(p_z, p_tmp)      //(fs2+p)/(fs2-p)

	z_a, ok := z.(numgo.Array)
	if ok {
		z_tmp := make(numgo.Array, len(z_a)) //fs2-z
		copy(z_tmp, z_a)
		z_a.AddScalar(fs2) //fs2+z
		z_tmp.SubByScalar(fs2)
		z_z, _ := z_a.DivArray(z_tmp) //(fs2+z)/(fs2-z)
		k_z_f := z_tmp.Product()
		k_z := k * real(ComplexDivByReal(k_z_f, ComplexArrProd(p_tmp))) //k*real(prod(fs2-z)/prod(fs2-p))

		for i := 0; i < degree; i++ {
			z_z = append(z_z, -1)
		}
		return z_z, p_z, k_z, nil
	} else {
		z_c, _ := z.([]complex128)
		z_tmp := ComplexArrSubByReal(fs2, z_c)       //fs2-z
		z_c = ComplexArrAddReal(fs2, z_c)            //fs2+z
		z_z := ComplexArrDiv(z_c, z_tmp)             //(fs2+z)/(fs2-z)
		k_z_c := ComplexArrProd(z_tmp)               //prod(fs2-z)
		k_z := k * real(k_z_c/ComplexArrProd(p_tmp)) //k*real(prod(fs2-z)/prod(fs2-p))

		for i := 0; i < degree; i++ {
			z_z = append(z_z, -1+0i)
		}
		return z_z, p_z, k_z, nil
	}
	return nil, nil, 0, err
}

/**
  z = atleast_1d(z)
  k = atleast_1d(k)
  if len(z.shape) > 1:
      temp = poly(z[0])
      b = zeros((z.shape[0], z.shape[1] + 1), temp.dtype.char)
      if len(k) == 1:
          k = [k[0]] * z.shape[0]
      for i in range(z.shape[0]):
          b[i] = k[i] * poly(z[i])
  else:
      b = k * poly(z)
  a = atleast_1d(poly(p))

  # Use real output if possible. Copied from numpy.poly, since
  # we can't depend on a specific version of numpy.
  if issubclass(b.dtype.type, numpy.complexfloating):
      # if complex roots are all complex conjugates, the roots are real.
      roots = numpy.asarray(z, complex)
      pos_roots = numpy.compress(roots.imag > 0, roots)
      neg_roots = numpy.conjugate(numpy.compress(roots.imag < 0, roots))
      if len(pos_roots) == len(neg_roots):
          if numpy.all(numpy.sort_complex(neg_roots) ==
                       numpy.sort_complex(pos_roots)):
              b = b.real.copy()

  if issubclass(a.dtype.type, numpy.complexfloating):
      # if complex roots are all complex conjugates, the roots are real.
      roots = numpy.asarray(p, complex)
      pos_roots = numpy.compress(roots.imag > 0, roots)
      neg_roots = numpy.conjugate(numpy.compress(roots.imag < 0, roots))
      if len(pos_roots) == len(neg_roots):
          if numpy.all(numpy.sort_complex(neg_roots) ==
                       numpy.sort_complex(pos_roots)):
              a = a.real.copy()

  return b, a
*/
func zpk2tf(z interface{}, p []complex128, k float64) (b, a numgo.Array) {

	return
}

/**
  if len(seq_of_zeros) == 0:
      return 1.0
  dt = seq_of_zeros.dtype
  a = ones((1,), dtype=dt)
  for k in range(len(seq_of_zeros)):
      a = NX.convolve(a, array([1, -seq_of_zeros[k]], dtype=dt),
                      mode='full')

  if issubclass(a.dtype.type, NX.complexfloating):
      # if complex roots are all complex conjugates, the roots are real.
      roots = NX.asarray(seq_of_zeros, complex)
      if NX.all(NX.sort(roots) == NX.sort(roots.conjugate())):
          a = a.real.copy()

  return a

------------------------------------------------------------
def convolve(a, v, mode='full'):
    a, v = array(a, copy=False, ndmin=1), array(v, copy=False, ndmin=1)
    if (len(v) > len(a)):
        a, v = v, a
    if len(a) == 0:
        raise ValueError('a cannot be empty')
    if len(v) == 0:
        raise ValueError('v cannot be empty')
    mode = _mode_from_name(mode)
    return multiarray.correlate(a, v[::-1], mode)
//File:      c:\users\wangjp\appdata\local\programs\python\python38\lib\site-packages\numpy\core\numeric.py
*/
func poly(seq_of_zeros interface{}) numgo.Array {

	return nil
}

/***********************************************
功能:N阶IIR数字滤波器设计,返回滤波器系数
输入:
	N: 滤波阶数
	Wn:自然频率数组，也称归一化的截止频率，fcf=截止频率*2/采样频率；
		如果是低通，高通滤波，fcf只有一个元素；
		如果是带通，带阻滤波，fcf数组有两个元素。
	btype:滤波类型，字符串；lp：表示低通；hp：表示高通；bp：表示带通；bs：表示带阻；
输出:b系数和a系数
说明:
编辑:wang_jp
时间:
***********************************************/
//iirfilter(N, Wn, rp=None, rs=None, btype='band', analog=False,ftype='butter', output='ba', fs=None)
func iirFilter(N int, Wn interface{}, btype string, fs ...float64) (b, a numgo.Array, err error) {
	z, p, k := buttap(N)

	var z_b []complex128 //带宽滤波时候的z参数
	switch strings.ToLower(btype) {
	case "lowpass", "low", "lp", "highpass", "high", "hp":
		warped := 1.0
		wn, ok := Wn.(float64)
		if !ok {
			err = fmt.Errorf("Wn should be of float64 type when 'lowpass' or 'highpass'.[当进行'低通'或者'高通'滤波时,参数 Wn 应该是 float64类型]")
		}
		_fs := 2.0
		if len(fs) > 0 {
			if fs[0] > 0 {
				_fs = fs[0]
				wn = 2 * wn / _fs
			}
		}
		warped = 2 * _fs * math.Tan(pi*wn/_fs)
		switch strings.ToLower(btype) {
		case "lowpass", "low", "lp":
			err = lp2lpZpk(z, p, &k, warped)
			if err != nil {
				return
			}
		case "highpass", "high", "hp":
			err = lp2hpZpk(z, p, &k, warped)
			if err != nil {
				return
			}
		}
		zo, po, ko, er := bilinearZpk(z, p, k, _fs)
		if er != nil {
			err = er
			return
		}
		b, a = zpk2tf(zo, po, ko)
	case "bandpass", "bp", "pass", "bandstop", "bs", "stop":
		wn, ok := Wn.(numgo.Array)
		if !ok {
			err = fmt.Errorf("Wn should be of [2]float64 type array when 'bandpass' or 'bandstop'.[当‘带通’或者‘带阻’滤波时,参数 Wn 应该是 [2]float64类型的数组]")
			return
		}
		_fs := 2.0
		if len(fs) > 0 {
			if fs[0] > 0 {
				_fs = fs[0]
				for i, v := range wn {
					wn[i] = 2 * v / _fs
				}
			}
		}
		if len(wn) < 2 {
			err = fmt.Errorf("Wn must specify start and stop frequencies for 'bandpass' or 'bandstop' filter.[设计'带通'或者'带阻'滤波器时,Wn 应该是开始频率和结束频率的数组]")
			return
		}
		var warped numgo.Array
		for _, v := range wn {
			warped = append(warped, 2*_fs*math.Tan(pi*v/_fs))
		}
		bw := warped[1] - warped[0]
		wo := math.Sqrt(warped[0] * warped[1])

		switch strings.ToLower(btype) {
		case "bandpass", "bp", "pass":
			z_b, p, k, err = lp2bpZpk(Reals2Complexs(z), p, k, wo, bw)
			if err != nil {
				return
			}
		case "bandstop", "bs", "stop":
			z_b, p, k, err = lp2bsZpk(Reals2Complexs(z), p, k, wo, bw)
			if err != nil {
				return
			}
		}
		zo, po, ko, er := bilinearZpk(z_b, p, k, _fs)
		if er != nil {
			err = er
			return
		}
		b, a = zpk2tf(zo, po, ko)
	}

	return
}
