package main

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/scgolang/dx7/sysex"
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage:\n")
	fmt.Fprintf(os.Stderr, "  %s [OPTIONS] [FILE]\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "If FILE is not provided, sysex data will be read from stdin.\n")
	fmt.Fprintf(os.Stderr, "OPTIONS\n")
	fmt.Fprintf(os.Stderr, "  -format          xml|json\n")
}

func main() {
	format := flag.String("format", "xml", "output format")
	flag.Parse()

	if *format != "xml" && *format != "json" {
		usage()
		os.Exit(1)
	}

	args := flag.Args()
	switch len(args) {
	case 0:
		// Read sysex data from stdin.
		if err := run(os.Stdin, os.Stdout, *format); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	case 1:
		// Read sysex data from a file.
		r, err := os.Open(args[0])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		if err := run(r, os.Stdout, *format); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	default:
		usage()
		os.Exit(1)
	}
}

// run is the heart of the program.
func run(r io.Reader, w io.Writer, format string) error {
	syx, err := sysex.New(r)
	if err != nil {
		return err
	}

	switch format {
	default:
		return fmt.Errorf("Unrecognized format: %s", format)
	case "xml":
		return xml.NewEncoder(w).Encode(syx)
	case "json":
		return json.NewEncoder(w).Encode(syx)
	}
}
