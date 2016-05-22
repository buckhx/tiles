package tiles

import (
	"fmt"
	"math/rand"
	"testing"
)

//TODO actually write a test
func TestTileRange(t *testing.T) {
	qks := []string{
		"0000",
		"0001",
		"0010",
		"0011",
		"0100",
		"0101",
		"1111",
	}
	idx := NewTileIndex()
	for i, qk := range qks {
		tile, _ := FromQuadkeyString(qk)
		idx.Add(tile, i)
	}
	c := 0
	for range idx.TileRange(1, 2) {
		//t.Logf("Tile %+v", tile.Quadkey())
		c++
	}
	if c != 5 {
		t.Error("TileRange should generate 5 tiles, got ", c)
	}
}

func TestTileIndex(t *testing.T) {
	idx := NewSuffixIndex()
	testIndex(t, TileIndex(idx))
}

func TestSuffixIndex(t *testing.T) {
	idx := NewSuffixIndex()
	testIndex(t, TileIndex(idx))
}

func testIndex(t *testing.T, idx TileIndex) {
	esb := FromCoordinate(40.7484, -73.9857, 18)
	sol := FromCoordinate(40.6892, -74.0445, 18)
	bbn := FromCoordinate(51.5007, -0.1246, 18)
	idx.Add(esb, "EmpireStateBuilding")
	idx.Add(sol, "StatueOfLiberty")
	idx.Add(bbn, "BigBen")
	nyc := Tile{X: 75, Y: 96, Z: 8}
	den := Tile{X: 106, Y: 194, Z: 9}
	switch {
	case len(idx.Values(esb)) != 1:
		t.Error("ESB: ", idx.Values(esb))
	case len(idx.Values(sol)) != 1:
		t.Error("SOL: ", idx.Values(sol))
	case len(idx.Values(nyc)) != 2:
		t.Error("NYC: ", idx.Values(nyc))
	case len(idx.Values(den)) != 0:
		t.Error("DEN: ", idx.Values(nyc))
	}
}

var vQ []interface{}

func BenchmarkKeysetValues(b *testing.B) {
	idx := &KeysetIndex{}
	mlat, mlon := 40.0, 73.0
	for i := 0; i < 10000; i++ {
		lat := mlat + rand.Float64()
		lon := mlon - rand.Float64()
		t := FromCoordinate(lat, lon, 18)
		idx.Add(t, fmt.Sprintf("%f,%f", lat, lon))
	}
	esb := Tile{X: 9649, Y: 12315, Z: 15}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		vQ = idx.Values(esb)
	}
}

func BenchmarkSuffixValues(b *testing.B) {
	idx := NewSuffixIndex() //NewTileIndex()
	mlat, mlon := 40.0, 73.0
	for i := 0; i < 10000; i++ {
		lat := mlat + rand.Float64()
		lon := mlon - rand.Float64()
		t := FromCoordinate(lat, lon, 18)
		idx.Add(t, fmt.Sprintf("%f,%f", lat, lon))
	}
	esb := Tile{X: 9649, Y: 12315, Z: 15}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		vQ = idx.Values(esb)
	}
}

func ExampleTileIndex() {
	idx := NewTileIndex()
	esb := FromCoordinate(40.7484, -73.9857, 18)
	sol := FromCoordinate(40.6892, -74.0445, 18)
	bbn := FromCoordinate(51.5007, -0.1246, 18)
	idx.Add(esb, "EmpireStateBuilding")
	idx.Add(sol, "StatueOfLiberty")
	idx.Add(bbn, "BigBen")
	nyc := Tile{X: 75, Y: 96, Z: 8}
	den := Tile{X: 106, Y: 194, Z: 9}
	fmt.Println("ESB Tile: ", idx.Values(esb))
	fmt.Println("SOL Tile: ", idx.Values(sol))
	fmt.Println("NYC Tile: ", idx.Values(nyc))    //contains both values!
	fmt.Println("DENVER Tile: ", idx.Values(den)) //contains no values!
	// Output: ESB Tile:  [EmpireStateBuilding]
	// SOL Tile:  [StatueOfLiberty]
	// NYC Tile:  [EmpireStateBuilding StatueOfLiberty]
	// DENVER Tile:  []
}
