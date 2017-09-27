package headers

import (
	"encoding/binary"
	"fmt"
	"io"
	"github.com/neofight/mobi/convert"
)

type PDB struct {
	Records []*PDBRecord
}

func ReadPDB(reader io.Reader) (*PDB, error) {

	var header [78]byte

	err := binary.Read(reader, binary.BigEndian, &header)

	if err != nil {
		return nil, fmt.Errorf("unable to read PDB header: %v", err)
	}

	recordCount := convert.FromUint16(header[76:78])

	records := make([]*PDBRecord, recordCount)

	for i := range records {

		records[i], err = ReadPDBRecord(reader)

		if err != nil {
			return nil, fmt.Errorf("unable to read PDB record %v: %v", i, err)
		}
	}

	var gapToData [2]byte

	err = binary.Read(reader, binary.BigEndian, &gapToData)

	if err != nil {
		return nil, fmt.Errorf("unable to read to end of PDB header: %v", err)
	}

	return &PDB{records}, nil
}
