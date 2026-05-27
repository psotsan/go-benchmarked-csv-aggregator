package aggregate

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

const separator = ","

func splitValidateLine(line string) ([]string, error) {
	var err error
	s := strings.Split(line, separator)

	if len(s) != 2 {
		err = fmt.Errorf("expected length 2. Got length %d", len(line))
		return nil, err
	}

	s[0] = strings.Trim(s[0], " ")
	s[1] = strings.Trim(s[1], " ")

	if s[0] == "" || s[1] == "" {
		err = fmt.Errorf("empty field(s) found")
		return nil, err
	}

	return s, nil
}

func AggregateLines(lines []string, errs *[]error) map[string]float64 {
	agg := make(map[string]float64)
	// errs := make([]error, 0)
	n := 0

	for _, l := range lines {
		n++
		s, err := splitValidateLine(l)
		if err != nil {
			*errs = append(*errs, fmt.Errorf("Line %d: %v", n, err))
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

func Process(r io.Reader) (map[string]float64, []error) {
	const lineLimit = 100_000
	var batch []string
	errs := make([]error, 0)
	scanner := bufio.NewScanner(r)

	n := 0
	for scanner.Scan() {
		n++
		if n > lineLimit {
			err := fmt.Errorf("Process: r exceeds maximum line limit")
			errs = append(errs, err)
			return nil, errs
		}
		batch = append(batch, scanner.Text())
	}

	ret := AggregateLines(batch, &errs)
	return ret, errs
}
