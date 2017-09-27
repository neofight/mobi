package headers

import (
	"encoding/binary"
	"fmt"
	"github.com/neofight/mobi/convert"
	"io"
)

type EXTHRecord struct {
	RecordType int
	RecordData []byte
}

func ReadEXTHRecord(reader io.Reader) (*EXTHRecord, error) {

	var header [8]byte

	err := binary.Read(reader, binary.BigEndian, &header)

	if err != nil {
		return nil, fmt.Errorf("unable to read EXTH record: %v", err)
	}

	recordLength := convert.FromUint32(header[4:8])

	recordData := make([]byte, recordLength-8)

	err = binary.Read(reader, binary.BigEndian, &recordData)

	if err != nil {
		return nil, fmt.Errorf("unable to read EXTH record data: %v", err)
	}

	return &EXTHRecord{
		RecordType: convert.FromUint32(header[0:4]),
		RecordData: recordData,
	}, nil
}
