package ccfpq

import (
	"fmt"
	ds "github.com/ciromdrs/graph-tools/datastructures"
	"github.com/ciromdrs/graph-tools/util"
	"io/ioutil"
	"strings"
)

const OPEN = "["
const CLOSE = "]"

type Grammar struct {
	Name        string
	Rules       map[ds.Vertex][][]ds.Vertex
	Alphabet    ds.VertexSet
	NonTerm     ds.VertexSet
	NestedExp   ds.VertexSet
	StartSymbol ds.Vertex
}

func NewGrammar(f Factory) *Grammar {
	return &Grammar{
		NonTerm:     f.NewVertexSet(),
		Alphabet:    f.NewVertexSet(),
		NestedExp:   f.NewVertexSet(),
		Rules:       make(map[ds.Vertex][][]ds.Vertex),
		StartSymbol: nil,
	}
}

func (g *Grammar) AddRule(lhs ds.Vertex, rhs []ds.Vertex, f Factory) {
	if g.StartSymbol == nil {
		g.StartSymbol = lhs
	}
	g.Rules[lhs] = append(g.Rules[lhs], rhs)
	g.NonTerm.Add(lhs)
	g.Alphabet.Remove(lhs)
	for _, s := range rhs {
		if g.NonTerm.Contains(s) {
			// do nothing
		} else if s == nil {
			panic("rhs should not contain nil symbols")
		} else if isNestedExp(s.Label()) {
			g.NestedExp.Add(s)
		} else {
			g.Alphabet.Add(s)
		}
	}
}

func (g *Grammar) Show() {
	fmt.Print("N = ")
	g.NonTerm.Show()
	fmt.Print("\nT = ")
	g.Alphabet.Show()
	fmt.Print("\nNE = ")
	g.NestedExp.Show()
	fmt.Println()

	for lhs := range g.Rules {
		for _, rhs := range g.Rules[lhs] {
			fmt.Print(lhs.ToString(), " ->")
			for _, s := range rhs {
				fmt.Print(" ", s.ToString())
			}
			fmt.Println()
		}
	}
}

func LoadGrammar(path string, f Factory) *Grammar {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		panic("Error openning file: " + path + "\n")
	}
	lines := strings.Split(string(data), "\n")
	g := NewGrammar(f)
	g.Name = util.GetFileName(path)
	for _, line := range lines {
		rule := parseExp(line, f)
		if len(rule) > 0 {
			lhs := rule[0]
			var rhs []ds.Vertex
			if len(rule) > 1 {
				rhs = rule[1:]
			}
			g.AddRule(lhs, rhs, f)
		}
	}
	return g
}

func parseExp(exp string, f Factory) []ds.Vertex {
	var rule []ds.Vertex
	if exp != "" {
		for _, str := range strings.Split(exp, " ") {
			rule = append(rule, f.NewPredicate(str))
		}
	}
	return rule
}

func isNestedExp(exp string) bool {
	return strings.HasPrefix(exp, OPEN) && strings.HasSuffix(exp, CLOSE)
}
