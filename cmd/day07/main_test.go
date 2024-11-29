package main

import (
	"reflect"
	"slices"
	"testing"
)

func compareHands(t *testing.T, hand1 string, hand2 string, cmpExp int) {
	h1 := Hand{hand1, 0}
	h2 := Hand{hand2, 0}
	cmp := ranker.Compare(h1, h2)
	if ((cmp < 0) != (cmpExp < 0)) || ((cmp > 0) != (cmpExp > 0)) {
		t.Errorf("Compare(%s, %s) expected %d, got %d", hand1, hand2, cmpExp, cmp)
	}
}

func TestCompare(t *testing.T) {
	type test struct {
		hand1, hand2 string
		cmpExp       int
	}

	for _, test := range []test{
		{"22627", "22992", -1},
		{"22992", "282A2", 1}, // full house vs 3 pair
		{"24224", "242K2", 1}, // full house vs 3 pair
		{"22427", "3QQ4Q", -1},
	} {
		compareHands(t, test.hand1, test.hand2, test.cmpExp)
	}
}

func TestSort(t *testing.T) {
	type test struct {
		hands  []Hand
		sorted []Hand
	}
	for _, test := range []test{
		{
			hands: []Hand{
				{"32T3K", 765},
				{"T55J5", 684},
				{"KK677", 28},
				{"KTJJT", 220},
				{"QQQJA", 483},
			},
			sorted: []Hand{
				{"32T3K", 765},
				{"KTJJT", 220},
				{"KK677", 28},
				{"T55J5", 684},
				{"QQQJA", 483},
			},
		},
		{
			hands: []Hand{
				{"22592", 0},
				{"22627", 0},
				{"24544", 0},
				{"24KKK", 0},
				{"26686", 0},
				{"26AAA", 0},
				{"242K2", 0},
				{"24224", 0},
				{"282A2", 0},
				{"2T2Q2", 0},
				{"22992", 0},
				{"2Q22Q", 0},
			},
			sorted: []Hand{
				{"22592", 0},
				{"22627", 0},
				{"242K2", 0},
				{"24544", 0},
				{"24KKK", 0},
				{"26686", 0},
				{"26AAA", 0},
				{"282A2", 0},
				{"2T2Q2", 0},
				{"22992", 0},
				{"24224", 0},
				{"2Q22Q", 0},
			},
		},
		/*
			{
				hands: []Hand{
					{"22627", 0},
					{"24544", 0},
					{"26686", 0},
					{"26AAA", 0},
					{"242K2", 0},
					{"24224", 0},
				},
				sorted: []Hand{
					{"22627", 0},
					{"24224", 0},
					{"242K2", 0},
					{"24544", 0},
					{"26686", 0},
					{"26AAA", 0},
				},
			}*/
	} {
		slices.SortFunc(test.hands, ranker.Compare)

		if !reflect.DeepEqual(test.hands, test.sorted) {
			t.Errorf("hands not sorted as expected: %+v", test.hands)
		}
	}
}

func TestRanks(t *testing.T) {
	type test struct {
		ranker HandRanker
		hand   string
		rank   HandRank
	}
	for _, test := range []test{
		{ranker, "12345", RankHighCard},
		{ranker, "12234", RankPair},
		{ranker, "11233", RankTwoPair},
		{ranker, "12223", RankTriple},
		{ranker, "12222", RankQuad},
		{ranker, "11111", RankPenta},
		{jokerRanker, "12345", RankHighCard},
		{jokerRanker, "12234", RankPair},
		{jokerRanker, "1223J", RankTriple},
		{jokerRanker, "J2233", RankFullHouse},
		{jokerRanker, "JJJ12", RankQuad},
		{jokerRanker, "11J1J", RankPenta},
	} {
		testRank := test.ranker.LookupRank(Hand{test.hand, 0})
		if test.rank != testRank {
			t.Errorf("hand %s: expected rank %d, got %d", test.hand, test.rank, testRank)
		}
	}
}
