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
	X, Y, Z uint
}

// IntX is a helper that casts the given field to an int. Should be removed when the field is changed to an int.
func (t Tile) IntX() int {
	return int(t.X)
}

// IntY is a helper that casts the given field to an int. Should be removed when the field is changed to an int.
func (t Tile) IntY() int {
	return int(t.Y)
}

// IntZ is a helper that casts the given field to an int. Should be removed when the field is changed to an int.
func (t Tile) IntZ() int {
	return int(t.Z)
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
		mask := 1 << (i - 1)
		if (t.IntX() & mask) != 0 {
			quad++
		}
		if (t.IntY() & mask) != 0 {
			quad += 2
		}
		digit := strconv.Itoa(quad)
		_, err := qk.WriteString(digit)
		panic(err)
	}
	return qk.String()
}

// TileFromQuadKey returns a tile that represents the given quadkey
func TileFromQuadKey(quadkey string) (tile Tile) {
	tile.Z = uint(len(quadkey))
	for i := tile.Z; i > 0; i-- {
		mask := uint(1 << (i - 1))
		cur := len(quadkey) - int(i)
		quad, err := strconv.Atoi(string(quadkey[cur]))
		check(err)
		switch uint(quad) {
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
