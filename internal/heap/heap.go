package heap

import (
	"encoding/binary"
)

type Address uint64

const addressSize = 8

type Heap struct {
	objects            map[Address]Object
	stackFrames        map[Address]StackFrame
	goroutines         map[Address]Goroutine
	stackFramePtrIndex map[Address][]Address
	byteOrder          binary.ByteOrder
}

type Object struct {
	Pointers []Address
	Addr     Address
	Size     uint64
}

type StackFrame struct {
	Pointers []Address
	FuncName string
	Size     uint64
	Addr     Address
}

type Goroutine struct{}

func New(byteOrder binary.ByteOrder) *Heap {
	return &Heap{
		objects:            map[Address]Object{},
		stackFrames:        map[Address]StackFrame{},
		goroutines:         map[Address]Goroutine{},
		stackFramePtrIndex: map[Address][]Address{},
		byteOrder:          byteOrder,
	}
}

func (h *Heap) WalkPointers(start Address, objectFn func(object Object)) {
	visited := map[Address]struct{}{start: {}}
	stack := make([]Address, 0, 64*1024/addressSize)
	stack = append(stack, start)

	for len(stack) > 0 {
		var addr Address
		addr, stack = stack[len(stack)-1], stack[:len(stack)-1]

		if nextObject, ok := h.objects[addr]; ok {
			for _, pointer := range nextObject.Pointers {
				if _, isVisited := visited[pointer]; isVisited {
					continue
				}

				object, ok := h.objects[pointer]
				if !ok {
					continue
				}

				objectFn(object)
				visited[object.Addr] = struct{}{}

				if len(object.Pointers) > 0 {
					stack = append(stack, pointer)
				}
			}
		}
	}
}
