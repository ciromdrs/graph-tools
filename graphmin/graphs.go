package graphmin

import (
	"fmt"
	ds "github.com/ciromdrs/graph-tools/datastructures"
)

type (
	// A Graph is a set of Edges
	Graph interface {
		Add(*Edge) bool
		Remove(*Edge) bool
		Contains(*Edge) bool
	}

	// HashGraph is a map-based Graph implementation.
	HashGraph struct {
		data *ds.MapSet
	}

	// An Edge keeps track of AugItems where it appears
	Edge struct {
		triple
		isNecessary  bool
		dependencies []itemPos
		exists       bool
	}

	triple struct {
		s, X, o ds.Vertex
	}

	itemPos struct {
		item *AugItem
		pos  int
	}
)

/* Edge methods and functions */

// newEdge creates an Edge object.
func newEdge(s, X, o ds.Vertex) *Edge {
	e := &Edge{
		isNecessary:  false,
		dependencies: nil,
		exists:       false,
	}
	e.triple = triple{s: s, X: X, o: o}
	return e
}

// Convert triple to string
func (t triple) String() string {
	return fmt.Sprintf("(%v, %v, %v)", t.s.String(), t.X.String(), t.o.String())
}

func (e *Edge) addDependency(item *AugItem, pos int) {
	for _, ip := range e.dependencies {
		if ip.item == item && ip.pos == pos {
			panic(fmt.Sprintf("Duplicated dependency (%v,%v).", item, pos))
		}
	}
	e.dependencies = append(e.dependencies, itemPos{item: item, pos: pos})
}

/* HashGraph methods and functions */

// newHashGraph creates a HashGraph object
func newHashGraph() *HashGraph {
	return &HashGraph{
		data: ds.NewMapSet(),
	}
}

// Add adds an edge to the graph. It returns a boolean value indicating
// whether the edge was in the graph.
func (g *HashGraph) Add(e *Edge) bool {
	return g.data.Add(e)
}

// Remove removes an edge from the graph. It returns a boolean value indicating
// whether the edge was in the graph.
func (g *HashGraph) Remove(e *Edge) bool {
	return g.data.Remove(e)
}

// Contains checks wether the graph contains the given edge
func (g *HashGraph) Contains(e *Edge) bool {
	return g.data.Contains(e)
}
