package ccfpq

import (
	"fmt"
	ds "github.com/ciromoraismedeiros/graph-tools/data_structures"
	"strconv"
)

const (
	SIMPLE_FACTORY = "simple_factory"
	SLICE_FACTORY  = "slice_factory"
)

var f Factory

type (
	Factory interface {
		NewGraph(string) ds.Graph
		NewVertex(string) ds.Vertex
		NewSuperVertex(ds.VertexSet) ds.SuperVertex
		NewPredicate(string) ds.Vertex
		NewVertexSet() ds.VertexSet
		NewObserversSet() observersSet
		NewRelationsSet() relationsSet
		Reset()
		NewSet() ds.Set
		nextSuperVertex() string
		NewMap(int) ds.Map
		Type() string
	}

	IDMap struct {
		size     int
		strtovex map[string]ds.BitVertex
		vextostr map[ds.BitVertex]string
	}

	BaseFactory struct {
		superVertexCount int
		VSize            int
		ESize            int
	}

	SimpleFactory struct {
		BaseFactory
	}

	SliceFactory struct {
		BaseFactory
		V *IDMap
		E *IDMap
	}
)

/* BaseFactory Methods */
func SetFactory(factoryType string, VAlloc, EAlloc int) {
	switch factoryType {
	case SIMPLE_FACTORY:
		f = NewSimpleFactory()
	case SLICE_FACTORY:
		f = NewSliceFactory(VAlloc, EAlloc)
	default:
		panic(fmt.Sprintf("Invalid factory type %s", factoryType))
	}
}

func (f *BaseFactory) nextSuperVertex() string {
	f.superVertexCount++
	return "super" + strconv.Itoa(f.superVertexCount)
}

/* SimpleFactory Functions and Methods */
func NewSimpleFactory() *SimpleFactory {
	return &SimpleFactory{}
}

func (f *SimpleFactory) NewGraph(name string) ds.Graph {
	g := ds.NewSimpleGraph()
	g.SetName(name)
	return g
}

func (f *SimpleFactory) NewVertex(label string) ds.Vertex {
	return ds.NewSimpleVertex(label)
}

func (f *SimpleFactory) NewPredicate(label string) ds.Vertex {
	return ds.NewSimpleVertex(label)
}

func (f *SimpleFactory) NewSuperVertex(vertices ds.VertexSet) ds.SuperVertex {
	vertex := f.NewVertex(f.nextSuperVertex())
	return ds.NewSuperVertex(vertex, vertices)
}

func (f *SimpleFactory) NewVertexSet() ds.VertexSet {
	return ds.NewSimpleVertexSet()
}

func (f *SimpleFactory) Reset() {
	f = NewSimpleFactory()
}

func (f *SimpleFactory) NewObserversSet() observersSet {
	return newMapObserversSet(f.VSize * f.ESize)
}

func (f *SimpleFactory) NewRelationsSet() relationsSet {
	return newMapRelationsSet(f.VSize, f.ESize)
}

func (f *SimpleFactory) NewSet() ds.Set {
	return ds.NewMapSet()
}

func (f *SimpleFactory) NewMap(preallocate int) ds.Map {
	return ds.NewSimpleMap()
}

func (f *SimpleFactory) Type() string {
	return SIMPLE_FACTORY
}

/* IDMap Functions and Methods */
func newIDMap(size int) *IDMap {
	return &IDMap{
		size:     0,
		strtovex: make(map[string]ds.BitVertex, size),
		vextostr: make(map[ds.BitVertex]string, size),
	}
}

func (m *IDMap) Add(id string) bool {
	added := false
	if _, exists := m.strtovex[id]; !exists {
		added = true
		v := ds.SliceIndexToVertex(m.size)
		m.size++ // increase after creating the new vertex
		m.strtovex[id] = v
		m.vextostr[v] = id
	}
	return added
}

func (m *IDMap) GetVertex(id string) (ds.BitVertex, *ds.NotFoundError) {
	v, exists := m.strtovex[id]
	if exists {
		return v, nil
	}
	return ds.BitVertex{}, &ds.NotFoundError{}
}

/* SliceFactory Functions and Methods */
func NewSliceFactory(VSize, ESize int) *SliceFactory {
	return &SliceFactory{
		BaseFactory: BaseFactory{
			VSize: VSize,
			ESize: ESize,
		},
		V: newIDMap(VSize),
		E: newIDMap(ESize),
	}
}

func (f *SliceFactory) NewGraph(name string) ds.Graph {
	return ds.NewBitGraph(name, f.VSize, f.ESize)
}

func (f *SliceFactory) NewVertex(label string) ds.Vertex {
	v, err := f.V.GetVertex(label)
	if err == nil {
		return v
	}
	if f.V.size < (f.VSize) {
		f.V.Add(label)
		v, _ := f.V.GetVertex(label)
		return v
	}
	panic("Capacity exceeded.")
}

func (f *SliceFactory) NewPredicate(label string) ds.Vertex {
	v, err := f.E.GetVertex(label)
	if err == nil {
		return v
	}
	if f.E.size < (f.ESize) {
		f.E.Add(label)
		v, _ = f.E.GetVertex(label)
		return v
	}
	panic("Capacity exceeded.")
}

func (f *SliceFactory) NewSuperVertex(vertices ds.VertexSet) ds.SuperVertex {
	vertex := f.NewVertex(f.nextSuperVertex())
	return ds.NewSuperVertex(vertex, vertices)
}

func (f *SliceFactory) NewVertexSet() ds.VertexSet {
	return ds.NewBitVertexSet()
}

func (f *SliceFactory) NewObserversSet() observersSet {
	return newSliceObserversSet(f.VSize, f.ESize)
}

func (f *SliceFactory) NewRelationsSet() relationsSet {
	return newSliceRelationsSet(f.VSize, f.ESize)
}

func (f *SliceFactory) Reset() {
	f = NewSliceFactory(f.VSize, f.ESize)
}

func (f *SliceFactory) NewSet() ds.Set {
	return ds.NewSliceSet(0)
}

func (f *SliceFactory) NewMap(preallocate int) ds.Map {
	return ds.NewSliceMap(preallocate)
}

func (f *SliceFactory) Type() string {
	return SLICE_FACTORY
}
