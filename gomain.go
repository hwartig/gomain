package main

import (
	"bufio"
	"encoding/csv"
	"github.com/oschwald/geoip2-golang"
	"golang.org/x/net/publicsuffix"
	"io"
	"log"
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

func AppendServerCountry(ipDB *geoip2.Reader) recordTransformer {
	return func(in []string) []string {
		ip, err := net.LookupIP(in[0])

		if err == nil {
			country, _ := ipDB.Country(ip[0])
			in = append(in, country.Country.IsoCode)
		} else {
			in = append(in, "")
		}
		return in
	}
}

func AppendTLD(in []string) []string {
	tld, _ := publicsuffix.PublicSuffix(in[0])

	return append(in, tld)
}

func chain(first, second recordTransformer) recordTransformer {
	return func(in []string) []string {
		return second(first(in))
	}
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
	db, err := geoip2.Open("GeoIP2-Country.mmdb")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	Process(
		os.Stdin,
		os.Stdout,
		chain(
			AppendServerCountry(db),
			AppendTLD,
		),
	)
}
