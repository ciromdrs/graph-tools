package graphmin

/* Package graphmin implements solvers for the Formal-Language-Constrained Graph
Minimization (FLGM) problem */

import (
	"github.com/ciromdrs/graph-tools/ccfpq"
	ds "github.com/ciromdrs/graph-tools/data_structures"
)

type (
	// A GraphMinimizer solves the FLGM problem
	GraphMinimizer interface {
		minimize(ccfpq.Grammar, Graph, []Query, []Edge) Graph
	}

	// An Edge keeps track of AugItems where it appears
	Edge interface {
	}

	// A Graph is a set of Edges
	Graph interface {
	}

	// An AugItem is an augmented trace item for solving the FLGM proble
	AugItem interface {
	}

	// A Query is a pair (vertex, symbol)
	Query struct {
		vertex ds.Vertex
		symbol ds.Vertex
	}
)
