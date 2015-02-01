package main

import (
	"bufio"
	"encoding/csv"
	"io"
	"os"
	"strings"
)

func Process(in io.Reader, out io.Writer) {
	tsvReader := csv.NewReader(in)
	tsvReader.Comma = '\t'

	outWriter := bufio.NewWriter(out)
	defer outWriter.Flush()

	for record, _ := tsvReader.Read(); record != nil; record, _ = tsvReader.Read() {
		outWriter.WriteString(strings.Join(record, "\t"))
		outWriter.WriteString("\n")
	}
}

func main() {
	Process(os.Stdin, os.Stdout)
}
