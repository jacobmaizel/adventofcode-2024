package main

import (
	"os"
	"testing"
)

func TestManualMovement(t *testing.T) {
	tests := []struct {
		grid        []string
		expectedAns int
	}{
		{
			[]string{
				"....#.....",
				".........#",
				"..........",
				"..#.......",
				".......#..",
				"..........",
				".#..^.....",
				"........#.",
				"#.........",
				"......#...",
			}, 41,
		},
	}

	for _, tt := range tests {
		m := NewMapFromGrid(tt.grid)

		m.guardUp()
		m.guardRight()
		m.guardDown()
		// m.printGrid()
		m.guardLeft()
		// m.printGrid()
		m.guardUp()
		// m.printGrid()
		m.guardRight()
		// m.printGrid()
		m.guardDown()
		// m.printGrid()
		m.guardLeft()
		// m.printGrid()

		m.guardUp()
		m.guardRight()
		// m.printGrid()
		m.guardDown()
		// m.printGrid()

		if tt.expectedAns != len(m.distinctPositions) {
			t.Fatalf("wrong!!! got=%d --- want=%d", len(m.distinctPositions), tt.expectedAns)
		}
	}
}

func TestAutoMovement(t *testing.T) {
	tests := []struct {
		grid        []string
		expectedAns int
	}{
		{
			[]string{
				"....#.....",
				".........#",
				"..........",
				"..#.......",
				".......#..",
				"..........",
				".#..^.....",
				"........#.",
				"#.........",
				"......#...",
			}, 41,
		},
	}

	for _, tt := range tests {
		m := NewMapFromGrid(tt.grid)

		m.moveGuard()

		if tt.expectedAns != len(m.distinctPositions) {
			t.Fatalf("wrong!!! got=%d --- want=%d", len(m.distinctPositions), tt.expectedAns)
		}
	}
}

/*
initial
569947948                1.904 ns/op           0 B/op          0 allocs/op
*/

func BenchmarkDay6p1(b *testing.B) {
	f, _ := os.Open("day6-aoc-input.txt")
	defer f.Close()

	m := NewMapFromReader(f)

	for i := 0; i < b.N; i++ {
		m.moveGuard()
	}
}
