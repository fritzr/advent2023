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

func (s Span) Intersect(t Span) Span {
	result := Span(s)
	if t[0] > s[0] {
		result[0] = t[0]
	}
	if t[1] < s[1] {
		result[1] = t[1]
	}
	if result[1] < result[0] {
		result[0] = result[1]
	}
	return result
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

func (s RangeSet[T]) GetRange(key int) *RangeResult[T] {
	index := s.bisect(key)
	if index < len(s.set) && s.set[index].span.Contains(key) {
		return &RangeResult[T]{s.set[index].span, &s.set[index].value}
	}
	return nil
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

// DoIntersectSet invokes a function on the intersection of a range with the set.
//
// If a call returns false, no more ranges are visited.
func (s RangeSet[T]) DoIntersect(t RangeSet[T], do func(ss, ts, ix Span, svalue, tvalue *T) bool) {
	if len(s.set) == 0 || len(t.set) == 0 {
		return
	}
	sIndex := s.bisect(t.set[0].span[0])
	tIndex := 0
	for sIndex < len(s.set) && tIndex < len(t.set) {
		sinfo := &s.set[sIndex]
		tinfo := &t.set[tIndex]
		if tinfo.span[0] >= s.Max() {
			break
		}
		intersection := sinfo.span.Intersect(tinfo.span)
		if intersection[0] != intersection[1] {
			if !do(sinfo.span, tinfo.span, intersection, &sinfo.value, &tinfo.value) {
				break
			}
		}
		if tinfo.span[1] > sinfo.span[1] {
			sIndex++
		} else {
			tIndex++
		}
	}
}

// IntersectSet intersects two sets and returns a new set with all intersecting regions.
func (s RangeSet[T]) Intersect(t RangeSet[T], combine func(ss, ts, sx Span, sval, tval *T) T) RangeSet[T] {
	result := RangeSet[T]{}
	s.DoIntersect(t, func(ss, ts, xs Span, svalue, tvalue *T) bool {
		result.set = append(result.set, spanValue[T]{xs, combine(ss, ts, xs, svalue, tvalue)})
		return true
	})
	return result
}

// Do invokes a function on all ranges in the set.
//
// If a call returns false, no more ranges are visited.
func (s RangeSet[T]) Do(do func(Span, *T) bool) {
	for _, spanInfo := range s.set {
		if !do(spanInfo.span, &spanInfo.value) {
			break
		}
	}
}

func (s RangeSet[T]) Min() int {
	if len(s.set) == 0 {
		return 0
	}
	return s.set[0].span[0]
}

func (s RangeSet[T]) MinValue() *T {
	if len(s.set) == 0 {
		return nil
	}
	return &s.set[0].value
}

func (s RangeSet[T]) Max() int {
	if len(s.set) == 0 {
		return 0
	}
	return s.set[len(s.set)-1].span[1]
}

func (s RangeSet[T]) MaxValue() *T {
	if len(s.set) == 0 {
		return nil
	}
	return &s.set[len(s.set)-1].value
}

func (s RangeSet[T]) Span() Span {
	return Span{s.Min(), s.Max()}
}
