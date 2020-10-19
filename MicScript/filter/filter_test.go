// filter project filter.go
package filter

import (
	"math"
	"testing"

	"github.com/bkzy-wangjp/MicEngine/MicScript/numgo"
)

/*
func TestNmoalizaton(t *testing.T) {
	tests := []struct {
		a numgo.Array
		b numgo.Array
	}{
		{
			numgo.Array{1, 2, 3, 4, 5},
			numgo.Array{2, 3, 4, 5, 6},
		},
		{
			numgo.Array{2, 2, 3, 4, 5},
			numgo.Array{2, 3, 4, 5, 6},
		},
	}
	for _, tt := range tests {
		normalization(tt.b, tt.a)
		t.Log(tt.a)
		t.Log(tt.b)
	}
}
*/
func TestInitZi(t *testing.T) {
	tests := []struct {
		a  numgo.Array
		b  numgo.Array
		zi numgo.Array
	}{
		{
			numgo.Array{1., -1.45424359, 0.57406192},
			numgo.Array{0.02995458, 0.05990916, 0.02995458},
			numgo.Array{0.97004542, -0.54410733},
		},
	}
	for _, tt := range tests {
		zi, err := initZi(tt.b, tt.a)
		if err != nil {
			t.Log(err)
		} else {
			if zi.IsEqual(tt.zi) == false {
				t.Error("错误:期望值:", tt.zi, "得到值:", zi)
			}
		}
	}
}

func TestInitTx(t *testing.T) {
	tests := []struct {
		X     numgo.Array
		nfilt int
		res   numgo.Array
	}{
		{
			numgo.Array{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
			5,
			numgo.Array{-14.0, -13.0, -12.0, -11.0, -10.0, -9.0, -8.0, -7.0, -6.0, -5.0, -4.0, -3.0, -2.0, -1.0, 0.0,
				1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0, 11.0, 12.0, 13.0, 14.0, 15.0, 16.0, 17.0, 18.0, 19.0, 20.0,
				21.0, 22.0, 23.0, 24.0, 25.0, 26.0, 27.0, 28.0, 29.0, 30.0, 31.0, 32.0, 33.0, 34.0, 35.0},
		},
	}
	for _, tt := range tests {
		tx := initTx(tt.X, tt.nfilt)
		if tx.IsEqual(tt.res) == false {
			t.Error("初始化Tx错误,期望值:", tt.res, "得到值:", tx)
		}
	}
}

func TestFilter(t *testing.T) {
	tests := []struct {
		a   numgo.Array
		b   numgo.Array
		X   numgo.Array
		res numgo.Array
	}{
		{
			numgo.Array{1.0, -0.66817864},
			numgo.Array{0.16591068, 0.16591068},
			numgo.Array{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
			numgo.Array{1.11174407, 2.07457773, 3.04970063, 4.03301343, 5.02176632,
				6.01410596, 7.00877004, 8.00487929, 9.00179257, 9.99900123,
				10.99604532, 11.99243773, 12.987584, 13.98068431, 14.97060169,
				15.95567469, 16.93344356, 17.90024498, 18.85060832, 19.77635425},
		},
		{
			numgo.Array{1., -1.45424359, 0.57406192},
			numgo.Array{0.02995458, 0.05990916, 0.02995458},
			numgo.Array{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
			numgo.Array{0.89899996, 1.93065909, 2.9565825, 3.97582069, 4.98859359,
				5.99580663, 6.99872006, 7.99875272, 8.9973915, 9.99617193,
				10.99669297, 12.00062758, 13.00968872, 14.02550955, 15.04939642,
				16.08191685, 17.12229509, 18.16760849, 19.2118143, 20.24469196},
		},
	}
	for _, tt := range tests {
		y, err := Filtfilt(tt.b, tt.a, tt.X)
		if err != nil {
			t.Log(err)
		} else {
			if y.IsEqual(tt.res) == false {
				t.Error("滤波错误,期望值:", tt.res, "\n得到值:", y)
			}
		}
	}

}

func TestButter(t *testing.T) {
	//var c1 complex128 = complex(0, 1.0)
	n := 2.0
	m := 5.3
	p := math.Pow(n, m)
	t.Log(p)
}
