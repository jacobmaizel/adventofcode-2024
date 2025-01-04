package main

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

var _ = fmt.Print

func Test_main(t *testing.T) {
	tests := []struct {
		name                            string
		input                           string
		expectedUniqueAntinodePositions int
	}{
		{
			name: "web case",
			input: `............
        ........0...
        .....0......
        .......0....
        ....0.......
        ......A.....
        ............
        ............
        ........A...
        .........A..
        ............
        ............`,

			expectedUniqueAntinodePositions: 14,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := strings.NewReader(tt.input)
			i := newInput(r)

			i.computePart1()

			// fmt.Println(i.String())

			// for k := range i.antinodePositions {
			// 	i.grid[k.row][k.col] = 'X'
			// }

			// fmt.Println(i.String())
			if len(i.antinodePositions) != tt.expectedUniqueAntinodePositions {
				t.Errorf("wrong ap count,got=%v want=%v", len(i.antinodePositions), tt.expectedUniqueAntinodePositions)
			}
		})
	}
}

func Test_mainPart2(t *testing.T) {
	tests := []struct {
		name                            string
		input                           string
		expectedUniqueAntinodePositions int
	}{
		{
			name: "web case",
			input: `............
        ........0...
        .....0......
        .......0....
        ....0.......
        ......A.....
        ............
        ............
        ........A...
        .........A..
        ............
        ............`,

			expectedUniqueAntinodePositions: 34,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := strings.NewReader(tt.input)
			i := newInput(r)

			i.computePart2()

			// fmt.Println(i.String())

			// for k := range i.antinodePositions {
			// 	i.grid[k.row][k.col] = '#'
			// }

			// fmt.Println(i.String())
			if len(i.antinodePositions) != tt.expectedUniqueAntinodePositions {
				t.Errorf("wrong ap count,got=%v want=%v", len(i.antinodePositions), tt.expectedUniqueAntinodePositions)
			}
		})
	}
}

func Test_distance(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		p1    position
		p2    position
		want  int
		want2 int
	}{
		// TODO: Add test cases.
		{
			name:  "t1",
			p1:    position{row: 1, col: 2},
			p2:    position{row: 1, col: 3},
			want:  0,
			want2: 1,
		},
		{
			name:  "t2",
			p1:    position{row: 1, col: 8},
			p2:    position{row: 2, col: 5},
			want:  1,
			want2: -3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got2 := distance(tt.p1, tt.p2)
			// TODO: update the condition below to compare got with tt.want.
			if got != tt.want || got2 != tt.want2 {
				t.Errorf("distance() = %v,%v, want %v,%v", got, got2, tt.want, tt.want2)
			}
		})
	}
}

func BenchmarkP1(b *testing.B) {
	f, _ := os.Open("aoc-day8-input.txt")
	defer f.Close()
	input := newInput(f)
	for i := 0; i < b.N; i++ {
		input.computePart1()
	}
}

func BenchmarkP2(b *testing.B) {
	f, _ := os.Open("aoc-day8-input.txt")
	defer f.Close()
	input := newInput(f)
	for i := 0; i < b.N; i++ {
		input.computePart2()
	}
}
