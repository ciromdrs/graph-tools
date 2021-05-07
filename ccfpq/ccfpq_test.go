package ccfpq

import (
	"encoding/csv"
	"fmt"
	ds "github.com/ciromdrs/graph-tools/datastructures"
	. "github.com/ciromdrs/graph-tools/util"
	"os"
	"strconv"
	"testing"
)

func TestSimpleEngine(t *testing.T) {
	if !testing.Short() {
		testDatabases(t, ds.SIMPLE_FACTORY)
	}
}

func TestSliceEngine(t *testing.T) {
	if !testing.Short() {
		testDatabases(t, ds.SLICE_FACTORY)
	}
}

func testDatabases(t *testing.T, factorytype string) {
	script, err := os.Open("testdata/script.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer script.Close()

	lines, err := csv.NewReader(script).ReadAll()
	if err != nil {
		t.Fatalf("Error reading script file: %v.", err)
	}

	for i, l := range lines {
		if len(l) < 3 {
			t.Errorf("Error reading line %d", i)
			continue
		}
		grammar := l[0]
		graph := l[1]
		expected, err := strconv.Atoi(l[2])
		if err != nil {
			t.Errorf("Error reading #results at line %d", i)
			continue
		}

		G, D, F := QuickLoad(grammar, graph, factorytype)
		Q := QueryAll(G, D)
		engine := NewTIEngine(G, D, Q, F)
		engine.Run()
		resCount := engine.CountResults()
		if resCount != expected {
			t.Errorf("Wrong results for grammar %s, graph %s: expected %d, "+
				"got %d", grammar, graph, expected, resCount)
		}
	}
}

func TestNonTerminalRelation(t *testing.T) {
	// TODO: Test SliceFactory
	G, D, F := QuickLoad("testdata/aSb.yrd", "testdata/graph-ti.txt", ds.SIMPLE_FACTORY)
	one := F.NewVertex("1")
	S := F.NewPredicate("S")
	a := F.NewPredicate("a")
	b := F.NewPredicate("b")
	Q := []pair{*newPair(one, S)}
	engine := NewTIEngine(G, D, Q, F)
	// TODO: test trace items content

	r := NewNonTerminalRelation(one, S, F)
	items := r.TraceItems(engine)
	Assert(t, items != nil && len(items) == 0,
		fmt.Sprintf("Wrong trace items for relation with nil rules. "+
			"Expected empty slice, got %v.", items))

	start := F.NewVertexSet()
	start.Add(one)
	var rule []ds.Vertex
	rule = append(rule, S, a, S, b)
	r.AddRule(start, rule[1:], engine)
	posets := make([]ds.VertexSet, len(rule))
	posets[0] = start
	for i := 1; i < len(posets); i++ {
		posets[i] = F.NewVertexSet()
	}
	want := F.NewTraceItem(rule, posets)

	items = r.TraceItems(engine)
	Assert(t, items != nil, "Expected non-nil TraceItems().")
	Assert(t, len(items) == 1, fmt.Sprintf("Expected length 1, got %v.",
		len(items)))
	Assert(t, items[0].Equals(want),
		fmt.Sprintf("Wrong trace item. Expected %v, got %v", want, items[0]))
}
