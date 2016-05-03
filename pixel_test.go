package tiles

import (
	"testing"
)

func TestPixelToCoords(t *testing.T) {
	coordTests := []struct {
		pixel  Pixel
		coords Coords
	}{
		{Pixel{6827, 12405, 7}, Coords{40.002372, -104.996338}},
	}
	errf := "Pixel%+v: %+v -> %+v"
	for _, test := range coordTests {
		coords := test.pixel.ToCoords()
		if coords.Equals(test.coords) {
			t.Errorf(errf, test.pixel, test.coords, coords)
		}
	}
}

func TestPixelToTile(t *testing.T) {
	coordTests := []struct {
		pixel Pixel
		tile  Tile
	}{
		{Pixel{6827, 12405, 7}, Tile{26, 48, 7}},
	}
	errf := "Pixel%+v: %+v -> %+v"
	for _, test := range coordTests {
		tile, tpixel := test.pixel.ToTile()
		_ = tpixel
		// t.Logf("%+v -> %+v %+v", test.pixel, tile, tpixel)
		// TODO test tpixel
		// pixel_test.go:33: {X:6827 Y:12405 Z:7} -> {X:26 Y:48 Z:7} {X:171 Y:117 Tile:0xc20801e040}
		if tile != test.tile {
			t.Errorf(errf, test.pixel, test.tile, tile)
		}
	}
}
