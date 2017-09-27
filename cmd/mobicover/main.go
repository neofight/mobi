package main

import (
	"github.com/neofight/mobi"
	"log"
	"os"
)

func main() {

	const path = "test.prc"

	book, err := mobi.Open(path)
	defer book.Close()

	if err != nil {
		log.Fatalf("failed to open %v: %v", path, err)
	}

	cover, err := book.Cover()

	if err != nil {
		log.Fatalf("failed to read the cover: %v", err)
	}

	coverFile, err := os.Create("cover.gif")
	defer coverFile.Close()

	if err != nil {
		log.Fatalf("failed to create cover.gif: %v", err)
	}

	coverFile.Write(cover)
}
