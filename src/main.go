package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"math"
	"math/big"
	"os"
	"path/filepath"
	"regexp"
	"runtime/pprof"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/golang/geo/s2"
)

// flags
var (
	inFile     = flag.String("f", "in.txt", "source file name")
	outFile    = flag.String("o", "out.txt", "outfile name")
	img        = flag.Bool("t", true, "produce route image")
	rate       = flag.Float64("r", 0.8, "decay value for SA proc")
	start      = flag.Int("s", 0, "index to rotate result to")
	meth       = flag.String("m", "auto", "opt method to use")
	anchor     = flag.String("a", "", "pass anchor coords for rotation")
	clusters   = flag.Int("cls", 0, "perform k-means clustering")
	format     = flag.Bool("fmt", true, "format output with headers and order")
	centers    = flag.Bool("ctr", false, "process centroids not locations")
	cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
)

// available methods in order of quality (0 is best)
var methOPT = map[int]string{
	-1: "auto",
	0:  "exh",
	1:  "opt",
	2:  "resOpt",
	3:  "bigOpt",
	4:  "nnMul",
	5:  "nn",
	6:  "none",
}

func main() {

	flag.Parse()

	// check method flag
	if !inOpt(*meth) {
		fmt.Printf("%q is not a valid method\n", *meth)
		fmt.Printf("valid methods: %s\n", dispMETH())
		return
	}

	// profiling start
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			panic(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	// grab current directory
	dir, err := os.Getwd()
	if err != nil {
		fmt.Printf("error getting current directory: %v\n", err)
	}

	// load and check file
	fmt.Println("loading file...")
	p, err := readFile(filepath.Join(dir, *inFile))
	if err != nil {
		fmt.Printf("error loading file: %v\n", err)
	}

	if len(p) > 0 {
		fmt.Printf("%v records loaded with starting tour length of: %.4f km\n", len(p), p.tourLen())
	} else {
		fmt.Println("empty points, quiting")
		return
	}

	// clustering interupt
	if *clusters > 0 {

		if *clusters >= len(p) {
			rVal := int(math.Pow(float64(len(p)), 0.33))
			fmt.Printf("error, %d nodes and asked for %d clusters\n", len(p), *clusters)
			fmt.Printf("choose a number < %d (recommended:%d)\n", len(p), rVal)
			return
		}

		fmt.Printf("finding %d clusters\n", *clusters)
		clsRes := p.kmeans(*clusters)
		clsPath := filepath.Join(dir, "clusters")
		if _, err := os.Stat(clsPath); os.IsNotExist(err) {
			os.Mkdir(clsPath, os.ModeDir)
		}

		// remove any output files already in directory
		if err := remFiles(clsPath); err != nil {
			fmt.Printf("error, could not remove files in clusters dir %v\n", err)
		}

		for i, v := range clsRes {
			clsCtr, clsDist := v.cls.centPnt()
			writeFile(v.cls, clsCtr, clsDist, filepath.Join(clsPath, "cls"+strconv.Itoa(i)+".txt"), *format)
		}

		fmt.Println("generating cluster map..")
		genClusters(clsRes, filepath.Join(clsPath, "clusters"))

		fmt.Println("skipping routing")
		return // don't perform routing if clustering is selected
	}

	// convert data to centroids
	if *centers {
		fmt.Println("creating centroid route")

		// aggregate to: label, <centroid>
		p = p.ctrAgg()

	}

	// process anchor flag
	if *anchor != "" {
		aPnt, err := anchToPnt(*anchor)
		if err != nil {
			fmt.Printf("error, could not parse %q: %v\n", *anchor, err)
			return
		}
		pNear, newStart := p.nearest(aPnt, false)
		*start = newStart //re-assign start flag

		fmt.Printf("using provided anchor:%v, at node:%d,%v\n", aPnt, newStart+1, pNear)
	}

	// choose method
	cnt := len(p)
	out := make(pnts, len(p))
	optDone := true

	s1 := time.Now()

	switch {
	case cnt < 4 || *meth == "none":
		fmt.Println("nothing to optimize")
		optDone = false
		out = p
	case *meth == "exh":
		quit := true
		out, quit = methodExh(p)
		if quit {
			return
		}
	case *meth == "opt":
		out = methodOpt(p, *rate, false, -1, true)
	case *meth == "resOpt":
		out = methodOpt(p.nna(*start), *rate, true, -1, true)
	case *meth == "bigOpt":
		out = methodOpt(p.nna(*start), *rate, true, 1, false)
	case *meth == "nn":
		out = methodNN(p, *start, false)
	case *meth == "nnMul":
		out = methodNN(p, *start, true)

	// auto
	case cnt < 11:
		out, _ = methodExh(p)
	case cnt <= 750: //max 7min
		out = methodOpt(p, *rate, false, -1, true)
	case cnt <= 3000: //max 8min
		out = methodOpt(p.nna(*start), *rate, true, -1, true)
	case cnt <= 10000:
		out = methodOpt(p.nna(*start), *rate, true, 1, false)
	case cnt > 10000:
		out = methodNN(p, *start, false)
	}

	if optDone {
		elap := time.Since(s1)
		fmt.Println("optimization took:", elap)
	}

	// rotate result to res[0] = start
	if *start != 0 {
		fmt.Printf("rotating result to node %d\n", *start+1)
	}

	var ix int
	for i, pnt := range out {
		if pnt == p[*start] {
			ix = i
			break
		}
	}
	out.rotIn(ix)

	if optDone {
		fmt.Printf("final tour length: %.4f km\n", out.tourLen())
	}
	fmt.Printf("writing results to %v\n", *outFile)

	// center point calc
	fmt.Printf("finding center point of data.. ")
	ctr, ctrDist := out.centPnt()
	fmt.Printf("{%.6f,%.6f}\t%.2fkm avg dist\n", ctr.lat, ctr.lon, ctrDist/float64(len(p)))

	if err := writeFile(out, ctr, ctrDist, filepath.Join(dir, *outFile), *format); err != nil {
		fmt.Printf("error writing file: %v\n", err)
		return
	}

	if *img {
		rName := (*outFile)[:strings.Index(*outFile, ".txt")]
		fmt.Println("generating route and center plot")
		if err := genRoute(out, ctr, rName+"_route"); err != nil {
			fmt.Printf("error building route: %v\n", err)
			return
		}
	}

}

// interact fucntions for methods selected in switch
func methodExh(p pnts) (pnts, bool) {

	fmt.Println("using exhaustive method")
	nodes := len(p)

	// provide warning for large sets
	if nodes > 11 {
		perms := factf(nodes)
		pPerSec := big.NewFloat(500000)
		secPerYear := big.NewFloat(float64(60 * 60 * 24 * 365))
		years := new(big.Float)
		years.Quo(perms, pPerSec).Quo(years, secPerYear)

		fmt.Printf("warning, aprox %8.4e years to calculate\n", years)
		buf := bufio.NewReader(os.Stdin)
		for {
			fmt.Print("continue (y/n): ")
			resp, err := buf.ReadBytes('\n')
			if err != nil {
				fmt.Println(err)
			} else if resp[0] == byte('y') {
				break
			} else if resp[0] == byte('n') {
				return pnts{}, true // quit state
			}
		}
	}

	return p.exh(), false
}

func methodOpt(p pnts, r float64, b bool, lim int, sa bool) pnts {
	str := "using"
	if b {
		str += " resticted (20 node)"
	}
	if lim != -1 {
		str += " " + strconv.Itoa(lim) + "-pass"
	}
	str += " 2-Opt"
	if sa {
		str += " with SA"
	} else {
		str += " without SA"
	}
	fmt.Println(str)

	return p.opt2SA(r, b, lim, sa)
}

func methodNN(p pnts, s int, m bool) pnts {

	if m {
		fmt.Println("using multi-start nearest neighbor")
		return p.nnaMul()
	}

	fmt.Printf("using nearest neighbor, starting at node: %d\n", s+1)
	return p.nna(s)

}

// file processing
func readFile(path string) (pnts, error) {

	csvFile, err := os.Open(path)
	if err != nil {
		return pnts{}, err
	}
	defer csvFile.Close()

	reader := csv.NewReader(csvFile)

	// configure csv reader
	reader.Comma = '\t' // set split token

	records, err := reader.ReadAll() // -> [][]string
	if err != nil {
		return pnts{}, err
	}
	if len(records) == 0 {
		return pnts{}, errors.New("empty input file")
	}
	records = records[1:] // drop header record

	// check for empty LatLon
	if err := checkEmp(records); err != nil {
		return pnts{}, err
	}

	p := make(pnts, len(records))

	for i, rec := range records {
		p[i].lab = strings.TrimSpace(rec[0])
		p[i].lat, err = strconv.ParseFloat(rec[1], 64)
		if err != nil {
			return pnts{}, err
		}
		p[i].lon, err = strconv.ParseFloat(rec[2], 64)
		if err != nil {
			return pnts{}, err
		}
	}

	// check records
	p, err = checkRec(p)
	if err != nil {
		return pnts{}, err
	}

	return p, nil
}

func checkEmp(recs [][]string) error {
	for i, rec := range recs {
		if strings.TrimSpace(rec[1]) == "" || strings.TrimSpace(rec[2]) == "" {
			return errors.New("missing LatLon! row: " + strconv.Itoa(i+2))
		}
	}
	return nil
}

func checkRec(recs pnts) (pnts, error) {

	// check populated
	for i, rec := range recs {
		if rec.lab == "" {
			return pnts{}, errors.New("unpopulated record! row: " + strconv.Itoa(i+2) + " column: label")
		}
		if rec.lat == 0 {
			return pnts{}, errors.New("unpopulated record! row: " + strconv.Itoa(i+2) + " column: lat")
		}
		if rec.lon == 0 {
			return pnts{}, errors.New("unpopulated record! row: " + strconv.Itoa(i+2) + " column: lon")
		}
	}

	// remove dups
	distVals := make(map[point]struct{})
	dups := []int{}

	tmp := pnts{}
	for i, rec := range recs {
		if _, ok := distVals[rec]; !ok {
			distVals[rec] = struct{}{}
			tmp = append(tmp, rec)
		} else {
			dups = append(dups, i+2)
		}
	}

	if len(recs)-len(tmp) == 1 {
		fmt.Printf("removed %d duplicate at row:%v\n", len(recs)-len(tmp), dups)
	} else if len(recs)-len(tmp) > 1 {
		fmt.Printf("removed %d duplicates at rows:%v\n", len(recs)-len(tmp), dups)
	}

	// check valid coords
	for i, rec := range tmp {
		chk := s2.LatLngFromDegrees(rec.lat, rec.lon)
		if !chk.IsValid() {
			return pnts{}, errors.New("invalid LatLng, row: " + strconv.Itoa(i+2))
		}
	}

	return tmp, nil
}

func writeFile(p pnts, c point, d float64, dest string, format bool) error {
	tour := make([][]string, len(p))

	if format {
		for i, loc := range p {
			tour[i] = []string{
				loc.lab,
				strconv.FormatFloat(loc.lat, 'f', 6, 64),
				strconv.FormatFloat(loc.lon, 'f', 6, 64),
				strconv.Itoa(i + 1),
			}
		}
		hdr := [][]string{
			{"center:",
				strconv.FormatFloat(c.lat, 'f', 6, 64),
				strconv.FormatFloat(c.lon, 'f', 6, 64),
				fmt.Sprintf("%.2f", d/float64(len(p))) + "km avg dist"},
			{},
			{"lab", "lat", "lon", "ord"},
		}
		tour = append(hdr, tour...)
	} else {
		for i, loc := range p {
			tour[i] = []string{
				loc.lab,
				strconv.FormatFloat(loc.lat, 'f', 6, 64),
				strconv.FormatFloat(loc.lon, 'f', 6, 64),
			}
		}
		hdr := [][]string{{"label", "lat", "lon"}}
		tour = append(hdr, tour...)

	}

	outFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer outFile.Close()

	outWriter := csv.NewWriter(outFile)

	// configure csv writer
	outWriter.Comma = '\t' // set split token

	outWriter.WriteAll(tour)
	if err := outWriter.Error(); err != nil {
		return err
	}

	return nil

}

// build sorted method string
func dispMETH() string {
	arrOrd := []int{}
	arrMETH := []string{}
	for v := range methOPT {
		arrOrd = append(arrOrd, v)
	}
	sort.Ints(arrOrd)
	for _, i := range arrOrd {
		arrMETH = append(arrMETH, methOPT[i])
	}
	return strings.Join(arrMETH, ", ")
}

func inOpt(inStr string) bool {
	for _, v := range methOPT {
		if inStr == v {
			return true
		}
	}
	return false
}

// convert flag string to point
func anchToPnt(inStr string) (point, error) {
	clnStr := strings.Replace(inStr, " ", "", -1)
	coords := strings.Split(clnStr, ",")

	lat, err := strconv.ParseFloat(coords[0], 64)
	if err != nil {
		return point{}, err
	}
	lon, err := strconv.ParseFloat(coords[1], 64)
	if err != nil {
		return point{}, err
	}

	chk := s2.LatLngFromDegrees(lat, lon)
	if !chk.IsValid() {
		return point{}, errors.New("invalid LatLng")
	}

	return point{lat, lon, "anchor"}, nil
}

// remove files in dir
func remFiles(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	fls, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range fls {
		isOut, _ := regexp.MatchString(`^cls\d+\.txt$`, name)
		isImg := name == "clusters.png"
		if isOut || isImg {
			err = os.Remove(filepath.Join(dir, name))
			if err != nil {
				return err
			}
		}

	}
	return nil
}

/*  ToDo
- use filepath package to make portable to other GOOS env

*/
