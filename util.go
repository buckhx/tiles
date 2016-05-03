package tiles

import (
	"fmt"
	"math"
)

// Earth Parameters
const (
	MinLat       float64 = -85.05112878
	MaxLat       float64 = 85.05112878
	MinLon       float64 = -180
	MaxLon       float64 = 180
	EarthRadiusM float64 = 6378137
)

// TileSize is the size in pixels of each tile. It can be tuned at the package level.
var TileSize uint = 256

// if val is outside of min-max range, clip it to min or max
// panic if min > max
func clip(val, min, max float64) float64 {
	if min > max {
		panic(fmt.Errorf("clip min %f > max %f", min, max))
	}
	return math.Min(math.Max(val, min), max)
}

// Gets the size of the x, y dimensions in pixels at the given zoom level
// x == y since the map is a square
func mapDimensions(zoom uint) uint {
	//TODO panic outside of zoom bounds
	return TileSize << zoom
}

// Gets the ground resoultion (meters/pixel) of the map at the lat and zoom
// TODO remove if unused
/*
func grndRes(lat float64, zoom uint) float64 {
	lat = clip(lat, MinLat, MaxLat)
	dim := float64(mapDimensions(zoom))
	return math.Cos(lat*math.Pi/180) * 2 * math.Pi * EarthRadiusM / dim
}
*/

// Gets the map scale at the lat, zoom & screen DPI expressed as the denominator N of the ratio 1 : N.
// TODO remove if unused
/*
func mapScale(lat float64, zoom, dpi uint) float64 {
	d := float64(dpi)
	return grndRes(lat, zoom) * d / 0.0254
}
*/

// method for approx float equality
func floatEquals(a, b float64) bool {
	var EPSILON = 0.00000001
	if (a-b) < EPSILON && (b-a) < EPSILON {
		return true
	}
	return false
}

func check(errs ...error) {
	for _, err := range errs {
		if err != nil {
			panic(err)
		}
	}
}
