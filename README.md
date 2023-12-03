# Advent of Code 2023

Advent of Code for 2023 written in Go.

Building:

```
go vet ./... && go build ./...
```

Running:

```
go run advent2023/cmd/dayXX [options]
```

Options (run with -h or -help to display usage):

```
Usage of advent2023/cmd/dayXX:
  -name string
    	input filename (relative to day directory) (default "input.txt")
  -path string
    	explicit input path (overrides -name)
```
