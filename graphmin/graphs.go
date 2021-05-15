package graphmin

import (
	"fmt"
	ds "github.com/ciromdrs/graph-tools/datastructures"
)

type (
	// A Graph is a set of Edges
	Graph interface {
		Add(*Edge)
		Get(ds.Vertex, ds.Vertex, ds.Vertex) *Edge
		Remove(ds.Vertex, ds.Vertex, ds.Vertex)
		Contains(ds.Vertex, ds.Vertex, ds.Vertex) bool
		Iterate() <-chan *Edge
		Size() int
		Equals(Graph) bool
		String() string
	}

	// HashGraph is a map-based Graph implementation.
	HashGraph struct {
		data *ds.SimpleMap
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
func newEdge(s, X, o ds.Vertex, isNecessary bool, dependencies []itemPos,
	exists bool) *Edge {
	e := &Edge{
		isNecessary:  isNecessary,
		dependencies: dependencies,
		exists:       exists,
	}
	e.triple = triple{s: s, X: X, o: o}
	return e
}

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
		data: ds.NewSimpleMap(),
	}
}

// Add uniquely adds an edge to the graph.
func (g *HashGraph) Add(e *Edge) {
	g.data.Set(e.triple.String(), e)
}

// Remove removes an edge from the graph.
func (g *HashGraph) Remove(s, X, o ds.Vertex) {
	t := triple{s: s, X: X, o: o}
	g.data.Remove(t.String())
}

// Contains checks wether the graph contains the given edge
func (g *HashGraph) Contains(s, X, o ds.Vertex) bool {
	t := triple{s: s, X: X, o: o}
	return g.data.Get(t.String()) != nil
}

// Iterate iterates over the edges.
func (g *HashGraph) Iterate() <-chan *Edge {
	ch := make(chan *Edge)
	go func() {
		for kv := range g.data.Iterate() {
			ch <- kv.Value.(*Edge)
		}
		defer close(ch)
	}()
	return ch
}

// Size returns the number of edges.
func (g *HashGraph) Size() int {
	return g.data.Size()
}

// Get returns the Edge object corresponding to (s,X,o).
func (g *HashGraph) Get(s, X, o ds.Vertex) *Edge {
	t := triple{s: s, X: X, o: o}
	return g.data.Get(t.String()).(*Edge)
}

// Equals checks whether graphs are equal.
func (g *HashGraph) Equals(other Graph) bool {
	if g.Size() != other.Size() {
		return false
	}
	for t := range g.Iterate() {
		if !other.Contains(t.s, t.X, t.o) {
			return false
		}
	}
	return true
}

func (g *HashGraph) String() string {
	return g.data.String()
}

// SimpleToHashGraph creates a HashGraph from a ds.SimpleGraph
func SimpleToHashGraph(simple *ds.SimpleGraph) *HashGraph {
	hash := newHashGraph()
	for t := range simple.Iterate() {
		hash.Add(newEdge(t[0], t[1], t[2], false, nil, true))
	}
	return hash
}
