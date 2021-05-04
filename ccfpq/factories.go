package ccfpq

import (
	"fmt"
	ds "github.com/ciromdrs/graph-tools/datastructures"
)

var f Factory

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

/* BaseFactory Methods */
func SetFactory(factoryType string, VAlloc, EAlloc int) {
	switch factoryType {
	case ds.SIMPLE_FACTORY:
		f = NewSimpleFactory()
	case ds.SLICE_FACTORY:
		f = NewSliceFactory(VAlloc, EAlloc)
	default:
		panic(fmt.Sprintf("Invalid factory type %s", factoryType))
	}
}

/* SimpleFactory Functions and Methods */
func NewSimpleFactory() *SimpleFactory {
	return &SimpleFactory{}
}

func (f *SimpleFactory) NewObserversSet() observersSet {
	return newMapObserversSet(f.VSize * f.ESize)
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
	return newSliceObserversSet(f.VSize, f.ESize)
}

func (f *SliceFactory) NewRelationsSet() relationsSet {
	return newSliceRelationsSet(f.VSize, f.ESize)
}
