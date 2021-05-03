package ccfpq

import (
	"fmt"
	"os"
	ds "rdf-ccfpq/go/data_structures"
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

// func testFactory(t *testing.T, f Factory) {
// 	testSet(f)
// 	f.Reset()
// 	testGraph(f)
// }

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
	if got := A.Get(0).(ds.Vertex); !got.Equals(zero) {
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
		if !B.Get(kv.Key).(ds.Vertex).Equals(kv.Value.(ds.Vertex)) {
			t.Fatalf("A != B")
		}
	}
}

func testSet(t *testing.T, f Factory) {
	A := f.NewSet()
	one := f.NewVertex("1")
	assert(A.Size() == 0, "Empty set should have Size 0", t)
	assert(!A.Contains(one), "Empty set should not contain 1.", t)

	A.Add(one)
	assert(A.Size() == 1, "A.Size() != 1", t)
	assert(A.Contains(one), "A should contain 1.", t)

	two := f.NewVertex("2")
	three := f.NewVertex("3")
	four := f.NewVertex("4")
	A.Add(two)
	A.Add(three)
	A.Add(four)
	assert(A.Size() == 4, "A.Size() != 4", t)

	B := f.NewSet()
	B.Add(one)
	B.Add(two)
	B.Add(three)
	B.Add(four)
	assert(A.Equals(B), "A should be equal to B.", t)

	A.Remove(four)
	assert(A.Size() == 3, "A.Size() != 3", t)
	assert(!A.Contains(four), "A should not contain 4.", t)
	assert(!A.Equals(B), "A should be different from B.", t)
}

func testGraph(t *testing.T, f Factory) {
	g := f.NewGraph("g")
	s := f.NewVertex("s")
	p := f.NewPredicate("p")
	o := f.NewVertex("o")

	assert(g.Size() == 0, "g.Size() != 0", t)
	assert(!g.Contains(s, p, o), "g should not contain (s,p,o).", t)

	g.Add(s, p, o)
	assert(g.Size() == 1, "g.Size() != 1", t)
	assert(g.VSize() == 2, fmt.Sprintln("g.VSize() ", (g.VSize()), " != 2"), t)
	assert(g.ESize() == 1, "g.ESize() != 1", t)
	assert(g.Contains(s, p, o), "g should contain (s,p,o).", t)
	assert(!g.Contains(o, p, s), "g should not contain (o,p,s).", t)
	g.Add(o, p, s)
	assert(g.Contains(o, p, s), "g should contain (o,p,s).", t)

	o2 := f.NewVertex("o2")
	g.Add(s, p, o2)
	assert(g.VSize() == 3, "g.VSize() != 3", t)
	set := f.NewVertexSet()
	set.Add(o2)
	set.Add(o)
	res := f.NewVertexSet()
	ds.ChanToSet(g.Objects(s, p), res)
	assert(res.Equals(set), "Wrong objects for (s,p).", t)

	{
		pairs := [][2]ds.Vertex{
			{s, o},
			{s, o2},
			{o, s},
		}
		res := ds.ChanToPairs(g.SubjectObjects(p))
		assert(len(res) == len(pairs), "len(g.SubjectObjects(*p)) != len(pairs).", t)
		equalPairs := true
		for _, p1 := range pairs {
			in := false
			for _, p2 := range res {
				in = in || (p1[0].Equals(p2[0]) && p1[1].Equals(p2[1]))
			}
			equalPairs = equalPairs && in
		}
		assert(equalPairs, "Wrong (subject,object) pairs for p.", t)
	}

	{
		all := f.NewVertexSet()
		all.Add(s)
		all.Add(o)
		assert(!all.Equals(g.AllNodes()), "o2 should not be in AllNodes.", t)
		all.Add(o2)
		assert(all.Equals(g.AllNodes()), "Wrong nodes.", t)
	}

	{
		triples := [][3]ds.Vertex{
			{s, p, o},
			{s, p, o2},
			{o, p, s},
		}
		res := ds.ChanToTriples(g.Iterate())
		assert(len(res) == len(triples), "len(g.Iterate()) != len(triples).", t)
		equal := true
		for _, t1 := range triples {
			in := false
			for _, t2 := range res {
				in = in || (t1[0].Equals(t2[0]) && t1[1].Equals(t2[1]) && t1[2].Equals(t2[2]))
			}
			equal = equal && in
		}
		assert(equal, "Wrong triples in g.", t)
	}

	{
		subjs := f.NewVertexSet()
		subjs.Add(s)
		subjs.Add(o)
		res := f.NewVertexSet()
		ds.ChanToSet(g.AllSubjects(), res)
		assert(res.Equals(subjs), "Wrong subjects.", t)
	}

	p2 := f.NewPredicate("p2")
	g.Add(s, p2, o)
	{
		assert(g.ESize() == 2, "g.ESize() != 2", t)
		preds := f.NewVertexSet()
		preds.Add(p)
		preds.Add(p2)
		res := f.NewVertexSet()
		ds.ChanToSet(g.Predicates(s), res)
		assert(res.Equals(preds), "Wrong predicates.", t)
	}

	{
		set := f.NewVertexSet()
		for i := 0; i < ds.CELL_SIZE*2; i++ {
			v := f.NewVertex(strconv.Itoa(i))
			set.Add(v)
		}
		assert(set.Size() == ds.CELL_SIZE*2, "set.Size() != CELL_SIZE*2", t)
	}
}

func assert(condition bool, errmsg string, t *testing.T) {
	if !condition {
		t.Fatalf(errmsg)
	}
}
