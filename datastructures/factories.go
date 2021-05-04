package datastructures

import (
	"fmt"
	"strconv"
)

const (
	SIMPLE_FACTORY = "simple_factory"
	SLICE_FACTORY  = "slice_factory"
)

var f Factory

type (
	Factory interface {
		NewGraph(string) Graph
		NewVertex(string) Vertex
		NewSuperVertex(VertexSet) SuperVertex
		NewPredicate(string) Vertex
		NewVertexSet() VertexSet
		Reset()
		NewSet() Set
		nextSuperVertex() string
		NewMap(int) Map
		Type() string
	}

	IDMap struct {
		size     int
		strtovex map[string]BitVertex
		vextostr map[BitVertex]string
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

func (f *SimpleFactory) NewGraph(name string) Graph {
	g := NewSimpleGraph()
	g.SetName(name)
	return g
}

func (f *SimpleFactory) NewVertex(label string) Vertex {
	return NewSimpleVertex(label)
}

func (f *SimpleFactory) NewPredicate(label string) Vertex {
	return NewSimpleVertex(label)
}

func (f *SimpleFactory) NewSuperVertex(vertices VertexSet) SuperVertex {
	vertex := f.NewVertex(f.nextSuperVertex())
	return NewSuperVertex(vertex, vertices)
}

func (f *SimpleFactory) NewVertexSet() VertexSet {
	return NewSimpleVertexSet()
}

func (f *SimpleFactory) Reset() {
	f = NewSimpleFactory()
}

func (f *SimpleFactory) NewSet() Set {
	return NewMapSet()
}

func (f *SimpleFactory) NewMap(preallocate int) Map {
	return NewSimpleMap()
}

func (f *SimpleFactory) Type() string {
	return SIMPLE_FACTORY
}

/* IDMap Functions and Methods */
func NewIDMap(size int) *IDMap {
	return &IDMap{
		size:     0,
		strtovex: make(map[string]BitVertex, size),
		vextostr: make(map[BitVertex]string, size),
	}
}

func (m *IDMap) Add(id string) bool {
	added := false
	if _, exists := m.strtovex[id]; !exists {
		added = true
		v := SliceIndexToVertex(m.size)
		m.size++ // increase after creating the new vertex
		m.strtovex[id] = v
		m.vextostr[v] = id
	}
	return added
}

func (m *IDMap) GetVertex(id string) (BitVertex, *NotFoundError) {
	v, exists := m.strtovex[id]
	if exists {
		return v, nil
	}
	return BitVertex{}, &NotFoundError{}
}

/* SliceFactory Functions and Methods */
func NewSliceFactory(VSize, ESize int) *SliceFactory {
	return &SliceFactory{
		BaseFactory: BaseFactory{
			VSize: VSize,
			ESize: ESize,
		},
		V: NewIDMap(VSize),
		E: NewIDMap(ESize),
	}
}

func (f *SliceFactory) NewGraph(name string) Graph {
	return NewBitGraph(name, f.VSize, f.ESize)
}

func (f *SliceFactory) NewVertex(label string) Vertex {
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

func (f *SliceFactory) NewPredicate(label string) Vertex {
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

func (f *SliceFactory) NewSuperVertex(vertices VertexSet) SuperVertex {
	vertex := f.NewVertex(f.nextSuperVertex())
	return NewSuperVertex(vertex, vertices)
}

func (f *SliceFactory) NewVertexSet() VertexSet {
	return NewBitVertexSet()
}

func (f *SliceFactory) Reset() {
	f = NewSliceFactory(f.VSize, f.ESize)
}

func (f *SliceFactory) NewSet() Set {
	return NewSliceSet(0)
}

func (f *SliceFactory) NewMap(preallocate int) Map {
	return NewSliceMap(preallocate)
}

func (f *SliceFactory) Type() string {
	return SLICE_FACTORY
}
