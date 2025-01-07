package main

import (
	"os"
	"strings"
	"testing"
)

func Test_rating_calc(t *testing.T) {
	tests := []struct {
		name           string // description of this test case
		input          string
		expectedRating int
	}{
		{
			name: "example",
			input: `..90..9
...1.98
...2..7
6543456
765.987
876....
987....`,
			expectedRating: 13,
		},
		{
			name: "example Test2 small",
			input: `.....0.
              ..4321.
              ..5..2.
              ..6543.
              ..7..4.
              ..8765.
              ..9....`,
			expectedRating: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := strings.NewReader(tt.input)

			in := newInput(r)

			// fmt.Println(in.String())
			// fmt.Println(in.trailheadPositions)
			res := in.calcRating()

			if res != tt.expectedRating {
				t.Errorf("got %d, want %d", res, tt.expectedRating)
			}
		})
	}
}

func Test_main(t *testing.T) {
	tests := []struct {
		name                         string // description of this test case
		input                        string
		expectedSumofTrailheadScores int
	}{
		{
			name: "example Test",
			input: `89010123
		78121874
		87430965
		96549874
		45678903
		32019012
		01329801
		10456732`,
			expectedSumofTrailheadScores: 36,
		},

		{
			name: "example Test3",
			input: `...0...
		...1...
		...2...
		6543456
		7.....7
		8.....8
		9.....9`,
			expectedSumofTrailheadScores: 2,
		},
		{
			name: "example Test2 small",
			input: `10..9..
		2...8..
		3...7..
		4567654
		...8..3
		...9..2
		.....01`,
			expectedSumofTrailheadScores: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := strings.NewReader(tt.input)

			in := newInput(r)

			// fmt.Println(in.String())
			// fmt.Println(in.trailheadPositions)
			res := in.calcScore()

			if res != tt.expectedSumofTrailheadScores {
				t.Errorf("got %d, want %d", res, tt.expectedSumofTrailheadScores)
			}
		})
	}
}

func Benchmark_part2(b *testing.B) {
	f, _ := os.Open("input-day10.txt")
	defer f.Close()

	in := newInput(f)

	for i := 0; i < b.N; i++ {
		in.calcRating()
	}
}

func Benchmark_part1(b *testing.B) {
	f, _ := os.Open("input-day10.txt")
	defer f.Close()

	in := newInput(f)

	for i := 0; i < b.N; i++ {
		in.calcScore()
	}
}
