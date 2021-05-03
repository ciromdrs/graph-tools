package ccfpq

import (
    ds "rdf-ccfpq/go/data_structures"
)

type (
    observersSet interface {
        add(ds.Vertex, ds.Vertex, *NodeSet)
        get(ds.Vertex, ds.Vertex) []*NodeSet
    }

    mapObserversSet struct {
        data *ds.SimpleMap
    }

    sliceObserversSet struct {
        data *ds.SliceMap
        ESize int
    }
)

/* mapObserversSet Methods and Functions */
func newMapObserversSet(size int) *mapObserversSet {
    return &mapObserversSet{data: f.NewMap(size).(*ds.SimpleMap)}
}

func (O *mapObserversSet) add(node, symbol ds.Vertex, set *NodeSet) {
    v := node.(ds.SimpleVertex)
    s := symbol.(ds.SimpleVertex)
    key := v.Label()+"|"+s.Label()
    var observers []*NodeSet
    var ok bool
    if observers, ok = O.data.Get(key).([]*NodeSet); ok {
        observers = append(observers, set)
    } else {
        observers = []*NodeSet{set}
    }
    O.data.Set(key, observers)
}

func (O *mapObserversSet) get(node, symbol ds.Vertex) []*NodeSet {
    v := node.(ds.SimpleVertex)
    s := symbol.(ds.SimpleVertex)
    key := v.Label()+"|"+s.Label()
    if nodesets, ok := O.data.Get(key).([]*NodeSet); ok {
        return nodesets
    }
    return nil
}


/* sliceObserversSet Methods and Functions */
func newSliceObserversSet(VSize, ESize int) *sliceObserversSet {
    total := VSize * ESize
    return &sliceObserversSet{
        data:  f.NewMap(total).(*ds.SliceMap),
        ESize: ESize,
    }
}

func (O *sliceObserversSet) add(node, symbol ds.Vertex, set *NodeSet) {
    v := node.(ds.BitVertex)
    s := symbol.(ds.BitVertex)
    i := (O.ESize * v.IndexInSlice()) + s.IndexInSlice()
    var observers []*NodeSet
    var ok bool
    if observers, ok = O.data.Get(i).([]*NodeSet); ok {
        observers = append(observers, set)
    } else {
        observers = []*NodeSet{set}
    }
    O.data.Set(i, observers)
}

func (O *sliceObserversSet) get(node, symbol ds.Vertex) []*NodeSet {
    v := node.(ds.BitVertex)
    s := symbol.(ds.BitVertex)
    i := (O.ESize * v.IndexInSlice()) + s.IndexInSlice()
    value := O.data.Get(i)
    if value != nil {
        return value.([]*NodeSet)
    }
    return nil
}
