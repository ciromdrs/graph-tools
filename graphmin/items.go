package graphmin

import (
	"fmt"
	ds "github.com/ciromdrs/graph-tools/datastructures"
)

type (
	// An AugItem is an augmented trace item for solving the FLGM problem.
	AugItem struct {
		rule  []ds.Vertex
		edges [][]*Edge
	}

	// AugItemSet is a set of AugItems.
	AugItemSet struct {
		data ds.Map
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

func newAugItemSet(f Factory, prealloc int) *AugItemSet {
	return &AugItemSet{
		data: f.NewMap(prealloc),
	}
}
