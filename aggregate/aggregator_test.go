package aggregate

import (
	"fmt"
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
			got, errs := AggregateLines(tt.lines)
			gotErr := len(errs)
			if gotErr != tt.errors {
				for _, err := range errs {
					fmt.Println(err)
				}
				t.Fatalf("AggregateLines(%q): expected %d errors, got %d errors", tt.name, tt.errors, gotErr)
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
