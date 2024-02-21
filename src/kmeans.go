package main

import (
	"fmt"
	"math/rand"
)

type cluster struct {
	ctr point
	cls pnts
}

// kmeans alg
// take list of points and partition into n clustered pnts
// https://en.wikipedia.org/wiki/K-means_clustering
func (ps *pnts) kmeans(cls int) []cluster {

	clsOut := make([]cluster, cls)

	// pick n random cetroids from pnts
	meanCtrs := ps.rndPoints(cls)
	for i := 0; i < cls; i++ {
		clsOut[i].ctr = meanCtrs[i]
	}

	loopCnt := 0
	for loopCnt < 100 {

		// assign each point to nearest centroid to create cluster
		clsOut = ps.asgnCtr(clsOut)

		// calculate new center of cluster
		oldCtrs := getCtrs(clsOut)
		for i := 0; i < cls; i++ {
			clsOut[i].ctr, _ = clsOut[i].cls.centPnt()
		}

		// loop until no center shift
		if compCtrs(oldCtrs, getCtrs(clsOut)) {
			return clsOut
		}
		loopCnt++
	}

	fmt.Println("clustering did not converge..")
	return clsOut

}

// pick random pnts
func (ps *pnts) rndPoints(n int) pnts {
	return ps.shuffle()[:n]
}

// shuffle a slice of pnts
func (ps *pnts) shuffle() pnts {
	ret := make(pnts, len(*ps))
	perm := rand.Perm(len(*ps))
	for i, ri := range perm {
		ret[i] = (*ps)[ri]
	}
	return ret
}

// assign points to nearest mean
func (ps *pnts) asgnCtr(c []cluster) []cluster {
	outCls := make([]cluster, len(c))
	ctrs := make(pnts, len(c))
	for i, v := range c {
		ctrs[i] = v.ctr
	}

	for _, v := range *ps {
		_, ix := ctrs.nearest(v, true)
		outCls[ix].ctr = ctrs[ix]
		outCls[ix].cls = append(outCls[ix].cls, v)
	}
	return outCls

}

func getCtrs(c []cluster) pnts {
	out := make(pnts, len(c))
	for i, v := range c {
		out[i] = v.ctr
	}
	return out
}

func compCtrs(preCtrs pnts, curCtrs pnts) bool {
	for i, v := range preCtrs {
		if v != curCtrs[i] {
			return false
		}
	}
	return true
}
