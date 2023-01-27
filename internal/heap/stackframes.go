package heap

import "github.com/alexey-medvedchikov/go-heapview/internal/heapfile"

type StackFrames struct {
	heap *Heap
}

func (h *Heap) StackFrames() StackFrames {
	return StackFrames{heap: h}
}

func (s StackFrames) Add(frame heapfile.StackFrame) {
	fr := StackFrame{
		FuncName: frame.FuncName,
		Size:     uint64(len(frame.Contents)),
		Addr:     Address(frame.Address),
	}

	for _, ptrOffset := range frame.PointerOffsets {
		ptr := Address(s.heap.byteOrder.Uint64(frame.Contents[ptrOffset:]))
		fr.Pointers = append(fr.Pointers, ptr)
		s.heap.stackFramePtrIndex[ptr] = append(s.heap.stackFramePtrIndex[ptr], fr.Addr)
	}

	s.heap.stackFrames[Address(frame.Address)] = fr
}

func (s StackFrames) HasAddress(addr Address) []StackFrame {
	var frames []StackFrame
	for _, frameAddr := range s.heap.stackFramePtrIndex[addr] {
		frames = append(frames, s.heap.stackFrames[frameAddr])
	}

	return frames
}
