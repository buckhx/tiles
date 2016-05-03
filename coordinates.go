package tiles

import (
	"fmt"
	"math"
)

// WGS-84 Lat Lon coordindates in degrees
type Coords struct {
	Lat, Lon float64
}

func (c Coords) Equals(that Coords) bool {
	eq := floatEquals(c.Lat, that.Lat)
	eq = eq && floatEquals(c.Lon, that.Lon)
	return eq
}

// Get the Pixel of the pixel at the zoom level
func (c Coords) ToPixel(zoom uint) Pixel {
	x := (c.Lon + 180) / 360.0
	sinLat := math.Sin(c.Lat * math.Pi / 180.0)
	y := 0.5 - math.Log((1+sinLat)/(1-sinLat))/(4*math.Pi)
	size := float64(mapDimensions(zoom))
	return Pixel{
		X: uint(clip(x*size+0.5, 0, size-1)),
		Y: uint(clip(y*size+0.5, 0, size-1)),
		Z: zoom,
	}

}

func (c Coords) String() string {
	return fmt.Sprintf("(%v, %v)", c.Lat, c.Lon)
}

// Coords that have been clipped to Max/Min Lat/Lon
func ClippedCoords(lat, lon float64) Coords {
	return Coords{
		Lat: clip(lat, MinLat, MaxLat),
		Lon: clip(lon, MinLon, MaxLon),
	}
}
