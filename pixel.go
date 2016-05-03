package tiles

import (
	"math"
)

// Pixel in a WGS84 Mercator map projection with a NW origin (0,0) of the projection
type Pixel struct {
	X, Y, Z uint
}

func (p Pixel) FloatX() float64 {
	return float64(p.X)
}

func (p Pixel) FloatY() float64 {
	return float64(p.Y)
}

func (p Pixel) ToCoords() Coords {
	size := float64(mapDimensions(p.Z))
	x := (clip(p.FloatX(), 0, size-1) / size) - 0.5
	y := 0.5 - (clip(p.FloatY(), 0, size-1) / size)
	lat := 90 - 360*math.Atan(math.Exp(-y*2*math.Pi))/math.Pi
	lon := 360.0 * x
	return ClippedCoords(lat, lon)
}

// Gets the tile that contains this pixel as well as the pixel within that tile.
func (p Pixel) ToTile() (tile Tile, tilePixel TilePixel) {
	tile = Tile{
		X: p.X / TileSize,
		Y: p.Y / TileSize,
		Z: p.Z,
	}
	tilePixel = TilePixel{
		X:    p.X % TileSize,
		Y:    p.Y % TileSize,
		Tile: &tile,
	}
	return
}

// Pixel whose origin (0,0) is NW corner of Tile
type TilePixel struct {
	X, Y uint
	Tile *Tile
}

func (p TilePixel) ToCoords() Coords {
	panic("TilePixel.ToCoords() Not Implemented")
	return Coords{}
}
