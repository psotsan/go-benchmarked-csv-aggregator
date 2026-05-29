package main

import (
	"os"
	"strings"
	"testing"
)

func TestRun(t *testing.T) {
	helpStr := "usage: " + os.Args[0] + " [options]\nOptions:\n  -help\n    \tshow help\n  -path string\n    \tfile path\n"

	rs := []struct {
		name     string
		args     []string
		r        strings.Reader
		w        strings.Builder
		errW     strings.Builder
		wWant    string
		errWWant string
	}{
		{
			name: "valid input",
			r: *strings.NewReader(`sales, 10
	purchases, 19
	sales, 3
	purchases, 11
	sales, 2`),
			wWant: `
Aggregation:
purchases: 30.0
sales: 15.0`,
		},
		{
			name:     "help flag",
			args:     []string{"-help"},
			errWWant: helpStr,
		},
		{
			name:     "double dash help flag",
			args:     []string{"--help"},
			errWWant: helpStr,
		},
		{
			name:     "invalid flag",
			args:     []string{"-hel"},
			errWWant: "flag provided but not defined: -hel\n" + helpStr,
		},
		{
			name: "input file",
			args: []string{"-path=test-file.txt"},
			wWant: `
Aggregation:
purchases: 17.0
sales: 70.0`,
		},
		{
			name:     "non-existant inputfile",
			args:     []string{"-path=inexistant-file.txt"},
			errWWant: "Could not open file inexistant-file.txt\n",
		},
		{
			name:  "empty input",
			wWant: "\nAggregation:\n",
		},
		{
			name: "3 invalid lines and 2 valid lines",
			args: []string{"-path=test-valid-and-invalid.txt"},
			wWant: `
Aggregation:
purchases: 20.0
sales: 10.0`,
			errWWant: "strconv.ParseFloat: parsing \"twenty\": invalid syntax\nLine 4: expected length 2. Got length 3\nstrconv.ParseFloat: parsing \"^23\": invalid syntax\n\n",
		},
	}

	for _, tt := range rs {
		t.Run(tt.name, func(t *testing.T) {
			run(tt.args, &tt.r, &tt.w, &tt.errW)
			if trimAll(tt.errW.String()) != trimAll(tt.errWWant) {
				t.Fatalf("run (%q): obtained error output:\n%q\nexpected:\n%q", tt.name, &tt.errW, tt.errWWant)
			}
			if trimAll(tt.w.String()) != trimAll(tt.wWant) {
				t.Fatalf("run (%q): obtained output\n%q\nexpected:\n%q", tt.name, tt.w.String(), tt.wWant)
			}
		})
	}
}
