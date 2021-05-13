package graphmin

import (
	ds "github.com/ciromdrs/graph-tools/datastructures"
)

type (
	// Factory interface.
	// TODO: Add preallocation size.
	Factory interface {
		ds.Factory
		NewAugItemSet(int) *AugItemSet
		NewEmptyPosets(int) []Graph
	}

	// HashFactory is a hashmap-based implementation of Factory.
	HashFactory struct {
		ds.SimpleFactory
	}
)

// NewHashFactory creates a HashFactory object.
func NewHashFactory() *HashFactory {
	return &HashFactory{}
}

// NewAugItemSet creates an AugItemSet object.
func (f *HashFactory) NewAugItemSet(prealloc int) *AugItemSet {
	return newAugItemSet(f, prealloc)
}

// NewEmptyPosets creates an array of empty posets.
func (f *HashFactory) NewEmptyPosets(length int) []Graph {
	posets := make([]Graph, length)
	for i := 0; i < length; i++ {
		posets[i] = newHashGraph()
	}
	return posets
}
