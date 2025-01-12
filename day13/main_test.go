package main

import (
	"fmt"
	"strings"
	"testing"
)

var _ = fmt.Println

func Test_part2(t *testing.T) {
	tests := []struct {
		name      string // description of this test case
		input     string
		expTokens int
	}{
		{
			name: "web",

			input: `Button A: X+94, Y+34
              Button B: X+22, Y+67
              Prize: X=8400, Y=5400

              Button A: X+26, Y+66
              Button B: X+67, Y+21
              Prize: X=12748, Y=12176

              Button A: X+17, Y+86
              Button B: X+84, Y+37
              Prize: X=7870, Y=6450

              Button A: X+69, Y+23
              Button B: X+27, Y+71
              Prize: X=18641, Y=10279`,
			expTokens: 480,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := strings.NewReader(tt.input)
			in := newInput(r)

			res := in.calcTokensP2()

			if res != tt.expTokens {
				t.Errorf("got %d, want %d", res, tt.expTokens)
			}
		})
	}
}

// func Test_main(t *testing.T) {
// 	tests := []struct {
// 		name      string // description of this test case
// 		input     string
// 		expTokens int
// 	}{
// 		{
// 			name: "web",

// 			input: `Button A: X+94, Y+34
//               Button B: X+22, Y+67
//               Prize: X=8400, Y=5400

//               Button A: X+26, Y+66
//               Button B: X+67, Y+21
//               Prize: X=12748, Y=12176

//               Button A: X+17, Y+86
//               Button B: X+84, Y+37
//               Prize: X=7870, Y=6450

//               Button A: X+69, Y+23
//               Button B: X+27, Y+71
//               Prize: X=18641, Y=10279`,
// 			expTokens: 480,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			r := strings.NewReader(tt.input)
// 			in := newInput(r)

// 			// res := in.part1()

// 			if res != tt.expTokens {
// 				t.Errorf("got %d, want %d", res, tt.expTokens)
// 			}
// 		})
// 	}
// }

// 1        1203075958 ns/op        523246128 B/op  19834368 allocs/op
// 54408             20134 ns/op           3 B/op          0 allocs/op
func BenchmarkP1(b *testing.B) {
	f := inputFile()
	in := newInput(f)

	for i := 0; i < b.N; i++ {
		in.calcTokensP2()
	}
}
