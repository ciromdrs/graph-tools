package datastructures

import (
	"fmt"
	"strconv"
)

const SELECTOR = 1 // bit selector

const CELL_SIZE = 32 // size of vertexBitMask
type vertexBitMask int32

type (
	/* Data structures */
	subjectsBitMask struct { // TODO: support memory pre-allocation too
		data []*predicatesBitMask
		size int
	}

	predicatesBitMask struct { // TODO: support memory pre-allocation too
		data []*BitVertexSet
		size int
	}

	BitVertex struct {
		index    int
		id       vertexBitMask
		bitIndex int
	}

	BitVertexSet struct { // TODO: support memory pre-allocation too
		BaseVertexSet
		data []vertexBitMask
		size int
	}

	BitGraph struct {
		BaseGraph
		data *subjectsBitMask
		V    *BitVertexSet
		E    *BitVertexSet
	}

	indexobjects struct {
		index int
		value *BitVertexSet
	}

	indexpredicates struct {
		index int
		value *predicatesBitMask
	}
)

/* Functions and Methods */

func SliceIndexToVertex(index int) BitVertex {
	return NewBitVertex(index/CELL_SIZE, SELECTOR<<(index%CELL_SIZE))
}

/* BitGraph Functions and Methods */
func (g *BitGraph) VSize() int {
	return g.V.Size()
}

func (g *BitGraph) ESize() int {
	return g.E.Size()
}

func (g *BitGraph) Iterate() <-chan [3]Vertex {
	ch := make(chan [3]Vertex)
	go func() {
		for iv := range g.data.iterate() {
			subjindex := iv.index
			predicates := iv.value
			for iv := range predicates.iterate() {
				predindex := iv.index
				objects := iv.value
				for o := range objects.Iterate() {
					s := SliceIndexToVertex(subjindex)
					p := SliceIndexToVertex(predindex)
					ch <- [3]Vertex{s, p, o}
				}
			}
		}
		defer close(ch)
	}()
	return ch
}

func (g *BitGraph) AllNodes() VertexSet {
	return g.V
}

func (g *BitGraph) AllSubjects() <-chan Vertex {
	ch := make(chan Vertex)
	go func() {
		for iv := range g.data.iterate() {
			sindex := iv.index
			preds := iv.value
			for iv := range preds.iterate() {
				objects := iv.value
				if objects.Size() > 0 {
					ch <- SliceIndexToVertex(sindex)
				}
			}
		}
		defer close(ch)
	}()
	return ch
}

func (g *BitGraph) SubjectObjects(p Vertex) <-chan [2]Vertex {
	pp := p.(BitVertex)
	ch := make(chan [2]Vertex)
	go func() {
		for iv := range g.data.iterate() {
			sindex := iv.index
			predicates := iv.value
			s := SliceIndexToVertex(sindex)
			for o := range predicates.data[pp.IndexInSlice()].Iterate() {
				ch <- [2]Vertex{s, o}
			}
		}
		defer close(ch)
	}()
	return ch
}

func (g *BitGraph) PredicateObjects(s Vertex) <-chan [2]Vertex {
	panic("Not implemented yet.")
}

func (g *BitGraph) Objects(s, p Vertex) <-chan Vertex {
	ss := s.(BitVertex)
	pp := p.(BitVertex)
	return g.data.data[ss.IndexInSlice()].data[pp.IndexInSlice()].Iterate()
}

func (g *BitGraph) Show() {
	for triple := range g.Iterate() {
		s := triple[0].(BitVertex)
		p := triple[1].(BitVertex)
		o := triple[2].(BitVertex)
		fmt.Printf("(%d,%8b) (%d,%8b) (%d,%8b)\n", s.index, s.id, p.index, p.id, o.index, o.id)
	}
}

func (g *BitGraph) Contains(s, p, o Vertex) bool {
	ss := s.(BitVertex)
	pp := p.(BitVertex)
	oo := o.(BitVertex)
	return g.data.contains(ss, pp, oo)
}

func newSubjectsBitMask(VSize, ESize int) *subjectsBitMask {
	return &subjectsBitMask{}
}

func (m *subjectsBitMask) contains(s, p, o BitVertex) bool {
	if s.IndexInSlice() >= len(m.data) {
		return false
	}
	return m.data[s.IndexInSlice()].contains(p, o)
}

func (m *subjectsBitMask) add(s, p, o BitVertex) bool {
	if s.IndexInSlice() >= len(m.data) {
		m.expand(s.IndexInSlice() + 1)
	}
	old := m.data[s.IndexInSlice()].size
	added := m.data[s.IndexInSlice()].add(p, o)
	if (old == 0) && added {
		m.size++
	}
	return added
}

func (m *subjectsBitMask) iterate() <-chan indexpredicates {
	ch := make(chan indexpredicates)
	go func() {
		remaining := m.size
		for i := 0; (i < len(m.data)) && (remaining > 0); i++ {
			predicates := m.data[i]
			if predicates.size > 0 {
				remaining--
				ch <- indexpredicates{index: i, value: predicates}
			}
		}
		defer close(ch)
	}()
	return ch
}

func (m *subjectsBitMask) expand(length int) {
	if length < len(m.data) {
		panic("Cannot expand subjectsBitMask. `length` is too small.")
	}
	oldSize := len(m.data)
	new := make([]*predicatesBitMask, length)
	copy(new, m.data)
	m.data = new
	for i := oldSize; i < length; i++ {
		m.data[i] = newPredicatesBitMask()
	}
}

func newPredicatesBitMask() *predicatesBitMask {
	return &predicatesBitMask{}
}

func (m *predicatesBitMask) contains(p, o BitVertex) bool {
	if p.IndexInSlice() >= len(m.data) {
		return false
	}
	return m.data[p.IndexInSlice()].Contains(o)
}

func (m *predicatesBitMask) add(p, o BitVertex) bool {
	if p.IndexInSlice() >= len(m.data) {
		m.expand(p.IndexInSlice() + 1)
	}
	old := m.data[p.IndexInSlice()].Size()
	added := m.data[p.IndexInSlice()].Add(o)
	if (old == 0) && added {
		m.size++
	}
	return added
}

func (m *predicatesBitMask) expand(length int) {
	if length < len(m.data) {
		panic("Cannot expand predicatesBitMask. `length` is too small.")
	}
	oldSize := len(m.data)
	new := make([]*BitVertexSet, length)
	copy(new, m.data)
	m.data = new
	for i := oldSize; i < length; i++ {
		m.data[i] = NewBitVertexSet()
	}
}

func (m *predicatesBitMask) iterate() <-chan indexobjects {
	ch := make(chan indexobjects)
	go func() {
		remaining := m.size
		for i := 0; (i < len(m.data)) && (remaining > 0); i++ {
			objects := m.data[i]
			if objects.Size() > 0 {
				remaining--
				ch <- indexobjects{index: i, value: objects}
			}
		}
		defer close(ch)
	}()
	return ch
}

func (g *BitGraph) Add(s, p, o Vertex) bool {
	ss := s.(BitVertex)
	pp := p.(BitVertex)
	oo := o.(BitVertex)
	if g.data.add(ss, pp, oo) {
		g.size++
		g.V.Add(s)
		g.E.Add(p)
		g.V.Add(o)
		return true
	}
	return false
}

func (g *BitGraph) LoadTriple([]string) {
	panic("Not implemented yet.")
}

func (g *BitGraph) Predicates(s Vertex) <-chan Vertex {
	ch := make(chan Vertex)
	go func() {
		sindex := s.(BitVertex).IndexInSlice()
		for iv := range g.data.data[sindex].iterate() {
			pindex := iv.index
			objects := iv.value
			if objects.Size() > 0 {
				ch <- SliceIndexToVertex(pindex)
			}
		}
		defer close(ch)
	}()
	return ch
}

func NewBitGraph(name string, VSize, ESize int) *BitGraph {
	g := &BitGraph{
		BaseGraph: BaseGraph{
			name: name,
		},
		V:    NewBitVertexSet(),
		E:    NewBitVertexSet(),
		data: newSubjectsBitMask(VSize, ESize),
	}
	return g
}

/* BitVertex Functions and Methods */
func NewBitVertex(index int, id vertexBitMask) BitVertex {
	good := false
	shifts := -1
	for i := 0; i < CELL_SIZE; i++ {
		if id == (SELECTOR << i) {
			good = true
			shifts = i + 1
			break
		}
	}
	if !good {
		fmt.Print("BitVertex id ", id, " ")
		panic("is not valid. It should be a power of 2.")
	}
	return BitVertex{
		index:    index,
		id:       id,
		bitIndex: shifts,
	}
}

func (a BitVertex) Equals(other Vertex) bool {
	b := other.(BitVertex)
	return (a.index == b.index) && (a.id == b.id)
}

func (a BitVertex) String() string {
	return fmt.Sprintf("(%d, %"+strconv.Itoa(CELL_SIZE)+"b)", a.index, a.id)
}

func (a BitVertex) IndexInSlice() int {
	return (a.index * CELL_SIZE) + a.bitIndex - 1
}

func (a BitVertex) Label() string {
	return a.String() // TODO: fetch label from IDMap
}

/* vertexBitMask Functions and Methods */
func (a vertexBitMask) contains(b vertexBitMask) bool {
	return (a & b) == b
}

func (vertices vertexBitMask) Iterate() <-chan vertexBitMask {
	ch := make(chan vertexBitMask)
	go func() {
		for j := 0; j < CELL_SIZE; j++ {
			v := vertexBitMask(SELECTOR << j)
			if vertices.contains(v) {
				ch <- v
			}
		}
		defer close(ch)
	}()
	return ch
}

/* BitVertexSet Functions and Methods */
func NewBitVertexSet() *BitVertexSet {
	return &BitVertexSet{}
}

func (m *BitVertexSet) Contains(v Vertex) bool {
	b := v.(BitVertex)
	if b.index >= len(m.data) {
		return false
	}
	return (m.data[b.index] & b.id) == b.id
}

func (s *BitVertexSet) Add(new Vertex) bool {
	v := new.(BitVertex)
	added := false
	if !s.Contains(v) {
		added = true
		s.size++
		if v.index >= len(s.data) {
			s.expand(v.index + 1)
		}
		s.data[v.index] = s.data[v.index] | v.id
	}
	return added
}

func (s *BitVertexSet) expand(length int) {
	if length < len(s.data) {
		panic("Cannot expand BitVertexSet. `length` is too small.")
	}
	new := make([]vertexBitMask, length)
	copy(new, s.data)
	s.data = new
}

func (s *BitVertexSet) Iterate() <-chan Vertex {
	ch := make(chan Vertex)
	go func() {
		found := 0
		for i, vertices := range s.data {
			for j := 0; j < CELL_SIZE; j++ {
				v := vertexBitMask(SELECTOR << j)
				if vertices.contains(v) {
					found++
					ch <- NewBitVertex(i, v)
					if found == s.Size() {
						defer close(ch)
						return
					}
				}
			}
		}
		defer close(ch)
	}()
	return ch
}

func (s *BitVertexSet) Remove(v Vertex) bool {
	b := v.(BitVertex)
	removed := s.Contains(b)
	if removed {
		s.data[b.index] = s.data[b.index] &^ b.id
		s.size--
	}
	return removed
}

func (s *BitVertexSet) Equals(other VertexSet) bool {
	if s.Size() != other.Size() {
		return false
	}
	for v := range s.Iterate() {
		if !other.Contains(v) {
			return false
		}
	}
	return true
}

func (s *BitVertexSet) Size() int {
	return s.size
}

func (s *BitVertexSet) Update(toAdd VertexSet) int {
	c := 0
	for v := range toAdd.Iterate() {
		if s.Add(v) {
			c++
		}
	}
	return c
}

func (s *BitVertexSet) String() string {
	out := "{ "
	for e := range s.Iterate() {
		out += e.String() + " "
	}
	out += "}"
	return out
}
