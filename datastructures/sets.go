package datastructures

import (
	"fmt"
)

type (
	SetElement interface {
		String() string
	}

	Set interface {
		Add(e SetElement) bool
		Contains(e SetElement) bool
		Equals(Set) bool
		Iterate() <-chan SetElement
		Remove(e SetElement) bool
		Show()
		String() string
		Size() int
		Update(Set) int
	}

	MapSet struct {
		data map[SetElement]bool
	}

	SliceSet struct {
		data []SliceSetElement
		size int
	}

	SliceSetElement interface {
		SetElement
		IndexInSlice() int
	}

	/* Errors */
	NotFoundError struct{}
)

/* NotFoundError Functions and Methods */
func (e *NotFoundError) Error() string {
	return fmt.Sprintf("Not found")
}

/* Common Set functions and methods */

func setString(s Set) string {
	out := "{ "
	for e := range s.Iterate() {
		out += e.String() + " "
	}
	out += "}"
	return out
}

// update adds all elements from `from` to `to`
func update(from, to Set) int {
	count := 0
	for e := range from.Iterate() {
		if to.Add(e) {
			count++
		}
	}
	return count
}

/* MapSet Functions and Methods */
func NewMapSet() *MapSet {
	return &MapSet{
		data: make(map[SetElement]bool),
	}
}

func (s *MapSet) Size() int {
	return len(s.data)
}

func (s *MapSet) Equals(other Set) bool {
	s2 := other.(*MapSet)
	if s.Size() != s2.Size() {
		return false
	}
	for e := range s2.Iterate() {
		if !s.Contains(e) {
			return false
		}
	}
	return true
}

func (s *MapSet) Add(e SetElement) bool {
	added := false
	if _, in := s.data[e]; !in {
		s.data[e] = true
		added = true
	}
	return added
}

func (s *MapSet) Contains(e SetElement) bool {
	if _, in := s.data[e]; in {
		return true
	}
	return false
}

func (s *MapSet) Remove(e SetElement) bool {
	removed := false
	if _, in := s.data[e]; in {
		delete(s.data, e)
		removed = true
	}
	return removed
}

func (s *MapSet) Iterate() <-chan SetElement {
	ch := make(chan SetElement)
	go func() {
		for e := range s.data {
			ch <- e
		}
		defer close(ch)
	}()
	return ch
}

// Update adds all elements in toAdd the set
func (s *MapSet) Update(toAdd Set) int {
	return update(toAdd, s)
}

func (s *MapSet) String() string {
	return setString(s)
}

func (s *MapSet) Show() {
	fmt.Println(s.String())
}

/* SliceSet Functions and Methods */
func NewSliceSet(preallocate int) *SliceSet {
	return &SliceSet{
		data: []SliceSetElement{}[:preallocate],
	}
}

func (s *SliceSet) Size() int {
	return s.size
}

func (s *SliceSet) Add(e SetElement) bool {
	added := false
	if !s.Contains(e) {
		added = true
		s.size++
		i := e.(SliceSetElement).IndexInSlice()
		if i >= len(s.data) {
			s.expand(i + 1)
		}
		s.data[i] = e.(SliceSetElement)
	}
	return added
}

func (s *SliceSet) expand(length int) {
	if length < len(s.data) {
		panic("Cannot expand BitVertexSet. `length` is too small.")
	}
	new := make([]SliceSetElement, length)
	copy(new, s.data)
	s.data = new
}

func (s *SliceSet) Contains(e SetElement) bool {
	se := e.(SliceSetElement)
	if se.IndexInSlice() >= len(s.data) {
		return false
	}
	return s.data[se.IndexInSlice()] != nil
}

func (s *SliceSet) Remove(e SetElement) bool {
	removed := false
	if s.Contains(e) {
		i := e.(SliceSetElement).IndexInSlice()
		s.data[i] = nil
		s.size--
		removed = true
	}
	return removed
}

func (s *SliceSet) Iterate() <-chan SetElement {
	ch := make(chan SetElement)
	go func() {
		for _, e := range s.data {
			ch <- e
		}
		defer close(ch)
	}()
	return ch
}

func (s *SliceSet) Equals(other Set) bool {
	s2 := other.(*SliceSet)
	if s.Size() != s2.Size() {
		return false
	}
	for i := range s.data {
		if s.data[i] != s2.data[i] {
			return false
		}
	}
	return true
}

// Update adds all elements in toAdd the set
func (s *SliceSet) Update(toAdd Set) int {
	return update(toAdd, s)
}

func (s *SliceSet) String() string {
	return setString(s)
}

func (s *SliceSet) Show() {
	fmt.Println(s.String())
}
