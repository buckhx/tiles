package tiles

import (
	"testing"
)

func TestClip(t *testing.T) {
	clipTests := []struct {
		val, min, max, out float64
	}{
		{0, 1, 5, 1},
		{9, 1, 5, 5},
		{3, 1, 5, 3},
	}
	errf := "clip() %+v -> %+v"
	for _, test := range clipTests {
		val := clip(test.val, test.min, test.max)
		if val != test.out {
			t.Errorf(errf, test, val)
		}
	}
}

func TestMapDimensions(t *testing.T) {
	mapDimTests := []struct {
		zoom, out int
	}{
		{0, 256},
		{1, 512},
		{15, 8388608},
	}
	errf := "mapDimensions() %+v -> %+v"
	for _, test := range mapDimTests {
		val := mapDimensions(test.zoom)
		if val != test.out {
			t.Errorf(errf, test, val)
		}
	}
}

/*
//TODO assert neither of these are used and remove
func TestGroundRes(t *testing.T) {
	lat := 40.0
	var zoom int = 7
	res := 936.86657226219847
	if out := grndRes(lat, zoom); !floatEquals(out, res) {
		t.Errorf("grndRes(%v, %v) -> %v not %v", lat, zoom, res, out)
	}
}

func TestMapScale(t *testing.T) {
	lat := 40.0
	var zoom int = 7
	var dpi int = 96
	scale := 3540913.0290224836
	if out := mapScale(lat, zoom, dpi); !floatEquals(out, scale) {
		t.Errorf("mapScale(%v, %v, %v) -> %v not %v", lat, zoom, dpi, scale, out)
	}
}
*/
