package headers

import (
	"encoding/binary"
	"fmt"
	"github.com/neofight/mobi/convert"
	"io"
)

type PalmDOC struct {
	TextLength int
}

func ReadPalmDOC(reader io.Reader) (*PalmDOC, error) {

	var header [16]byte

	err := binary.Read(reader, binary.BigEndian, &header)

	if err != nil {
		return nil, fmt.Errorf("unable to read PalmDOC header: %v", err)
	}

	return &PalmDOC{
		TextLength: convert.FromUint32(header[4:8]),
	}, nil
}
