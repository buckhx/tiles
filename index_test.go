package tiles

import (
	"bytes"
	"fmt"
	"index/suffixarray"
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
	idx := NewTileIndex()
	testIndex(t, idx)
}

func TestKeysetIndex(t *testing.T) {
	idx := &KeysetIndex{}
	testIndex(t, TileIndex(idx))
}

func TestSuffixIndex(t *testing.T) {
	idx := NewSuffixIndex()
	testIndex(t, TileIndex(idx))
}

func TestPrefixes(t *testing.T) {
	keys := [][]byte{
		[]byte("0123"),
		[]byte("01234"),
		[]byte("0231"),
		[]byte("01412"),
		[]byte("3023"),
	}
	d := []byte{0}
	b := bytes.Join(keys, d)              //join w/ zeros
	data := bytes.Join([][]byte{d, d}, b) //pad w/ zeros
	idx := suffixarray.New(data)
	if len(prefixes(idx, data, []byte("01"))) != 3 {
		t.Error("prefixes() did not return correct keys")
	}
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

var bV []interface{}

func BenchmarkKeysetValues(b *testing.B) {
	idx := &KeysetIndex{}
	hydrateIndex(idx)
	idx.sort()
	esb := Tile{X: 9649, Y: 12315, Z: 15}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		bV = idx.Values(esb)
	}
}

func BenchmarkSuffixValues(b *testing.B) {
	idx := NewSuffixIndex()
	hydrateIndex(idx)
	idx.sort()
	esb := Tile{X: 9649, Y: 12315, Z: 15}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		bV = idx.Values(esb)
	}
}

func hydrateIndex(idx TileIndex) {
	mlat, mlon := 40.7, -73.9
	for i := 0; i < 10000; i++ {
		lat := mlat + 0.1*rand.Float64()
		lon := mlon - 0.1*rand.Float64()
		t := FromCoordinate(lat, lon, 18)
		idx.Add(t, fmt.Sprintf("%f,%f", lat, lon))
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
