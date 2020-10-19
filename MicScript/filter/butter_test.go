package filter

import (
	"math"
	"testing"

	"github.com/bkzy-wangjp/MicEngine/MicScript/numgo"
)

func TestButtap(t *testing.T) {
	tests := []struct {
		n int
		z numgo.Array
		p []complex128
		k float64
	}{
		{0, numgo.Array{}, []complex128{}, 1},
		{1, numgo.Array{}, []complex128{-1 - 0i}, 1},
		{2, numgo.Array{}, []complex128{(-0.7071067811865476 + 0.7071067811865475i), (-0.7071067811865476 - 0.7071067811865475i)}, 1},
		{3, numgo.Array{}, []complex128{(-0.5000000000000001 + 0.8660254037844386i), (-1 + 0i), (-0.5000000000000001 - 0.8660254037844386i)}, 1},
		{4, numgo.Array{}, []complex128{(-0.38268343236508984 + 0.9238795325112867i), (-0.9238795325112867 + 0.3826834323650898i), (-0.9238795325112867 - 0.3826834323650898i), (-0.38268343236508984 - 0.9238795325112867i)}, 1},
		{5, numgo.Array{}, []complex128{(-0.30901699437494745 + 0.9510565162951535i), (-0.8090169943749473 + 0.5877852522924731i), (-1 + 0i), (-0.8090169943749473 - 0.5877852522924731i), (-0.30901699437494745 - 0.9510565162951535i)}, 1},
	}
	for _, tt := range tests {
		z, p, k := buttap(tt.n)
		if !tt.z.IsEqual(z) {
			t.Error("z错误:期望值:", tt.z, "得到值:", z)
		}
		if !ComplexArrsIsEqual(tt.p, p) {
			t.Error("p错误:期望值:", tt.p, "得到值:", p)
		}
		if math.Abs(tt.k-k) > 1e-5 {
			t.Error("k错误:期望值:", tt.k, "得到值:", k)
		}
	}
}

func TestLp2lpZpk(t *testing.T) {
	tests := []struct {
		z  numgo.Array
		p  []complex128
		k  float64
		wo float64
		oz numgo.Array
		op []complex128
		ok float64
	}{
		{numgo.Array{}, []complex128{}, 1, 1.6568542494923801,
			numgo.Array{}, []complex128{}, 1},
		{numgo.Array{}, []complex128{-1 - 0i}, 1, 1.6568542494923801,
			numgo.Array{}, []complex128{(-1.6568542494923801 + 0i)}, 1.6568542494923801},
		{numgo.Array{}, []complex128{(-0.7071067811865476 + 0.7071067811865475i), (-0.7071067811865476 - 0.7071067811865475i)}, 1, 1.6568542494923801,
			numgo.Array{}, []complex128{(-1.17157287525381 + 1.1715728752538097i), (-1.17157287525381 - 1.1715728752538097i)}, 2.7451660040609585},
		{numgo.Array{}, []complex128{(-0.5000000000000001 + 0.8660254037844386i), (-1 + 0i), (-0.5000000000000001 - 0.8660254037844386i)}, 1, 1.6568542494923801,
			numgo.Array{}, []complex128{(-0.8284271247461903 + 1.4348778704286014i), (-1.6568542494923801 + 0i), (-0.8284271247461903 - 1.4348778704286014i)}, 4.5483399593904155},
		{numgo.Array{}, []complex128{(-0.38268343236508984 + 0.9238795325112867i), (-0.9238795325112867 + 0.3826834323650898i), (-0.9238795325112867 - 0.3826834323650898i), (-0.38268343236508984 - 0.9238795325112867i)}, 1, 1.6568542494923801,
			numgo.Array{}, []complex128{(-0.6340506711244289 + 1.530733729460359i), (-1.530733729460359 + 0.6340506711244288i), (-1.530733729460359 - 0.6340506711244288i), (-0.6340506711244289 - 1.530733729460359i)}, 7.53593638985201},
		{numgo.Array{}, []complex128{(-0.30901699437494745 + 0.9510565162951535i), (-0.8090169943749473 + 0.5877852522924731i), (-1 + 0i), (-0.8090169943749473 - 0.5877852522924731i), (-0.30901699437494745 - 0.9510565162951535i)}, 1, 1.6568542494923801,
			numgo.Array{}, []complex128{(-0.5119961202954946 + 1.5757620305310442i), (-1.3404232450416844 + 0.973874493049735i), (-1.6568542494923801 + 0i), (-1.3404232450416844 - 0.973874493049735i), (-0.5119961202954946 - 1.5757620305310442i)}, 12.485948231430568},
	}
	for _, tt := range tests {
		err := lp2lpZpk(tt.z, tt.p, &tt.k, tt.wo)
		if err != nil {
			t.Log(err)
		} else {
			if !tt.z.IsEqual(tt.oz) {
				t.Error("z错误:期望值:", tt.oz, "得到值:", tt.z)
			}
			if !ComplexArrsIsEqual(tt.p, tt.op) {
				t.Error("p错误:期望值:", tt.op, "得到值:", tt.p)
			}
			if math.Abs(tt.k-tt.ok) > 1e-5 {
				t.Error("k错误:期望值:", tt.ok, "得到值:", tt.k)
			}
		}
	}
}

func TestLp2hpZpk(t *testing.T) {
	tests := []struct {
		z  numgo.Array
		p  []complex128
		k  float64
		wo float64
		oz numgo.Array
		op []complex128
		ok float64
	}{
		{numgo.Array{}, []complex128{}, 1, 1.6568542494923801,
			numgo.Array{}, []complex128{}, 1},
		{numgo.Array{}, []complex128{-1 - 0i}, 1, 1.6568542494923801,
			numgo.Array{}, []complex128{(-1.6568542494923801 + 0i)}, 1},
		{numgo.Array{}, []complex128{(-0.7071067811865476 + 0.7071067811865475i), (-0.7071067811865476 - 0.7071067811865475i)}, 1, 1.6568542494923801,
			numgo.Array{}, []complex128{(-1.17157287525381 - 1.1715728752538097i), (-1.17157287525381 + 1.1715728752538097i)}, 1},
		{numgo.Array{}, []complex128{(-0.5000000000000001 + 0.8660254037844386i), (-1 + 0i), (-0.5000000000000001 - 0.8660254037844386i)}, 1, 1.6568542494923801,
			numgo.Array{}, []complex128{(-0.8284271247461903 - 1.4348778704286014i), (-1.6568542494923801 - 0i), (-0.8284271247461903 + 1.4348778704286014i)}, 1},
		{numgo.Array{}, []complex128{(-0.38268343236508984 + 0.9238795325112867i), (-0.9238795325112867 + 0.3826834323650898i), (-0.9238795325112867 - 0.3826834323650898i), (-0.38268343236508984 - 0.9238795325112867i)}, 1, 1.6568542494923801,
			numgo.Array{}, []complex128{(-0.6340506711244289 - 1.530733729460359i), (-1.530733729460359 - 0.6340506711244288i), (-1.530733729460359 + 0.6340506711244288i), (-0.6340506711244289 + 1.530733729460359i)}, 1},
		{numgo.Array{}, []complex128{(-0.30901699437494745 + 0.9510565162951535i), (-0.8090169943749473 + 0.5877852522924731i), (-1 + 0i), (-0.8090169943749473 - 0.5877852522924731i), (-0.30901699437494745 - 0.9510565162951535i)}, 1, 1.6568542494923801,
			numgo.Array{}, []complex128{(-0.5119961202954946 - 1.5757620305310442i), (-1.3404232450416844 - 0.973874493049735i), (-1.6568542494923801 + 0i), (-1.3404232450416844 + 0.973874493049735i), (-0.5119961202954946 + 1.5757620305310442i)}, 1},
	}
	for _, tt := range tests {
		err := lp2hpZpk(tt.z, tt.p, &tt.k, tt.wo)
		if err != nil {
			t.Log(err)
		} else {
			if !tt.z.IsEqual(tt.oz) {
				t.Error("z错误:期望值:", tt.oz, "得到值:", tt.z)
			}
			if !ComplexArrsIsEqual(tt.p, tt.op) {
				t.Error("p错误:期望值:", tt.op, "得到值:", tt.p)
			}
			if math.Abs(tt.k-tt.ok) > 1e-5 {
				t.Error("k错误:期望值:", tt.ok, "得到值:", tt.k)
			}
		}
	}
}

func TestLp2bpZpk(t *testing.T) {
	tests := []struct {
		z  numgo.Array
		p  []complex128
		k  float64
		wo float64
		bw float64
		oz []complex128
		op []complex128
		ok float64
	}{
		{numgo.Array{}, []complex128{}, 1, 1.6568542494923801, 1,
			[]complex128{}, []complex128{}, 1},
		{numgo.Array{}, []complex128{-1 - 0i}, 1, 1.6568542494923801, 1,
			[]complex128{0. + 0.i}, []complex128{-0.5 - 1.57960945i, -0.5 + 1.57960945i}, 1},
		{numgo.Array{}, []complex128{(-0.7071067811865476 + 0.7071067811865475i), (-0.7071067811865476 - 0.7071067811865475i)}, 1, 1.6568542494923801, 1,
			[]complex128{0. + 0.i, 0. + 0.i},
			[]complex128{-0.27818715 - 1.30501409i, -0.27818715 + 1.30501409i,
				-0.42891963 + 2.01212087i, -0.42891963 - 2.01212087i}, 1},
		{numgo.Array{}, []complex128{(-0.5000000000000001 + 0.8660254037844386i), (-1 + 0i), (-0.5000000000000001 - 0.8660254037844386i)}, 1, 1.6568542494923801, 1,
			[]complex128{0. + 0.i, 0. + 0.i, 0. + 0.i},
			[]complex128{-0.18614736 - 1.2623466i, -0.5 - 1.57960945i,
				-0.18614736 + 1.2623466i, -0.31385264 + 2.128372i,
				-0.5 + 1.57960945i, -0.31385264 - 2.128372i}, 1},
		{numgo.Array{}, []complex128{(-0.38268343236508984 + 0.9238795325112867i), (-0.9238795325112867 + 0.3826834323650898i), (-0.9238795325112867 - 0.3826834323650898i), (-0.38268343236508984 - 0.9238795325112867i)}, 1, 1.6568542494923801, 1,
			[]complex128{0. + 0.i, 0. + 0.i, 0. + 0.i, 0. + 0.i},
			[]complex128{-0.13965717 - 1.24821052i, -0.40681994 - 1.4122254i,
				-0.40681994 + 1.4122254i, -0.13965717 + 1.24821052i,
				-0.24302627 + 2.17209005i, -0.5170596 + 1.79490883i,
				-0.5170596 - 1.79490883i, -0.24302627 - 2.17209005i}, 1},
		{numgo.Array{}, []complex128{(-0.30901699437494745 + 0.9510565162951535i), (-0.8090169943749473 + 0.5877852522924731i), (-1 + 0i), (-0.8090169943749473 - 0.5877852522924731i), (-0.30901699437494745 - 0.9510565162951535i)}, 1, 1.6568542494923801, 1,
			[]complex128{0. + 0.i, 0. + 0.i, 0. + 0.i, 0. + 0.i, 0. + 0.i},
			[]complex128{-0.11172534 - 1.24180999i, -0.3317974 - 1.34109932i,
				-0.5 - 1.57960945i, -0.3317974 + 1.34109932i,
				-0.11172534 + 1.24180999i, -0.19729166 + 2.1928665i,
				-0.4772196 + 1.92888457i, -0.5 + 1.57960945i,
				-0.4772196 - 1.92888457i, -0.19729166 - 2.1928665i}, 1},
	}
	for i, tt := range tests {
		iz := Reals2Complexs(tt.z)
		z, p, k, err := lp2bpZpk(iz, tt.p, tt.k, tt.wo, tt.bw)
		if err != nil {
			t.Log(err)
		} else {
			if !ComplexArrsIsEqual(z, tt.oz) {
				t.Error(i, "z错误:\n期望值:", tt.oz, "\n得到值:", z)
			}
			if !ComplexArrsIsEqual(p, tt.op) {
				t.Error(i, "p错误:\n期望值:", tt.op, "\n得到值:", p)
			}
			if math.Abs(k-tt.ok) > 1e-5 {
				t.Error(i, "k错误:\n期望值:", tt.ok, "\n得到值:", k)
			}
		}
	}
}

func TestLp2bsZpk(t *testing.T) {
	tests := []struct {
		z  numgo.Array
		p  []complex128
		k  float64
		wo float64
		bw float64
		oz []complex128
		op []complex128
		ok float64
	}{
		{numgo.Array{}, []complex128{}, 1, 1.6568542494923801, 1,
			[]complex128{}, []complex128{}, 1},
		{numgo.Array{}, []complex128{-1 - 0i}, 1, 1.6568542494923801, 1,
			[]complex128{0. + 1.65685425i, 0. - 1.65685425i}, []complex128{-0.5 - 1.57960945i, -0.5 + 1.57960945i}, 1},
		{numgo.Array{}, []complex128{(-0.7071067811865476 + 0.7071067811865475i), (-0.7071067811865476 - 0.7071067811865475i)}, 1, 1.6568542494923801, 1,
			[]complex128{0. + 1.65685425i, 0. + 1.65685425i, 0. - 1.65685425i, 0. - 1.65685425i},
			[]complex128{-0.27818715 + 1.30501409i, -0.27818715 - 1.30501409i,
				-0.42891963 - 2.01212087i, -0.42891963 + 2.01212087i}, 1},
		{numgo.Array{}, []complex128{(-0.5000000000000001 + 0.8660254037844386i), (-1 + 0i), (-0.5000000000000001 - 0.8660254037844386i)}, 1, 1.6568542494923801, 1,
			[]complex128{0. + 1.65685425i, 0. + 1.65685425i, 0. + 1.65685425i, 0. - 1.65685425i, 0. - 1.65685425i, 0. - 1.65685425i},
			[]complex128{-0.18614736 + 1.2623466i, -0.5 - 1.57960945i,
				-0.18614736 - 1.2623466i, -0.31385264 - 2.128372i,
				-0.5 + 1.57960945i, -0.31385264 + 2.128372i}, 1},
		{numgo.Array{}, []complex128{(-0.38268343236508984 + 0.9238795325112867i), (-0.9238795325112867 + 0.3826834323650898i), (-0.9238795325112867 - 0.3826834323650898i), (-0.38268343236508984 - 0.9238795325112867i)}, 1, 1.6568542494923801, 1,
			[]complex128{0. + 1.65685425i, 0. + 1.65685425i, 0. + 1.65685425i, 0. + 1.65685425i,
				0. - 1.65685425i, 0. - 1.65685425i, 0. - 1.65685425i, 0. - 1.65685425i},
			[]complex128{-0.13965717 + 1.24821052i, -0.40681994 + 1.4122254i,
				-0.40681994 - 1.4122254i, -0.13965717 - 1.24821052i,
				-0.24302627 - 2.17209005i, -0.5170596 - 1.79490883i,
				-0.5170596 + 1.79490883i, -0.24302627 + 2.17209005i}, 1},
		{numgo.Array{}, []complex128{(-0.30901699437494745 + 0.9510565162951535i), (-0.8090169943749473 + 0.5877852522924731i), (-1 + 0i), (-0.8090169943749473 - 0.5877852522924731i), (-0.30901699437494745 - 0.9510565162951535i)}, 1, 1.6568542494923801, 1,
			[]complex128{0. + 1.65685425i, 0. + 1.65685425i, 0. + 1.65685425i, 0. + 1.65685425i,
				0. + 1.65685425i, 0. - 1.65685425i, 0. - 1.65685425i, 0. - 1.65685425i,
				0. - 1.65685425i, 0. - 1.65685425i},
			[]complex128{-0.11172534 + 1.24180999i, -0.3317974 + 1.34109932i,
				-0.5 - 1.57960945i, -0.3317974 - 1.34109932i,
				-0.11172534 - 1.24180999i, -0.19729166 - 2.1928665i,
				-0.4772196 - 1.92888457i, -0.5 + 1.57960945i,
				-0.4772196 + 1.92888457i, -0.19729166 + 2.1928665i}, 1},
	}
	for i, tt := range tests {
		iz := Reals2Complexs(tt.z)
		z, p, k, err := lp2bsZpk(iz, tt.p, tt.k, tt.wo, tt.bw)
		if err != nil {
			t.Log(err)
		} else {
			if !ComplexArrsIsEqual(z, tt.oz) {
				t.Error(i, "z错误:\n期望值:", tt.oz, "\n得到值:", z)
			}
			if !ComplexArrsIsEqual(p, tt.op) {
				t.Error(i, "p错误:\n期望值:", tt.op, "\n得到值:", p)
			}
			if math.Abs(k-tt.ok) > 1e-5 {
				t.Error(i, "k错误:\n期望值:", tt.ok, "\n得到值:", k)
			}
		}
	}
}

func TestShiftband(t *testing.T) {
	tests := []struct {
		in  []complex128
		wo  float64
		res []complex128
	}{
		{[]complex128{-0.35355339 - 0.35355339i, -0.35355339 + 0.35355339i}, 1.6568542494923801,
			[]complex128{-0.27818715 + 1.30501409i, -0.27818715 - 1.30501409i,
				-0.42891963 - 2.01212087i, -0.42891963 + 2.01212087i}},
		{[]complex128{-0.25 - 0.4330127i, -0.5 - 0.i, -0.25 + 0.4330127i}, 1.6568542494923801,
			[]complex128{-0.18614736 + 1.2623466i, -0.5 - 1.57960945i,
				-0.18614736 - 1.2623466i, -0.31385264 - 2.128372i,
				-0.5 + 1.57960945i, -0.31385264 + 2.128372i}},
	}
	for i, tt := range tests {
		res := shiftband(tt.in, tt.wo)
		if !ComplexArrsIsEqual(tt.res, res) {
			t.Error(i, "p错误:\n期望值:", tt.res, "\n得到值:", res)
		}
	}
}