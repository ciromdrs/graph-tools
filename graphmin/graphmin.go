package graphmin

/* Package graphmin implements solvers for the Formal-Language-Constrained Graph
Minimization (FLGM) problem */

import (
	"github.com/ciromdrs/graph-tools/ccfpq"
	ds "github.com/ciromdrs/graph-tools/datastructures"
)

type (
	// A GraphMinimizer solves the FLGM problem
	GraphMinimizer interface {
		minimize(ccfpq.Grammar, Graph, []Query, []Edge) Graph
	}

	// A Query is a pair (vertex, symbol)
	Query struct {
		vertex ds.Vertex
		symbol ds.Vertex
	}
)
