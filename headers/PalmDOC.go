package headers

import (
	"encoding/binary"
	"fmt"
	"io"
)

type PalmDOC struct {
}

func ReadPalmDOC(reader io.Reader) (*PalmDOC, error) {

	var header [16]byte

	err := binary.Read(reader, binary.BigEndian, &header)

	if err != nil {
		return nil, fmt.Errorf("unable to read PalmDOC header: %v", err)
	}

	return &PalmDOC{}, nil
}
