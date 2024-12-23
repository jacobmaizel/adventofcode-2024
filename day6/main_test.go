package main

import (
	"fmt"
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

func TestObsPlacementLoopCount(t *testing.T) {
	tests := []struct {
		grid                   []string
		expectedObstPlacements int
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
			}, 6,
		},
	}

	for i, tt := range tests {
		n := fmt.Sprintf("idx %d", i)

		t.Run(n, func(t *testing.T) {
			m := NewMapFromGrid(tt.grid)

			res := m.simulateDifferentObstructionPositions()

			if tt.expectedObstPlacements != res {
				t.Fatalf("wrong placement count: want=%d, got=%d\n", tt.expectedObstPlacements, res)
			}
		})
	}
}

func TestGuardLoop(t *testing.T) {
	tests := []struct {
		grid       []string
		expectLoop bool
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
			}, false,
		},
		{
			[]string{
				"....#.....",
				".........#",
				"..........",
				"..#.......",
				".......#..",
				"..........",
				".#.#^.....",
				"........#.",
				"#.........",
				"......#...",
			}, true,
		},
	}

	for i, tt := range tests {
		n := fmt.Sprintf("idx %d", i)

		t.Run(n, func(t *testing.T) {
			m := NewMapFromGrid(tt.grid)

			m.moveGuard()

			if m.guardLeavingGrid && m.stuckInLoop && tt.expectLoop == m.stuckInLoop {
				t.Fatal("cant leave and also be in a loop..")
			}
			if tt.expectLoop != m.stuckInLoop {
				t.Fatalf("mismatch expect loop: got=%t, want=%t\n", m.stuckInLoop, tt.expectLoop)
			}
		})
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

/*
initial
1        4134710709 ns/op        5821483544 B/op    2209030 allocs/op
*/


func BenchmarkDay6p2(b *testing.B) {
	f, _ := os.Open("day6-aoc-input.txt")
	defer f.Close()

	m := NewMapFromReader(f)

	for i := 0; i < b.N; i++ {
		m.simulateDifferentObstructionPositions()
	}
}
