package main

import (
	"fmt"
	"os"
	"strings"
	"testing"
	// "os"
	// "strings"
	// "testing"
)

var _ = fmt.Printf

// func Test_checksum(t *testing.T) {
// 	tests := []struct {
// 		name   string // description of this test case
// 		input  string
// 		expect int
// 	}{
// 		// TODO: Add test cases.
// 		{name: "example", input: "2333133121414131402", expect: 2858},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			i := newInput(strings.NewReader(tt.input))
// 			// fmt.Println(string(i.diskMap))
// 			i.expandDiskMap()
// 			var res strings.Builder

// 			for _, b := range i.expandedMap {
// 				res.WriteString(b.String())
// 			}

// 			fmt.Println(res.String())

// 			i.compressDiskMap()

// 			fmt.Println(i.String())
// 			fmt.Println(i.input)

// 			e := i.checkSum()

// 			if tt.expect != e {
// 				t.Fatalf("\nchecksum() = %d\nwant=         %d\n", e, tt.expect)
// 			}
// 		})
// 	}
// }

func Test_swap_block(t *testing.T) {
	tests := []struct {
		name   string // description of this test case
		input  string
		start  int
		end    int
		expect string
	}{
		// {name: "example", input: "2333133121414131402", expect: "00...111...2...333.44.5555.6666.777.888899"},
		// {name: "example", input: "2333133121414131402", expect: "00333111...2.......44.5555.6666.777.888899"},
		{name: "example2", input: "2333133121414131402", expect: "0099.111...2...333.44.5555.6666.777.8888.."},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := newInput(strings.NewReader(tt.input))
			// fmt.Println(string(i.diskMap))

			i.buildFileSystem()

			i.swapBlocks(2, 2, len(i.String())-1)

			if tt.expect != i.String() {
				t.Fatalf("\nafter swap = %s\nwant=           %s\n", i.String(), tt.expect)
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

			i.buildFileSystem()

			if tt.expect != i.String() {
				t.Fatalf("\nexpandDiskmap = %s\nwant=           %s\n", i.String(), tt.expect)
			}
		})
	}
}

// // different compression rules for part 2
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

			i.buildFileSystem()
			i.compressFileSystem()

			if tt.expect != i.String() {
				t.Fatalf("\ncompressFile = %s\nwant=          %s\n", i.String(), tt.expect)
			}
		})
	}
}

func Benchmark_p2(b *testing.B) {
	f, _ := os.Open("aoc-day9-input.txt")

	defer f.Close()

	for idx := 0; idx < b.N; idx++ {
		i := newInput(f)
		i.checkSum()
	}
}
