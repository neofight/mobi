package convert

import "encoding/binary"

func FromUint16(value []byte) int {
	return int(binary.BigEndian.Uint16(value))
}

func FromUint32(value []byte) int {
	return int(binary.BigEndian.Uint32(value))
}
