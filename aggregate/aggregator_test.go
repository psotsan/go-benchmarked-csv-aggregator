package aggregate

import (
	"fmt"
	"strings"
	"testing"
)

func TestAggregateLines(t *testing.T) {
	als := []struct {
		name   string
		lines  []string
		want   map[string]float64
		errors int
	}{
		{
			name:  "single valid line",
			lines: []string{"sales,10"},
			want:  map[string]float64{"sales": 10},
		},
		{
			name:  "single valid line with spaces",
			lines: []string{"  sales  , 10 "},
			want:  map[string]float64{"sales": 10},
		},
		{
			name: "3 valid lines with same key",
			lines: []string{
				"sales, 10",
				"sales, 3",
				"sales, 2",
			},
			want: map[string]float64{"sales": 15},
		},
		{
			name: "5 valid lines with 2 keys",
			lines: []string{
				"sales, 10",
				"purchases, 19",
				"sales, 3",
				"purchases, 11",
				"sales, 2",
			},
			want: map[string]float64{
				"sales":     15,
				"purchases": 30,
			},
		},
		{
			name: "3 invalid float lines",
			lines: []string{
				"sales, twenty",
				"sales, 7a6",
				"purchases, ^23",
			},
			want:   map[string]float64{},
			errors: 3,
		},
		{
			name: "3 invalid float lines and 2 valid lines",
			lines: []string{
				"purchases, 20",
				"sales, twenty",
				"sales, 10",
				"sales, 7a6",
				"purchases, ^23",
			},
			want: map[string]float64{
				"purchases": 20,
				"sales":     10,
			},
			errors: 3,
		},
		{
			name: "3 valid lines and 2 invalid length lines",
			lines: []string{
				"purchases, 20, 50",
				"sales, 3",
				"purchases, 20",
				"sales, 7",
				"purchases",
			},
			want: map[string]float64{
				"purchases": 20,
				"sales":     10,
			},
			errors: 2,
		},
		{
			name: "3 valid lines and 2 invalid lines with empty fields",
			lines: []string{
				",50",
				"sales, 3",
				"purchases, 20",
				"sales, 7",
				"purchases,",
			},
			want: map[string]float64{
				"purchases": 20,
				"sales":     10,
			},
			errors: 2,
		},
	}

	for _, tt := range als {
		t.Run(tt.name, func(t *testing.T) {
			sErrs := make([]error, 0)
			got := AggregateLines(tt.lines, &sErrs)
			gotErr := len(sErrs)
			if gotErr != tt.errors {
				for _, err := range sErrs {
					fmt.Println(err)
				}
				t.Fatalf("AggregateLines(%q): expected %d errors, got %d errors", tt.name, tt.errors, gotErr)
			}

			if len(got) != len(tt.want) {
				t.Fatalf("AggregateLines(%q): got and want lengths differ\nlen(got): %d\nlen(want): %d",
					tt.name, len(got), len(tt.want))
			}

			for key, val := range got {
				if tt.want[key] != val {
					t.Errorf("AggregateLines(%q): expected %s -> %.1f, got %s -> %.1f",
						tt.name, key, tt.want[key], key, val)
				}
			}
		})
	}
}

func TestProcess(t *testing.T) {
	var r1 *strings.Reader = strings.NewReader(`sales, 10
	purchases, 19
	sales, 3
	purchases, 11
	sales, 2`)

	var r2 *strings.Reader = strings.NewReader(`purchases, 20
	sales, twenty
	sales, 10
	sales, 15, 80
	purchases, ^23`)

	var r3 *strings.Reader = strings.NewReader(strings.Repeat("sales, 10\n", lineLimit+1))

	ps := []struct {
		name    string
		r       *strings.Reader
		want    string
		wantMap map[string]float64
		errors  int
	}{
		{
			name: "5 valid lines with 2 keys",
			r:    r1,
			want: "purchases: 30.0\nsales: 15.0\n",
			wantMap: map[string]float64{
				"sales":     15,
				"purchases": 30,
			},
		},
		{
			name: "3 invalid  lines and 2 valid lines",
			r:    r2,
			want: "purchases: 20.0\nsales: 10.0\n",
			wantMap: map[string]float64{
				"purchases": 20,
				"sales":     10,
			},
			errors: 3,
		},
		{
			name:    "limit exceeded",
			r:       r3,
			wantMap: nil,
			errors:  1,
		},
	}

	for _, tt := range ps {
		t.Run(tt.name, func(t *testing.T) {
			var sb, errSb strings.Builder
			var gotErrW int
			got, errs := Process(tt.r, &sb, &errSb)
			gotErr := len(errs)
			errW := strings.Trim(errSb.String(), "\n")

			if errSb.Len() == 0 {
				gotErrW = 0
			} else {
				gotErrW = len(strings.Split(errW, "\n"))
			}

			if gotErr != tt.errors {
				for _, err := range errs {
					fmt.Println(err)
				}
				t.Fatalf("Process(%q): expected %d returned errors, got %d errors", tt.name, tt.errors, gotErr)
			}

			if gotErrW != tt.errors {
				fmt.Printf("%v", errW)
				t.Fatalf("Process(%q): expected %d errors in Writer, got %d errors", tt.name, tt.errors, gotErrW)
			}

			if len(got) != len(tt.wantMap) {
				t.Fatalf("Process(%q): got and want lengths differ\nlen(got): %d\nlen(want): %d",
					tt.name, len(got), len(tt.wantMap))
			}

			for key, val := range got {
				if tt.wantMap[key] != val {
					t.Errorf("Process(%q): expected %s -> %.1f, got %s -> %.1f",
						tt.name, key, tt.wantMap[key], key, val)
				}
			}

			if sb.String() != tt.want {
				t.Errorf("Process(%q): expected %s\ngot %s", tt.name, tt.want, sb.String())
			}
		})
	}
}
