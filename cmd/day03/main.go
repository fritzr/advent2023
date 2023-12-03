package main

import (
	"advent2023/util"
	"fmt"
	"log"
	"strconv"
	"strings"
	"unicode"
)

func clampLo(value, lower int) int {
	if value < lower {
		value = lower
	}
	return value
}

func clampHi(value, upper int) int {
	if value > upper {
		value = upper
	}
	return value
}

func findAdjacent(grid []string, row int, col int, colEnd int, search func(byte) bool) [2]int {
	rowStart := clampLo(row-1, 0)
	rowEnd := clampHi(row+2, len(grid))
	row = rowStart
	colStart := clampLo(col-1, 0)
	colEnd = clampHi(colEnd+1, len(grid[0]))
	// fmt.Printf("searching:")
	for row = rowStart; row < rowEnd; row++ {
		for col := colStart; col < colEnd; col++ {
			// fmt.Printf(" %d,%d", row, col)
			if search(grid[row][col]) {
				// fmt.Printf("\n           found adjacent '%c' at (%d, %d)\n", grid[row][col], row, col)
				return [2]int{row, col}
			}
		}
		// fmt.Printf("\n          ")
	}
	// fmt.Printf(" not found\n")
	return [2]int{-1, -1}
}

func parseNumber(s string, begin int) (int, int) {
	end := strings.IndexFunc(s[begin:], func(r rune) bool { return !unicode.IsDigit(r) })
	if end < 0 {
		end = len(s)
	} else {
		end += begin
	}
	num, _ := strconv.Atoi(s[begin:end])
	return num, end
}

type gearInfo struct {
	count int
	ratio int
}

func isSymbol(b byte) bool {
	return b != '.' && !unicode.IsDigit(rune(b))
}

func main() {
	lines, err := util.ReadInputLines(3)
	if err != nil {
		log.Fatalf("%s", err)
	}

	gears := map[[2]int]*gearInfo{}
	sum := 0
	sumGears := uint64(0)
	for row, line := range lines {
		/*
			if row > 0 {
				fmt.Printf("    %s\n", lines[row-1])
			}
			fmt.Printf("==> %s\n", line)
			if row < len(line)-1 {
				fmt.Printf("    %s\n", lines[row+1])
			}
		*/
		for col := 0; col < len(line); col++ {
			if !unicode.IsDigit(rune(line[col])) {
				continue
			}
			number, colEnd := parseNumber(line, col)
			// fmt.Printf("number %d in row %d spans [%d, %d)\n", number, row, col, colEnd)
			found := findAdjacent(lines, row, col, colEnd, isSymbol)
			col = colEnd
			if found[0] < 0 {
				continue
			}
			sum += number
			if lines[found[0]][found[1]] == '*' {
				info, ok := gears[found]
				if !ok {
					info = &gearInfo{count: 0, ratio: number}
					gears[found] = info
				}
				info.count += 1
				switch count := info.count; count {
				case 1: // first number, remember to produce ratio
					info.ratio = number
				case 2: // just enough numbers, compute ratio and add to sum
					info.ratio *= number
					sumGears += uint64(info.ratio)
				case 3: // too many adjacent numbers, remove ratio from sum
					sumGears -= uint64(info.ratio)
				}
			}
		}
	}
	fmt.Println(sum)
	fmt.Println(sumGears)
}

// vim: set ts=2 sw=2:
