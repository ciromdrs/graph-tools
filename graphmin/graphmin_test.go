package graphmin

import (
	ds "github.com/ciromdrs/graph-tools/data_structures"
	"testing"
)

func TestNewEdge(t *testing.T) {
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
