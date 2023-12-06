package util

import (
	"fmt"
	"sort"
)

type Span [2]int
type spanValue[T any] struct {
	span  Span
	value T
}

func (s Span) Contains(value int) bool {
	return s[0] <= value && value < s[1]
}

func (s Span) Overlaps(t Span) bool {
	return !(t[1] <= s[0] || t[0] >= s[1])
}

func (s Span) String() string {
	return fmt.Sprintf("[%d,%d)", s[0], s[1])
}

type RangeSet[T any] struct {
	set []spanValue[T]
}

// bisectRow returns the index of the next gridInfo for a given column
func (s RangeSet[T]) bisect(value int) int {
	return sort.Search(len(s.set), func(index int) bool {
		return s.set[index].span.Contains(value) || s.set[index].span[0] > value
	})
}

func (s *RangeSet[T]) Add(span Span, value T) *T {
	index := s.bisect(span[0])
	if index == len(s.set) {
		s.set = append(s.set, spanValue[T]{span, value})
	} else {
		s.set = append(s.set[:index+1], s.set[index:]...)
		s.set[index] = spanValue[T]{span, value}
	}
	return &s.set[index].value
}

func (s RangeSet[T]) Get(key int) *T {
	index := s.bisect(key)
	if index < len(s.set) && s.set[index].span.Contains(key) {
		return &s.set[index].value
	}
	return nil
}

type RangeResult[T any] struct {
	Span  Span
	Value *T
}

func (s RangeSet[T]) Intersect(span Span) []RangeResult[T] {
	index := s.bisect(span[0])
	results := make([]RangeResult[T], 0)
	for index < len(s.set) && s.set[index].span.Overlaps(span) {
		results = append(results, RangeResult[T]{s.set[index].span, &s.set[index].value})
	}
	return results
}
