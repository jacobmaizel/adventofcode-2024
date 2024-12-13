package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
)

func TestMain(t *testing.T) {
	input := []string{
		"MMMSXXMASM",
		"MSAMXMSMSA",
		"AMXSXMAAMM",
		"MSAMASMSMX",
		"XMASAMXAMM",
		"XXAMMXXAMA",
		"SMSMSASXSS",
		"SAXAMASAAA",
		"MAMMMXMMMM",
		"MXMXAXMASX",
	}

	expect := 18

	im := NewInputGrid(input)

	// im.printGrid()

	// fmt.Printf("ROWS: %d, COLS: %d\n", im.rows, im.cols)
	res := im.WordSearch()

	if res != 18 {
		t.Fatalf("Wrong answer. got=%d,want=%d\n", res, expect)
	}
}

func TestSeekVertical(t *testing.T) {
	tests := []struct {
		input        []string
		expectedXmas int
	}{
		{
			input: []string{
				"SSSS",
				"AAAA",
				"MMMM",
				"XXXX",
			},
			expectedXmas: 4,
		},
		{
			input: []string{
				"XXXX",
				"MMMM",
				"AAAA",
				"SSSS",
				"ASDF",
			},
			expectedXmas: 4,
		},
	}

	for i, tt := range tests {
		n := fmt.Sprintf("idx=%d", i)

		t.Run(n, func(t *testing.T) {
			im := NewInputGrid(tt.input)

			// im.printGrid()
			total := 0
			for x := range im.rows {
				for y := range im.cols {
					total += im.FindVertical(x, y)
				}
			}

			if total != tt.expectedXmas {
				t.Fatalf("wrong vertical total got=%d, want=%d", total, tt.expectedXmas)
			}
		})
	}
}

func TestSeekHorizontal(t *testing.T) {
	tests := []struct {
		input        []string
		expectedXmas int
	}{
		{
			input: []string{
				"XMAS",
				"SAMX",
				"SAMX",
				"XXXX",
			},
			expectedXmas: 3,
		},
		{
			input: []string{
				"XXXX",
				"MMMM",
				"AMAA",
				"SSSS",
				"ASDF",
			},
			expectedXmas: 0,
		},
	}

	for i, tt := range tests {
		n := fmt.Sprintf("idx=%d", i)

		t.Run(n, func(t *testing.T) {
			im := NewInputGrid(tt.input)

			// im.printGrid()
			total := 0
			for x := range im.rows {
				for y := range im.cols {
					total += im.FindHorizontal(x, y)
				}
			}

			if total != tt.expectedXmas {
				t.Fatalf("wrong horizontal total got=%d, want=%d", total, tt.expectedXmas)
			}
		})
	}
}

func TestSeekDiagonal(t *testing.T) {
	tests := []struct {
		input        []string
		startPoint   []int
		expectedXmas int
	}{
		{
			input: []string{
				"SXXXXXS",
				"XAXXXAX",
				"XXMXMXX",
				"XXXXXXX",
				"XXMXMXX",
				"XAXXXAX",
				"SXXXXXS",
			},
			startPoint:   []int{0, 0},
			expectedXmas: 4,
		},
		{
			// testing up boundary
			input: []string{
				"XAXXXAX",
				"XXMXMXX",
				"XXXXXXX",
				"XXMXMXX",
				"XAXXXAX",
				"SXXXXXS",
			},
			startPoint:   []int{0, 0},
			expectedXmas: 2,
		},
	}

	for i, tt := range tests {
		n := fmt.Sprintf("idx=%d", i)

		t.Run(n, func(t *testing.T) {
			im := NewInputGrid(tt.input)

			// im.printGrid()
			total := 0
			for x := range im.rows {
				for y := range im.cols {
					total += im.FindDiagonal(x, y)
				}
			}

			if total != tt.expectedXmas {
				t.Fatalf("wrong horizontal total got=%d, want=%d", total, tt.expectedXmas)
			}
		})
	}
}

func BenchmarkWordSearch(b *testing.B) {
	/*
		   initial:  789   1279141 ns/op  484962 B/op                        60602 allocs/op
		   skipping non x:  1932   622397 ns/op                 236697 B/op  29580 allocs/op
		   early return for mismatch byte:  2901  412033 ns/op  119592 B/op  14944 allocs/op
			 fixes for diagonal early return: 4490  262942 ns/op      23 B/op  0 allocs/op
	*/
	file, _ := os.Open("aoc-day4-input.txt")
	defer file.Close()
	str, _ := io.ReadAll(file)
	lines := strings.Split(string(str), "\n")
	im := NewInputGrid(lines[:len(lines)-1])

	for i := 0; i < b.N; i++ {
		im.WordSearch()
	}
}
