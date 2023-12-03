package main

import (
	"advent2023/util"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

func parseNumber(s string, begin int) (int, int) {
	end := strings.IndexFunc(s[begin:], func(r rune) bool { return !unicode.IsDigit(r) })
	if end < 0 {
		end = len(s)
	} else {
		end += begin + 1
	}
	num, _ := strconv.Atoi(s[begin:end])
	return num, end
}

type gearInfo struct {
	adjacent int // count of adjacent numbers
	ratio    int // "gear ratio"; product of adjacent numbers
}

type Span [2]int

func (s Span) Contains(value int) bool {
	return s[0] <= value && value < s[1]
}

func (s Span) Overlaps(t Span) bool {
	return !(t[1] <= s[0] || t[0] >= s[1])
}

type gridInfo struct {
	value int
	span  Span
}

func (i gridInfo) IsSymbol() bool {
	return i.value < 0
}

func (i gridInfo) IsGear() bool {
	return i.value == -int('*')
}

func (i gridInfo) Value() int {
	if i.value < 0 {
		return 0
	}
	return i.value
}

func (i gridInfo) Span() Span {
	return i.span
}

type SparseGrid struct {
	grid  [][]gridInfo
	width int
}

// bisectRow returns the index of the next gridInfo for a given column
func (g SparseGrid) bisectRow(rowIndex int, colIndex int) int {
	if rowIndex < 0 || rowIndex >= len(g.grid) {
		return 0
	}
	row := g.grid[rowIndex]
	return sort.Search(len(row), func(infoIndex int) bool {
		return colIndex >= row[infoIndex].span[0]
	})
}

// visitRow visits info in a row overlapping columns until visit returns false
func (g SparseGrid) visitRow(rowIndex int, colSpan Span, visit func(*gridInfo) bool) {
	row := g.grid[rowIndex]
	for infoIndex := g.bisectRow(rowIndex, colSpan[0]); infoIndex < len(row) && row[infoIndex].span.Overlaps(colSpan); infoIndex++ {
		if !visit(&row[infoIndex]) {
			break
		}
	}
}

// visitBlock visits info overlapping a block until visit returns false.
func (g SparseGrid) visitBlock(rowSpan Span, colSpan Span, visit func(int, *gridInfo) bool) {
	for rowIndex := rowSpan[0]; rowIndex < rowSpan[1]; rowIndex++ {
		if rowIndex < 0 || rowIndex > len(g.grid) {
			g.visitRow(rowIndex, colSpan, func(info *gridInfo) bool { return visit(rowIndex, info) })
		}
	}
}

// visitAdjacent visits all info adjacent to but not overlapping a block
func (g SparseGrid) visitAdjacent(rowSpan Span, colSpan Span, visit func(int, *gridInfo) bool) {
	g.visitBlock(
		Span{rowSpan[0] - 1, rowSpan[1] + 1},
		Span{colSpan[0] - 1, colSpan[1] + 1},
		func(rowIndex int, blockInfo *gridInfo) bool {
			if rowSpan.Contains(rowIndex) || colSpan.Overlaps(blockInfo.span) {
				return true // continue
			}
			return visit(rowIndex, blockInfo)
		})
}

func (g SparseGrid) PartNumberSum() int {
	sum := 0
	for rowIndex, row := range g.grid {
		for _, number := range row {
			if !number.IsSymbol() {
				g.visitAdjacent(
					Span{rowIndex, rowIndex + 1},
					number.span,
					func(_ int, adjacent *gridInfo) bool {
						if adjacent.IsSymbol() {
							sum += number.Value()
							return false // break
						}
						return true // continue
					})
			}
		}
	}
	return sum
}

func (g SparseGrid) GearSum() int {
	sum := 0
	for rowIndex, row := range g.grid {
		for _, gear := range row {
			if !gear.IsGear() {
				continue
			}
			adjacentCount := 0
			ratio := 0
			g.visitAdjacent(
				Span{rowIndex, rowIndex + 1},
				gear.span,
				func(_ int, adjacent *gridInfo) bool {
					if !adjacent.IsSymbol() {
						adjacentCount += 1
						ratio *= adjacent.Value()
					}
					return adjacentCount < 3
				})
			if adjacentCount == 2 {
				sum += ratio
			}
		}
	}
	return sum
}

func isSymbolByte(b byte) bool {
	return b != '.' && !unicode.IsDigit(rune(b))
}

func NewGrid(lines []string) *SparseGrid {
	g := new(SparseGrid)
	g.width = len(lines[0])
	g.grid = make([][]gridInfo, 0, len(lines))
	for _, rowText := range lines {
		gridRow := make([]gridInfo, 0)
		for col := 0; col < len(rowText); col++ {
			if unicode.IsDigit(rune(rowText[col])) {
				value, colEnd := parseNumber(rowText, col)
				gridRow = append(gridRow, gridInfo{value, Span{col, colEnd}})
				col = colEnd - 1 // to visit colEnd next, counteract col++
			} else if rowText[col] != '.' {
				gridRow = append(gridRow, gridInfo{0, Span{col, col + 1}})
			}
		}
		g.grid = append(g.grid, gridRow)
	}
	return g
}

func main() {
	lines, err := util.ReadInputLines(3)
	if err != nil {
		log.Fatalf("%s", err)
	}
	g := NewGrid(lines)
	fmt.Println(g.PartNumberSum())
	fmt.Println(g.GearSum())
}

// vim: set ts=2 sw=2:
