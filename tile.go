// Package tiles is a collection of conversion utilities to go between geo/pixel/tile/quadkey space
// This package uses WGS84 coordinates and a mercator projection
// There is also a TileIndex which can be used to store data in a single place and aggregate when needed
package tiles

import (
	"bytes"
	"errors"
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

// Quadkey returns the string representation of a Bing Maps quadkey. See more https://msdn.microsoft.com/en-us/library/bb259689.aspx
// Panics if the tile is invalid or if it can't write to the internal buffer
func (t Tile) Quadkey() Quadkey {
	var qk bytes.Buffer
	for i := t.Z; i > 0; i-- {
		q := 0
		m := 1 << uint(i-1)
		if (t.X & m) != 0 {
			q++
		}
		if (t.Y & m) != 0 {
			q += 2
		}
		// strconv.Itoa(q) was the bottleneck
		var d byte
		switch q {
		case 0:
			d = '0'
		case 1:
			d = '1'
		case 2:
			d = '2'
		case 3:
			d = '3'
		default:
			panic("Invalid tile.Quadkey()")
		}
		_ = qk.WriteByte(d)
	}
	return Quadkey(qk.String())
}

// FromQuadkeyString returns a tile that represents the given quadkey string. Returns an error if quadkey string is invalid.
func FromQuadkeyString(qk string) (tile Tile, err error) {
	tile.Z = len(qk)
	for i := tile.Z; i > 0; i-- {
		mask := 1 << uint(i-1)
		c := len(qk) - i
		q := qk[c]
		switch q {
		case '0':
			break
		case '1':
			tile.X |= mask
		case '2':
			tile.Y |= mask
		case '3':
			tile.X |= mask
			tile.Y |= mask
		default:
			err = errors.New("Invalid Quadkey " + qk)
			tile = Tile{} // zero tile
			return
		}
	}
	return
}

// FromCoordinate take float lat/lons and a zoom and return a tile
// Clips the coordinates if they are outside of Min/MaxLat/Lon
func FromCoordinate(lat, lon float64, zoom int) Tile {
	c := ClippedCoords(lat, lon)
	p := c.ToPixel(zoom)
	t, _ := p.ToTile()
	return t
}
