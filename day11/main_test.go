package main

import (
	// "fmt"
	"strings"
	"testing"
)

func Test_main(t *testing.T) {
	tests := []struct {
		name     string // description of this test case
		input    string
		blinks   int
		expected string
	}{
		{name: "example", input: "0 1 10 99 999", expected: "1 2024 1 0 9 9 2021976", blinks: 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := strings.NewReader(tt.input)

			i := newInput(f)

			// fmt.Println("stones", i.stones)

			i.simulateBlinks(tt.blinks)

			expSpl := strings.Split(tt.expected, " ")

			if len(i.stones) != len(expSpl) {
				t.Errorf("wrong len of stones after blinks, got=%d, want=%d", len(i.stones), len(tt.expected))
			}

			for idx := range expSpl {
				if i.stones[idx] != expSpl[idx] {
					t.Errorf("mismatch at idx=%d, got=%s, want=%s", idx, i.stones[idx], expSpl[idx])
				}
			}
		})
	}
}


/*
initial, 25 blinks
33          34984295 ns/op        60714923 B/op   708332 allocs/op
*/

func Benchmark_p1(b *testing.B) {
	// f, _ := os.Open("input.txt")
	// defer f.Close()
	ip := "2 77706 5847 9258441 0 741 883933 12"
	for n := 0; n < b.N; n++ {
		i := newInput(strings.NewReader(ip))
		i.simulateBlinks(25)
	}
}
