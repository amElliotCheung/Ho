package interpreter

import "encoding/binary"

type Frame struct {
	fn *CompiledFunction
	ip int
	bp int
}

func NewFrame(fn *CompiledFunction, bp int) *Frame {
	return &Frame{fn: fn, ip: 0, bp: bp}
}

func (f *Frame) readUint8(offset int) int {
	return int(f.fn.Instructions[f.ip+offset])
}

func (f *Frame) readUint16(offset int) int {
	idx := binary.BigEndian.Uint16(f.fn.Instructions[f.ip+offset:])
	return int(idx)
}
