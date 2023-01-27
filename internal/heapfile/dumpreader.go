package heapfile

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"reflect"
)

// DumpReader is used to parse heap dump file. You can set handlers for each record type via On* fields.
//
//	If the field is not set, the record will be skipped.
type DumpReader struct {
	OnObjectFn           func(record Object) error
	OnOtherRootFn        func(record OtherRoot) error
	OnTypeDescFn         func(record TypeDesc) error
	OnGoroutineFn        func(record Goroutine) error
	OnStackFrameFn       func(record StackFrame) error
	OnDumpParamsFn       func(record DumpParams) error
	OnFinalizerFn        func(record Finalizer) error
	OnItabFn             func(record Itab) error
	OnOSThreadFn         func(record OSThread) error
	OnMemStatsFn         func(record MemStats) error
	OnQueuedFinalizerFn  func(record Finalizer) error
	OnDataSegmentFn      func(record Segment) error
	OnBSSSegmentFn       func(record Segment) error
	OnDeferFn            func(record Defer) error
	OnPanicFn            func(record Panic) error
	OnAllocProfileFn     func(record AllocProfile) error
	OnAllocStackSampleFn func(record AllocStackSample) error
}

type Reader interface {
	io.ByteReader
	io.Reader
}

// Read parses heap dump. On every record it will invoke a certain On* function. The return error is either an error
//
//		from parser itself or propagated from callback. Current version of heap dump supported is 1.7.
//	 Read https://github.com/golang/go/wiki/heapdump15-through-heapdump17 for the details.
func (d DumpReader) Read(r Reader) error {
	if err := readMagic17(r); err != nil {
		return err
	}

	for {
		err := d.readRecord(r)
		if err == io.EOF {
			return nil
		} else if err != nil {
			return err
		}
	}
}

func readMagic17(r io.Reader) error {
	magic := []byte("go1.7 heap dump\n")

	buf := make([]byte, len(magic))
	if _, err := io.ReadFull(r, buf); err != nil {
		return err
	}

	if !bytes.Equal(buf, magic) {
		return fmt.Errorf("unknown format: %s", buf)
	}

	return nil
}

func (d DumpReader) readRecord(r Reader) error {
	recordType, err := binary.ReadUvarint(r)
	if err != nil {
		return err
	}

	switch recordType {
	case 0:
		return io.EOF
	case 1:
		var record Object
		if err := decodeObject(r, &record); err != nil {
			return err
		}
		if d.OnObjectFn != nil {
			return d.OnObjectFn(record)
		}
	case 2:
		var record OtherRoot
		if err := decodeRecord(r, &record); err != nil {
			return err
		}
		if d.OnOtherRootFn != nil {
			return d.OnOtherRootFn(record)
		}
	case 3:
		var record TypeDesc
		if err := decodeRecord(r, &record); err != nil {
			return err
		}
		if d.OnTypeDescFn != nil {
			return d.OnTypeDescFn(record)
		}
	case 4:
		var record Goroutine
		if err := decodeRecord(r, &record); err != nil {
			return err
		}
		if d.OnGoroutineFn != nil {
			return d.OnGoroutineFn(record)
		}
	case 5:
		var record StackFrame
		if err := decodeRecord(r, &record); err != nil {
			return err
		}
		if d.OnStackFrameFn != nil {
			return d.OnStackFrameFn(record)
		}
	case 6:
		var record DumpParams
		if err := decodeRecord(r, &record); err != nil {
			return err
		}
		if d.OnDumpParamsFn != nil {
			return d.OnDumpParamsFn(record)
		}
	case 7:
		var record Finalizer
		if err := decodeRecord(r, &record); err != nil {
			return err
		}
		if d.OnFinalizerFn != nil {
			return d.OnFinalizerFn(record)
		}
	case 8:
		var record Itab
		if err := decodeRecord(r, &record); err != nil {
			return err
		}
		if d.OnItabFn != nil {
			return d.OnItabFn(record)
		}
	case 9:
		var record OSThread
		if err := decodeRecord(r, &record); err != nil {
			return err
		}
		if d.OnOSThreadFn != nil {
			return d.OnOSThreadFn(record)
		}
	case 10:
		var record MemStats
		if err := decodeRecord(r, &record); err != nil {
			return err
		}
		if d.OnMemStatsFn != nil {
			return d.OnMemStatsFn(record)
		}
	case 11:
		var record Finalizer
		if err := decodeRecord(r, &record); err != nil {
			return err
		}
		if d.OnQueuedFinalizerFn != nil {
			return d.OnQueuedFinalizerFn(record)
		}
	case 12:
		var record Segment
		if err := decodeRecord(r, &record); err != nil {
			return err
		}
		if d.OnDataSegmentFn != nil {
			return d.OnDataSegmentFn(record)
		}
	case 13:
		var record Segment
		if err := decodeRecord(r, &record); err != nil {
			return err
		}
		if d.OnBSSSegmentFn != nil {
			return d.OnBSSSegmentFn(record)
		}
	case 14:
		var record Defer
		if err := decodeRecord(r, &record); err != nil {
			return err
		}
		if d.OnDeferFn != nil {
			return d.OnDeferFn(record)
		}
	case 15:
		var record Panic
		if err := decodeRecord(r, &record); err != nil {
			return err
		}
		if d.OnPanicFn != nil {
			return d.OnPanicFn(record)
		}
	case 16:
		var record AllocProfile
		if err := decodeRecord(r, &record); err != nil {
			return err
		}
		if d.OnAllocProfileFn != nil {
			return d.OnAllocProfileFn(record)
		}
	case 17:
		var record AllocStackSample
		if err := decodeRecord(r, &record); err != nil {
			return err
		}
		if d.OnAllocStackSampleFn != nil {
			return d.OnAllocStackSampleFn(record)
		}
	default:
		return fmt.Errorf("unknown record type: %d", recordType)
	}

	return nil
}

func decodeObject(r Reader, dst *Object) error {
	var err error
	dst.Address, err = binary.ReadUvarint(r)
	if err != nil {
		return err
	}

	dst.Contents, err = readBytes(r)
	if err != nil {
		return err
	}

	dst.PointerOffsets, err = readFieldList(r)
	return err
}

func decodeRecord(r Reader, dst any) error {
	dstType := reflect.TypeOf(dst)

	if dstType.Kind() != reflect.Pointer {
		return fmt.Errorf("dst must be pointer, got: %s", dstType.Kind())
	}

	if dstType.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("dst must be pointer to struct, got: %s", dstType.Elem().Kind())
	}

	dstVal := reflect.ValueOf(dst).Elem()
	for i := 0; i < dstType.Elem().NumField(); i++ {
		fieldType := dstType.Elem().Field(i).Type

		if fieldType.Kind() == reflect.Uint64 {
			val, err := binary.ReadUvarint(r)
			if err != nil {
				return err
			}
			dstVal.Field(i).Set(reflect.ValueOf(val))
			continue
		}

		if fieldType.Kind() == reflect.String {
			val, err := readString(r)
			if err != nil {
				return err
			}
			dstVal.Field(i).Set(reflect.ValueOf(val))
			continue
		}

		if fieldType.Kind() == reflect.Bool {
			val, err := readBool(r)
			if err != nil {
				return err
			}
			dstVal.Field(i).Set(reflect.ValueOf(val))
			continue
		}

		if fieldType.Kind() == reflect.Slice {
			if fieldType.Elem().Kind() == reflect.Uint64 {
				val, err := readFieldList(r)
				if err != nil {
					return err
				}
				dstVal.Field(i).Set(reflect.ValueOf(val))
				continue
			}

			if fieldType.Elem().Kind() == reflect.Uint8 {
				val, err := readBytes(r)
				if err != nil {
					return err
				}
				dstVal.Field(i).Set(reflect.ValueOf(val))
				continue
			}

			if fieldType.Elem() == reflect.TypeOf(Frame{}) {
				val, err := readStackFrames(r)
				if err != nil {
					return err
				}
				dstVal.Field(i).Set(reflect.ValueOf(val))
				continue
			}
		}

		if fieldType.Kind() == reflect.Array && fieldType.Elem().Kind() == reflect.Uint64 && fieldType.Len() == 256 {
			val, err := read256Uint64(r)
			if err != nil {
				return err
			}
			dstVal.Field(i).Set(reflect.ValueOf(val))
			continue
		}

		return fmt.Errorf("unknown field type: %s", fieldType)
	}

	return nil
}

func readBool(r io.ByteReader) (bool, error) {
	result, err := binary.ReadUvarint(r)
	if err != nil {
		return false, err
	}

	if result == 0 {
		return false, nil
	}

	if result == 1 {
		return true, nil
	}

	return false, fmt.Errorf("unknown bool value: %d", result)
}

func readBytes(r Reader) ([]byte, error) {
	length, err := binary.ReadUvarint(r)
	if err != nil {
		return nil, err
	}

	buf := make([]byte, length)
	bytesRead, err := io.ReadFull(r, buf)
	if err != nil {
		return nil, err
	}

	if uint64(bytesRead) != length {
		return nil, fmt.Errorf("string record unexpected end after %d want %d", bytesRead, length)
	}

	return buf, nil
}

func readString(r Reader) (string, error) {
	b, err := readBytes(r)
	return string(b), err
}

func readFieldList(r io.ByteReader) ([]uint64, error) {
	var fieldOffsets []uint64

	for {
		fieldKind, err := binary.ReadUvarint(r)
		if err != nil {
			return nil, err
		}

		if fieldKind == 0 {
			break
		}

		if fieldKind == 1 {
			fieldOffset, err := binary.ReadUvarint(r)
			if err != nil {
				return nil, err
			}
			fieldOffsets = append(fieldOffsets, fieldOffset)
		} else {
			return nil, fmt.Errorf("unknown field kind: %d", fieldKind)
		}
	}

	return fieldOffsets, nil
}

func read256Uint64(r io.ByteReader) ([256]uint64, error) {
	var result [256]uint64

	for i := 0; i < 256; i++ {
		var err error
		result[i], err = binary.ReadUvarint(r)
		if err != nil {
			return result, err
		}
	}

	return result, nil
}

func readStackFrames(r Reader) ([]Frame, error) {
	length, err := binary.ReadUvarint(r)
	if err != nil {
		return nil, err
	}

	result := make([]Frame, length)
	for i := uint64(0); i < length; i++ {
		result[i].FuncName, err = readString(r)
		if err != nil {
			return nil, err
		}

		result[i].FuncName, err = readString(r)
		if err != nil {
			return nil, err
		}

		result[i].Line, err = binary.ReadUvarint(r)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}
