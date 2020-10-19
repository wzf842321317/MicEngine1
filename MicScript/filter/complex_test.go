// filter project filter.go
package filter

import (
	"testing"
)

func TestRealDivComplex(t *testing.T) {
	tests := []struct {
		r   float64
		c   complex128
		res complex128
	}{
		{1.0, 2 + 1i, 0.4 - 0.2i},
		{3.0, 1i, -3i},
		{123.0, 20 + 0i, 6.15 + 0i},
	}
	for _, tt := range tests {
		res := RealDivComplex(tt.r, tt.c)
		if res != tt.res {
			t.Error("错误:期望值:", tt.res, "得到值:", res)
		}
	}
}

func TestRealMulComplex(t *testing.T) {
	tests := []struct {
		r   float64
		c   complex128
		res complex128
	}{
		{1.0, 2 + 1i, 2 + 1i},
		{3.0, 1i, 3i},
		{123.0, 20 + 0i, 2460 + 0i},
	}
	for _, tt := range tests {
		res := RealMulComplex(tt.r, tt.c)
		if res != tt.res {
			t.Error("错误:期望值:", tt.res, "得到值:", res)
		}
	}
}

func TestComplexArrMulReal(t *testing.T) {
	tests := []struct {
		comps []complex128
		r     float64
		res   []complex128
	}{
		{[]complex128{(-0.5000000000000001 + 0.8660254037844386i), (-1 + 0i), (-0.5000000000000001 - 0.8660254037844386i)}, 1.0,
			[]complex128{(-0.5000000000000001 + 0.8660254037844386i), (-1 + 0i), (-0.5000000000000001 - 0.8660254037844386i)}},
		{[]complex128{(-0.5000000000000001 + 0.8660254037844386i), (-1 + 0i), (-0.5000000000000001 - 0.8660254037844386i)}, -1.0,
			[]complex128{(0.5000000000000001 - 0.8660254037844386i), (1 - 0i), (0.5000000000000001 + 0.8660254037844386i)}},
		{[]complex128{(-0.5000000000000001 + 0.8660254037844386i), (-1 + 0i), (-0.5000000000000001 - 0.8660254037844386i)}, -5.0,
			[]complex128{(2.5000000000000004 - 4.330127018922193i), (5 - 0i), (2.5000000000000004 + 4.330127018922193i)}},
	}
	for _, tt := range tests {
		res := ComplexArrMulReal(tt.comps, tt.r)
		if !ComplexArrsIsEqual(res, tt.res) {
			t.Error("错误:期望值:", tt.res, "得到值:", res)
		}
	}
}

func TestComplexArrProd(t *testing.T) {
	tests := []struct {
		comps []complex128
		res   complex128
	}{
		{[]complex128{}, 1 + 0i},
		{[]complex128{(-0.5000000000000001 + 0.8660254037844386i), (-1 + 0i), (-0.5000000000000001 - 0.8660254037844386i)}, -1 + 0i},
		{[]complex128{-0.82842712 - 1.43487787i, -1.65685425 - 0.i, -0.82842712 + 1.43487787i}, -4.548339945716935 + 2.220446049250313e-16i},
		{[]complex128{-0.63405067 - 1.53073373i, -1.53073373 - 0.63405067i, -1.53073373 + 0.63405067i, -0.63405067 + 1.53073373i}, 7.535936389852008 + 4.440892098500626e-16i},
	}
	for _, tt := range tests {
		res := ComplexArrProd(tt.comps)
		if !ComplexsIsEqual(res, tt.res) {
			t.Error("错误:期望值:", tt.res, "得到值:", res)
		}
	}
}
