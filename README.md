[![Go Report Card](https://goreportcard.com/badge/github.com/psotsan/go-benchmarked-csv-aggregator)](https://goreportcard.com/report/github.com/psotsan/go-benchmarked-csv-aggregator)
# CSV Aggregator

Go CLI tool to aggregate numeric values from a two‑column CSV file.  
Reads from stdin or a file, groups by the first column, and sums the second column.

---

## Quick start

```bash
# Build
go build -o csv-aggregator

# Run with stdin
echo "sales,10\npurchases,20\nsales,5" | ./csv-aggregator

# Run with a file
./csv-aggregator -path data.csv
```

---

## File format

- *Two columns separated by a comma: key,value*

- *Empty lines and lines starting with # are ignored.*

- *Keys are trimmed of spaces/tabs.*

- *Values can be integers or decimals.*

- *Lines with wrong format (wrong number of columns or non‑numeric value) are skipped – errors go to stderr.*

---

## Example input
```
sales, 10
purchases, 20
sales, 5
purchases, abc   # invalid, will be skipped
```

## Flags
- -path     Read from file instead of stdin
- -help     Show usage help

---

## Example output
```bash
$ ./csv-aggregator -path data.csv

Aggregation:
purchases: 20.0
sales: 15.0
```

---

## Tests and benchmarks
```bash
# Run all tests
go test -v ./...

# Run tests with coverage
go test -cover ./...

# Run benchmarks (100k lines)
go test -bench=. ./aggregate
```