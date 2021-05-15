package graphmin

import (
	"fmt"
	"github.com/ciromdrs/graph-tools/ccfpq"
	ds "github.com/ciromdrs/graph-tools/datastructures"
	. "github.com/ciromdrs/graph-tools/util"
	"testing"
)

func TestConverter(t *testing.T) {
	G, D, F := ccfpq.QuickLoad("../ccfpq/testdata/aSb.yrd",
		"../ccfpq/testdata/graph-ti.txt",
		ds.SIMPLE_FACTORY)

	S := F.NewPredicate("S")
	a := F.NewPredicate("a")
	b := F.NewPredicate("b")
	one := F.NewVertex("1")
	two := F.NewVertex("2")
	three := F.NewVertex("3")
	four := F.NewVertex("4")
	five := F.NewVertex("5")

	Q := []ccfpq.Query{F.NewQuery(one, S)}

	engine := ccfpq.NewTIEngine(G, D, Q, F)
	engine.Run()

	hashGraph := SimpleToHashGraph(D.(*ds.SimpleGraph))
	hashGraph.Add(newEdge(two, S, two, false, nil, true))
	hashGraph.Add(newEdge(two, S, three, false, nil, true))
	hashGraph.Add(newEdge(two, S, four, false, nil, true))
	hashGraph.Add(newEdge(one, S, one, false, nil, true))
	hashGraph.Add(newEdge(two, S, two, false, nil, true))
	hashGraph.Add(newEdge(three, S, three, false, nil, true))

	rule0 := []ds.Vertex{S, a, S, b}
	rule1 := []ds.Vertex{S}
	hashFactory := NewHashFactory()
	item0 := newAugItem(rule0, hashFactory.NewEmptyPosets(len(rule0)))
	item1 := newAugItem(rule1, hashFactory.NewEmptyPosets(len(rule1)))
	item0.AddEdge(hashGraph.Get(one, S, one), 0)
	item0.AddEdge(hashGraph.Get(two, S, two), 0)
	item0.AddEdge(hashGraph.Get(three, S, three), 0)
	item0.AddEdge(hashGraph.Get(one, a, two), 1)
	item0.AddEdge(hashGraph.Get(one, a, three), 1)
	item0.AddEdge(hashGraph.Get(two, a, three), 1)
	item0.AddEdge(hashGraph.Get(two, S, two), 2)
	item0.AddEdge(hashGraph.Get(two, S, three), 2)
	item0.AddEdge(hashGraph.Get(two, S, four), 2)
	item0.AddEdge(hashGraph.Get(three, S, three), 2)
	item0.AddEdge(hashGraph.Get(three, b, four), 3)
	item0.AddEdge(hashGraph.Get(four, b, five), 3)
	item1.AddEdge(newEdge(one, S, one, false, nil, true), 0)
	item1.AddEdge(newEdge(two, S, two, false, nil, true), 0)
	item1.AddEdge(newEdge(three, S, three, false, nil, true), 0)

	want := hashFactory.NewAugItemSet(0)
	want.Add(item0)
	want.Add(item1)
	converter := newConverter(hashFactory)
	traceItems := engine.R.TraceItems(F)
	// augitems := converter.Convert(traceItems, hashGraph)
	AssertPanic(t, func() { converter.Convert(traceItems, hashGraph) },
		fmt.Sprintf("Should not run not implemented method."))
	/*Assert(t, augitems.Equals(want),
	fmt.Sprintf("Wrong augitems. Want %v, got %v.", want.String(),
		augitems.String()))*/
}
