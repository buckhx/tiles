package tiles

import (
	"fmt"
	"math"
)

const maxChildrenDiff = 8

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

// Children returns a slice of the Quadkeys in the next level of this tree
func (q Quadkey) Children() []Quadkey {
	return []Quadkey{
		q + "0",
		q + "1",
		q + "2",
		q + "3",
	}
}

// ChildrenAt returns a slice of the Quadkeys in the specified level of this tree.
// The number of recursively generated children is 4^(currentLevel - z).
// For technical reasons, z can only be at max 8 levels more precise than the current level.
func (q Quadkey) ChildrenAt(z int) ([]Quadkey, error) {
	diff := z - q.Level()
	cap := int(math.Pow(4, float64(diff)))

	if z <= q.Level() {
		return make([]Quadkey, 0), nil
	} else if z > ZMax {
		return make([]Quadkey, 0), fmt.Errorf("level has to be less than %d", ZMax)
	} else if diff > maxChildrenDiff {
		return make([]Quadkey, 0), fmt.Errorf("children level has to be less than %d", q.Level()+maxChildrenDiff+1)
	}

	if diff == 1 {
		return q.Children(), nil
	}

	children := make([]Quadkey, 0, cap)
	for _, c := range q.Children() {
		add, err := c.ChildrenAt(z)
		if err != nil {
			return make([]Quadkey, 0), err
		}
		children = append(children, add...)
	}

	return children, nil
}

// ToTile returns the Tile represented by this Quadkey
func (q Quadkey) ToTile() Tile {
	t, err := FromQuadkeyString(string(q))
	check(err)
	return t
}
