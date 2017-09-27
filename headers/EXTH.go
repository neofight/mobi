package headers

import (
	"encoding/binary"
	"fmt"
	"github.com/neofight/mobi/convert"
	"io"
)

type EXTH struct {
	Records []*EXTHRecord
}

func ReadEXTH(reader io.Reader) (*EXTH, error) {

	var header [12]byte

	err := binary.Read(reader, binary.BigEndian, &header)

	if err != nil {
		return nil, fmt.Errorf("unable to read EXTH header: %v", err)
	}

	recordCount := convert.FromUint32(header[8:12])

	records := make([]*EXTHRecord, recordCount)

	for i := range records {

		records[i], err = ReadEXTHRecord(reader)

		if err != nil {
			return nil, fmt.Errorf("unable to read EXTH record %v: %v", i, err)
		}
	}

	headerLength := convert.FromUint32(header[4:8])

	paddingLength := 4 - (headerLength % 4)

	if paddingLength != 4 {

		skip := make([]byte, paddingLength)

		binary.Read(reader, binary.BigEndian, &skip)

		if err != nil {
			return nil, fmt.Errorf("unable to read to end of EXTH record: %v", err)
		}
	}

	return &EXTH{records}, nil
}
