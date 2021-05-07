package ccfpq

import (
	ds "github.com/ciromdrs/graph-tools/datastructures"
)

type (
	Factory interface {
		ds.Factory
		NewObserversSet() observersSet
		NewRelationsSet() relationsSet
		NewTraceItem(ds.VertexSet, []ds.Vertex) *TraceItem
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
func (f *SimpleFactory) NewTraceItem(start ds.VertexSet,
	rule []ds.Vertex) *TraceItem {
	return newTraceItem(start, rule, f)
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
func (f *SliceFactory) NewTraceItem(start ds.VertexSet,
	rule []ds.Vertex) *TraceItem {
	return newTraceItem(start, rule, f)
}

// newTraceItem creates a TraceItem object. It is a behavior common to Simple
// and Slice factories
func newTraceItem(start ds.VertexSet, rule []ds.Vertex, f Factory) *TraceItem {
	ti := &TraceItem{
		rule:   rule,
		posets: make([]ds.VertexSet, len(rule)),
	}
	ti.posets[0] = start
	for i := 1; i < len(ti.posets); i++ {
		ti.posets[i] = f.NewVertexSet()
	}
	return ti
}
