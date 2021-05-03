package data_structures

import (
    "fmt"
    "io/ioutil"
    "rdf-ccfpq/go/util"
    "strings"
)

/* Interfaces */
type Vertex interface {
    ToString() string
    Label() string
    Equals(other Vertex) bool
}

type VertexSet interface {
    Add(Vertex) bool
    Contains(Vertex) bool
    Equals(VertexSet) bool
    Iterate() <-chan Vertex
    Remove(Vertex) bool
    Show()
    Size() int
    Update(VertexSet) int
}

type Graph interface {
    Add(Vertex, Vertex, Vertex) bool
    AllSubjects() <-chan Vertex
    AllNodes() VertexSet
    Contains(Vertex, Vertex, Vertex) bool
    Iterate() <-chan [3]Vertex
    Name() string
    Objects(Vertex, Vertex) <-chan Vertex
    Predicates(Vertex) <-chan Vertex
    PredicateObjects(Vertex) <-chan [2]Vertex
    SetName(string)
    SubjectObjects(Vertex) <-chan [2]Vertex
    Show()
    Size() int
    VSize() int
    ESize() int
}

/* Data Structures */
type SuperVertex struct {
    Vertex   Vertex
    Vertices VertexSet
}

type BaseVertexSet struct{}

/* BaseGraph Functions and Methods */
type BaseGraph struct {
    name string
    size int
}

func (g *BaseGraph) Name() string {
    return g.name
}

func (g *BaseGraph) SetName(name string) {
    g.name = name
}

func (g *BaseGraph) Size() int {
    return g.size
}

func (g *BaseGraph) Show() {
    for triple := range g.Iterate() {
        fmt.Println(triple)
    }
}

func (g *BaseGraph) Load(path string) {
    data, err := ioutil.ReadFile(path)
    if err != nil {
        panic("Error openning file: " + path + "\n")
    }
    lines := strings.Split(string(data), "\n")
    g.name = util.GetFileName(path)
    for _, line := range lines {
        if line != "" {
            triple := strings.Split(line, " ")
            if len(triple) != 3 {
                panic("Error reading line: " + string(line))
            }
            g.LoadTriple(triple)
        }
    }
}

func (g *BaseGraph) Iterate() <-chan [3]Vertex {
    panic("Abstract method")
}

func (g *BaseGraph) LoadTriple(triple []string) {
    panic("Abstract method")
}

/* BaseVertexSet Functions and Methods */
func (s *BaseVertexSet) Equals(other VertexSet) bool {
    panic("Abstract method.")
}

func (s *BaseVertexSet) Show() {
    fmt.Print("{ ")
    for v := range s.Iterate() {
        fmt.Print(v.ToString(), " ")
    }
    fmt.Println("}")
}

func (s *BaseVertexSet) Add(v Vertex) bool {
    panic("Abstract method.")
}

func (s *BaseVertexSet) Size() int {
    panic("Abstract method.")
}

func (s *BaseVertexSet) Update(toAdd VertexSet) int {
    panic("Abstract method.")
}

func (s *BaseVertexSet) Iterate() <-chan Vertex {
    panic("Abstract method.")
}

func ChanToSet(ch <-chan Vertex, s VertexSet) {
    for e := range ch {
        s.Add(e)
    }
}

func ChanToPairs(ch <-chan [2]Vertex) [][2]Vertex {
    var pairs [][2]Vertex
    for e := range ch {
        pairs = append(pairs, e)
    }
    return pairs
}

func ChanToTriples(ch <-chan [3]Vertex) [][3]Vertex {
    var triples [][3]Vertex
    for e := range ch {
        triples = append(triples, e)
    }
    return triples
}

/* SuperVertex Functions and Methods */
func NewSuperVertex(vertex Vertex, vertices VertexSet) SuperVertex {
    return SuperVertex{
        Vertex: vertex,
        Vertices: vertices,
    }
}

func (v SuperVertex) ToString() string {
    return v.Vertex.ToString()
}

func (v SuperVertex) Label() string {
    return v.Vertex.Label()
}

func (v SuperVertex) Equals(other Vertex) bool {
    v2, ok := other.(SuperVertex)
    if !ok {
        return false
    }
    return v.Vertices.Equals(v2.Vertices)
}
