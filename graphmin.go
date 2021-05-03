package graphmin

/* Package graphmin implements solvers for the Formal-Language-Constrained Graph
Minimization (FLGM) problem */

// GraphMinimizer solves the FLGM problem
type GraphMinimizer interface {
	minimize(Grammar, Graph, []Query, []Edge) Graph
}
