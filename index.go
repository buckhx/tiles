package tiles

import (
	"bytes"
	"index/suffixarray"
	"sort"
	"sync"
)

const (
	zero = byte('\x00')
)

// TileIndex stores indexes values by tile.
// If a deep level of tile is added and a shallower one is requested, the values are aggregated up.
type TileIndex interface {
	TileRange(zmin, zmax int) <-chan Tile
	Values(t Tile) (vals []interface{})
	Add(t Tile, val ...interface{})
}

// NewTileIndex returns the default TileIndex
func NewTileIndex() TileIndex {
	return &KeysetIndex{}
}

// KeysetIndex is a TileIndex implementation that uses a sorted keyset.
// A trie would be more efficient, but KeysetIndex mirrors the range queries of boltdb which could be dropped in if the entire index won't fit in memory.
// KeysetIndex is thread safe
type KeysetIndex struct {
	// Implementation uses a sorted keyset.
	// A trie would be more efficient, but
	sorted bool
	keys   []qkey
	values [][]interface{}
	sync.RWMutex
}

// TileRange returns a channel of all tiles in the index in the zoom range
// If zmax is greater than the deepest tile level, the deepest tile level returns
// Acquires a readlock for duration of returned channel being open
func (idx *KeysetIndex) TileRange(zmin, zmax int) <-chan Tile {
	idx.sort()
	tiles := make(chan Tile, 1<<10)
	go func() {
		defer close(tiles)
		idx.RLock()
		defer idx.RUnlock()
		for i := 0; i < len(idx.keys)-1; i++ {
			qmax := idx.keys[i].qk.Level()
			for z := zmin; z <= zmax && z <= qmax; z++ {
				q := idx.keys[i].qk.Parent(z)
				n := idx.keys[i+1].qk
				if !n.HasParent(q) {
					tiles <- q.ToTile()
				}
			}
		}
		q := idx.keys[len(idx.keys)-1].qk
		for z := zmin; z <= zmax && z <= len(q); z++ {
			tiles <- q.Parent(z).ToTile()
		}
	}()
	return tiles
}

// Values returns a list of values aggregated under the requested tile
func (idx *KeysetIndex) Values(t Tile) (vals []interface{}) {
	idx.sort()
	idx.RLock()
	defer idx.RUnlock()
	qk := t.Quadkey()
	i := idx.search(qk)
	for i < len(idx.keys) {
		n := idx.keys[i]
		if n.qk == qk || n.qk.HasParent(qk) {
			vals = append(vals, idx.values[n.v]...)
		}
		i++
	}
	return
}

// Add adds a value, but will not be indexed
func (idx *KeysetIndex) Add(t Tile, val ...interface{}) {
	idx.Lock()
	defer idx.Unlock()
	idx.values = append(idx.values, val)
	qk := qkey{qk: t.Quadkey(), v: len(idx.values) - 1}
	idx.keys = append(idx.keys, qk)
	idx.sorted = false
}

// sorts the tiles, nothing happens if the sorted flag is set
func (idx *KeysetIndex) sort() {
	if !idx.sorted {
		idx.Lock()
		sort.Sort(byQk(idx.keys))
		idx.sorted = true
		idx.Unlock()
	}
}

func (idx *KeysetIndex) search(qk Quadkey) int {
	return sort.Search(len(idx.keys), func(i int) bool { return idx.keys[i].qk >= qk })
}

type qkey struct {
	qk Quadkey
	v  int
}

type byQk []qkey

func (q byQk) Len() int           { return len(q) }
func (q byQk) Swap(i, j int)      { q[i], q[j] = q[j], q[i] }
func (q byQk) Less(i, j int) bool { return q[i].qk < q[j].qk }

//SuffixIndex is a TileIndex that uses a suffixarray to lookup values
//It IS NOT currently safe for concurrent access.
type SuffixIndex struct {
	// \x00 joined string of keys for suffixarray
	indexed []byte
	index   *suffixarray.Index
	tiles   map[Quadkey][]interface{}
}

//NewSuffixIndex returns a new SuffixIndex
func NewSuffixIndex() *SuffixIndex {
	return &SuffixIndex{
		tiles: make(map[Quadkey][]interface{}),
	}
}

//TileRange returns all the tiles available in this index.
//It currently DOES NOT return unique values
func (idx *SuffixIndex) TileRange(zmin, zmax int) <-chan Tile {
	tiles := make(chan Tile, 1<<10)
	go func() {
		defer close(tiles)
		seen := make(map[Tile]struct{}, len(idx.tiles)*(zmax-zmin+1))
		for k := range idx.tiles {
			for z := zmin; z <= zmax; z++ {
				t := k[:z].ToTile()
				if _, ok := seen[t]; !ok {
					tiles <- t
					seen[t] = struct{}{}
				}
			}
		}
	}()
	return tiles
}

//Values returns all the values aggregated under the given tile
func (idx *SuffixIndex) Values(t Tile) (vals []interface{}) {
	idx.sort()
	qk := t.Quadkey()
	keys := prefixes(idx.index, idx.indexed, []byte(qk))
	for _, k := range keys {
		qk := Quadkey(k)
		//fmt.Println(qk)
		vals = append(vals, idx.tiles[qk]...)
	}
	return
}

//Add adds a tile and values associated with it
func (idx *SuffixIndex) Add(t Tile, v ...interface{}) {
	// Set index to nil b/c adding invalidates index
	idx.index = nil
	qk := t.Quadkey()
	idx.tiles[qk] = append(idx.tiles[qk], v...)
}

func (idx *SuffixIndex) sort() {
	if idx.index == nil {
		keys := make([][]byte, len(idx.tiles))
		i := 0
		for k := range idx.tiles {
			keys[i] = []byte(k)
			i++
		}
		d := []byte{zero}
		b := bytes.Join(keys, d)                    //join w/ zeros
		idx.indexed = bytes.Join([][]byte{d, d}, b) //pad w/ zeros
		idx.index = suffixarray.New(idx.indexed)
	}
}

//prefixes assumes a \x00 delimited data with \x00 padding
func prefixes(idx *suffixarray.Index, data, q []byte) (keys [][]byte) {
	for _, i := range idx.Lookup(q, -1) {
		if data[i-1] == zero { //if previous byte is zero, it's a prefix
			var l int
			for l = i; data[l] != zero; l++ {
				// iterating until end of word
			}
			k := data[i:l]
			keys = append(keys, k)
		}
	}
	return
}
