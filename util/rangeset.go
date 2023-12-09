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

func (s spanValue[T]) String() string {
	return fmt.Sprintf("%s=%v", s.span, s.value)
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
func bisect[T any](set []spanValue[T], value int) int {
	return sort.Search(len(set), func(index int) bool {
		return set[index].span.Contains(value) || set[index].span[0] > value
	})
}

func add[T any](set *[]spanValue[T], span Span, value T) *T {
	index := bisect(*set, span[0])
	if index == len(*set) {
		*set = append(*set, spanValue[T]{span, value})
	} else {
		*set = append((*set)[:index+1], (*set)[index:]...)
		(*set)[index] = spanValue[T]{span, value}
	}
	return &(*set)[index].value
}

func (s *RangeSet[T]) Add(span Span, value T) *T {
	return add(&s.set, span, value)
}

func extend[T any](s *[]spanValue[T], span Span, value T) *spanValue[T] {
	*s = append(*s, spanValue[T]{span, value})
	return &(*s)[len(*s)-1]
}

func (s RangeSet[T]) GetRange(key int) *RangeResult[T] {
	index := bisect(s.set, key)
	if index < len(s.set) && s.set[index].span.Contains(key) {
		return &RangeResult[T]{s.set[index].span, &s.set[index].value}
	}
	return nil
}

func get[T any](set []spanValue[T], key int) *T {
	index := bisect(set, key)
	if index < len(set) && set[index].span.Contains(key) {
		return &set[index].value
	}
	return nil
}

func (s RangeSet[T]) Get(key int) *T {
	return get(s.set, key)
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
	sIndex := bisect(s.set, t.set[0].span[0])
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

type CombineFunc[T any] func(svalue, tvalue *T) T

// Cover computes the cover set of two range sets.
//
// This is every range in the union of the two sets, where the
// intersections are combined with the combine function.
//
//	         __________  __________
//	______
//	      ______
//	           ______
//	                ______
//	                      ______
func (s RangeSet[T]) DoCover(t RangeSet[T], combine CombineFunc[T], visit func(Span, T) bool) {
	if len(s.set) == 0 {
		return
	}
	if len(t.set) == 0 {
		return
	}
	sIndex := 0
	tIndex := 0
	sweep := min(s.set[0].span[0], t.set[0].span[0])

	// sweep line algorithm:
	//   S = sweep line
	//   L = first = span with least lower bound
	//   R = second = span with greatest lower bound
	// min{S.max, T.max} when comparing spans S (from s.set) and T (from t.set).
	for sIndex < len(s.set) && tIndex < len(t.set) {
		first := &s.set[sIndex]
		second := &t.set[tIndex]
		if second.span[0] < first.span[0] {
			second, first = first, second
		}

		// skip gaps
		if sweep < first.span[0] {
			sweep = first.span[0]
		}

		// check intersection
		if second.span[0] < first.span[1] {
			// add [S,Rl] if not empty
			if sweep != second.span[0] {
				visit(Span{sweep, second.span[0]}, first.value)
			}
			// add [Rl, min{Lr,Rr}]; advance span with min{Lr,Rr}
			sweep = first.span[1]
			if second.span[1] < sweep {
				sweep = second.span[1]
			}
			visit(Span{second.span[0], sweep}, combine(&first.value, &second.value))
		} else {
			// no intersection: add [S,Lr]
			visit(Span{sweep, first.span[1]}, first.value)
			sweep = first.span[1]
		}
		// advance span with min{Lr,Rr}
		if s.set[sIndex].span[1] < t.set[tIndex].span[1] {
			sIndex++
		} else {
			tIndex++
		}
	}
	for _, set := range [][]spanValue[T]{s.set[sIndex:], t.set[tIndex:]} {
		if len(set) == 0 {
			continue
		}
		info := &set[0]
		// skip gap
		if sweep <= info.span[0] {
			sweep = info.span[0]
		}
		visit(Span{sweep, info.span[1]}, info.value)
		for index := 1; index < len(set); index++ {
			visit(set[index].span, set[index].value)
		}
	}
}

func (s RangeSet[T]) Cover(t RangeSet[T], combine CombineFunc[T]) RangeSet[T] {
	cover := RangeSet[T]{}
	s.DoCover(t, combine, func(s Span, value T) bool {
		extend(&cover.set, s, value)
		return true
	})
	return cover
}

// IntersectSet intersects two sets and returns a new set with all intersecting regions.
func (s RangeSet[T]) Intersect(t RangeSet[T], combine CombineFunc[T]) RangeSet[T] {
	result := RangeSet[T]{}
	s.DoIntersect(t, func(ss, ts, xs Span, svalue, tvalue *T) bool {
		result.set = append(result.set, spanValue[T]{xs, combine(svalue, tvalue)})
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

func (s RangeSet[T]) String() string {
	str := "{"
	for index, spanValue := range s.set {
		str += fmt.Sprintf(" [%d]=%s", index, spanValue)
	}
	str += " }"
	return str
}

type RangeMap RangeSet[int]

// Combine range-maps.
func (s RangeMap) CombineMap(t RangeMap) RangeMap {
	result := RangeMap{}
	if len(s.set) == 0 {
		return t
	}
	if len(t.set) == 0 {
		return s
	}
	sweep := Min(RangeSet[int](s).Min(), RangeSet[int](t).Min())
	sIndex := 0
	tIndex := 0
	for sIndex < len(s.set) && tIndex < len(t.set) {
		sspan := s.set[sIndex].span
		if sweep > sspan[0] {
			sspan[0] = sweep
		}
		tspan := t.set[tIndex].span
		if sweep > tspan[0] {
			tspan[0] = sweep
		}

		if tspan[0] < sspan[0] {
			sweep = Min(tspan[1], sspan[0])
			extend(&result.set, Span{tspan[0], sweep}, t.set[tIndex].value)
		}

		// Map contiguous subset of domain
		if sspan.Contains(sweep) {
			svalue := s.set[sIndex].value
			mapped := Span{sweep + svalue, svalue + sspan[1]}
			tMapIndex := bisect(t.set, mapped[0])
			if tMapIndex != len(t.set) {
				tinfo := &t.set[tMapIndex]
				value := svalue + tinfo.value
				if mapped[0] < tinfo.span[0] {
					if tinfo.span[0] < mapped[1] {
						extend(&result.set, Span{sweep, sweep + tinfo.span[0] - mapped[0]}, svalue)
						sweep += tinfo.span[0] - mapped[0]
						mapped[0] = tinfo.span[0]
					} else {
						value = svalue
					}
				} else if tinfo.span[1] < mapped[1] {
					mapped[1] = tinfo.span[1]
				}
				delta := mapped[1] - mapped[0]
				extend(&result.set, Span{sweep, sweep + delta}, value)
				sweep += delta
			} else {
				extend(&result.set, Span{sweep, sspan[1]}, svalue)
				sweep = sspan[1]
			}
		}

		if sweep >= sspan[1] {
			sIndex++
		}
		if sweep >= tspan[1] {
			tIndex++
		}
	}
	for ; sIndex < len(s.set); sIndex++ {
		extend(&result.set, s.set[sIndex].span, s.set[sIndex].value)
	}
	for ; tIndex < len(t.set); tIndex++ {
		extend(&result.set, t.set[tIndex].span, t.set[tIndex].value)
	}
	return result
}

func (s *RangeMap) Add(span Span, value int) {
	add(&s.set, span, value)
}

func (s RangeMap) Reduce(maps []RangeMap) RangeMap {
	result := s
	for _, rangeMap := range maps {
		newMap := result.CombineMap(rangeMap)
		fmt.Printf("map(\n     %s,\n     %s\n  => %s\n)\n", result, rangeMap, newMap)
		result = newMap
	}
	return result
}

func (s RangeMap) Map(value int) int {
	index := bisect(s.set, value)
	if index == len(s.set) || value < s.set[index].span[0] {
		return value
	}
	return value + s.set[index].value
}

func (s RangeMap) Maps(value int) bool {
	index := bisect(s.set, value)
	if index == len(s.set) || value < s.set[index].span[0] {
		return false
	}
	return true
}

func (s RangeMap) String() string {
	str := "{"
	for index, spanValue := range s.set {
		str += fmt.Sprintf(" [%d]=%s%+d=>%s", index, spanValue.span, spanValue.value,
			Span{spanValue.span[0] + spanValue.value, spanValue.span[1] + spanValue.value})
	}
	str += " }"
	return str
}

func (s RangeMap) Count() int {
	return len(s.set)
}
