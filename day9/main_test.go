package main

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

var _ = fmt.Printf

func Test_checksum(t *testing.T) {
	tests := []struct {
		name   string // description of this test case
		input  string
		expect int
	}{
		// TODO: Add test cases.
		{name: "example", input: "2333133121414131402", expect: 2858},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := newInput(strings.NewReader(tt.input))
			// fmt.Println(string(i.diskMap))
			i.expandDiskMap()
			var res strings.Builder

			for _, b := range i.expandedMap {
				res.WriteString(b.String())
			}

			// fmt.Println(res.String())

			i.compressDiskMap()

			var re2 strings.Builder

			for _, b := range i.expandedMap {
				re2.WriteString(b.String())
			}

			// fmt.Println(re2.String())

			e := i.checkSum()

			if tt.expect != e {
				t.Fatalf("\nchecksum() = %d\nwant=         %d\n", e, tt.expect)
			}
		})
	}
}

func Test_expand_blocks(t *testing.T) {
	tests := []struct {
		name   string // description of this test case
		input  string
		expect string
	}{
		// TODO: Add test cases.
		{name: "example", input: "2333133121414131402", expect: "00...111...2...333.44.5555.6666.777.888899"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := newInput(strings.NewReader(tt.input))
			// fmt.Println(string(i.diskMap))

			i.expandDiskMap()

			var res strings.Builder

			for _, b := range i.expandedMap {
				res.WriteString(b.String())
			}

			r := res.String()
			if tt.expect != r {
				t.Fatalf("\nexpandDiskmap = %s\nwant=           %s\n", r, tt.expect)
			}
		})
	}
}

// different compression rules for part 2
func Test_compress_filePart2(t *testing.T) {
	tests := []struct {
		name   string // description of this test case
		input  string
		expect string
	}{
		// TODO: Add test cases.
		{name: "example", input: "2333133121414131402", expect: "00992111777.44.333....5555.6666.....8888.."},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := newInput(strings.NewReader(tt.input))
			// fmt.Println(string(i.diskMap))

			i.expandDiskMap()
			i.compressDiskMap()

			var res strings.Builder

			for _, b := range i.expandedMap {
				res.WriteString(b.String())
			}
			r := res.String()
			if tt.expect != r {
				t.Fatalf("\ncompressFile = %s\nwant=          %s\n", r, tt.expect)
			}
		})
	}
}

// func Test_compress_file(t *testing.T) {
// 	tests := []struct {
// 		name   string // description of this test case
// 		input  string
// 		expect string
// 	}{
// 		// TODO: Add test cases.
// 		{name: "example", input: "2333133121414131402", expect: "0099811188827773336446555566.............."},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			i := newInput(strings.NewReader(tt.input))
// 			// fmt.Println(string(i.diskMap))

// 			e := i.compressDiskMap()

// 			var res strings.Builder

// 			for _, b := range e {
// 				res.WriteString(b.String())
// 			}
// 			r := res.String()
// 			if tt.expect != r {
// 				t.Fatalf("\ncompressFile = %s\nwant=          %s\n", r, tt.expect)
// 			}
// 		})
// 	}
// }

/*
initial correct answer
355           3116640 ns/op         9083822 B/op       10030 allocs/op
932           1323666 ns/op         2263125 B/op       38810 allocs/op
*/
func Benchmark_expand_disk_map(b *testing.B) {
	f, _ := os.Open("aoc-day9-input.txt")

	defer f.Close()

	i := newInput(f)

	for idx := 0; idx < b.N; idx++ {
		i.expandDiskMap()
	}
}
