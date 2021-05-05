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

		G, D := QuickLoad(grammar, graph, factorytype)
		Q := QueryAll(G, D)
		R, _, _ := Run(D, G, Q)
		resCount := 0
		for _, p := range Q {
			var node ds.Vertex
			if super, isSuperVertex := p.node.(ds.SuperVertex); isSuperVertex {
				node = super.Vertex
			} else {
				node = p.node
			}
			resCount += R.get(node, p.symbol).Objects().Size()
		}
		if resCount != expected {
			t.Errorf("Wrong results for grammar %s, graph %s: expected %d, "+
				"got %d", grammar, graph, expected, resCount)
		}
	}
}

func TestNonTerminalRelation(t *testing.T) {
	SetFactory(ds.SIMPLE_FACTORY, 0, 0)
	x := f.NewVertex("x")
	a := f.NewVertex("a")
	b := f.NewVertex("b")
	S := f.NewVertex("S")

	NEW = ds.NewMapSet()

	r := NewNonTerminalRelation(x, S)
	items := r.TraceItems()
	Assert(t, items != nil && len(items) == 0,
		fmt.Sprintf("Wrong trace items for relation with nil rules. "+
			"Expected empty slice, got %v.", items))

	start := f.NewVertexSet()
	start.Add(x)
	var rule []ds.Vertex
	rule = append(rule, S, a, S, b)
	r.AddRule(start, rule[1:], nil)
	want := NewTraceItem(start, rule)

	items = r.TraceItems()
	Assert(t, items != nil, "Expected non-nil TraceItems().")
	Assert(t, len(items) == 1, fmt.Sprintf("Expected length 1, got %v.",
		len(items)))
	Assert(t, items[0].Equals(want),
		fmt.Sprintf("Wrong trace item. Expected %v, got %v", want, items[0]))
}
