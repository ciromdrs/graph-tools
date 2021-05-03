package data_structures

import "fmt"

type (
    MapKey interface {}

    MapValue interface {}

    Map interface {
        Set(MapKey, MapValue)
        Get(MapKey) MapValue
        Remove(MapKey)
        Iterate() <- chan KeyValue
        Show()
    }

    KeyValue struct {
        Key   MapKey
        Value MapValue
    }

    SimpleMap struct {
        data map[MapKey]MapValue
    }

    SliceMap struct {
        data []MapValue
        ZeroValue MapValue
    }
)

/* SimpleMap Functions and Methods */
func NewSimpleMap() *SimpleMap {
    return &SimpleMap{data: make(map[MapKey]MapValue)}
}

func (m *SimpleMap) Set(key MapKey, value MapValue) {
    m.data[key] = value
}

func (m *SimpleMap) Get(key MapKey) MapValue {
    return m.data[key]
}

func (m *SimpleMap) Remove(key MapKey) {
    delete(m.data, key)
}

func (m *SimpleMap) Iterate() <- chan KeyValue {
    ch := make(chan KeyValue)
    go func() {
        for k, v := range m.data {
            ch <- KeyValue{
                Key: k,
                Value: v,
            }
        }
        defer close(ch)
    }()
    return ch
}

func (m *SimpleMap) Show() {
    fmt.Print("{")
    for kv := range m.Iterate() {
        fmt.Print(kv.Key,"-->",kv.Value,", ")
    }
    fmt.Println("}")
}

/* SliceMap Functions and Methods */
func NewSliceMap(capacity int) *SliceMap {
    return &SliceMap{data: make([]MapValue, capacity)}
}

func (m *SliceMap) Set(key MapKey, value MapValue) {
    i := key.(int)
    if i >= len(m.data) {
        m.expand(i+1)
    }
    m.data[i] = value
}

func (m *SliceMap) Get(key MapKey) MapValue {
    i := key.(int)
    if i >= len(m.data) {
        return nil
    }
    return m.data[i]
}

func (m *SliceMap) Remove(key MapKey) {
    if i := key.(int); i < len(m.data) {
        m.data[i] = nil
    }
}

func (m *SliceMap) expand(length int) {
    if length < len(m.data) {
        panic("Cannot expand SliceMap. `length` is too small.")
    }
    newlength := length
    for ; newlength < length; newlength = newlength << 1 {}
    new := make([]MapValue, newlength)
    copy(new, m.data)
    m.data = new
}

func (m *SliceMap) Iterate() <- chan KeyValue {
    ch := make(chan KeyValue)
    go func() {
        for i, v := range m.data {
            if v != nil {
                ch <- KeyValue{
                    Key: i,
                    Value: v,
                }
            }
        }
        defer close(ch)
    }()
    return ch
}

func (m *SliceMap) Show() {
    fmt.Print("{")
    for kv := range m.Iterate() {
        fmt.Print(kv.Key,"-->",kv.Value)
    }
    fmt.Println("}")
}
