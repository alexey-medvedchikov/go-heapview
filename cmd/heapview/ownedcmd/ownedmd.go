package ownedcmd

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"errors"
	"io"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/alexey-medvedchikov/go-heapview/internal/fileutils"
	"github.com/alexey-medvedchikov/go-heapview/internal/heap"
	"github.com/alexey-medvedchikov/go-heapview/internal/heapfile"
)

func Command() *cli.Command {
	return &cli.Command{
		Name: "owned",
		Action: func(c *cli.Context) error {
			fpath := c.Args().Get(0)
			return fileutils.WithFileOpened(fpath, func(fp *os.File) error {
				return ownedAction(fp)
			}, os.O_RDONLY, 0640)
		},
		Usage: "Show pointers that are rooted in stack frames together with the statistics",
	}
}

type frame struct {
	Address  heap.Address
	FuncName string
}

type pointer struct {
	Address    heap.Address
	Size       int
	OwnedSize  int
	OwnedCount int
	Frames     []frame
}

func ownedAction(r io.Reader) error {
	encoder := json.NewEncoder(os.Stdout)

	h, err := readDump(r)
	if err != nil {
		return err
	}

	return h.Objects().Walk(func(object heap.Object) error {
		frames := h.StackFrames().HasAddress(object.Addr)

		if len(frames) > 0 {
			stats := h.Objects().Stats(object.Addr)

			p := pointer{
				Address:    object.Addr,
				Size:       int(object.Size),
				OwnedSize:  int(stats.OwnedSize),
				OwnedCount: int(stats.OwnedCount),
			}

			for _, fr := range frames {
				p.Frames = append(p.Frames, frame{
					Address:  fr.Addr,
					FuncName: fr.FuncName,
				})
			}

			if err := encoder.Encode(p); err != nil {
				return err
			}
		}

		return nil
	})
}

func readDump(r io.Reader) (*heap.Heap, error) {
	var h *heap.Heap

	endiannessUnknownErr := errors.New("DumpParams missing, endianness unknown")

	reader := heapfile.DumpReader{
		OnDumpParamsFn: func(record heapfile.DumpParams) error {
			var byteOrder binary.ByteOrder = binary.LittleEndian
			if record.BigEndian {
				byteOrder = binary.BigEndian
			}
			h = heap.New(byteOrder)
			return nil
		},
		OnObjectFn: func(record heapfile.Object) error {
			if h == nil {
				return endiannessUnknownErr
			}
			h.Objects().Add(record)
			return nil
		},
		OnStackFrameFn: func(record heapfile.StackFrame) error {
			if h == nil {
				return endiannessUnknownErr
			}
			h.StackFrames().Add(record)
			return nil
		},
		OnGoroutineFn: func(record heapfile.Goroutine) error {
			if h == nil {
				return endiannessUnknownErr
			}
			h.Goroutines().Add(record)
			return nil
		},
	}

	buffered := bufio.NewReader(r)
	return h, reader.Read(buffered)
}
