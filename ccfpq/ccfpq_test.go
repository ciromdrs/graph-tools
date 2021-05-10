package ccfpq

import (
	"encoding/csv"
	"fmt"
	ds "github.com/ciromdrs/graph-tools/datastructures"
	. "github.com/ciromdrs/graph-tools/util"
	"os"
	"strconv"
	"strings"
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

func TestSimpleNonTerminalRelation(t *testing.T) {
	testNonTerminalRelation(t, ds.SIMPLE_FACTORY)
}

func TestSliceNonTerminalRelation(t *testing.T) {
	testNonTerminalRelation(t, ds.SLICE_FACTORY)
}

func testNonTerminalRelation(t *testing.T, factoryType string) {
	G, D, F := QuickLoad("testdata/aSb.yrd", "testdata/graph-ti.txt",
		factoryType)
	one := F.NewVertex("1")
	S := F.NewPredicate("S")
	Q := []Query{F.NewQuery(one, S)}
	engine := NewTIEngine(G, D, Q, F)

	r := NewNonTerminalRelation(one, S, F)
	items := r.TraceItems(engine)
	Assert(t, items != nil && len(items) == 0,
		fmt.Sprintf("Wrong trace items for relation with nil rules. "+
			"Expected empty slice, got %v.", items))

	var want []*TraceItem
	want = append(want, traceItemFromString("S 1 a 2,3 S 2,3,4 b 4,5", F))
	want = append(want, traceItemFromString("S 1", F))
	want = append(want, traceItemFromString("S 2 a 3 S 3 b 4", F))
	want = append(want, traceItemFromString("S 2", F))
	want = append(want, traceItemFromString("S 3 a  S  b ", F))
	want = append(want, traceItemFromString("S 3", F))
	engine.Run()
	items = engine.R.TraceItems(F)
	Assert(t, items != nil, "Expected non-nil trace items.")
	Assert(t, len(items) == len(want),
		fmt.Sprintf("Expected %v trace items, got %v.", len(want), len(items)))
	for _, it := range want {
		found := false
		for _, other := range items {
			found = found || it.Equals(other)
		}
		msg := fmt.Sprintf("Wrong trace items. Item %v not found.\n",
			it.String())
		msg += "items = ["
		for _, it := range items {
			msg += it.String() + ", "
		}
		msg += "]\nwant = "
		for _, it := range want {
			msg += it.String() + ", "
		}
		msg += "]"
		Assert(t, found, msg)
	}
}

// traceItemFromString is a shorthand for building TraceItems.
func traceItemFromString(str string, factory Factory) *TraceItem {
	parts := strings.Split(str, " ")
	var rule []ds.Vertex
	var posets []ds.VertexSet
	for i, p := range parts {
		if i%2 == 0 {
			// symbol
			rule = append(rule, factory.NewPredicate(p))
		} else {
			// poset
			poset := factory.NewVertexSet()
			for _, e := range strings.Split(p, ",") {
				if e != "" {
					poset.Add(factory.NewVertex(e))
				}
			}
			posets = append(posets, poset)
		}
	}
	return factory.NewTraceItem(rule, posets)
}
