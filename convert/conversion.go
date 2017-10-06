package convert

import (
	"bytes"
	"encoding/binary"
	"io"
)

func FromUint16(value []byte) int {
	return int(binary.BigEndian.Uint16(value))
}

func FromUint32(value []byte) int {
	return int(binary.BigEndian.Uint32(value))
}

func FromLZ77(text []byte) []byte {

	var reader = bytes.NewReader(text)

	var buffer [4096]byte
	var pos int

	for {
		if pos == 4096 {
			break
		}

		c, err := reader.ReadByte()

		if err == io.EOF {
			break
		}

		switch {

		// 0x00: "1 literal" copy that byte unmodified to the decompressed stream.
		case c == 0x00:
			buffer[pos] = c
			pos++

		// 0x09 to 0x7f: "1 literal" copy that byte unmodified to the decompressed stream.
		case c >= 0x09 && c <= 0x7f:
			buffer[pos] = c
			pos++

		// 0x01 to 0x08: "literals": the byte is interpreted as a count from 1 to 8, and that many literals are copied
		// unmodified from the compressed stream to the decompressed stream.
		case c >= 0x01 && c <= 0x08:
			length := int(c)
			for i := 0; i < length; i++ {
				c, err = reader.ReadByte()
				buffer[pos] = c
				pos++
			}

		// 0x80 to 0xbf: "length, distance" pair: the 2 leftmost bits of this byte ('10') are discarded, and the
		// following 6 bits are combined with the 8 bits of the next byte to make a 14 bit "distance, length" item.
		// Those 14 bits are broken into 11 bits of distance backwards from the current location in the uncompressed
		// text, and 3 bits of length to copy from that point (copying n+3 bytes, 3 to 10 bytes).
		case c >= 0x80 && c <= 0xbf:
			c2, _ := reader.ReadByte()

			distance := (int(c&0x3F)<<8 | int(c2)) >> 3
			length := int(c2&0x07) + 3

			start := pos - distance

			for i := 0; i < length; i++ {
				c = buffer[start+i]
				buffer[pos] = c
				pos++
			}

		// 0xc0 to 0xff: "byte pair": this byte is decoded into 2 characters: a space character, and a letter formed
		// from this byte XORed with 0x80.
		case c >= 0xc0:
			buffer[pos] = ' '
			pos++
			buffer[pos] = c ^ 0x80
			pos++
		}
	}

	return buffer[:pos]
}
