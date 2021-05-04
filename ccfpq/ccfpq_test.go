package ccfpq

import (
	"encoding/csv"
	ds "github.com/ciromdrs/graph-tools/datastructures"
	"os"
	"strconv"
	"testing"
)

func TestSimpleEngine(t *testing.T) {
	testDatabases(t, ds.SIMPLE_FACTORY)
}

func TestSliceEngine(t *testing.T) {
	testDatabases(t, ds.SLICE_FACTORY)
}

func testDatabases(t *testing.T, factorytype string) {
	script, err := os.Open("script.csv")
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
