package main

import (
	"image/color"
	"math/rand"
	"os"
	"strconv"

	"github.com/flopp/go-staticmaps"
	"github.com/fogleman/gg"
	"github.com/golang/geo/s2"
)

// plot points with highlighted center
func genPoints(p pnts, c point, nm string) error {
	ctx := sm.NewContext()
	ctx.SetSize(800, 600)
	for _, loc := range p {
		ctx.AddMarker(sm.NewMarker(s2.LatLngFromDegrees(loc.lat, loc.lon), color.RGBA{255, 51, 51, 0xff}, 10.0))
	}
	ctx.AddMarker(sm.NewMarker(s2.LatLngFromDegrees(c.lat, c.lon), color.RGBA{10, 10, 255, 0xff}, 12.0))

	img, err := ctx.Render()
	if err != nil {
		return err
	}

	if _, err := os.Stat(nm + ".png"); !os.IsNotExist(err) {
		os.Remove(nm + ".png")
	}

	if err := gg.SavePNG(nm+".png", img); err != nil {
		return err
	}

	return nil
}

// build route from ordered points
func genRoute(p pnts, c point, nm string) error {
	ctx := sm.NewContext()
	ctx.SetSize(800, 600)
	for _, loc := range p {
		ctx.AddMarker(sm.NewMarker(s2.LatLngFromDegrees(loc.lat, loc.lon), color.RGBA{255, 51, 51, 0xff}, 10.0))
	}
	ctx.AddMarker(sm.NewMarker(s2.LatLngFromDegrees(c.lat, c.lon), color.RGBA{10, 10, 255, 0xff}, 12.0))

	path := make([]s2.LatLng, len(p))
	for i, loc := range p {
		path[i] = s2.LatLngFromDegrees(loc.lat, loc.lon)
	}

	ctx.AddPath(sm.NewPath(path, color.RGBA{155, 51, 255, 0xff}, 3.0))

	img, err := ctx.Render()
	if err != nil {
		return err
	}

	if _, err := os.Stat(nm + ".png"); !os.IsNotExist(err) {
		os.Remove(nm + ".png")
	}

	if err := gg.SavePNG(nm+".png", img); err != nil {
		return err
	}

	return nil
}

// plot clusters with highlighted center
func genClusters(c []cluster, nm string) error {

	ctx := sm.NewContext()
	ctx.SetSize(800, 600)
	mkr := color.RGBA{10, 10, 255, 0xff}

	//contruct palette
	del := len(palette) - len(c)
	newPal := make([]color.RGBA, len(palette))
	if del < 0 {
		copy(newPal, palette)

		for len(newPal) < len(c) {
			rnd := randColor()
			if !inPal(rnd) && rnd != mkr {
				newPal = append(newPal, rnd)
			}
		}
	} else {
		newPal = shufflePal(palette)
		newPal = newPal[:len(c)]
	}

	for i, cls := range c {
		clr := newPal[i]
		for _, loc := range cls.cls {
			ctx.AddMarker(sm.NewMarker(s2.LatLngFromDegrees(loc.lat, loc.lon), clr, 10.0))
		}
		ctrPoint := sm.NewMarker(s2.LatLngFromDegrees(cls.ctr.lat, cls.ctr.lon), mkr, 12.0)
		ctrPoint.Label = strconv.Itoa(i)
		ctrPoint.LabelColor = color.RGBA{0xfe, 0xfe, 0xfa, 0xff}
		ctx.AddMarker(ctrPoint)
	}

	img, err := ctx.Render()
	if err != nil {
		return err
	}

	if _, err := os.Stat(nm + ".png"); !os.IsNotExist(err) {
		os.Remove(nm + ".png")
	}

	if err := gg.SavePNG(nm+".png", img); err != nil {
		return err
	}

	return nil
}

func randColor() color.RGBA {
	return color.RGBA{
		uint8(rand.Intn(256)),
		uint8(rand.Intn(256)),
		uint8(rand.Intn(256)),
		0xff,
	}
}

func inPal(c color.RGBA) bool {
	for _, v := range palette {
		if v == c {
			return true
		}
	}
	return false
}

var palette = []color.RGBA{
	color.RGBA{0xfd, 0xe, 0x35, 0xff},
	color.RGBA{0xff, 0x60, 0x37, 0xff},
	color.RGBA{0xff, 0x35, 0x5e, 0xff},
	color.RGBA{0xff, 0x99, 0x66, 0xff},
	color.RGBA{0xff, 0xcc, 0x33, 0xff},
	color.RGBA{0xcc, 0xff, 0x0, 0xff},
	color.RGBA{0x66, 0xff, 0x66, 0xff},
	color.RGBA{0xee, 0x34, 0xd2, 0xff},
	color.RGBA{0xff, 0xff, 0x66, 0xff},
	color.RGBA{0x9c, 0x51, 0xb6, 0xff},
}

func shufflePal(p []color.RGBA) []color.RGBA {
	ret := make([]color.RGBA, len(p))
	perm := rand.Perm(len(p))
	for i, ri := range perm {
		ret[i] = (p)[ri]
	}
	return ret
}
