package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/psotsan/go-benchmarked-csv-aggregator/aggregate"
)

type config struct {
	path string
	help bool
}

func trimAll(s string) string {
	s = strings.Trim(s, " ")
	s = strings.Trim(s, "\t")
	s = strings.Trim(s, "\n")
	return s
}

func parseArgs(c *config, args []string, errW io.Writer) (*flag.FlagSet, error) {
	name := os.Args[0]
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	fs.SetOutput(errW)

	fs.StringVar(&c.path, "path", "", "file path")
	fs.BoolVar(&c.help, "help", false, "show help")

	fs.Usage = func() {
		fmt.Fprintf(errW, "usage: %s [options]\n", name)
		fmt.Fprintln(errW, "Options:")
		fs.PrintDefaults()
	}

	err := fs.Parse(args)
	return fs, err
}

func run(args []string, r io.Reader, w, errW io.Writer) int {
	var c config

	fs, err := parseArgs(&c, args, errW)
	if err != nil {
		return 1
	}

	if c.help {
		fs.Usage()
		return 0
	}

	if c.path != "" {
		r, err = os.Open(c.path)
		if err != nil {
			e := fmt.Errorf("Could not open file %s", c.path)
			fmt.Fprintln(errW, e)
			return 1
		}
	}

	var sb, errSb strings.Builder

	aggregate.Process(r, &sb, &errSb)
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Aggregation:")
	fmt.Fprintln(w, sb.String())
	fmt.Fprintln(errW, errSb.String())

	return 0
}

func main() {
	args := os.Args[1:]
	r := os.Stdin
	w := os.Stdout
	errW := os.Stderr

	exit := run(args, r, w, errW)

	os.Exit(exit)
}
