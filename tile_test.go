package tiles

import (
	"fmt"
	"testing"
)

func TestTileToPixel(t *testing.T) {
	tileTests := []struct {
		tile  Tile
		pixel Pixel
	}{
		{Tile{26, 48, 7}, Pixel{6656, 12288, 7}},
		//{Tile{26, 48, 7}, Pixel{6827, 12405, 7}},
	}
	errf := "Tile%+v: %+v -> %+v"
	for _, test := range tileTests {
		pixel := test.tile.ToPixel()
		if pixel != test.pixel {
			t.Errorf(errf, test.tile, test.pixel, pixel)
		}
	}
}

func TestTileToQuadkey(t *testing.T) {
	tileTests := []struct {
		tile    Tile
		quadkey string
	}{
		{Tile{26, 48, 7}, "0231010"},
	}
	errf := "Tile%+v: %+v -> %+v"
	for _, test := range tileTests {
		qk := test.tile.QuadKey()
		if qk != test.quadkey {
			t.Errorf(errf, test.tile, test.quadkey, qk)
		}
	}
}

func TestTileFromQuadkey(t *testing.T) {
	tileTests := []struct {
		quadkey string
		tile    Tile
	}{
		{"0231010", Tile{26, 48, 7}},
	}
	errf := "QuadKey%+v: %+v -> %+v"
	for _, test := range tileTests {
		tile := TileFromQuadKey(test.quadkey)
		if tile != test.tile {
			t.Errorf(errf, test.quadkey, test.tile, tile)
		}
	}
}

func ExampleCoordinateToTile() {
	esbLat := 40.7484
	esbLon := -73.9857
	tile, _ := CoordinateToTile(esbLat, esbLon, 18)
	fmt.Println(tile)
}
