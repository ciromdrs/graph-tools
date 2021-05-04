package datastructures

import (
	"fmt"
	. "github.com/ciromdrs/graph-tools/util"
	"os"
	"strconv"
	"testing"
)

var (
	slf *SliceFactory
	spf *SimpleFactory
)

func TestMain(m *testing.M) {
	slf = NewSliceFactory(200, 2)
	spf = NewSimpleFactory()
	code := m.Run()
	os.Exit(code)
}

func TestSliceMap(t *testing.T) {
	slf.Reset()
	testMap(t, slf)
}

func TestSimpleMap(t *testing.T) {
	spf.Reset()
	testMap(t, spf)
}

func TestSliceSet(t *testing.T) {
	slf.Reset()
	testSet(t, slf)
}

func TestSimpleSet(t *testing.T) {
	spf.Reset()
	testSet(t, spf)
}

func TestSliceGraph(t *testing.T) {
	slf.Reset()
	testGraph(t, slf)
}

func TestSimpleGraph(t *testing.T) {
	spf.Reset()
	testGraph(t, spf)
}

func testMap(t *testing.T, f Factory) {
	A := f.NewMap(0)
	zero := f.NewVertex("0")

	A.Set(0, zero)
	if got := A.Get(0).(Vertex); !got.Equals(zero) {
		t.Fatalf("A should contain 0.")
	}

	one := f.NewVertex("1")
	two := f.NewVertex("2")
	three := f.NewVertex("3")
	A.Set(1, one)
	A.Set(2, two)
	A.Set(3, three)
	A.Remove(3)

	if got := A.Get(3); got != nil {
		t.Fatalf("A should not contain 3.")
	}

	B := f.NewMap(0)
	B.Set(0, zero)
	B.Set(1, one)
	B.Set(2, two)
	for kv := range A.Iterate() {
		if !B.Get(kv.Key).(Vertex).Equals(kv.Value.(Vertex)) {
			t.Fatalf("A != B")
		}
	}
}

func testSet(t *testing.T, f Factory) {
	A := f.NewSet()
	one := f.NewVertex("1")
	Assert(t, A.Size() == 0, "Empty set should have Size 0")
	Assert(t, !A.Contains(one), "Empty set should not contain 1.")

	A.Add(one)
	Assert(t, A.Size() == 1, "A.Size() != 1")
	Assert(t, A.Contains(one), "A should contain 1.")

	two := f.NewVertex("2")
	three := f.NewVertex("3")
	four := f.NewVertex("4")
	A.Add(two)
	A.Add(three)
	A.Add(four)
	Assert(t, A.Size() == 4, "A.Size() != 4")

	B := f.NewSet()
	B.Add(one)
	B.Add(two)
	B.Add(three)
	B.Add(four)
	Assert(t, A.Equals(B), "A should be equal to B.")

	A.Remove(four)
	Assert(t, A.Size() == 3, "A.Size() != 3")
	Assert(t, !A.Contains(four), "A should not contain 4.")
	Assert(t, !A.Equals(B), "A should be different from B.")
}

func testGraph(t *testing.T, f Factory) {
	g := f.NewGraph("g")
	s := f.NewVertex("s")
	p := f.NewPredicate("p")
	o := f.NewVertex("o")

	Assert(t, g.Size() == 0, "g.Size() != 0")
	Assert(t, !g.Contains(s, p, o), "g should not contain (s,p,o).")

	g.Add(s, p, o)
	Assert(t, g.Size() == 1, "g.Size() != 1")
	Assert(t, g.VSize() == 2, fmt.Sprintln("g.VSize() ", (g.VSize()), " != 2"))
	Assert(t, g.ESize() == 1, "g.ESize() != 1")
	Assert(t, g.Contains(s, p, o), "g should contain (s,p,o).")
	Assert(t, !g.Contains(o, p, s), "g should not contain (o,p,s).")
	g.Add(o, p, s)
	Assert(t, g.Contains(o, p, s), "g should contain (o,p,s).")

	o2 := f.NewVertex("o2")
	g.Add(s, p, o2)
	Assert(t, g.VSize() == 3, "g.VSize() != 3")
	set := f.NewVertexSet()
	set.Add(o2)
	set.Add(o)
	res := f.NewVertexSet()
	ChanToSet(g.Objects(s, p), res)
	Assert(t, res.Equals(set), "Wrong objects for (s,p).")

	{
		pairs := [][2]Vertex{
			{s, o},
			{s, o2},
			{o, s},
		}
		res := ChanToPairs(g.SubjectObjects(p))
		Assert(t, len(res) == len(pairs), "len(g.SubjectObjects(*p)) != len(pairs).")
		equalPairs := true
		for _, p1 := range pairs {
			in := false
			for _, p2 := range res {
				in = in || (p1[0].Equals(p2[0]) && p1[1].Equals(p2[1]))
			}
			equalPairs = equalPairs && in
		}
		Assert(t, equalPairs, "Wrong (subject,object) pairs for p.")
	}

	{
		all := f.NewVertexSet()
		all.Add(s)
		all.Add(o)
		Assert(t, !all.Equals(g.AllNodes()), "o2 should not be in AllNodes.")
		all.Add(o2)
		Assert(t, all.Equals(g.AllNodes()), "Wrong nodes.")
	}

	{
		triples := [][3]Vertex{
			{s, p, o},
			{s, p, o2},
			{o, p, s},
		}
		res := ChanToTriples(g.Iterate())
		Assert(t, len(res) == len(triples), "len(g.Iterate()) != len(triples).")
		equal := true
		for _, t1 := range triples {
			in := false
			for _, t2 := range res {
				in = in || (t1[0].Equals(t2[0]) && t1[1].Equals(t2[1]) && t1[2].Equals(t2[2]))
			}
			equal = equal && in
		}
		Assert(t, equal, "Wrong triples in g.")
	}

	{
		subjs := f.NewVertexSet()
		subjs.Add(s)
		subjs.Add(o)
		res := f.NewVertexSet()
		ChanToSet(g.AllSubjects(), res)
		Assert(t, res.Equals(subjs), "Wrong subjects.")
	}

	p2 := f.NewPredicate("p2")
	g.Add(s, p2, o)
	{
		Assert(t, g.ESize() == 2, "g.ESize() != 2")
		preds := f.NewVertexSet()
		preds.Add(p)
		preds.Add(p2)
		res := f.NewVertexSet()
		ChanToSet(g.Predicates(s), res)
		Assert(t, res.Equals(preds), "Wrong predicates.")
	}

	{
		set := f.NewVertexSet()
		for i := 0; i < CELL_SIZE*2; i++ {
			v := f.NewVertex(strconv.Itoa(i))
			set.Add(v)
		}
		Assert(t, set.Size() == CELL_SIZE*2, "set.Size() != CELL_SIZE*2")
	}
}
