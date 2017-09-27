package mobi

import (
	"encoding/binary"
	"fmt"
	"github.com/neofight/mobi/convert"
	"github.com/neofight/mobi/headers"
	"os"
)

type Book struct {
	file          *os.File
	pdbHeader     *headers.PDB
	palmDOCHeader *headers.PalmDOC
	mobiHeader    *headers.MOBI
	exthHeader    *headers.EXTH
}

func Open(path string) (*Book, error) {

	var book Book

	var err error

	book.file, err = os.Open(path)

	if err != nil {
		return nil, fmt.Errorf("unable to open %v: %v", path, err)
	}

	book.pdbHeader, err = headers.ReadPDB(book.file)

	if err != nil {
		return nil, fmt.Errorf("unable to read PDB header: %v", err)
	}

	book.palmDOCHeader, err = headers.ReadPalmDOC(book.file)

	if err != nil {
		return nil, fmt.Errorf("unable to read PalmDOC header: %v", err)
	}

	book.mobiHeader, err = headers.ReadMOBI(book.file)

	if err != nil {
		return nil, fmt.Errorf("unable to read MOBI header: %v", err)
	}

	if book.mobiHeader.EXTHHeaderPresent {

		book.exthHeader, err = headers.ReadEXTH(book.file)

		if err != nil {
			return nil, fmt.Errorf("unable to read EXTH header: %v", err)
		}
	}

	return &book, nil
}

func (mobiFile Book) Close() error {
	return mobiFile.file.Close()
}

func (mobiFile Book) Cover() ([]byte, error) {

	for _, r := range mobiFile.exthHeader.Records {

		if r.RecordType == 201 {
			coverIndex := mobiFile.mobiHeader.FirstImageIndex + convert.FromUint32(r.RecordData)

			record := mobiFile.pdbHeader.Records[coverIndex]
			nextRecord := mobiFile.pdbHeader.Records[coverIndex+1]

			coverOffset := record.RecordDataOffset
			coverSize := nextRecord.RecordDataOffset - coverOffset

			_, err := mobiFile.file.Seek(int64(coverOffset), 0)

			if err != nil {
				return nil, fmt.Errorf("unable to find cover: %v", err)
			}

			cover := make([]byte, coverSize)

			err = binary.Read(mobiFile.file, binary.BigEndian, &cover)

			if err != nil {
				return nil, fmt.Errorf("unable to read cover: %v", err)
			}

			return cover, nil
		}
	}

	return nil, nil
}
