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

func parseArgs(c *config, args []string) (*flag.FlagSet, error) {
	name := os.Args[0]
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	fs.StringVar(&c.path, "path", "", "file path")
	fs.BoolVar(&c.help, "help", false, "show help")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage %s [options]\n", name)
		fmt.Fprintln(os.Stderr, "Options:")
		fs.PrintDefaults()
	}

	err := fs.Parse(args)
	return fs, err
}

func main() {
	var c config
	args := os.Args[1:]

	fs, err := parseArgs(&c, args)

	if c.help {
		fs.Usage()
		os.Exit(0)
	}

	if err != nil {
		fmt.Println(err)
		fs.Usage()
		os.Exit(1)
	}

	var r *os.File = os.Stdin
	if c.path != "" {
		r, err = os.Open(c.path)
		if err != nil {
			e := fmt.Errorf("Could not open file %s", c.path)
			fmt.Fprintln(os.Stderr, e)
			os.Exit(1)
		}
	}

	var sb, errSb strings.Builder

	aggregate.Process(r, &sb, &errSb)
	fmt.Println()
	fmt.Println("Aggregation:")
	fmt.Println(sb.String())
}
