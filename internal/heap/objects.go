package heap

import "github.com/alexey-medvedchikov/go-heapview/internal/heapfile"

type Objects struct {
	heap *Heap
}

func (h *Heap) Objects() Objects {
	return Objects{heap: h}
}

func (o Objects) Walk(fn func(object Object) error) error {
	for _, object := range o.heap.objects {
		if err := fn(object); err != nil {
			return err
		}
	}

	return nil
}

func (o Objects) Add(object heapfile.Object) {
	obj := Object{
		Addr: Address(object.Address),
		Size: uint64(len(object.Contents)),
	}

	for _, ptrOffset := range object.PointerOffsets {
		ptr := o.heap.byteOrder.Uint64(object.Contents[ptrOffset:])
		obj.Pointers = append(obj.Pointers, Address(ptr))
	}

	o.heap.objects[Address(object.Address)] = obj
}

type ObjectStats struct {
	OwnedSize  uint64
	OwnedCount uint64
}

func (o Objects) Stats(addr Address) ObjectStats {
	var ownedSize uint64
	var ownedCount uint64

	o.heap.WalkPointers(addr, func(object Object) {
		ownedSize += object.Size
		ownedCount += 1
	})

	return ObjectStats{
		OwnedSize:  ownedSize,
		OwnedCount: ownedCount,
	}
}
