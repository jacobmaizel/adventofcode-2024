package main

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

// the AFTER portion of a rule CAME BEFORE the before portion of the rule
func TestMain(t *testing.T) {
	tests := []struct {
		input string
		// expectValid         bool
		expectedMidpointSum int
	}{
		{"82|32\n94|34\n\n\n54,82,18,29,94,32,34", 29},
		// {"32|82\n94|34\n\n\n54,82,18,29,94,32,34", 0},
	}

	for i, tt := range tests {

		n := fmt.Sprintf("idx=%d", i)
		t.Run(n, func(t *testing.T) {
			r := strings.NewReader(tt.input)
			inp := NewInputFromReader(r)

			res := inp.CalcP1()

			if res != tt.expectedMidpointSum {
				t.Fatalf("Wrong answer, want=%d - got=%d", tt.expectedMidpointSum, res)
			}
		})

	}
}

func TestPart2(t *testing.T) {
	tests := []struct {
		inputRows           [][]int
		brokenRules         []Rule
		expectedMidpointSum int
	}{
		{[][]int{
			{75, 97, 47, 61, 53},
			{61, 13, 29},
			{97, 13, 75, 29, 47},
		}, []Rule{
			{before: 47, after: 53},
			{before: 97, after: 13},
			{before: 97, after: 61},
			{before: 97, after: 47},
			{before: 75, after: 29},
			{before: 61, after: 13},
			{before: 75, after: 53},
			{before: 29, after: 13},
			{before: 97, after: 29},
			{before: 53, after: 29},
			{before: 61, after: 53},
			{before: 97, after: 53},
			{before: 61, after: 29},
			{before: 47, after: 13},
			{before: 75, after: 47},
			{before: 97, after: 75},
			{before: 47, after: 61},
			{before: 75, after: 61},
			{before: 47, after: 29},
			{before: 75, after: 13},
			{before: 53, after: 13},
		}, 123},
	}

	for i, tt := range tests {

		n := fmt.Sprintf("idx=%d", i)
		t.Run(n, func(t *testing.T) {
			inp := NewInputFromParts(tt.inputRows, tt.brokenRules)

			res := inp.CalcP2()

			if res != tt.expectedMidpointSum {
				t.Fatalf("Wrong answer, want=%d - got=%d", tt.expectedMidpointSum, res)
			}
		})

	}
}

// the AFTER portion of a rule CAME BEFORE the before portion of the rule
func TestAfterRuleBreak(t *testing.T) {
	tests := []struct {
		input               string
		expectedMidpointSum int
	}{
		{"18|82\n\n\n54,82,18,29,94,32,34", 0},
	}

	for i, tt := range tests {

		n := fmt.Sprintf("idx=%d", i)
		t.Run(n, func(t *testing.T) {
			r := strings.NewReader(tt.input)
			inp := NewInputFromReader(r)

			res := inp.CalcP1()

			if res != tt.expectedMidpointSum {
				t.Fatalf("Wrong answer, want=%d - got=%d", tt.expectedMidpointSum, res)
			}
		})

	}
}

// the BEFORE portion of the rule CAME AFTER the after portion of the rule
func TestBeforeRuleBreak(t *testing.T) {
	tests := []struct {
		input               string
		expectValid         bool
		expectedMidpointSum int
	}{
		{"29|82\n\n\n54,82,18,29,94,32,34", false, 0},
	}

	for i, tt := range tests {

		n := fmt.Sprintf("idx=%d", i)
		t.Run(n, func(t *testing.T) {
			r := strings.NewReader(tt.input)
			inp := NewInputFromReader(r)

			res := inp.CalcP1()

			if res == 0 && tt.expectValid {
				t.Fatal("Expected valid...")
			}

			if res != tt.expectedMidpointSum {
				t.Fatalf("Wrong answer, want=%d - got=%d", tt.expectedMidpointSum, res)
			}
		})

	}
}

/*
initial:  429  2802167 ns/op         3535525 B/op      20912 allocs/op
with ptrs: 42 29188391 ns/op        33794528 B/op    2019002 allocs/op
removing unused afterRules: 662           1773943 ns/op         1823045 B/op      10720 allocs/op
using map of page: rules: 3013            386086 ns/op          110297 B/op        523 allocs/op
remove dupe strconvs: 3819            325321 ns/op           38009 B/op        233 allocs/op
*/

func BenchmarkDay5p1(b *testing.B) {
	f, _ := os.Open("aoc-day5-input.txt")
	defer f.Close()

	inp := NewInputFromReader(f)

	for i := 0; i < b.N; i++ {
		inp.CalcP1()
	}
}

/*
Part 2
initial
384           3070162 ns/op         1845874 B/op      10896 allocs/op

remove some unused arrays
546           2239621 ns/op         1026635 B/op       9215 allocs/op
*/

func BenchmarkDay5p2(b *testing.B) {
	f, _ := os.Open("aoc-day5-input.txt")
	defer f.Close()

	inp := NewInputFromReader(f)

	for i := 0; i < b.N; i++ {
		inp.CalcP2()
	}
}
