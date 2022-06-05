package interpreter

import "encoding/binary"

type Instructions []byte

func (i Instructions) readUint8(pos int) int {
	return int(i[pos])
}

func (i Instructions) readUint16(pos int) int {
	result := binary.BigEndian.Uint16(i[pos:])
	return int(result)
}
