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
	// A CFPQEngine is a Context-Free Path Query evaluation engine
	CFPQEngine interface {
		Run() (time.Duration, uint64)
		Graph() ds.Graph
		Grammar() *Grammar
		Query() []pair
		Factory() Factory
		CountResults() int
	}

	// TIEngine is the Trace-Item-based CFPQ evaluation engine
	TIEngine struct {
		graph   ds.Graph
		grammar *Grammar
		query   []pair
		f       Factory
		R       relationsSet
		NEW     ds.Set
		O       observersSet
	}

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

/* NodeSet Methods and Functions */
func NewNodeSet(f Factory) *NodeSet {
	return &NodeSet{
		nodes: f.NewVertexSet(),
		new:   f.NewVertexSet(),
	}
}

func (ns *NodeSet) Mark(f Factory) {
	for k := range ns.new.Iterate() {
		ns.nodes.Add(k)
	}
	ns.new = f.NewVertexSet()
}

func (ns *NodeSet) String() string {
	return "&{ " + ns.prev.predicate.String() + " " + ns.nodes.String() + " " +
		ns.next.predicate.String() + "}"
}

func (engine *TIEngine) GetOrCreate(node, label ds.Vertex, G *Grammar) Relation {
	r := engine.R.get(node, label)
	if r == nil {
		if label == nil {
			panic("Should not create empty relations")
			// nested expressions
		} else if isNestedExp(label.Label()) {
			r2 := NewNestedRelation(node, label, engine.Factory())
			subexp := strings.TrimPrefix(label.Label(), OPEN)
			subexp = strings.TrimSuffix(subexp, CLOSE)
			r2.SetRule(parseExp(subexp, engine.Factory()), engine)
			r = r2
			// nonterminals
		} else if G.NonTerm.Contains(label) {
			r2 := NewNonTerminalRelation(node, label, engine.Factory())
			startVertices := engine.Factory().NewVertexSet()
			startVertices.Add(node)
			for _, rule := range G.Rules[label] {
				r2.AddRule(startVertices, rule, engine)
			}
			r = r2
			// terminals
		} else if G.Alphabet.Contains(label) {
			// do nothing (do not delete this if clause)
			// expressions
		} else {
			r2 := NewExpressionRelation(node, label, engine.Factory())
			startVertices := engine.Factory().NewVertexSet()
			startVertices.Add(node)
			r2.SetRule(startVertices, parseExp(label.Label(), engine.Factory()),
				engine)
			r = r2
		}
		if r != nil {
			engine.R.set(node, label, r) // do not add nil relations
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

// NewTIEngine creates a new TIEngine
func NewTIEngine(grammar *Grammar, graph ds.Graph, query []pair,
	factory Factory) *TIEngine {
	engine := &TIEngine{
		graph:   graph,
		grammar: grammar,
		query:   query,
		f:       factory,
		R:       factory.NewRelationsSet(),
		NEW:     ds.NewMapSet(),
		O:       factory.NewObserversSet(),
	}
	return engine
}

func (engine *TIEngine) BuildBaseGraph(graph ds.Graph, grammar *Grammar) ds.Graph {
	D := engine.Factory().NewGraph(graph.Name())
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

func (engine *TIEngine) processNew(nodeSet *NodeSet, G *Grammar) {
	if symbol := nodeSet.next; symbol != nil {
		// update symbol's object set and NEW
		new := engine.Factory().NewVertexSet()
		for a := range nodeSet.new.Iterate() {
			var destinations ds.VertexSet = engine.Factory().NewVertexSet()
			if symbol.predicate == nil {
				panic("symbol.predicate should not be nil.")
			} else {
				if r := engine.GetOrCreate(a, symbol.predicate, G); r != nil {
					destinations = r.Objects(engine.Factory())
				}
				if !G.Alphabet.Contains(symbol.predicate) {
					engine.O.add(a, symbol.predicate, symbol.objNodeSet)
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
			engine.NEW.Add(symbol.objNodeSet)
		}
	} else {
		// updates relation's objects and notifies O
		new := engine.Factory().NewVertexSet()
		if nodeSet.relation.IsNested() &&
			nodeSet.relation.Objects(engine.Factory()).Size() == 0 {
			new.Add(nodeSet.relation.Node())
		} else {
			for a := range nodeSet.new.Iterate() {
				if !nodeSet.relation.Objects(engine.Factory()).Contains(a) {
					new.Add(a)
				}
			}
		}
		if new.Size() > 0 {
			nodeSet.relation.AddObjects(new)
			for _, o := range engine.O.get(nodeSet.relation.Node(),
				nodeSet.relation.Label()) {
				for n := range new.Iterate() {
					if !o.new.Contains(n) {
						o.new.Add(n)
						engine.NEW.Add(o)
					}
				}
			}
		}
	}
	nodeSet.Mark(engine.Factory())
}

func (engine *TIEngine) pickAndRemove() *NodeSet {
	for e := range engine.NEW.Iterate() {
		ns := e.(*NodeSet)
		engine.NEW.Remove(ns)
		return ns
	}
	return nil
}

func (engine *TIEngine) Run() (time.Duration, uint64) {
	os.Setenv("GOGC", "off")
	// allocateMemory(15000)
	runtime.GC()
	// var startmem runtime.MemStats
	// runtime.ReadMemStats(&startmem)

	startusr, startsys := util.GetTime()

	engine.R = engine.Factory().NewRelationsSet()
	engine.O = engine.Factory().NewObserversSet()
	engine.NEW = ds.NewMapSet()

	// Creating terminal relations
	for s := range engine.Graph().AllSubjects() {
		for p := range engine.Graph().Predicates(s) {
			objects := engine.Factory().NewVertexSet()
			ds.ChanToSet(engine.Graph().Objects(s, p), objects)
			engine.R.set(s, p, NewTerminalRelation(s, p, objects))
		}
	}

	// Initializing non-terminal relations
	for _, p := range engine.Query() {
		var r *NonTerminalRelation
		var node ds.Vertex
		label := p.symbol
		var startVertices ds.VertexSet

		if super, isSuperVertex := p.node.(ds.SuperVertex); isSuperVertex {
			node = super.Vertex
			startVertices = super.Vertices
		} else {
			node = p.node
			startVertices = engine.Factory().NewVertexSet()
			startVertices.Add(node)
		}
		r = NewNonTerminalRelation(node, label, engine.Factory())
		for _, rule := range engine.Grammar().Rules[label] {
			r.AddRule(startVertices, rule, engine)
		}
		engine.R.set(node, label, r)
	}

	for engine.NEW.Size() > 0 {
		newNodeSet := engine.pickAndRemove()
		engine.processNew(newNodeSet, engine.Grammar())
	}

	endusr, endsys := util.GetTime()
	usrtime := time.Duration(endusr - startusr)
	systime := time.Duration(endsys - startsys)

	runtime.GC()
	var endmem runtime.MemStats
	runtime.ReadMemStats(&endmem)
	memory := endmem.Alloc
	os.Setenv("GOGC", "100")
	return usrtime + systime, memory
}

// QuickLoad instantiates a factory and loads a grammar and a graph
func QuickLoad(grammarfile string, graphfile string,
	factoryType string) (*Grammar, ds.Graph, Factory) {
	graph, grammar := LoadInfo(grammarfile, graphfile)
	VAlloc := graph.VSize()
	EAlloc := grammar.Alphabet.Size() + grammar.NonTerm.Size() +
		grammar.NestedExp.Size() + 1
	f := NewFactory(factoryType, VAlloc, EAlloc)
	D := f.NewGraph(graph.Name())
	for triple := range graph.Iterate() {
		s := f.NewVertex(triple[0].Label())
		p := f.NewPredicate(triple[1].Label())
		o := f.NewVertex(triple[2].Label())
		D.Add(s, p, o)
	}
	G := LoadGrammar(grammarfile, f)
	return G, D, f
}

func LoadInfo(grammarfile string, graphfile string) (ds.Graph, *Grammar) {
	f := NewSimpleFactory()
	grammar := LoadGrammar(grammarfile, f)
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

// NewFactory returns a new Factory object. VAlloc and EAlloc are used in the
// SliceFactory for memory pre-allocation.
func NewFactory(factoryType string, VAlloc, EAlloc int) Factory {
	var f Factory
	switch factoryType {
	case ds.SIMPLE_FACTORY:
		f = NewSimpleFactory()
	case ds.SLICE_FACTORY:
		f = NewSliceFactory(VAlloc, EAlloc)
	default:
		panic(fmt.Sprintf("Invalid factory type %s", factoryType))
	}
	return f
}

func (engine *TIEngine) Graph() ds.Graph {
	return engine.graph
}

func (engine *TIEngine) Grammar() *Grammar {
	return engine.grammar
}

func (engine *TIEngine) Query() []pair {
	return engine.query
}

func (engine *TIEngine) Factory() Factory {
	return engine.f
}

// CountResults returns the number of results for the query
func (engine *TIEngine) CountResults() int {
	resCount := 0
	for _, p := range engine.Query() {
		var node ds.Vertex
		if super, isSuperVertex := p.node.(ds.SuperVertex); isSuperVertex {
			node = super.Vertex
		} else {
			node = p.node
		}
		resCount += engine.R.get(node, p.symbol).Objects(engine.Factory()).Size()
	}
	return resCount
}
