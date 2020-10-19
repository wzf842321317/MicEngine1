package numgo

import (
	"testing"
)

func TestMul(t *testing.T) {
	a := Array{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}
	tests := []struct {
		arr Array
		n   int
		fv  string
	}{
		{a, 1, "mean"},
		{a, 4, "median"},
		{a, 3, "mode"},
	}

	for _, tt := range tests {
		res := tt.arr.MoveWindowFilter(tt.n, tt.fv)
		t.Log(res)
	}
}

func TestReverse(t *testing.T) {
	tests := []struct {
		arr Array
		res Array
	}{
		{Array{1, 2, 3, 4, 5, 6, 7, 8}, Array{8, 7, 6, 5, 4, 3, 2, 1}},
	}

	for _, tt := range tests {
		tt.arr.Reverse()
		if tt.arr.IsEqual(tt.res) == false {
			t.Error("错误:期望值:", tt.res, "得到值:", tt.arr)
		}
	}
}

func TestCompanion(t *testing.T) {
	tests := []struct {
		arr Array
		res Matrix
	}{
		{Array{1, 2}, Matrix{Array{-2}}},
		{Array{1, -10, 31, -30}, Matrix{Array{10., -31., 30.}, Array{1, 0, 0}, Array{0, 1, 0}}},
	}

	for _, tt := range tests {
		mt, err := tt.arr.Companion()
		if err != nil {
			t.Log(err)
		} else if mt.IsEqual(tt.res) == false {
			t.Error("错误:期望值:", tt.res, "得到值:", tt.arr)
		}
	}
}
