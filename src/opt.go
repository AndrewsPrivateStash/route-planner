package main

import (
	"math"
	"math/big"
	"math/rand"
	"sort"
)

// cartesean coordinats
type cart struct {
	x float64
	y float64
	z float64
}

// a single geo-point with label
type point struct {
	lat float64
	lon float64
	lab string
}

// degrees to radians
func (p *point) dToR() point {
	return point{
		p.lat / 180 * math.Pi,
		p.lon / 180 * math.Pi,
		p.lab,
	}
}

// polar to cartesean coords
func (p *point) polToCart() cart {
	radPnt := p.dToR()
	return cart{
		math.Cos(radPnt.lat) * math.Cos(radPnt.lon),
		math.Cos(radPnt.lat) * math.Sin(radPnt.lon),
		math.Sin(radPnt.lat),
	}
}

// a list of points
type pnts []point

// remove element from pnts (in place)
func (ps *pnts) rem(ind int) {
	*ps = append((*ps)[:ind], (*ps)[ind+1:]...) //maybe mem leak
}

// swap two nodes (return copy)
func (ps *pnts) swap(ix1 int, ix2 int) pnts {
	t := make(pnts, len(*ps))
	copy(t, *ps)
	t[ix1], t[ix2] = t[ix2], t[ix1]
	return t
}

// 2opt swap (return copy)
// [0,i) + rev[i,j] + (j,oo)
func (ps *pnts) optSwap(ix1 int, ix2 int) pnts {
	t := append(pnts{}, (*ps)[0:ix1]...)
	t = append(t, rev((*ps)[ix1:ix2+1])...)
	t = append(t, (*ps)[ix2+1:]...)
	return t
}

// rotate points (return copy)
func (ps *pnts) rot(ix int) pnts {
	ix = ix % len(*ps)
	return append((*ps)[ix:], (*ps)[0:ix]...)
}

// rotate points (in place)
func (ps *pnts) rotIn(ix int) {
	ix = ix % len(*ps)
	copy(*ps, append((*ps)[ix:], (*ps)[0:ix]...))
}

// sort points on label (return copy)
func (ps *pnts) sortLab() pnts {
	out := make(pnts, len(*ps))
	copy(out, *ps)

	sort.Slice(out, func(i, j int) bool {
		return out[i].lab < out[j].lab
	})
	return out
}

// given point, find nearest point from pnts
func (ps *pnts) nearest(p point, dup bool) (point, int) {
	min := math.MaxFloat64
	var best point
	var index int
	for i, loc := range *ps {
		if loc != p || dup { // records assumed distinct
			h := haver(loc, p)
			if h < min {
				min = h
				best = loc
				index = i
			}
		}
	}
	return best, index
}

// length of tour (km)
func (ps *pnts) tourLen() float64 {
	var tourDist float64
	for i := 0; i < len(*ps)-1; i++ {
		tourDist += haver((*ps)[i], (*ps)[i+1])
	}
	tourDist += haver((*ps)[len(*ps)-1], (*ps)[0])

	return tourDist
}

// ordered length of tour (km)
func (ps *pnts) oTourLen(ord []int) float64 {
	var tourDist float64
	for i := 0; i < len(*ps)-1; i++ {
		tourDist += haver((*ps)[ord[i]], (*ps)[ord[i+1]])
	}
	tourDist += haver((*ps)[ord[len(ord)-1]], (*ps)[ord[0]])

	return tourDist
}

// central point
//http://www.geomidpoint.com/calculation.html
type pntDist struct {
	p    point
	dist float64
}

func (ps *pnts) centPnt() (point, float64) {
	const floatErrorMax = 1e-6
	var xb, yb, zb float64

	for _, loc := range *ps {
		tmp := loc.polToCart()
		xb += tmp.x
		yb += tmp.y
		zb += tmp.z
	}
	xb /= float64(len(*ps))
	yb /= float64(len(*ps))
	zb /= float64(len(*ps))

	hyp := math.Sqrt(xb*xb + yb*yb)
	latOut := math.Atan2(zb, hyp) * 180 / math.Pi
	lonOut := math.Atan2(yb, xb) * 180 / math.Pi
	ctr := point{latOut, lonOut, ""}

	// iterative improvement
	minDist := func(c point, p pnts) float64 {
		var sum float64
		for _, loc := range p {
			sum += haver(c, loc)
		}
		return sum
	}(ctr, *ps)

	sumDistCh := func(c point, p pnts, ch chan<- pntDist, d chan<- bool) {
		var sum float64
		for _, loc := range p {
			sum += haver(c, loc)
		}
		ch <- pntDist{c, sum}
		d <- true

	}

	incr := 1e-1 // starting offset; about 10km
	var adjVals = [3]float64{-incr, 0, incr}
	var swapped bool

	for incr > floatErrorMax {
		cells := make(chan pntDist, 8) // return channel for routines
		done := make(chan bool)        // track individual completion of routines

		// build 8 alternate points in radius and calc dist sum
		for i, dlat := range adjVals {
			for j, dlon := range adjVals {
				if i == 0 && j == 0 {
					continue
				}
				tmp := point{ctr.lat + dlat, ctr.lon + dlon, ""}
				go sumDistCh(tmp, *ps, cells, done)
			}
		}

		// swap current best if any of the 8 are better
		swapped = false
		procCnt, readCnt := 0, 0
		for procCnt < 8 || readCnt < 8 {
			select {
			case <-done:
				procCnt++
				if procCnt == 8 {
					close(cells)
				}

			case v, ok := <-cells:
				readCnt++
				if v.dist < minDist && ok {
					minDist, ctr = v.dist, v.p
					swapped = true
				}
			}

		}
		close(done)

		// reduce radius if no improvement found
		if !swapped {
			incr /= 2
			adjVals[0], adjVals[1], adjVals[2] = -incr, 0, incr
		}
	}

	return ctr, minDist
}

// aggreagte centers accross common labels
// used for routing groups versus locations
func (ps *pnts) ctrAgg() pnts {

	// sort on label
	sorted := ps.sortLab()

	// store centers in new slice
	var out pnts
	prev := 0
	for i := 0; i < len(sorted); i++ {
		if i+1 == len(sorted) || sorted[i].lab != sorted[i+1].lab {
			tmp := make(pnts, (i-prev)+1)

			if i+1 == len(sorted) {
				copy(tmp, sorted[prev:i+1])
			} else {
				copy(tmp, sorted[prev:])
			}

			if len(tmp) == 1 {
				tmp[0].lab = sorted[i].lab
				out = append(out, tmp[0])
			} else {
				ctr, _ := tmp.centPnt()
				ctr.lab = sorted[i].lab
				out = append(out, ctr)
			}

			prev = i + 1
		}

	}

	return out
}

// nearest neighbor algorithm
func (ps *pnts) nna(start int) pnts {

	// ordered out-slice
	ordPnts := make(pnts, 1, len(*ps))
	ordPnts[0] = (*ps)[start]

	// copy of in-slice as container
	q := make(pnts, len(*ps))
	copy(q, *ps)
	q.rem(start) // remove starting point

	cnt := len(q)
	for i := 0; i < cnt; i++ {
		p, ix := q.nearest(ordPnts[len(ordPnts)-1], false)
		ordPnts = append(ordPnts, p)
		q.rem(ix)

	}

	return ordPnts

}

// nearest neighbor multi-start (try all starting nodes)
func (ps *pnts) nnaMul() pnts {

	res := make(pnts, len(*ps))
	copy(res, *ps)
	min := ps.tourLen()

	for i := 0; i < len(*ps); i++ {
		iter := ps.nna(i)
		tl := iter.tourLen()
		if tl < min {
			res = iter
			min = tl
		}
	}

	return res
}

// exhaustive search ## don't use > 11 nodes! ##
func (ps *pnts) exh() pnts {
	bestOrd := make([]int, len(*ps))
	outPnts := make(pnts, len(*ps))
	minTour := (*ps).tourLen()
	ord := make([]int, len(*ps))
	for i := 0; i < len(*ps); i++ {
		ord[i] = i
	}

	copy(bestOrd, ord)
	for n := ord; n != nil; n = nextPerm(n) {
		t := ps.oTourLen(n)
		if t < minTour {
			minTour = t
			bestOrd = n
		}
	}

	for i := 0; i < len(*ps); i++ {
		outPnts[i] = (*ps)[bestOrd[i]]
	}

	return outPnts
}

// 2-opt
// https://en.wikipedia.org/wiki/2-opt
func (ps *pnts) opt2SA(rate float64, big bool, lim int, sa bool) pnts {

	// set starting values
	stTour := append(pnts{}, *ps...) // use given order to start
	stLen := stTour.tourLen()
	bestLen := stLen
	bestTour := append(pnts{}, stTour...)

	// rand.Seed(42) // for SA RNG

	curTemp := 100.0 // no need to param

	var iters, cnt, cnt2 int
	upd := true

	// main loop; do until no swaps can be made
	for upd {

		if lim != -1 { // allow for single pass runs
			if iters >= lim {
				break
			}
		}

		upd = false
		for i := 0; i < len(stTour)-2; i++ {
			for j := i + 2; j < len(stTour); j++ { // +2 to skip connected nodes

				if big {
					if j-i > 23 { //restrict search to 20 nodes forward
						continue
					}
				}

				// perform swap
				tmp := bestTour.optSwap(i, j)
				len := tmp.tourLen()

				if len < bestLen {
					bestLen = len
					bestTour = tmp
					cnt++
					upd = true
				} else if curTemp > 1 && sa { // SA effect
					if saProb(bestLen, len, curTemp) > rand.Float64() {
						bestLen = len
						bestTour = tmp
						cnt2++
						curTemp *= rate
						upd = true
					}
				}
			}

		}
		iters++
	}
	return bestTour
}

// SA Probabilty function
func saProb(old float64, new float64, temp float64) float64 {
	return math.Exp((old - new) / temp)
}

// haversine dist function
// https://en.wikipedia.org/wiki/Haversine_formula
func haver(pnt1, pnt2 point) float64 {
	p1 := pnt1.dToR()
	p2 := pnt2.dToR()
	const R = 6378.1 //earth equatorial radius (km)
	return 2 * R * math.Asin(math.Sqrt(
		sqr(math.Sin((p2.lat-p1.lat)/2))+
			math.Cos(p1.lat)*math.Cos(p2.lat)*
				sqr(math.Sin((p2.lon-p1.lon)/2))))

}

// a faster distance estimate for iteration (relative)
// https://math.stackexchange.com/questions/29157/how-do-i-convert-the-distance-between-two-lat-long-points-into-feet-meters
func fastDist(pnt1, pnt2 point) float64 {
	avgLat := ((pnt1.lat + pnt2.lat) / 2.0) / 180 * math.Pi
	dlat := pnt1.lat - pnt2.lat
	dlon := (pnt1.lon - pnt2.lon) * math.Cos(avgLat)
	return dlat*dlat + dlon*dlon
}

// find next permutation in lexographical order
// https://en.wikipedia.org/wiki/Permutation#Generation_in_lexicographic_order
func nextPerm(arr []int) []int {

	a := make([]int, len(arr))
	copy(a, arr)

	// Find the largest index k such that a[k] < a[k + 1]
	maxK := -1
	for i := 0; i < len(a)-1; i++ {
		if a[i] < a[i+1] {
			maxK = i
		}
	}
	if maxK == -1 { // no more permutations
		return nil
	}

	// Find the largest index l greater than k such that a[k] < a[l]
	maxL := -1
	for i := maxK + 1; i < len(a); i++ {
		if a[maxK] < a[i] {
			maxL = i
		}
	}

	// Swap the value of a[k] with that of a[l]
	a[maxK], a[maxL] = a[maxL], a[maxK]

	// Reverse the sequence from a[k + 1] up to and including the final element a[n]
	revTmp := a[maxK+1:] //still points to 'a', reverse in place
	for left, right := 0, len(revTmp)-1; left < right; left, right = left+1, right-1 {
		revTmp[left], revTmp[right] = revTmp[right], revTmp[left]
	}

	return a
}

// reverse point elements in range (return copy)
func rev(s pnts) pnts {
	t := make(pnts, len(s))
	copy(t, s)
	for left, right := 0, len(t)-1; left < right; left, right = left+1, right-1 {
		t[left], t[right] = t[right], t[left]
	}
	return t
}

// arbitrary factorial
func fact(n int) *big.Int {
	if n == 0 {
		return big.NewInt(int64(1))
	}

	val, counter := big.NewInt(int64(n)), big.NewInt(int64(n))
	dec := big.NewInt(int64(-1))

	for i := n; i > 1; i-- {
		counter.Add(counter, dec)
		val.Mul(val, counter)
	}
	return val
}

func factf(n int) *big.Float {
	if n == 0 {
		return big.NewFloat(float64(1))
	}
	val, counter := big.NewFloat(float64(n)), big.NewFloat(float64(n))
	dec := big.NewFloat(float64(-1))

	for i := n; i > 1; i-- {
		counter.Add(counter, dec)
		val.Mul(val, counter)
	}
	return val
}

func myP(x float64, y int) float64 {
	var val = x
	for i := 1; i < y; i++ {
		val *= x
	}
	return val
}

func sqr(x float64) float64 {
	return x * x
}

/* ToDo
- explore concurent structures for speed {exh and 2opt}  (involved)
- remove root from haver for speed
- implement spcialized math functions for {sin, cos}, to avoid switching done in math package
- build distance matrix [][]float64 to speed up haver

*/
