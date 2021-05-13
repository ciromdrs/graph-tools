package graphmin

import (
	"fmt"
	ds "github.com/ciromdrs/graph-tools/datastructures"
	. "github.com/ciromdrs/graph-tools/util"
	"testing"
)

func TestAugItem(t *testing.T) {
	var f Factory = NewHashFactory()
	S := ds.NewSimpleVertex("S")
	a := ds.NewSimpleVertex("a")
	b := ds.NewSimpleVertex("b")
	c := ds.NewSimpleVertex("c")
	s := ds.NewSimpleVertex("s")
	o := ds.NewSimpleVertex("o")
	e1 := newEdge(s, a, o)

	rule := []ds.Vertex{S, a, b, c}
	posets := f.NewEmptyPosets(len(rule))
	item := newAugItem(rule, posets)
	if item.Rule[0] != S || item.Rule[1] != a || item.Rule[2] != b ||
		item.Rule[3] != c {
		t.Fatalf("Wrong rule. Expected %v, got %v", rule, item.Rule)
	}
	if len(item.Posets) < 3 {
		t.Fatalf("Expected edges of length 3, got %v", len(item.Posets))
	}

	AssertPanic(t, func() { item.AddEdge(e1, 0) },
		fmt.Sprintf("Should not add inexistent edge %v.", e1))
	e1.exists = true
	AssertPanic(t, func() { item.AddEdge(e1, 0) },
		fmt.Sprintf("Should not add edge %v with wrong predicate.", e1))
	e1.exists = true
	item.AddEdge(e1, 1)
	Assert(t, item.Posets[1].Contains(s, a, o),
		fmt.Sprintf("Error adding edge. item.Posets[1] should contain %v.", e1))
	{
		want := itemPos{item: item, pos: 0}
		if e1.dependencies[0].item != item || e1.dependencies[0].pos != 1 {
			t.Fatalf("Eror adding dependency. Expected %v, got %v",
				want, e1.dependencies[0])
		}
	}
	AssertPanic(t, func() { item.AddEdge(e1, 0) },
		fmt.Sprintf("Should not add duplicated dependency %v %d.", e1, 0))
	e2 := newEdge(s, b, o)
	e2.exists = true
	AssertPanic(t, func() { item.AddEdge(e2, 0) },
		"Should not add edge with wrong predicate b.")
}

func TestAugItemSet(t *testing.T) {
	f := NewHashFactory()
	f.NewAugItemSet(0)
}
