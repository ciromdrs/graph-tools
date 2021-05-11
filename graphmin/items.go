package graphmin

import (
	"fmt"
	ds "github.com/ciromdrs/graph-tools/datastructures"
)

type (
	// An AugItem is an augmented trace item for solving the FLGM problem.
	AugItem struct {
		Rule   []ds.Vertex
		Posets [][]*Edge
	}

	// AugItemSet is a set of AugItems.
	AugItemSet struct {
		data ds.Map
	}
)

func newAugItem(rule []ds.Vertex) *AugItem {
	posets := make([][]*Edge, len(rule))
	return &AugItem{
		Rule:   rule,
		Posets: posets,
	}
}

func (item *AugItem) AddEdge(e *Edge, pos int) {
	if !e.exists {
		panic(fmt.Sprintf("Edge %v does not exist.", e))
	}
	if !e.X.Equals(item.Rule[pos]) {
		panic(fmt.Sprintf("Wrong predicate. Expected %v, got %v.",
			item.Rule[pos], e.X))
	}
	item.Posets[pos] = append(item.Posets[pos], e)
	e.addDependency(item, pos)
}

func newAugItemSet(f Factory, prealloc int) *AugItemSet {
	return &AugItemSet{
		data: f.NewMap(prealloc),
	}
}
