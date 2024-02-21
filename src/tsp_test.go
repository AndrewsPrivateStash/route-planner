package main

import (
	"fmt"
	"math"
	"math/big"
	"testing"
)

const floatErrorMax = 1e-6
const dat1 = `C:\Users\Andrew Pfaendler\Desktop\Code\go-work\src\tssOpt\routeplanner\testData\test.dat`
const dat2 = `C:\Users\Andrew Pfaendler\Desktop\Code\go-work\src\tssOpt\routeplanner\testData\test2.dat`
const dat3 = `C:\Users\Andrew Pfaendler\Desktop\Code\go-work\src\tssOpt\routeplanner\testData\ctrTest.dat`

// test dToR
func TestDtR(t *testing.T) {

	var cases = []pnts{
		pnts{point{0, 0, ""}, point{0, 0, ""}},
		pnts{point{45, 0, ""}, point{0.785398, 0, ""}},
		pnts{point{90, 0, ""}, point{1.570796, 0, ""}},
		pnts{point{135, 0, ""}, point{2.356194, 0, ""}},
		pnts{point{180, 0, ""}, point{3.141593, 0, ""}},
		pnts{point{0, 45, ""}, point{0, 0.785398, ""}},
		pnts{point{45, 45, ""}, point{0.785398, 0.785398, ""}},
		pnts{point{90, 45, ""}, point{1.570796, 0.785398, ""}},
		pnts{point{135, 45, ""}, point{2.356194, 0.785398, ""}},
		pnts{point{180, 45, ""}, point{3.141593, 0.785398, ""}},
		pnts{point{0, 90, ""}, point{0, 1.570796, ""}},
		pnts{point{45, 90, ""}, point{0.785398, 1.570796, ""}},
		pnts{point{90, 90, ""}, point{1.570796, 1.570796, ""}},
		pnts{point{135, 90, ""}, point{2.356194, 1.570796, ""}},
		pnts{point{180, 90, ""}, point{3.141593, 1.570796, ""}},
		pnts{point{0, 135, ""}, point{0, 2.356194, ""}},
		pnts{point{45, 135, ""}, point{0.785398, 2.356194, ""}},
		pnts{point{90, 135, ""}, point{1.570796, 2.356194, ""}},
		pnts{point{135, 135, ""}, point{2.356194, 2.356194, ""}},
		pnts{point{180, 135, ""}, point{3.141593, 2.356194, ""}},
		pnts{point{0, 180, ""}, point{0, 3.141593, ""}},
		pnts{point{45, 180, ""}, point{0.785398, 3.141593, ""}},
		pnts{point{90, 180, ""}, point{1.570796, 3.141593, ""}},
		pnts{point{135, 180, ""}, point{2.356194, 3.141593, ""}},
		pnts{point{180, 180, ""}, point{3.141593, 3.141593, ""}},
	}
	for _, tst := range cases {
		tstPrime := tst[0].dToR()
		latDiff := math.Abs(tstPrime.lat - tst[1].lat)
		lonDiff := math.Abs(tstPrime.lon - tst[1].lon)
		if latDiff > floatErrorMax || lonDiff > floatErrorMax {
			t.Errorf("%v DtR--> %v does not match %v", tst[0], tstPrime, tst[1])
		}
	}

}

// test fact
func TestFact(t *testing.T) {
	var cases = [][]int{
		{0, 1}, {1, 1}, {2, 2},
		{3, 6}, {4, 24}, {5, 120},
		{6, 720}, {7, 5040}, {8, 40320},
		{9, 362880}, {10, 3628800}, {11, 39916800},
	}
	for _, tst := range cases {
		val := fact(tst[0])
		cor := big.NewInt(int64(tst[1]))
		if val.Cmp(cor) != 0 {
			t.Errorf("fact(%d) expected %d received %d", tst[0], tst[1], val)
		}
	}

}

// test factf
func TestFactf(t *testing.T) {
	var cases = [][]int{
		{0, 1}, {1, 1}, {2, 2},
		{3, 6}, {4, 24}, {5, 120},
		{6, 720}, {7, 5040}, {8, 40320},
		{9, 362880}, {10, 3628800}, {11, 39916800},
	}
	for _, tst := range cases {
		val := factf(tst[0])
		cor := big.NewFloat(float64(tst[1]))
		if val.Cmp(cor) != 0 {
			t.Errorf("fact(%d) expected %f received %f", tst[0], cor, val)
		}
	}

}

// test myP
type powers struct {
	base float64
	exp  int
	cor  float64
}

func TestMyP(t *testing.T) {
	var cases = []powers{
		{0, 5, 0},
		{1, 5, 1},
		{-1, 5, -1},
		{2, 5, 32},
		{1.4, 5, 5.37824},
		{-1.4, 5, -5.37824},
		{3.14159265, 5, 306.019684},
	}
	for _, tst := range cases {
		val := myP(tst.base, tst.exp)
		diff := math.Abs(val - tst.cor)
		if diff > floatErrorMax {
			t.Errorf("myP(%f,%d) expected %f received %f", tst.base, tst.exp, tst.cor, val)
		}
	}
}

// test sqr
func TestSqr(t *testing.T) {
	var cases = [][]float64{
		{0, 0},
		{1, 1},
		{-1, 1},
		{1.4, 1.96},
		{-1.4, 1.96},
		{3.14159265, 9.869604},
	}
	for _, tst := range cases {
		val := sqr(tst[0])
		diff := math.Abs(val - tst[1])
		if diff > floatErrorMax {
			t.Errorf("sqr(%f) expected %f received %f", tst[0], tst[1], val)
		}
	}
}

// test haver
type pairs struct {
	p1  point
	p2  point
	cor float64
}

func TestHaver(t *testing.T) {
	var cases = []pairs{
		{point{0, 0, ""}, point{0, 0, ""}, 0},
		{point{45, 45, ""}, point{45, 45, ""}, 0},
		{point{45, -45, ""}, point{45, 45, ""}, 6679.130701},
		{point{45, 45, ""}, point{-45, 45, ""}, 10018.696052},
		{point{84.9999744, -135.0006867, "NP"}, point{-72.2940075, 0.6939949, "SP"}, 18418.845891},
		{point{45.5428626, -122.794813, "OR"}, point{42.752916, -71.5669218, "NH"}, 4033.73557971},
		{point{45.5214857, -122.8324972, "AP"}, point{45.5744697, -122.566121, "PR"}, 21.587481},
	}
	for _, tst := range cases {
		val := haver(tst.p1, tst.p2)
		diff := math.Abs(val - tst.cor)
		if diff > floatErrorMax {
			t.Errorf("haver(%v,%v) expected %f received %f", tst.p1, tst.p2, tst.cor, val)
		}
	}

}

func BenchmarkHaver(b *testing.B) {
	var cases = pnts{point{45.5428626, -122.794813, "OR"}, point{42.752916, -71.5669218, "NH"}}
	for n := 0; n < b.N; n++ {
		haver(cases[0], cases[1])
	}
}

// test tourLen
type tours struct {
	pnts
	len float64
}

func TestTourLen(t *testing.T) {
	var cases = []tours{
		{pnts{point{0, 0, ""}, point{0, 0, ""}, point{0, 0, ""}}, 0},
		{pnts{point{45, 45, ""}, point{45, 45, ""}, point{45, 45, ""}}, 0},
		{pnts{point{0, 0, ""}, point{1, 0, ""}, point{0, 1, ""}}, 380.0623139},
		{pnts{point{0, 0, ""}, point{45, 45, ""}, point{89, 89, ""}}, 21625.632737},
		{pnts{point{0, 0, ""}, point{-45, -45, ""}, point{-89, -89, ""}}, 21625.632737},
	}
	for _, tst := range cases {
		val := tst.tourLen()
		diff := math.Abs(val - tst.len)
		if diff > floatErrorMax {
			t.Errorf("tourLen(%v) expected %f received %f", tst.pnts, tst.len, val)
		}
	}

}

// test oTourLen
type oTours struct {
	ord []int
	cor float64
}

func TestOTourLen(t *testing.T) {
	var tour = pnts{
		point{0, 0, ""},
		point{1, 1, ""},
		point{-1, 1, ""},
		point{0, 2, ""},
	}

	var cases = []oTours{
		{[]int{0, 1, 2, 3}, 760.1246278},
		{[]int{0, 1, 3, 2}, 629.6984954},
		{[]int{0, 2, 1, 3}, 760.1246278},
		{[]int{0, 2, 3, 1}, 629.6984954},
		{[]int{0, 3, 2, 1}, 760.1246278},
		{[]int{0, 3, 1, 2}, 760.1246278},
	}
	for _, tst := range cases {
		val := tour.oTourLen(tst.ord)
		diff := math.Abs(val - tst.cor)
		if diff > floatErrorMax {
			t.Errorf("oTourLen(%v) expected %f received %f", tst.ord, tst.cor, val)
		}
	}

}

// test nearest
type nst struct {
	fnd point
	dat pnts
	res point
}

func TestNearest(t *testing.T) {
	var cases = []nst{
		{point{0, 0, "1"}, pnts{point{0, 0, "1"}, point{0, 0, "2"}, point{0, 0, "3"}}, point{0, 0, "2"}},
		{point{0, 0, "1"}, pnts{point{0, 0, "1"}, point{2, 4, "2"}, point{-5, 5, "3"}}, point{2, 4, "2"}},
		{point{-5, 5, "3"}, pnts{point{0, 0, "1"}, point{2, 4, "2"}, point{-5, 5, "3"}}, point{0, 0, "1"}},
		{point{-8, 8, ""}, pnts{point{0, 0, "1"}, point{2, 4, "2"}, point{-5, 5, "3"}}, point{-5, 5, "3"}},
	}
	for _, tst := range cases {
		val, _ := tst.dat.nearest(tst.fnd, false)
		if val != tst.res {
			t.Errorf("nearest(%v) for %v, expected %v received %v", tst.fnd, tst.dat, tst.res, val)
		}
	}
	// TODO add true flag test
}

// test rem
func TestRem(t *testing.T) {
	var cases = [][]pnts{
		{ //drop ix 0
			pnts{point{0, 0, ""}, point{2, 4, ""}, point{-5, 5, ""}},
			pnts{point{2, 4, ""}, point{-5, 5, ""}},
		},
		{ //drop ix 1
			pnts{point{0, 0, ""}, point{2, 4, ""}, point{-5, 5, ""}},
			pnts{point{0, 0, ""}, point{-5, 5, ""}},
		},
		{ //drop ix 2
			pnts{point{0, 0, ""}, point{2, 4, ""}, point{-5, 5, ""}},
			pnts{point{0, 0, ""}, point{2, 4, ""}},
		},
	}

	comp := func(v pnts, c pnts) bool {
		if len(v) != len(c) {
			return true
		}
		for i := 0; i < len(v); i++ {
			if v[i] != c[i] {
				return true
			}
		}
		return false
	}

	for i, tst := range cases {
		val := make(pnts, len(tst[0]))
		copy(val, tst[0])
		val.rem(i)
		if comp(val, tst[1]) {
			t.Errorf("rem(%d) for %v, expected %v received %v", i, tst[0], tst[1], val)
		}
	}
}

// test swap
type swp struct {
	st   pnts
	ed   pnts
	i, j int
}

func TestSwap(t *testing.T) {
	var cases = []swp{
		{ //swap ix 0,1
			pnts{point{0, 0, ""}, point{2, 4, ""}, point{-5, 5, ""}},
			pnts{point{2, 4, ""}, point{0, 0, ""}, point{-5, 5, ""}},
			0, 1,
		},
		{ //swap ix 0,2
			pnts{point{0, 0, ""}, point{2, 4, ""}, point{-5, 5, ""}},
			pnts{point{-5, 5, ""}, point{2, 4, ""}, point{0, 0, ""}},
			0, 2,
		},
		{ //swap ix 1,2
			pnts{point{0, 0, ""}, point{2, 4, ""}, point{-5, 5, ""}},
			pnts{point{0, 0, ""}, point{-5, 5, ""}, point{2, 4, ""}},
			1, 2,
		},
	}

	comp := func(v pnts, c pnts) bool {
		if len(v) != len(c) {
			return true
		}
		for i := 0; i < len(v); i++ {
			if v[i] != c[i] {
				return true
			}
		}
		return false
	}

	for _, tst := range cases {
		val := tst.st.swap(tst.i, tst.j)
		if comp(val, tst.ed) {
			t.Errorf("swap(%d, %d) for %v, expected %v received %v", tst.i, tst.j, tst.st, tst.ed, val)
		}
	}
}

// test rev
func TestRev(t *testing.T) {
	var cases = [][]pnts{
		{
			pnts{point{0, 0, ""}, point{2, 4, ""}, point{-5, 5, ""}},
			pnts{point{-5, 5, ""}, point{2, 4, ""}, point{0, 0, ""}},
		},
		{
			pnts{point{2, 4, ""}, point{-5, 5, ""}, point{0, 0, ""}},
			pnts{point{0, 0, ""}, point{-5, 5, ""}, point{2, 4, ""}},
		},
		{
			pnts{point{-5, 5, ""}, point{0, 0, ""}, point{2, 4, ""}},
			pnts{point{2, 4, ""}, point{0, 0, ""}, point{-5, 5, ""}},
		},
	}

	comp := func(v pnts, c pnts) bool {
		if len(v) != len(c) {
			return true
		}
		for i := 0; i < len(v); i++ {
			if v[i] != c[i] {
				return true
			}
		}
		return false
	}

	for _, tst := range cases {
		val := rev(tst[0])
		if comp(val, tst[1]) {
			t.Errorf("rev(%v), expected %v received %v", tst[0], tst[1], val)
		}
	}
}

// test rot
func TestRot(t *testing.T) {
	var cases = []swp{
		{ //rotate to ix 0
			pnts{point{0, 0, ""}, point{2, 4, ""}, point{-5, 5, ""}},
			pnts{point{0, 0, ""}, point{2, 4, ""}, point{-5, 5, ""}},
			0, 0,
		},
		{ //rotate to ix 1
			pnts{point{0, 0, ""}, point{2, 4, ""}, point{-5, 5, ""}},
			pnts{point{2, 4, ""}, point{-5, 5, ""}, point{0, 0, ""}},
			1, 0,
		},
		{ //rotate to ix 2
			pnts{point{0, 0, ""}, point{2, 4, ""}, point{-5, 5, ""}},
			pnts{point{-5, 5, ""}, point{0, 0, ""}, point{2, 4, ""}},
			2, 0,
		},
		{ //rotate to ix 4
			pnts{point{0, 0, ""}, point{2, 4, ""}, point{-5, 5, ""}},
			pnts{point{2, 4, ""}, point{-5, 5, ""}, point{0, 0, ""}},
			4, 0,
		},
	}
	comp := func(v pnts, c pnts) bool {
		if len(v) != len(c) {
			return true
		}
		for i := 0; i < len(v); i++ {
			if v[i] != c[i] {
				return true
			}
		}
		return false
	}

	for _, tst := range cases {
		val := tst.st.rot(tst.i)
		val2 := make(pnts, len(tst.st))
		copy(val2, tst.st)
		val2.rotIn(tst.i)

		if comp(val, tst.ed) {
			t.Errorf("rot(%d) for %v, expected %v received %v", tst.i, tst.st, tst.ed, val)
		}
		if comp(val2, tst.ed) {
			t.Errorf("rotIn(%d) for %v, expected %v received %v", tst.i, tst.st, tst.ed, val2)
		}
	}

}

// test nextPerm
func TestNextPerm(t *testing.T) {
	// next case
	var cases = [][]int{
		{0, 1, 2, 3, 4, 5}, {0, 1, 2, 3, 5, 4},
		{5, 4, 3, 2, 1, 0}, {},
	}
	comp := func(v []int, c []int) bool {
		if len(v) != len(c) {
			return true
		}
		for i := 0; i < len(v); i++ {
			if v[i] != c[i] {
				return true
			}
		}
		return false
	}
	for i := 0; i < len(cases); i += 2 {
		val := nextPerm(cases[i])
		if comp(val, cases[i+1]) {
			t.Errorf("nextPerm(%v), expected %v received %v", cases[i], cases[i+1], val)
		}
	}

	// perm count case
	cases = [][]int{
		{0, 1}, {2},
		{0, 1, 2}, {6},
		{0, 1, 2, 3}, {24},
		{0, 1, 2, 3, 4}, {120},
		{0, 1, 2, 3, 4, 5}, {720},
	}
	perms := func(lst []int) int {
		cnt := 0
		for n := lst; n != nil; n = nextPerm(n) {
			cnt++
		}
		return cnt
	}
	for i := 0; i < len(cases); i += 2 {
		val := perms(cases[i])
		if val != cases[i+1][0] {
			t.Errorf("nextPermCnt(%v), expected %d received %d", cases[i], cases[i+1][0], val)
		}
	}

}

// test anchToPnt
type ancType struct {
	passed string
	res    point
	errRes bool
}

func TestAnchToPnt(t *testing.T) {
	var cases = []ancType{
		{"47.782816,-122.343771", point{47.782816, -122.343771, "anchor"}, false},
		{"47.782816, -122.343771", point{47.782816, -122.343771, "anchor"}, false},
		{" 47.782816,  -122.343771 ", point{47.782816, -122.343771, "anchor"}, false},
		{"47.782816, -422.343771", point{}, true},
	}

	for _, tst := range cases {
		val, err := anchToPnt(tst.passed)
		if val != tst.res || (err == nil && tst.errRes) || (err != nil && !tst.errRes) {
			t.Errorf("anchToPnt(%q), expected %v and err=%t received %v and %v", tst.passed, tst.res, tst.errRes, val, err)
		}
	}

}

type centerPnt struct {
	vals pnts
	res  point
}

func TestCentPnt(t *testing.T) {
	var cases = []centerPnt{
		{pnts{point{0, 0, ""}, point{0, 1, ""}, point{1, 0, ""}}, point{0.211341, 0.211342, ""}},
		{pnts{point{0, 0, ""}, point{0, 0, ""}, point{0, 0, ""}}, point{0, 0, ""}},
		{pnts{point{48.775831, -122.444898, ""}, point{47.441353, -122.200362, ""}, point{48, -122, ""}}, point{48, -122.000001, ""}},
	}

	for _, tst := range cases {
		val, _ := tst.vals.centPnt()
		latDiff := math.Abs(val.lat - tst.res.lat)
		lonDiff := math.Abs(val.lon - tst.res.lon)
		if latDiff > floatErrorMax || lonDiff > floatErrorMax {
			t.Errorf("centerPnt of %v, expected %v, received %v", tst.vals, tst.res, val)
		}

	}

}

func BenchmarkCenter(b *testing.B) {

	// load file
	p, err := readFile(dat1)
	if err != nil {
		fmt.Printf("error loading file: %v\n", err)
		return
	}

	for n := 0; n < b.N; n++ {
		p.centPnt()
	}
}

func TestShuffle(t *testing.T) {
	var cases = []pnts{
		{point{0, 0, ""},
			point{0, 1, ""},
			point{1, 0, ""},
			point{2, 2, ""},
			point{4, 4, ""},
			point{5, 5, ""},
			point{7, 7, ""},
		},
	}

	for _, tst := range cases {
		val := tst.shuffle()
		if len(val) != len(tst) {
			t.Errorf("for %v, expected length %d and got %d\n", tst, len(tst), len(val))
		}

	}

}

func TestKmeans(t *testing.T) {
	// load test file
	p, err := readFile(dat2)
	if err != nil {
		fmt.Printf("error loading file: %v\n", err)
		return
	}

	count := func(c []cluster) int {
		cnt := 0
		for _, v := range c {
			cnt += len(v.cls)
		}
		return cnt
	}

	val := p.kmeans(5)
	if count(val) != len(p) {
		t.Errorf("expected %d vals, and got %d", len(p), count(val))
	}

	//genClusters(val, "clusterTestPNG")

}

type pntsInt struct {
	vals pnts
	cnt  int
}

func TestCtrAgg(t *testing.T) {
	var cases = []pntsInt{
		{pnts{point{0, 0, "A"}}, 1},
		{pnts{point{0, 0, "A"}, point{0, 1, "B"}, point{1, 0, "B"}, point{0.211341, 0.211342, "A"}}, 2},
		{pnts{point{0, 0, "A"}, point{0, 0, "A"}, point{0, 0, "B"}, point{0, 0, "B"}}, 2},
		{pnts{point{48.775831, -122.444898, "A"}, point{47.441353, -122.200362, "B"}, point{48, -122, "C"}, point{48, -122.000001, "D"}}, 4},
	}

	for _, tst := range cases {
		val := tst.vals.ctrAgg()
		if len(val) != tst.cnt {
			t.Errorf("ctrAgg of %v, expected cnt %v, received %v", tst.vals, tst.cnt, len(val))
		}

	}

}

// test optSwap
// test nna
// test nnaMul
// test exh
// test opt2SA
// test saProb
// test fastDist
