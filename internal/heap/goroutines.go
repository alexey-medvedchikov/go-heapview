package heap

import "github.com/alexey-medvedchikov/go-heapview/internal/heapfile"

type Goroutines struct {
	heap *Heap
}

func (h *Heap) Goroutines() Goroutines {
	return Goroutines{heap: h}
}

func (g Goroutines) Add(record heapfile.Goroutine) {
	g.heap.goroutines[Address(record.Frame)] = Goroutine{}
}
