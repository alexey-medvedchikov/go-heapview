package heapfile

type Object struct {
	Address        uint64
	Contents       []byte
	PointerOffsets []uint64
}

type TypeDesc struct {
	Address   uint64
	Size      uint64
	Name      string
	IsPointer bool
}

type Goroutine struct {
	DescAddress      uint64
	StackTop         uint64
	ID               uint64
	GoStmtLocation   uint64
	Status           uint64
	IsSystem         bool
	IsBackground     bool
	WaitingSinceNano uint64
	WaitReason       string
	Frame            uint64
	OsThreadDesc     uint64
	TopDefer         uint64
	TopPanic         uint64
}

type StackFrame struct {
	Address        uint64
	Depth          uint64
	ChildPointer   uint64
	Contents       []byte
	EntryPC        uint64
	CurrentPC      uint64
	ContinuationPC uint64
	FuncName       string
	PointerOffsets []uint64
}

type DumpParams struct {
	BigEndian       bool
	PointerSize     uint64
	HeapStartAddr   uint64
	HeapEndAddr     uint64
	Arch            string
	GoExperimentEnv string
	NCPU            uint64
}

type Itab struct {
	Address      uint64
	TypeDescAddr uint64
}

type Finalizer struct {
	Address     uint64
	FuncPointer uint64
	EntryPC     uint64
	ArgType     uint64
	ObjType     uint64
}

type OSThread struct {
	Address uint64
	ID      uint64
	OSID    uint64
}

type Segment struct {
	Address        uint64
	Contents       []byte
	PointerOffsets []uint64
}

type MemStats struct {
	Alloc        uint64
	TotalAlloc   uint64
	Sys          uint64
	Lookups      uint64
	Mallocs      uint64
	Frees        uint64
	HeapAlloc    uint64
	HeapSys      uint64
	HeapIdle     uint64
	HeapInuse    uint64
	HeapReleased uint64
	HeapObjects  uint64
	StackInuse   uint64
	StackSys     uint64
	MSpanInuse   uint64
	MSpanSys     uint64
	MCacheInuse  uint64
	MCacheSys    uint64
	BuckHashSys  uint64
	GCSys        uint64
	OtherSys     uint64
	NextGC       uint64
	LastGC       uint64
	PauseTotalNs uint64
	PauseNs      [256]uint64
	NumGC        uint64
}

type OtherRoot struct {
	Description string
	Pointer     uint64
}

type Defer struct {
	Address   uint64
	Goroutine uint64
	Argp      uint64
	PC        uint64
	FuncVal   uint64
	EntryPC   uint64
	NextDefer uint64
}

type Panic struct {
	Address      uint64
	Goroutine    uint64
	Type         uint64
	Data         uint64
	DeferPointer uint64
	NextPanic    uint64
}

type Frame struct {
	FuncName string
	FileName string
	Line     uint64
}

type AllocProfile struct {
	ID          uint64
	Size        uint64
	StackFrames []Frame
	Allocs      uint64
	Frees       uint64
}

type AllocStackSample struct {
	Address uint64
	ID      uint64
}
