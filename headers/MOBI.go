package headers

import (
	"encoding/binary"
	"fmt"
	"github.com/neofight/mobi/convert"
	"io"
)

type MOBI struct {
	EXTHHeaderPresent bool
	FirstContentIndex int
	FirstImageIndex   int
	FirstNonBookIndex int
}

func ReadMOBI(reader io.Reader) (*MOBI, error) {

	var header [232]byte

	err := binary.Read(reader, binary.BigEndian, &header)

	if err != nil {
		return nil, fmt.Errorf("unable to read MOBI header: %v", err)
	}

	headerLength := convert.FromUint32(header[4:8])

	skip := make([]byte, headerLength-232)

	err = binary.Read(reader, binary.BigEndian, &skip)

	if err != nil {
		return nil, fmt.Errorf("unable to read to end of MOBI header: %v", err)
	}

	return &MOBI{
		EXTHHeaderPresent: (convert.FromUint32(header[112:116]) & 0x40) != 0,
		FirstContentIndex: convert.FromUint16(header[176:178]),
		FirstImageIndex:   convert.FromUint32(header[92:96]),
		FirstNonBookIndex: convert.FromUint32(header[64:68]),
	}, nil
}
