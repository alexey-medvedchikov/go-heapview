package dumpcmd

import (
	"bufio"
	"encoding/json"
	"io"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/alexey-medvedchikov/go-heapview/internal/fileutils"
	"github.com/alexey-medvedchikov/go-heapview/internal/heapfile"
)

func Command() *cli.Command {
	return &cli.Command{
		Name: "dump",
		Action: func(c *cli.Context) error {
			fpath := c.Args().Get(0)
			return fileutils.WithFileOpened(fpath, func(fp *os.File) error {
				return dumpAction(fp)
			}, os.O_RDONLY, 0640)
		},
		Usage: "Output contents of the heap dump file in a newline-delimited JSON format",
	}
}

func dumpAction(r io.Reader) error {
	encoder := json.NewEncoder(os.Stdout)

	type record struct {
		Type   string
		Record any
	}

	reader := heapfile.DumpReader{
		OnObjectFn: func(v heapfile.Object) error {
			return encoder.Encode(record{Type: "Object", Record: any(v)})
		},
		OnOtherRootFn: func(v heapfile.OtherRoot) error {
			return encoder.Encode(record{Type: "OtherRoot", Record: any(v)})
		},
		OnTypeDescFn: func(v heapfile.TypeDesc) error {
			return encoder.Encode(record{Type: "TypeDesc", Record: any(v)})
		},
		OnGoroutineFn: func(v heapfile.Goroutine) error {
			return encoder.Encode(record{Type: "Goroutine", Record: any(v)})
		},
		OnStackFrameFn: func(v heapfile.StackFrame) error {
			return encoder.Encode(record{Type: "StackFrame", Record: any(v)})
		},
		OnDumpParamsFn: func(v heapfile.DumpParams) error {
			return encoder.Encode(record{Type: "DumpParams", Record: any(v)})
		},
		OnFinalizerFn: func(v heapfile.Finalizer) error {
			return encoder.Encode(record{Type: "Finalizer", Record: any(v)})
		},
		OnItabFn: func(v heapfile.Itab) error {
			return encoder.Encode(record{Type: "Itab", Record: any(v)})
		},
		OnOSThreadFn: func(v heapfile.OSThread) error {
			return encoder.Encode(record{Type: "OSThread", Record: any(v)})
		},
		OnMemStatsFn: func(v heapfile.MemStats) error {
			return encoder.Encode(record{Type: "MemStats", Record: any(v)})
		},
		OnQueuedFinalizerFn: func(v heapfile.Finalizer) error {
			return encoder.Encode(record{Type: "QueuedFinalizer", Record: any(v)})
		},
		OnDataSegmentFn: func(v heapfile.Segment) error {
			return encoder.Encode(record{Type: "DataSegment", Record: any(v)})
		},
		OnBSSSegmentFn: func(v heapfile.Segment) error {
			return encoder.Encode(record{Type: "BSSSegment", Record: any(v)})
		},
		OnDeferFn: func(v heapfile.Defer) error {
			return encoder.Encode(record{Type: "Defer", Record: any(v)})
		},
		OnPanicFn: func(v heapfile.Panic) error {
			return encoder.Encode(record{Type: "Panic", Record: any(v)})
		},
		OnAllocProfileFn: func(v heapfile.AllocProfile) error {
			return encoder.Encode(record{Type: "AllocProfile", Record: any(v)})
		},
		OnAllocStackSampleFn: func(v heapfile.AllocStackSample) error {
			return encoder.Encode(record{Type: "AllocStackSample", Record: any(v)})
		},
	}

	buffered := bufio.NewReader(r)
	return reader.Read(buffered)
}
