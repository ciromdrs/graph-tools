package graphmin

/* Package graphmin implements solvers for the Formal-Language-Constrained Graph
Minimization (FLGM) problem */

import (
	"fmt"
	"github.com/ciromdrs/graph-tools/ccfpq"
	ds "github.com/ciromdrs/graph-tools/datastructures"
)

type (
	// A GraphMinimizer solves the FLGM problem
	GraphMinimizer interface {
		minimize(ccfpq.Grammar, Graph, []Query, []Edge) Graph
	}

	// An AugItem is an augmented trace item for solving the FLGM problem
	AugItem struct {
		rule  []ds.Vertex
		edges [][]*Edge
	}

	// A Query is a pair (vertex, symbol)
	Query struct {
		vertex ds.Vertex
		symbol ds.Vertex
	}
)

func newAugItem(rule []ds.Vertex) *AugItem {
	edges := make([][]*Edge, len(rule))
	return &AugItem{
		rule:  rule,
		edges: edges,
	}
}

func (item *AugItem) addEdge(e *Edge, pos int) {
	if !e.exists {
		panic(fmt.Sprintf("Edge %v does not exist.", e))
	}
	if !e.X.Equals(item.rule[pos]) {
		panic(fmt.Sprintf("Wrong predicate. Expected %v, got %v.",
			item.rule[pos], e.X))
	}
	item.edges[pos] = append(item.edges[pos], e)
	e.addDependency(item, pos)
}
