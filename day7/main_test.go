package main

import (
	"os"
	"strings"
	"testing"
)

func TestAocWebCasesP2(t *testing.T) {
	input := []struct {
		input  string
		expAns int
	}{
		{
			input: `190: 10 19
      3267: 81 40 27
      83: 17 5
      156: 15 6
      7290: 6 8 6 15
      161011: 16 10 13
      192: 17 8 14
      21037: 9 7 18 13
      292: 11 6 16 20`,
			expAns: 11387,
		},
	}

	for _, tt := range input {
		r := strings.NewReader(tt.input)

		i := NewInput(r)

		res := i.RunPart2()
		if res != tt.expAns {
			t.Fatalf("wrong answer, got=%d ... want=%d", res, tt.expAns)
		}

	}
}

func TestAocWebCasesP1(t *testing.T) {
	input := []struct {
		input  string
		expAns int
	}{
		{
			input: `190: 10 19
      3267: 81 40 27
      83: 17 5
      156: 15 6
      7290: 6 8 6 15
      161011: 16 10 13
      192: 17 8 14
      21037: 9 7 18 13
      292: 11 6 16 20`,
			expAns: 3749,
		},
	}

	for _, tt := range input {
		r := strings.NewReader(tt.input)

		i := NewInput(r)

		res := i.RunPart1()
		if res != tt.expAns {
			t.Fatalf("wrong answer, got=%d ... want=%d", res, tt.expAns)
		}

	}
}

/*
initial
10         111399208 ns/op        181197540 B/op   6038014 allocs/op
*/
func BenchmarkP1(b *testing.B) {
	f, _ := os.Open("aoc-day7-input.txt")
	defer f.Close()

	in := NewInput(f)

	for i := 0; i < b.N; i++ {
		in.RunPart1()
	}
}

// initial 1        10551440583 ns/op       18398080328 B/op        433981256 allocs/op

func BenchmarkP2(b *testing.B) {
	f, _ := os.Open("aoc-day7-input.txt")
	defer f.Close()

	in := NewInput(f)

	for i := 0; i < b.N; i++ {
		in.RunPart2()
	}
}
