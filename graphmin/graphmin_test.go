package graphmin

import (
	"fmt"
	"github.com/ciromdrs/graph-tools/ccfpq"
	ds "github.com/ciromdrs/graph-tools/datastructures"
	. "github.com/ciromdrs/graph-tools/util"
	"testing"
)

func TestEdge(t *testing.T) {
	s := ds.NewSimpleVertex("s")
	p := ds.NewSimpleVertex("p")
	o := ds.NewSimpleVertex("o")
	_ = newEdge(s, p, o, false, nil, true)
	// TODO: Write actual tests for Edges
}

func TestHashGraph(t *testing.T) {
	var g Graph = newHashGraph()
	Assert(t, g.Size() == 0,
		fmt.Sprintf("Wrong g.Size(). Expected 0, got %v.", g.Size()))

	s := ds.NewSimpleVertex("s")
	p := ds.NewSimpleVertex("p")
	o := ds.NewSimpleVertex("o")
	e1 := newEdge(s, p, o, false, nil, true)

	Assert(t, !g.Contains(s, p, o), "Should not contain edge.")
	Assert(t, g.Size() == 0,
		fmt.Sprintf("Wrong g.Size(). Expected 0, got %v.", g.Size()))

	g.Add(e1)
	Assert(t, g.Size() == 1,
		fmt.Sprintf("Wrong g.Size(). Expected 1, got %v.", g.Size()))
	Assert(t, g.Contains(s, p, o), "Should contain edge.")

	g.Remove(s, p, o)
	Assert(t, !g.Contains(s, p, o), "Should not contain removed edge.")
	Assert(t, g.Size() == 0,
		fmt.Sprintf("Wrong g.Size(). Expected 0, got %v.", g.Size()))

	e2 := newEdge(s, p, ds.NewSimpleVertex("o2"), false, nil, true)
	g.Add(e1)
	g.Add(e2)
	Assert(t, g.Size() == 2,
		fmt.Sprintf("Wrong g.Size(). Expected 2, got %v.", g.Size()))
	for e := range g.Iterate() {
		Assert(t, e == e1 || e == e2, "Wrong edges in iteration.")
	}

	g2 := newHashGraph()
	g2.Add(e1)
	g2.Add(e2)
	Assert(t, g.Equals(g2), "Graphs should be equal.")
	g2.Remove(s, p, o)
	Assert(t, !g.Equals(g2), "Graphs should not be equal.")
	g2.Add(newEdge(o, p, s, false, nil, true))
	Assert(t, !g.Equals(g2), "Graphs should not be equal.")

}

func TestSimpleToHashGraphConversion(t *testing.T) {
	var databases []string
	databases = append(databases, "../ccfpq/testdata/atom-primitive.txt")
	if !testing.Short() {
		databases = append(databases,
			"../ccfpq/testdata/biomedical-mesure-primitive.txt",
			"../ccfpq/testdata/foaf.txt",
			"../ccfpq/testdata/funding.txt",
			"../ccfpq/testdata/generations.txt",
			"../ccfpq/testdata/people_pets.txt",
			"../ccfpq/testdata/pizza.txt",
			"../ccfpq/testdata/skos.txt",
			"../ccfpq/testdata/travel.txt",
			"../ccfpq/testdata/univ-bench.txt",
			"../ccfpq/testdata/wine.txt",
		)
	}

	for _, graph := range databases {
		_, simpleGraph, _ := ccfpq.QuickLoad("../ccfpq/testdata/sc.yrd",
			graph, ds.SIMPLE_FACTORY)
		hashGraph := SimpleToHashGraph(simpleGraph.(*ds.SimpleGraph))
		Assert(t, hashGraph.Size() == simpleGraph.Size(),
			fmt.Sprintf("Wrong hashGraph.Size(). Expected %v, got %v.",
				simpleGraph.Size(),
				hashGraph.Size()))
		for t1 := range simpleGraph.Iterate() {
			found := false
			for e := range hashGraph.Iterate() {
				t2 := e.triple
				if t1[0] == t2.s && t1[1] == t2.X && t1[2] == t2.o {
					found = true
					break
				}
			}
			Assert(t, found, fmt.Sprintf("Missing triple %v in hashGraph.", t1))
		}
	}
}
