package tiles

import (
	"testing"
)

//TODO actually write a test
func TestTileIndex(t *testing.T) {
	qks := []string{
		"0000",
		"0001",
		"0010",
		"0011",
		"0100",
		"0101",
		"1111",
	}
	idx := &TileIndex{}
	for i, qk := range qks {
		idx.Add(TileFromQuadKey(qk), i)
	}
	for tile := range idx.TileRange(1, 2) {
		t.Logf("Tile %+v", tile.QuadKey())
	}
}
