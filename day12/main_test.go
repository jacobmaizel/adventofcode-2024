package main

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

var _ = fmt.Println

func Test_main(t *testing.T) {
	tests := []struct {
		name          string // description of this test case
		input         string
		expectedPrice int
	}{
		{
			name: "web",
			input: `RRRRIICCFF
              RRRRIICCCF
              VVRRRCCFFF
              VVRCCCJFFF
              VVVVCJJCFE
              VVIVCCJJEE
              VVIIICJJEE
              MIIIIIJJEE
              MIIISIJEEE
              MMMISSJEEE`,
			expectedPrice: 1930,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := strings.NewReader(tt.input)

			i := newInput(r)

			plots := make(plotLookup)

			i.partOneQueueBFS(plots)
			a := i.calcPart1FenceCost(plots)

			if a != tt.expectedPrice {
				t.Errorf("Got=%d, want=%d", a, tt.expectedPrice)
			}
		})
	}
}

func Test_p2(t *testing.T) {
	tests := []struct {
		name          string // description of this test case
		input         string
		expectedPrice int
	}{
		{
			name: "web",
			input: `RRRRIICCFF
              RRRRIICCCF
              VVRRRCCFFF
              VVRCCCJFFF
              VVVVCJJCFE
              VVIVCCJJEE
              VVIIICJJEE
              MIIIIIJJEE
              MIIISIJEEE
              MMMISSJEEE`,
			expectedPrice: 1206,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := strings.NewReader(tt.input)

			i := newInput(r)

			plots := make(plotLookup)

			i.partOneQueueBFS(plots)
			a := i.calcPart2FenceCost(plots)

			if a != tt.expectedPrice {
				t.Errorf("Got=%d, want=%d", a, tt.expectedPrice)
			}
		})
	}
}

/*
part 1 floodfill initial dfs ->
319           3499944 ns/op         2379585 B/op    3264 allocs/op
*/

func BenchmarkPart1DFS(b *testing.B) {
	f, _ := os.Open("input.txt")
	defer f.Close()

	in := newInput(f)
	for i := 0; i < b.N; i++ {
		plots := make(plotLookup)
		in.partOneRecursiveDFS(plots)
		// in.calcPart1FenceCost(plots)
	}
}

func BenchmarkPart1BFS(b *testing.B) {
	f, _ := os.Open("input.txt")
	defer f.Close()

	in := newInput(f)
	for i := 0; i < b.N; i++ {
		plots := make(plotLookup)
		in.partOneQueueBFS(plots)
	}
}

func BenchmarkPart1Calc(b *testing.B) {
	f, _ := os.Open("input.txt")
	defer f.Close()
	in := newInput(f)

	plots := make(plotLookup)
	in.partOneQueueBFS(plots)

	for i := 0; i < b.N; i++ {
		in.calcPart1FenceCost(plots)
	}
}

func BenchmarkPart2Calc(b *testing.B) {
	f, _ := os.Open("input.txt")
	defer f.Close()
	in := newInput(f)

	plots := make(plotLookup)
	in.partOneQueueBFS(plots)

	for i := 0; i < b.N; i++ {
		in.calcPart2FenceCost(plots)
	}
}
