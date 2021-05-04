package ccfpq

import (
	"fmt"
	ds "github.com/ciromdrs/graph-tools/datastructures"
	"github.com/ciromdrs/graph-tools/util"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
	"time"
)

type (
	NodeSet struct {
		nodes    ds.VertexSet
		new      ds.VertexSet
		prev     *Symbol
		next     *Symbol
		relation Relation
	}

	pair struct {
		node   ds.Vertex
		symbol ds.Vertex
	}
)

var (
	R   relationsSet
	NEW ds.Set
	O   observersSet
)

/* NodeSet Methods and Functions */
func NewNodeSet() *NodeSet {
	return &NodeSet{
		nodes: f.NewVertexSet(),
		new:   f.NewVertexSet(),
	}
}

func (ns *NodeSet) Mark() {
	for k := range ns.new.Iterate() {
		ns.nodes.Add(k)
	}
	ns.new = f.NewVertexSet()
}

func GetOrCreate(node, label ds.Vertex, G *Grammar) Relation {
	r := R.get(node, label)
	if r == nil {
		if label.Equals(epsilon()) {
			panic("Should not create empty relations")
			// nested expressions
		} else if isNestedExp(label.Label()) {
			r2 := NewNestedRelation(node, label)
			subexp := strings.TrimPrefix(label.Label(), OPEN)
			subexp = strings.TrimSuffix(subexp, CLOSE)
			r2.SetRule(parseExp(subexp))
			r = r2
			// nonterminals
		} else if G.NonTerm.Contains(label) {
			r2 := NewNonTerminalRelation(node, label)
			startVertices := f.NewVertexSet()
			startVertices.Add(node)
			for _, rule := range G.Rules[label] {
				r2.AddRule(startVertices, rule, G)
			}
			r = r2
			// terminals
		} else if G.Alphabet.Contains(label) {
			// do nothing (do not delete this if clause)
			// expressions
		} else {
			r2 := NewExpressionRelation(node, label)
			startVertices := f.NewVertexSet()
			startVertices.Add(node)
			r2.SetRule(startVertices, parseExp(label.Label()))
			r = r2
		}
		if r != nil {
			R.set(node, label, r) // do not add nil relations
		}
	}
	return r
}

/* pair Methods and Functions */
func newPair(node, symbol ds.Vertex) *pair {
	return &pair{
		node:   node,
		symbol: symbol,
	}
}

/* Graph Parsing Functions */
func BuildBaseGraph(graph ds.Graph, grammar *Grammar) ds.Graph {
	D := f.NewGraph(graph.Name())
	D.SetName(graph.Name())
	for pElem := range grammar.Alphabet.Iterate() {
		p := pElem.(ds.Vertex)
		for pair := range graph.SubjectObjects(p) {
			s := pair[0]
			o := pair[1]
			D.Add(s, p, o)
		}
	}
	return D
}

func AddNew(nodeSet *NodeSet) {
	NEW.Add(nodeSet)
}

func processNew(nodeSet *NodeSet, G *Grammar) {
	if symbol := nodeSet.next; symbol != nil {
		// update symbol's object set and NEW
		new := f.NewVertexSet()
		for a := range nodeSet.new.Iterate() {
			var destinations ds.VertexSet = f.NewVertexSet()
			if symbol.predicate.Equals(epsilon()) {
				destinations.Add(a)
			} else {
				if r := GetOrCreate(a, symbol.predicate, G); r != nil {
					destinations = r.Objects()
				}
				if !G.Alphabet.Contains(symbol.predicate) {
					O.add(a, symbol.predicate, symbol.objNodeSet)
				}
			}
			for b := range destinations.Iterate() {
				if !symbol.objNodeSet.nodes.Contains(b) {
					new.Add(b)
				}
			}
		}
		if new.Size() > 0 {
			symbol.objNodeSet.new.Update(new)
			NEW.Add(symbol.objNodeSet)
		}
	} else {
		// updates relation's objects and notifies O
		new := f.NewVertexSet()
		if nodeSet.relation.IsNested() && nodeSet.relation.Objects().Size() == 0 {
			new.Add(nodeSet.relation.Node())
		} else {
			for a := range nodeSet.new.Iterate() {
				if !nodeSet.relation.Objects().Contains(a) {
					new.Add(a)
				}
			}
		}
		if new.Size() > 0 {
			nodeSet.relation.AddObjects(new)
			for _, o := range O.get(nodeSet.relation.Node(),
				nodeSet.relation.Label()) {
				for n := range new.Iterate() {
					if !o.new.Contains(n) {
						o.new.Add(n)
						AddNew(o)
					}
				}
			}
		}
	}
	nodeSet.Mark()
}

func pickAndRemove() *NodeSet {
	for e := range NEW.Iterate() {
		ns := e.(*NodeSet)
		NEW.Remove(ns)
		return ns
	}
	return nil
}

func Run(D ds.Graph, G *Grammar, Q []pair) (relationsSet, time.Duration, uint64) {
	os.Setenv("GOGC", "off")
	// allocateMemory(15000)
	runtime.GC()
	// var startmem runtime.MemStats
	// runtime.ReadMemStats(&startmem)

	startusr, startsys := util.GetTime()

	R = f.NewRelationsSet()
	O = f.NewObserversSet()
	NEW = ds.NewMapSet()

	// Creating terminal relations
	for s := range D.AllSubjects() {
		for p := range D.Predicates(s) {
			objects := f.NewVertexSet()
			ds.ChanToSet(D.Objects(s, p), objects)
			R.set(s, p, NewTerminalRelation(s, p, objects))
		}
	}

	// Initializing non-terminal relations
	for _, p := range Q {
		var r *NonTerminalRelation
		var node ds.Vertex
		label := p.symbol
		var startVertices ds.VertexSet

		if super, isSuperVertex := p.node.(ds.SuperVertex); isSuperVertex {
			node = super.Vertex
			startVertices = super.Vertices
		} else {
			node = p.node
			startVertices = f.NewVertexSet()
			startVertices.Add(node)
		}
		r = NewNonTerminalRelation(node, label)
		for _, rule := range G.Rules[label] {
			r.AddRule(startVertices, rule, G)
		}
		R.set(node, label, r)
	}

	for NEW.Size() > 0 {
		newNodeSet := pickAndRemove()
		processNew(newNodeSet, G)
	}

	endusr, endsys := util.GetTime()
	usrtime := time.Duration(endusr - startusr)
	systime := time.Duration(endsys - startsys)

	runtime.GC()
	var endmem runtime.MemStats
	runtime.ReadMemStats(&endmem)
	memory := endmem.Alloc
	os.Setenv("GOGC", "100")
	return R, usrtime + systime, memory
}

func QuickLoad(grammarfile string, graphfile string, factoryType string) (*Grammar, ds.Graph) {
	graph, grammar := LoadInfo(grammarfile, graphfile)
	VAlloc := graph.VSize()
	EAlloc := grammar.Alphabet.Size() + grammar.NonTerm.Size() +
		grammar.NestedExp.Size() + 1
	SetFactory(factoryType, VAlloc, EAlloc)
	D := f.NewGraph(graph.Name())
	for triple := range graph.Iterate() {
		s := f.NewVertex(triple[0].Label())
		p := f.NewPredicate(triple[1].Label())
		o := f.NewVertex(triple[2].Label())
		D.Add(s, p, o)
	}
	G := LoadGrammar(grammarfile)
	return G, D
}

func LoadInfo(grammarfile string, graphfile string) (ds.Graph, *Grammar) {
	f = NewSimpleFactory()
	grammar := LoadGrammar(grammarfile)
	data, err := ioutil.ReadFile(graphfile)
	if err != nil {
		panic("Error openning file: " + graphfile + "\n")
	}
	lines := strings.Split(string(data), "\n")
	graph := f.NewGraph(util.GetFileName(graphfile)).(*ds.SimpleGraph)
	for _, line := range lines {
		if line != "" {
			triple := strings.Split(line, " ")
			if len(triple) != 3 {
				fmt.Println("Error reading line:\n", line)
				return nil, nil
			}
			p := f.NewVertex(triple[1])
			if grammar.Alphabet.Contains(p) {
				graph.LoadTriple(triple)
			}
		}
	}
	return graph, grammar
}
