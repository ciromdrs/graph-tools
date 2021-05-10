package graphmin

import (
	ds "github.com/ciromdrs/graph-tools/datastructures"
)

type (
	// Factory interface.
	Factory interface {
		ds.Factory
		NewAugItemSet(int) *AugItemSet
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
