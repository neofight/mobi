package main

import (
	"github.com/neofight/mobi"
	"log"
	"os"
)

func main() {

	const path = "test.prc"

	book, err := mobi.Open(path)

	if err != nil {
		log.Fatalf("failed to open %v: %v", path, err)
	}

	defer book.Close()

	text, err := book.Text()

	if err != nil {
		log.Fatalf("failed to read the text: %v", err)
	}

	textFile, err := os.Create("text.txt")

	if err != nil {
		log.Fatalf("failed to create cover.gif: %v", err)
	}

	defer textFile.Close()

	textFile.Write([]byte(text))
}
