package ccfpq

import (
	ds "github.com/ciromdrs/graph-tools/datastructures"
)

type (
	Factory interface {
		ds.Factory
		NewObserversSet() observersSet
		NewRelationsSet() relationsSet
		NewTraceItem([]ds.Vertex, []ds.VertexSet) *TraceItem
		NewQuery(node, label ds.Vertex) Query
	}

	baseFactory struct{}

	SimpleFactory struct {
		baseFactory
		ds.SimpleFactory
	}

	SliceFactory struct {
		baseFactory
		ds.SliceFactory
	}
)

// NewQuery creates a Query object.
func (f *baseFactory) NewQuery(node, label ds.Vertex) Query {
	return newQuery(node, label)
}

// newTraceItem creates a TraceItem object.
func (f *baseFactory) NewTraceItem(rule []ds.Vertex, posets []ds.VertexSet) *TraceItem {
	ti := &TraceItem{
		Rule:   rule,
		Posets: posets,
	}
	return ti
}

/* SimpleFactory Functions and Methods */
func NewSimpleFactory() *SimpleFactory {
	return &SimpleFactory{}
}

func (f *SimpleFactory) NewObserversSet() observersSet {
	return newMapObserversSet(f.VSize*f.ESize, f)
}

func (f *SimpleFactory) NewRelationsSet() relationsSet {
	return newMapRelationsSet(f.VSize, f.ESize)
}

/* SliceFactory Functions and Methods */
func NewSliceFactory(VSize, ESize int) *SliceFactory {
	slf := ds.NewSliceFactory(VSize, ESize)
	new := &SliceFactory{SliceFactory: *slf}
	new.V = ds.NewIDMap(VSize)
	new.E = ds.NewIDMap(ESize)
	return new

}

func (f *SliceFactory) NewObserversSet() observersSet {
	return newSliceObserversSet(f.VSize, f.ESize, f)
}

func (f *SliceFactory) NewRelationsSet() relationsSet {
	return newSliceRelationsSet(f.VSize, f.ESize)
}
