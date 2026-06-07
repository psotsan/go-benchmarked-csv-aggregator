package aggregate

import (
	"bufio"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
)

const (
	separator = ","
	lineLimit = 100_000
)

func splitValidateLine(line string) ([]string, error) {
	var err error
	s := strings.Split(line, separator)

	if len(s) != 2 {
		err = fmt.Errorf("expected length 2. Got length %d", len(s))
		return nil, err
	}

	s[0] = strings.Trim(s[0], " ")
	s[0] = strings.Trim(s[0], "\t")
	s[1] = strings.Trim(s[1], " ")
	s[1] = strings.Trim(s[1], "\t")

	if s[0] == "" || s[1] == "" {
		err = fmt.Errorf("empty field(s) found")
		return nil, err
	}

	return s, nil
}

func keysToStrSlice(m map[string]float64) []string {
	var s []string

	for key := range m {
		s = append(s, key)
	}

	return s
}

func aggregateLines(lines []string, errs *[]error) map[string]float64 {
	// processes a batch of lines, skipping empty lines and comments (#).
	// For each valid CSV line, it accumulates the numeric value under the key.
	// Errors are appended to the provided slice.
	agg := make(map[string]float64)
	n := 0

	for _, l := range lines {
		n++

		l = strings.TrimLeft(l, " ")
		l = strings.TrimLeft(l, "\t")
		l = strings.TrimLeft(l, " ")

		if l == "" || strings.HasPrefix(l, "#") {
			continue
		}

		s, err := splitValidateLine(l)
		if err != nil {
			*errs = append(*errs, fmt.Errorf("line %d: %v", n, err))
			continue
		}

		key := strings.Trim(s[0], " ")
		val, err := strconv.ParseFloat(s[1], 64)
		if err != nil {
			*errs = append(*errs, err)
			continue
		}

		if _, ok := agg[key]; ok {
			agg[key] += val
			continue
		}

		agg[key] = val
	}

	return agg
}

// Process reads the input and aggregates the lines. It writes to stdout, stderr and returns the map of values and slice of errors.
func Process(r io.Reader, w io.Writer, errW io.Writer) (map[string]float64, []error) {
	var batch []string
	errs := make([]error, 0)
	scanner := bufio.NewScanner(r)

	n := 0
	for scanner.Scan() {
		n++
		if n > lineLimit {
			err := fmt.Errorf("Process: r exceeds maximum line limit")
			errs = append(errs, err)
			if _, err := fmt.Fprintln(errW, err); err != nil {
				errs = append(errs, fmt.Errorf("failed to write result: %w", err))
			}
			return nil, errs
		}
		batch = append(batch, scanner.Text())
	}

	ret := aggregateLines(batch, &errs)

	for _, err := range errs {
		if _, e := fmt.Fprintln(errW, err); e != nil {
			errs = append(errs, fmt.Errorf("failed to write result: %w", e))
		}
	}

	sl := keysToStrSlice(ret)
	sort.Strings(sl)

	for _, s := range sl {
		if _, err := fmt.Fprintf(w, "%s: %.1f\n", s, ret[s]); err != nil {
			errs = append(errs, fmt.Errorf("failed to write result: %w", err))
		}
	}

	return ret, errs
}
