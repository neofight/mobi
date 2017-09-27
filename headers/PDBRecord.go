package headers

import (
	"encoding/binary"
	"fmt"
	"io"
	"github.com/neofight/mobi/convert"
)

type PDBRecord struct {
	RecordDataOffset int
}

func ReadPDBRecord(reader io.Reader) (*PDBRecord, error) {

	var header [8]byte

	err := binary.Read(reader, binary.BigEndian, &header)

	if err != nil {
		return nil, fmt.Errorf("unable to read PDB record: %v", err)
	}

	return &PDBRecord{
		RecordDataOffset: convert.FromUint32(header[0:4]),
	}, nil
}
