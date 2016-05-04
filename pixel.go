package tiles

import (
	"math"
)

// Pixel in a WGS84 Mercator map projection with a NW origin (0,0) of the projection
type Pixel struct {
	X, Y, Z int
}

func (p Pixel) floatX() float64 {
	return float64(p.X)
}

func (p Pixel) floatY() float64 {
	return float64(p.Y)
}

// ToCoords converts to WGS84 coordaintes
func (p Pixel) ToCoords() Coords {
	size := float64(mapDimensions(p.Z))
	x := (clip(p.floatX(), 0, size-1) / size) - 0.5
	y := 0.5 - (clip(p.floatY(), 0, size-1) / size)
	lat := 90 - 360*math.Atan(math.Exp(-y*2*math.Pi))/math.Pi
	lon := 360.0 * x
	return ClippedCoords(lat, lon)
}

// ToTile gets the tile that contains this pixel as well as the offset pixel within that tile.
func (p Pixel) ToTile() (tile Tile, offset TilePixel) {
	tile = Tile{
		X: p.X / TileSize,
		Y: p.Y / TileSize,
		Z: p.Z,
	}
	offset = TilePixel{
		X:    p.X % TileSize,
		Y:    p.Y % TileSize,
		Tile: &tile,
	}
	return
}

// TilePixel is a pixel whose origin (0,0) is NW corner of Tile referenced in to tile field
type TilePixel struct {
	X, Y int
	Tile *Tile
}

func (p TilePixel) toCoords() Coords {
	panic("TilePixel.ToCoords() Not Implemented")
	//return Coords{}
}
