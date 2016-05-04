package tiles

import (
	"testing"
)

func TestQuadkeyHasParent(t *testing.T) {
	tests := []struct {
		q Quadkey
		p Quadkey
		t bool
	}{
		{"", "", false},
		{"0", "001", false},
		{"012", "0", true},
		{"0", "0", false},
		{"012123123", "012", true},
	}
	errf := "Quadkey(%q).HasParent(%q) -> %+v"
	for _, test := range tests {
		if test.q.HasParent(test.p) != test.t {
			t.Errorf(errf, test.q, test.p, !test.t)
		}
	}
}

func TestQuadkeyParent(t *testing.T) {
	tests := []struct {
		q Quadkey
		z int
		p Quadkey
	}{
		{"", 0, ""},
		{"0", 0, ""},
		{"012", 1, "0"},
		{"012", 3, "012"},
	}
	errf := "Quadkey(%q).Parent(%d) -> %q"
	for _, test := range tests {
		p := test.q.Parent(test.z)
		if p != test.p {
			t.Errorf(errf, test.q, test.z, p)
		}
	}
}

func TestQuadkeyLevel(t *testing.T) {
	tests := []struct {
		q Quadkey
		z int
	}{
		{"", 0},
		{"0", 1},
		{"012", 3},
	}
	errf := "Quadkey(%q).Level() -> %d"
	for _, test := range tests {
		z := test.q.Level()
		if z != test.z {
			t.Errorf(errf, test.q, z)
		}
	}
}

func TestQuadkeyChildren(t *testing.T) {
	tests := []struct {
		q Quadkey
		c []Quadkey
	}{
		{"", []Quadkey{"0", "1", "2", "3"}},
		{"0", []Quadkey{"00", "01", "02", "03"}},
	}
	errf := "Quadkey(%q).Children() -> %+v"
	for _, test := range tests {
		c := test.q.Children()
		if !qkSliceEqual(c, test.c) {
			t.Errorf(errf, test.q, c)
		}
	}
}

func TestQuadkeyToTile(t *testing.T) {
	tests := []struct {
		q Quadkey
		t Tile
	}{
		{"", Tile{X: 0, Y: 0, Z: 0}},
		{"0", Tile{X: 0, Y: 0, Z: 1}},
		{"0231010301", Tile{X: 213, Y: 388, Z: 10}},
	}
	errf := "Quadkey(%q).ToTile() -> %+v"
	for _, test := range tests {
		tile := test.q.ToTile()
		if tile != test.t {
			t.Errorf(errf, test.q, tile)
		}
	}
}

func qkSliceEqual(x, y []Quadkey) bool {
	if len(x) != len(y) {
		return false
	}
	for i, v := range x {
		if y[i] != v {
			return false
		}
	}
	return true
}
