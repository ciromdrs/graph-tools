package data_structures

import (
    "fmt"
)

/* Data structures */
type SimpleVertexSet struct {
    BaseVertexSet
    mapset *MapSet
}

type PredicatesMap map[Vertex]*SimpleVertexSet
type SubjectsMap map[Vertex]PredicatesMap

type SimpleGraph struct {
    BaseGraph
    data SubjectsMap
}

type SimpleVertex struct {
    label string
}

/* Functions and Methods */

/* SimpleVertex Methods and Functions */
func NewSimpleVertex(label string) SimpleVertex {
    return SimpleVertex{label: label}
}

func (v SimpleVertex) ToString() string {
    return v.Label()
}

func (v SimpleVertex) Label() string {
    return v.label
}

func (v SimpleVertex) Equals(other Vertex) bool {
    return v.label == other.(SimpleVertex).label
}

/* SimpleVertexSet Methods and Functions */
func NewSimpleVertexSet() *SimpleVertexSet {
    return &SimpleVertexSet{
        mapset: NewMapSet(),
    }
}

func (s *SimpleVertexSet) Show() {
    s.mapset.Show()
}

func (s *SimpleVertexSet) Size() int {
    return s.mapset.Size()
}

func (s *SimpleVertexSet) Add(e Vertex) bool {
    return s.mapset.Add(e)
}

func (s *SimpleVertexSet) Contains(e Vertex) bool {
    return s.mapset.Contains(e)
}

func (s *SimpleVertexSet) Remove(e Vertex) bool {
    return s.mapset.Remove(e)
}

func (s *SimpleVertexSet) Iterate() <-chan Vertex {
    ch := make(chan Vertex)
    go func() {
        for e := range s.mapset.Iterate() {
            ch <- e.(Vertex)
        }
        defer close(ch)
    }()
    return ch
}

func (s *SimpleVertexSet) Update(toAdd VertexSet) int {
    c := 0
    for v := range toAdd.Iterate() {
        if s.Add(v) {
            c++
        }
    }
    return c
}

func (s *SimpleVertexSet) Equals(other VertexSet) bool {
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

/* SimpleGraph Methods and Functions */
func NewSimpleGraph() *SimpleGraph {
    g := &SimpleGraph{data: make(SubjectsMap)}
    return g
}

func (g *SimpleGraph) VSize() int {
    return g.AllNodes().Size()
}

func (g *SimpleGraph) ESize() int {
    all := NewSimpleVertexSet()
    for s := range g.data {
        for p := range g.Predicates(s) {
            all.Add(p)
        }
    }
    return all.Size()
}

func (g *SimpleGraph) Iterate() <-chan [3]Vertex {
    ch := make(chan [3]Vertex)
    go func() {
        for s, predicates := range g.data {
            for p, objects := range predicates {
                for o := range objects.Iterate() {
                    ch <- [3]Vertex{s, p, o}
                }
            }
        }
        defer close(ch)
    }()
    return ch
}

func (smap SubjectsMap) add(s, p, o Vertex) bool {
    if _, subjectExists := smap[s]; !subjectExists {
        smap[s] = make(PredicatesMap)
    }
    return smap[s].add(p, o)
}

func (pmap PredicatesMap) add(p, o Vertex) bool {
    if _, predicateExists := pmap[p]; !predicateExists {
        pmap[p] = NewSimpleVertexSet()
    }
    return pmap[p].Add(o)
}

func (g *SimpleGraph) LoadTriple(triple []string) {
    s := NewSimpleVertex(triple[0])
    p := NewSimpleVertex(triple[1])
    o := NewSimpleVertex(triple[2])
    g.Add(s, p, o)
}

func (g *SimpleGraph) Add(s, p, o Vertex) bool {
    added := g.data.add(s, p, o)
    if added {
        g.size++
    }
    return added
}

func (g *SimpleGraph) Show() {
    for triple := range g.Iterate() {
        s := triple[0]
        p := triple[1]
        o := triple[2]
        fmt.Printf("%s %s %s\n", s, p, o)
    }
}

func (g *SimpleGraph) PredicateObjects(s Vertex) <-chan [2]Vertex {
    ch := make(chan [2]Vertex)
    go func() {
        for p, _ := range g.data[s] {
            for o := range g.data[s][p].Iterate() {
                ch <- [2]Vertex{p, o}
            }
        }
        defer close(ch)
    }()
    return ch
}

func (g *SimpleGraph) Objects(s, p Vertex) <-chan Vertex {
    return g.data[s][p].Iterate()
}

func (g *SimpleGraph) SubjectObjects(p Vertex) <-chan [2]Vertex {
    ch := make(chan [2]Vertex)
    go func() {
        for s, _ := range g.data {
            for o := range g.data[s][p].Iterate() {
                ch <- [2]Vertex{s, o}
            }
        }
        defer close(ch)
    }()
    return ch
}

func (g *SimpleGraph) AllNodes() VertexSet {
    all := NewSimpleVertexSet()
    for s, _ := range g.data {
        all.Add(s)
        for p, _ := range g.data[s] {
            for o := range g.data[s][p].Iterate() {
                all.Add(o.(Vertex))
            }
        }
    }
    return all
}

func (g *SimpleGraph) AllSubjects() <-chan Vertex {
    ch := make(chan Vertex)
    go func() {
        for s, _ := range g.data {
            ch <- s
        }
        defer close(ch)
    }()
    return ch
}

func (g *SimpleGraph) Predicates(s Vertex) <-chan Vertex {
    ch := make(chan Vertex)
    go func() {
        for p, _ := range g.data[s] {
            ch <- p
        }
        defer close(ch)
    }()
    return ch
}

func (g *SimpleGraph) Contains(s, p, o Vertex) bool {
    if _, ok := g.data[s]; ok {
        if _, ok := g.data[s][p]; ok {
            if g.data[s][p].Contains(o) {
                return true
            }
        }
    }
    return false
}

func Union(g1, g2 SimpleGraph) *SimpleGraph {
    u := NewSimpleGraph()

    for s, predicates := range g1.data {
        for p, objects := range predicates {
            for o := range objects.Iterate() {
                u.Add(s, p, o.(Vertex))
            }
        }
    }

    for s, predicates := range g2.data {
        for p, objects := range predicates {
            for o := range objects.Iterate() {
                u.Add(s, p, o.(Vertex))
            }
        }
    }
    return u
}
