package tiles

import (
	"fmt"
	"math"
)

// Coordinate is a simple struct for hold WGS-84 Lat Lon coordinates in degrees
type Coordinate struct {
	Lat, Lon float64
}

// Equals checks if these coords are equal avoiding some float precision
func (c Coordinate) Equals(that Coordinate) bool {
	eq := floatEquals(c.Lat, that.Lat)
	eq = eq && floatEquals(c.Lon, that.Lon)
	return eq
}

// ToPixel gets the Pixel of the coord at the zoom level
func (c Coordinate) ToPixel(zoom int) Pixel {
	x := (c.Lon + 180) / 360.0
	sinLat := math.Sin(c.Lat * math.Pi / 180.0)
	y := 0.5 - math.Log((1+sinLat)/(1-sinLat))/(4*math.Pi)
	size := float64(mapDimensions(zoom))
	return Pixel{
		X: int(clip(x*size+0.5, 0, size-1)),
		Y: int(clip(y*size+0.5, 0, size-1)),
		Z: zoom,
	}

}

func (c Coordinate) String() string {
	return fmt.Sprintf("(%v, %v)", c.Lat, c.Lon)
}

// ClippedCoords that have been clipped to Max/Min Lat/Lon
// This can be used as a constructor to assert bad values will be clipped
func ClippedCoords(lat, lon float64) Coordinate {
	return Coordinate{
		Lat: clip(lat, MinLat, MaxLat),
		Lon: clip(lon, MinLon, MaxLon),
	}
}
