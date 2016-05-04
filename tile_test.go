package tiles_test

import (
	"fmt"
	"testing"

	"github.com/buckhx/tiles"
)

func TestTileToPixel(t *testing.T) {
	tileTests := []struct {
		tile  tiles.Tile
		pixel tiles.Pixel
	}{
		{tiles.Tile{X: 26, Y: 48, Z: 7}, tiles.Pixel{X: 6656, Y: 12288, Z: 7}},
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
		tile    tiles.Tile
		quadkey tiles.Quadkey
	}{
		{tiles.Tile{X: 26, Y: 48, Z: 7}, tiles.Quadkey("0231010")},
	}
	errf := "Tile%+v: %+v -> %+v"
	for _, test := range tileTests {
		qk := test.tile.Quadkey()
		if qk != test.quadkey {
			t.Errorf(errf, test.tile, test.quadkey, qk)
		}
	}
}

func TestTileFromQuadkey(t *testing.T) {
	tileTests := []struct {
		quadkey tiles.Quadkey
		tile    tiles.Tile
	}{
		{tiles.Quadkey("0231010"), tiles.Tile{X: 26, Y: 48, Z: 7}},
	}
	errf := "QuadKey%+v: %+v -> %+v"
	for _, test := range tileTests {
		tile := test.quadkey.ToTile()
		if tile != test.tile {
			t.Errorf(errf, test.quadkey, test.tile, tile)
		}
	}
}

func ExampleFromCoordinate() {
	esbLat := 40.7484
	esbLon := -73.9857
	tile := tiles.FromCoordinate(esbLat, esbLon, 18)
	fmt.Println(tile)
}
