package graphmin

import (
	"fmt"
	ds "github.com/ciromdrs/graph-tools/datastructures"
)

type (
	// An AugItem is an augmented trace item for solving the FLGM problem.
	AugItem struct {
		Rule   []ds.Vertex
		Posets []Graph // Posets are sets of edges, which, in turn, are graphs
	}

	// AugItemSet is a set of AugItems.
	AugItemSet struct {
		data ds.Map
	}
)

func newAugItem(rule []ds.Vertex, posets []Graph) *AugItem {
	if len(rule) != len(posets) {
		panic("Rule and Posets must be of same length.")
	}
	return &AugItem{
		Rule:   rule,
		Posets: posets,
	}
}

// AddEdge adds a connection edge to an item and its corresponding dependency.
func (aug *AugItem) AddEdge(e *Edge, pos int) {
	if !e.exists {
		panic(fmt.Sprintf("Edge %v does not exist.", e))
	}
	if !e.X.Equals(aug.Rule[pos]) {
		panic(fmt.Sprintf("Wrong predicate. Expected %v, got %v.",
			aug.Rule[pos], e.X))
	}
	aug.Posets[pos].Add(e)
	e.addDependency(aug, pos)
}

// Equals checks whether two items are equal by comparing values, not pointers.
func (aug *AugItem) Equals(other *AugItem) bool {
	if len(aug.Rule) != len(other.Rule) {
		return false
	}
	for i := range aug.Rule {
		if aug.Rule[i] != other.Rule[i] {
			return false
		}
	}
	for i := range aug.Posets {
		if !aug.Posets[i].Equals(other.Posets[i]) {
			return false
		}
	}
	return true
}

func newAugItemSet(f Factory, prealloc int) *AugItemSet {
	return &AugItemSet{
		data: f.NewMap(prealloc),
	}
}

// Add adds an item or replaces an old one, in case it exists.
func (s *AugItemSet) Add(new *AugItem) {
	var items []*AugItem
	if got, ok := s.data.Get(new.Rule[0]).([]*AugItem); ok {
		items = got
	}
	// try to replace
	replaced := false
	for i, old := range items {
		if len(old.Rule) == len(new.Rule) {
			equal := true
			for j := range old.Rule {
				equal = equal && old.Rule[j].Equals(new.Rule[j])
			}
			if equal {
				items[i] = new
				replaced = true
				break
			}
		}
	}
	// if not replaced, add new
	if !replaced {
		items = append(items, new)
	}
	s.data.Set(new.Rule[0], items)
}

// Get returns the AugItem for the given rule.
func (s *AugItemSet) Get(rule []ds.Vertex) *AugItem {
	items := s.GetAll(rule[0])
	for _, it := range items {
		if len(it.Rule) == len(rule) {
			equal := true
			for i := range it.Rule[1:] {
				equal = equal && it.Rule[i] == rule[i]
			}
			if equal {
				return it
			}
		}
	}
	return nil
}

// GetAll returns all AugItems for the given lhs.
func (s *AugItemSet) GetAll(lhs ds.Vertex) []*AugItem {
	return s.data.Get(lhs).([]*AugItem)
}

// Size returns the number of AugItems.
func (s *AugItemSet) Size() int {
	size := 0
	for kv := range s.data.Iterate() {
		size += len(kv.Value.([]*AugItem))
	}
	return size
}
