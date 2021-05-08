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

func newEdge(s, X, o ds.Vertex) *Edge {
	e := &Edge{
		isNecessary:  false,
		dependencies: nil,
		exists:       false,
	}
	e.triple = triple{s: s, X: X, o: o}
	return e
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
func (g *HashGraph) Add(*Edge) bool {
	panic("Not implemented yet")
}

// Remove removes an edge from the graph. It returns a boolean value indicating
// whether the edge was in the graph.
func (g *HashGraph) Remove(*Edge) bool {
	panic("Not implemented yet")
}

// Contains checks wether the graph contains the given edge
func (g *HashGraph) Contains(*Edge) bool {
	panic("Not implemented yet")
}
