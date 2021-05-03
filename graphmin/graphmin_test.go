package graphmin

import (
	ds "github.com/ciromdrs/graph-tools/data_structures"
	"testing"
)

func TestEdge(t *testing.T) {
	s := ds.NewSimpleVertex("s")
	p := ds.NewSimpleVertex("p")
	o := ds.NewSimpleVertex("o")
	e := newEdge(s, p, o)
	if e.isNecessary {
		t.Fatalf("New edges should not be necessary.")
	}
	if e.exists {
		t.Fatalf("New edges should not exist.")
	}
	if e.dependencies != nil {
		t.Fatalf("New edges should have no dependencies.")
	}
}

func TestAugItem(t *testing.T) {
	S := ds.NewSimpleVertex("S")
	a := ds.NewSimpleVertex("a")
	b := ds.NewSimpleVertex("b")
	c := ds.NewSimpleVertex("c")
	item := newAugItem(S, []ds.Vertex{a, b, c})
	if item.lhs != S {
		t.Fatalf("Wrong lhs. Expected %v, got %v", S, item.lhs)
	}
	if item.rhs[0] != a || item.rhs[1] != b || item.rhs[2] != c {
		t.Fatalf("Wrong rhs. Expected %v %v %v, got %v", a, b, c, item.rhs)
	}
	if len(item.edges) < 3 {
		t.Fatalf("Expected edges of length 3, got %v", len(item.edges))
	}
}
