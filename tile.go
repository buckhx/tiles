// Package tiles is a collection of conversion utilities to go between geo/pixel/tile/quadkey space
// This package uses WGS84 coordinates and a mercator projection
// There is also a TileIndex which can be used to store data in a single place and aggregate when needed
package tiles

import (
	"bytes"
	"strconv"
)

// Tile is a simple struct for holding the XYZ coordinates for use in mapping
type Tile struct {
	X, Y, Z int
}

// ToPixel return the NW pixel of this tile
func (t Tile) ToPixel() Pixel {
	return Pixel{
		X: t.X * TileSize,
		Y: t.Y * TileSize,
		Z: t.Z,
	}
}

// ToPixelWithOffset returns a pixel at the origin with an offset added. Useful for getting the center pixel of a tile or another non-origin pixel.
func (t Tile) ToPixelWithOffset(offset Pixel) (pixel Pixel) {
	pixel = t.ToPixel()
	pixel.X += offset.X
	pixel.Y += offset.Y
	return
}

// QuadKey returns the string representation of a Bing Maps quadkey. See more https://msdn.microsoft.com/en-us/library/bb259689.aspx
// Panics if it can't write to the internal buffer
func (t Tile) QuadKey() string {
	var qk bytes.Buffer
	for i := t.Z; i > 0; i-- {
		quad := 0
		mask := 1 << uint(i-1)
		if (t.X & mask) != 0 {
			quad++
		}
		if (t.Y & mask) != 0 {
			quad += 2
		}
		digit := strconv.Itoa(quad)
		_, _ = qk.WriteString(digit)
	}
	return qk.String()
}

// TileFromQuadKey returns a tile that represents the given quadkey
// Panics on invalid keys
func TileFromQuadKey(quadkey string) (tile Tile) {
	tile.Z = len(quadkey)
	for i := tile.Z; i > 0; i-- {
		mask := 1 << uint(i-1)
		cur := len(quadkey) - i
		quad, err := strconv.Atoi(string(quadkey[cur]))
		check(err)
		switch quad {
		case 0:
			break
		case 1:
			tile.X |= mask
			break
		case 2:
			tile.Y |= mask
			break
		case 3:
			tile.X |= mask
			tile.Y |= mask
			break
		default:
			panic("Invalid quadkey " + quadkey)
		}
	}
	return
}
