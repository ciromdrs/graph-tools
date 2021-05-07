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
	}

	SimpleFactory struct {
		ds.SimpleFactory
	}

	SliceFactory struct {
		ds.SliceFactory
	}
)

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

// NewTraceItem returns a new TraceItem object
func (f *SimpleFactory) NewTraceItem(rule []ds.Vertex,
	posets []ds.VertexSet) *TraceItem {
	return newTraceItem(rule, posets)
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

// NewTraceItem returns a new TraceItem object
func (f *SliceFactory) NewTraceItem(rule []ds.Vertex,
	posets []ds.VertexSet) *TraceItem {
	return newTraceItem(rule, posets)
}

// newTraceItem creates a TraceItem object. It is a behavior common to Simple
// and Slice factories
func newTraceItem(rule []ds.Vertex, posets []ds.VertexSet) *TraceItem {
	ti := &TraceItem{
		rule:   rule,
		posets: posets,
	}
	return ti
}
