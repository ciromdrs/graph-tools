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
	e1 := newEdge(s, a, o, false, nil, true)

	rule := []ds.Vertex{S, a, b, c}
	AssertPanic(t, func() { newAugItem(rule, f.NewEmptyPosets(len(rule)-1)) },
		"Should not allow different legth.")
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
	e2 := newEdge(s, b, o, false, nil, true)
	e2.exists = true
	AssertPanic(t, func() { item.AddEdge(e2, 0) },
		"Should not add edge with wrong predicate b.")
	posets = f.NewEmptyPosets(len(rule))
	item2 := newAugItem(rule, posets)
	item2.AddEdge(e1, 1)
	Assert(t, item.Equals(item2),
		fmt.Sprintf("Wrong item. Want %v, got %v.", item2, item))
	item2 = newAugItem([]ds.Vertex{S}, f.NewEmptyPosets(1))
	Assert(t, !item.Equals(item2), "Items of length should be different.")
	item2 = newAugItem([]ds.Vertex{S, c, b, a}, f.NewEmptyPosets(4))
	Assert(t, !item.Equals(item2), "Items' rule should be different.")
	item2 = newAugItem(rule, f.NewEmptyPosets(len(rule)))
	item2.AddEdge(e2, 2)
	Assert(t, !item.Equals(item2), "Items should be different.")
}

func TestAugItemSet(t *testing.T) {
	var f Factory = NewHashFactory()
	set := f.NewAugItemSet(0)
	S := ds.NewSimpleVertex("S")
	a := ds.NewSimpleVertex("a")
	b := ds.NewSimpleVertex("b")
	c := ds.NewSimpleVertex("c")
	one := ds.NewSimpleVertex("1")
	two := ds.NewSimpleVertex("2")
	three := ds.NewSimpleVertex("3")

	rule := []ds.Vertex{S, a, b, c}
	it := newAugItem(rule, f.NewEmptyPosets(len(rule)))
	set.Add(it)
	Assert(t, set.Size() == 1,
		fmt.Sprintf("Wrong Size(). Want 1, got %v.", set.Size()))

	it2 := newAugItem(rule, f.NewEmptyPosets(len(rule)))
	e1 := newEdge(one, a, two, false, nil, true)
	it2.AddEdge(e1, 1)
	set.Add(it2)
	Assert(t, set.Size() == 1,
		fmt.Sprintf("Wrong Size(). Want 1, got %v.", set.Size()))
	Assert(t, set.Get(rule).Equals(it2), "Items should be equal.")
	Assert(t, set.Get([]ds.Vertex{S}) == nil, "Wrong item. Expected nil.")
	all := set.GetAll(rule[0])
	Assert(t, len(all) == 1,
		fmt.Sprintf("Wrong length. Want 1, got %v.", len(all)))
	Assert(t, all[0].Equals(it2),
		fmt.Sprintf("Wrong items. Expected { %v }, got %v", it2, all))

	set2 := f.NewAugItemSet(0)
	Assert(t, !set.Equals(set2), "AugItemSets should not be equal.")
	it3 := newAugItem(rule, f.NewEmptyPosets(len(rule)))
	e2 := newEdge(one, a, three, false, nil, true)
	it3.AddEdge(e2, 1)
	Assert(t, !set.Equals(set2), "AugItemSets should not be equal.")
	set2 = f.NewAugItemSet(0)
	it3 = newAugItem(rule, f.NewEmptyPosets(len(rule)))
	it3.AddEdge(e1, 1)
	set2.Add(it3)
	Assert(t, set.Equals(set2), "AugItemSets should be equal.")
}
