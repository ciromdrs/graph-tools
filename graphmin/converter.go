package graphmin

import (
	// "fmt"
	"github.com/ciromdrs/graph-tools/ccfpq"
)

type (
	// Converter converts trace items to augmented trace items.
	Converter struct {
		factory Factory
	}
)

// newHashConverter creates a HashConverter object.
func newConverter(factory Factory) *Converter {
	return &Converter{factory: factory}
}

// Convert converts TraceItems to AugItems.
func (conv *Converter) Convert(traceItems []*ccfpq.TraceItem,
	graph Graph) *AugItemSet {
	panic("Not implemented yet.")
}
