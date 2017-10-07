package mobi

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/neofight/mobi/convert"
	"github.com/neofight/mobi/headers"
	"golang.org/x/net/html"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"
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

func (mobiFile Book) Markup() (string, error) {

	startIndex := mobiFile.mobiHeader.FirstContentIndex
	endIndex := mobiFile.mobiHeader.FirstNonBookIndex - 1

	text := make([]byte, 0)

	for index := startIndex; index <= endIndex; index++ {

		record := mobiFile.pdbHeader.Records[index]
		nextRecord := mobiFile.pdbHeader.Records[index+1]

		recordOffset := record.RecordDataOffset
		recordSize := nextRecord.RecordDataOffset - recordOffset

		_, err := mobiFile.file.Seek(int64(recordOffset), 0)

		if err != nil {
			return "", fmt.Errorf("unable to find text: %v", err)
		}

		recordData := make([]byte, recordSize)

		err = binary.Read(mobiFile.file, binary.BigEndian, &recordData)

		if err != nil {
			return "", fmt.Errorf("unable to read text: %v", err)
		}

		recordText := convert.FromLZ77(recordData)

		text = append(text, recordText...)
	}

	text = text[:mobiFile.palmDOCHeader.TextLength]

	if !utf8.Valid(text) {
		return "", errors.New("unable to decompress text")
	}

	return string(text), nil
}

func (mobiFile Book) Text() (string, error) {

	markup, err := mobiFile.Markup()

	if err != nil {
		return "", fmt.Errorf("unable to read markup: %v", err)
	}

	pos, err := getTOCPosition(markup)

	if err != nil {
		return "", fmt.Errorf("unable to locate TOC: %v", err)
	}

	bookmarks, err := parseTOC(markup[pos:])

	text := make([]string, 0)

	for i := range bookmarks {

		start := bookmarks[i]
		var end int

		if i < len(bookmarks)-1 {
			end = bookmarks[i+1]
		} else {
			end = pos
		}

		paragraphs, err := parseChapter(markup[start:end])

		if err != nil {
			return "", fmt.Errorf("unable to parse chapter: %v", err)
		}

		text = append(text, paragraphs...)
	}

	return strings.Join(text, "\n\n"), nil
}

func getTOCPosition(markup string) (int, error) {

	htmlReader := strings.NewReader(markup)

	tokenizer := html.NewTokenizer(htmlReader)

	for {
		tokenType := tokenizer.Next()

		switch {
		case tokenType == html.ErrorToken:
			return 0, fmt.Errorf("unable to find reference element")
		case tokenType == html.SelfClosingTagToken:
			token := tokenizer.Token()

			if token.Data == "reference" {
				filepos, err := attr(token, "filepos")

				if err != nil {
					return 0, errors.New("filepos attribute missing")
				}

				pos, err := strconv.Atoi(filepos)

				if err != nil {
					return 0, errors.New("filepos attribute invalid")
				}

				return pos, nil
			}
		}
	}
}

func parseTOC(markup string) ([]int, error) {

	toc := make([]int, 0)

	htmlReader := strings.NewReader(markup)

	tokenizer := html.NewTokenizer(htmlReader)

	for {
		tokenType := tokenizer.Next()

		switch {
		case tokenType == html.ErrorToken:
			return toc[1:], nil
		case tokenType == html.StartTagToken:
			token := tokenizer.Token()

			if token.Data == "a" {
				filepos, err := attr(token, "filepos")

				if err != nil {
					continue
				}

				pos, err := strconv.Atoi(filepos)

				if err != nil {
					return nil, errors.New("filepos attribute invalid")
				}

				toc = append(toc, pos)
			}
		}
	}
}

func parseChapter(markup string) ([]string, error) {

	paragraphs := make([]string, 0)

	htmlReader := strings.NewReader(markup)

	tokenizer := html.NewTokenizer(htmlReader)

	for {
		tokenType := tokenizer.Next()

		switch {
		case tokenType == html.ErrorToken:
			return paragraphs, nil
		case tokenType == html.TextToken:
			token := tokenizer.Token()

			if len(strings.TrimSpace(token.Data)) > 0 {
				paragraphs = append(paragraphs, strings.TrimSpace(token.Data))
			}
		}
	}
}

func attr(t html.Token, name string) (string, error) {
	for _, a := range t.Attr {
		if a.Key == name {
			return a.Val, nil
		}
	}

	return "", fmt.Errorf("attribute %v not found", name)
}
