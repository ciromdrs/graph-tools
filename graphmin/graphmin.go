package graphmin

/* Package graphmin implements solvers for the Formal-Language-Constrained Graph
Minimization (FLGM) problem */

import (
	"fmt"
	"github.com/ciromdrs/graph-tools/ccfpq"
	ds "github.com/ciromdrs/graph-tools/data_structures"
)

type (
	// A GraphMinimizer solves the FLGM problem
	GraphMinimizer interface {
		minimize(ccfpq.Grammar, Graph, []Query, []Edge) Graph
	}

	// A Graph is a set of Edges
	Graph interface {
	}

	// An AugItem is an augmented trace item for solving the FLGM problem
	AugItem struct {
		lhs   ds.Vertex
		rhs   []ds.Vertex
		edges [][]*Edge
	}

	// A Query is a pair (vertex, symbol)
	Query struct {
		vertex ds.Vertex
		symbol ds.Vertex
	}

	itemPos struct {
		item *AugItem
		pos  int
	}

	// An Edge keeps track of AugItems where it appears
	Edge struct {
		s, X, o      ds.Vertex
		isNecessary  bool
		dependencies []itemPos
		exists       bool
	}
)

func newEdge(s, X, o ds.Vertex) *Edge {
	return &Edge{
		s:            s,
		X:            X,
		o:            o,
		isNecessary:  false,
		dependencies: nil,
		exists:       false,
	}
}

func (e *Edge) addDependency(item *AugItem, pos int) {
	e.dependencies = append(e.dependencies, itemPos{item: item, pos: pos})
}

func newAugItem(lhs ds.Vertex, rhs []ds.Vertex) *AugItem {
	edges := make([][]*Edge, len(rhs))
	return &AugItem{
		lhs:   lhs,
		rhs:   rhs,
		edges: edges,
	}
}

func (item *AugItem) addEdge(e *Edge, pos int) {
	if !e.exists {
		panic(fmt.Sprintf("Edge %v does not exist.", e))
	}
	if !e.X.Equals(item.rhs[pos]) {
		panic(fmt.Sprintf("Wrong predicate. Expected %v, got %v.",
			item.rhs[pos], e.X))
	}
	item.edges[pos] = append(item.edges[pos], e)
	e.addDependency(item, pos)
}
