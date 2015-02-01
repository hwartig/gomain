package main

import (
	"bufio"
	"encoding/csv"
	"io"
	"net"
	"os"
	"strings"
)

type recordTransformer func(in []string) []string

func AppendIP(in []string) []string {
	// TODO: find out how to mock LookupIP so we can test this
	ip, err := net.LookupIP(in[0])

	if err == nil {
		in = append(in, ip[0].String())
	} else {
		in = append(in, "")
	}
	return in
}

func Process(in io.Reader, out io.Writer, fn recordTransformer) {
	tsvReader := csv.NewReader(in)
	tsvReader.Comma = '\t'

	outWriter := bufio.NewWriter(out)
	defer outWriter.Flush()

	for record, _ := tsvReader.Read(); record != nil; record, _ = tsvReader.Read() {
		record = fn(record)

		outWriter.WriteString(strings.Join(record, "\t"))
		outWriter.WriteString("\n")
	}
}

func main() {
	Process(os.Stdin, os.Stdout, AppendIP)
}
