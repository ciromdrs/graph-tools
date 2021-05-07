package ccfpq

import (
	"fmt"
	. "github.com/ciromdrs/graph-tools/util"
	"testing"
)

func TestSimpleGrammar(t *testing.T) {
	testGrammar(t, NewSimpleFactory())
}

func TestSliceGrammar(t *testing.T) {
	testGrammar(t, NewSliceFactory(10, 10))
}

func testGrammar(t *testing.T, f Factory) {
	gramfile := "sc.yrd"
	G := LoadGrammar("testdata/"+gramfile, f)
	Assert(t, G.Name == gramfile,
		fmt.Sprintf("Wrong grammar name. Want %v, got %v", gramfile, G.Name))
	Assert(t, G.NonTerm.Size() == 2,
		fmt.Sprintf("Wrong number of non-terminals. Want 2, got %v.",
			G.NonTerm.Size()))
	Assert(t, G.Alphabet.Size() == 2,
		fmt.Sprintf("Wrong number of terminals. Want 2, got %v.",
			G.Alphabet.Size()))
	// TODO: test Rules, NestedExp, StartSymbol
}
