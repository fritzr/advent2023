package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

type Span [2]int

func (s Span) Contains(value int) bool {
	return s[0] <= value && value < s[1]
}

func (s Span) Overlaps(t Span) bool {
	return !(t[1] <= s[0] || t[0] >= s[1])
}

func (s Span) String() string {
	return fmt.Sprintf("[%d,%d)", s[0], s[1])
}

type RangeSet struct {
	set []Span
}

// bisectRow returns the index of the next gridInfo for a given column
func (s RangeSet) bisect(value int) int {
	return sort.Search(len(s.set), func(index int) bool {
		return s.set[index].span.Contains(value) || s.set[index].span[0] > value
	})
}

func (s RangeSet) Insert(span Span) {
	index := s.bisect(span[0])
	if index == len(s.set) {
		s.set = append(s.set, span)
	}
	s.set = append(s.set[:index], span, s.set[index:])
}
