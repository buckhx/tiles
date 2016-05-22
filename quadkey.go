package tiles

// Quadkey represents a Bing Maps quadkey
// It can also be used as a quadtree data structure
type Quadkey string

// HasParent returns a true if o is a parent of q.
// If q == o, it return false
func (q Quadkey) HasParent(o Quadkey) bool {
	z := len(o)
	if len(q) <= z {
		return false
	}
	return q[:z] == o
}

// Parent returns the parent of the object at the given level z.
// If level invalid (<0 || > q.Level()) it panics
func (q Quadkey) Parent(z int) Quadkey {
	return q[:z]
}

// Level returns the depth of the quadkey in the tree structure
func (q Quadkey) Level() int {
	return len(q)
}

// Children returns a slice of the the Quadkeys in the next level of this tree
func (q Quadkey) Children() []Quadkey {
	return []Quadkey{
		q + "0",
		q + "1",
		q + "2",
		q + "3",
	}
}

// ToTile returns the Tile represented by this Quadkey
func (q Quadkey) ToTile() Tile {
	t, err := FromQuadkeyString(string(q))
	check(err)
	return t
}
