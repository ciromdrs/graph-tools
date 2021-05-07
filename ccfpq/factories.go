package ccfpq

import (
	ds "github.com/ciromdrs/graph-tools/datastructures"
)

type (
	Factory interface {
		ds.Factory
		NewObserversSet() observersSet
		NewRelationsSet() relationsSet
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
