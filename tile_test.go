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

var (
	// These are globals to make sure that the compiler doesn't skip benchmarks
	bT tiles.Tile
	bQ tiles.Quadkey
)

func BenchmarkTileFromCoordinate(b *testing.B) {
	var t tiles.Tile
	z := 18
	lat := 40.7484
	lon := -73.9857
	for i := 0; i < b.N; i++ {
		t = tiles.FromCoordinate(lat, lon, z)
	}
	bT = t
}

func BenchmarkTileFromQuadkey(b *testing.B) {
	var t tiles.Tile
	qk := "032010110132023321"
	for i := 0; i < b.N; i++ {
		t, _ = tiles.FromQuadkeyString(qk)
	}
	bT = t

}

func BenchmarkQuadkeyFromTile(b *testing.B) {
	var q tiles.Quadkey
	t := tiles.Tile{X: 77197, Y: 98526, Z: 18}
	for i := 0; i < b.N; i++ {
		q = t.Quadkey()
	}
	bQ = q
}

func ExampleFromCoordinate() {
	esbLat := 40.7484
	esbLon := -73.9857
	tile := tiles.FromCoordinate(esbLat, esbLon, 18)
	fmt.Println(tile)
}
