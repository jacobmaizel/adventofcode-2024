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
		expTotal int
	}{
		{name: "example", input: "125 17", expected: "1 2024 1 0 9 9 2021976", blinks: 25, expTotal: 55312},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := strings.NewReader(tt.input)

			i := newInput(f)

			res := i.simulateBlinks(tt.blinks)
			if res != tt.expTotal {
				t.Fatal("wrong answer", res)
			}
		})
	}
}

func Benchmark_p1_25(b *testing.B) {
	// f, _ := os.Open("input.txt")
	// defer f.Close()
	ip := "2 77706 5847 9258441 0 741 883933 12"
	for n := 0; n < b.N; n++ {
		i := newInput(strings.NewReader(ip))
		i.simulateBlinks(25)
	}
}

func Benchmark_p1_75(b *testing.B) {
	// f, _ := os.Open("input.txt")
	// defer f.Close()
	ip := "2 77706 5847 9258441 0 741 883933 12"
	i := newInput(strings.NewReader(ip))
	for n := 0; n < b.N; n++ {
		i.simulateBlinks(75)
	}
}
