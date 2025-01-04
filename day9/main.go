package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
)

/*
Day 9
*/

type input struct {
	diskMap []byte
}

func newInput(r io.Reader) *input {
	v, _ := io.ReadAll(r)
	return &input{v}
}

type file struct {
	val byte // the value at that index from the original disk map
	id  int  // the original index that this came from in the input disk map
}

type emptySpace struct{}

type block interface {
	block()
	String() string
}

func (f file) block()       {}
func (f emptySpace) block() {}

const (
	FILE_MODE = iota
	FREE_SPACE_MODE
)

var EMPTY_SPACE_BLOCK = emptySpace{}

func (e emptySpace) String() string {
	return "."
}

func (f file) String() string {
	s := strconv.Itoa(f.id)
	return s
}

func (f file) FullString() string {
	return fmt.Sprintf("[fileId=%d, val=%s]", f.id, string(f.val))
}

func (i *input) expandDiskMap() []block {
	currFileId := 0

	var buff []block

	for diskIdx := range i.diskMap {
		diskmapVal := i.diskMap[diskIdx]

		switch diskIdx % 2 {
		case FILE_MODE:
			// fmt.Println("File mode:", diskIdx%2, " FileID:", currFileId)

			loops, err := strconv.Atoi(string(diskmapVal))

			f := newFile(diskmapVal, currFileId)

			if err != nil {
				// fmt.Println("Error converting to int at idx file mode", diskIdx, string(diskmapVal))
			}

			// fmt.Printf("inserting %s %d times\n", f, loops)
			for range loops {
				buff = append(buff, f)
			}
			currFileId++
		case FREE_SPACE_MODE:
			// fmt.Println("Freespace mode", string(diskmapVal))
			loops, err := strconv.Atoi(string(diskmapVal))
			if err != nil {
				// fmt.Println("Error converting to int at idx free space mode", diskIdx, string(diskmapVal))
			}
			for range loops {
				buff = append(buff, EMPTY_SPACE_BLOCK)
			}
		}
	}

	return buff
}

func newFile(val byte, id int) file {
	return file{val, id}
}

func (i *input) compressDiskMap() []block {
	expandedDiskMap := i.expandDiskMap()

	// fmt.Println(string(ex[0:500]))
	// fmt.Println(string(ex))
	start, end := 0, len(expandedDiskMap)-1

	for start <= end {
		// move start up until a free space
		for expandedDiskMap[start] != EMPTY_SPACE_BLOCK && start < end {
			// fmt.Printf("Moving start for free space: idx=%d, \nval=%s\n\n", start, string(ex))
			start++
		}
		// fmt.Printf("Swapping start and end: start=%d(%s), end=%d(%s)\n", start, string(ex[start]), end, string(ex[end]))
		// swap start and end idxs
		expandedDiskMap[start], expandedDiskMap[end] = expandedDiskMap[end], expandedDiskMap[start]
		end--
	}

	return expandedDiskMap
}

func (i *input) checkSum() int {
	ans := 0

	com := i.compressDiskMap()
	// fmt.Println(string(com[0:500]))
	for idx, val := range com {
		if _, ok := val.(emptySpace); ok {
			break
		}

		s, err := strconv.Atoi(val.String())
		if err != nil {
			fmt.Printf("Err at idx=%d, val=%s, s=%d, ans=%d\n", idx, val.String(), s, ans)
			// break
		}
		// fmt.Printf("idx=%d, val=%s, s=%d, ans=%d\n", idx, string(val), s, ans)

		if idx < 100 {
			// fmt.Printf("Adding %d * %d = %d, new total=%d\n", idx, s, idx*s, ans+idx*s)
		}

		ans += idx * s

	}
	return ans
}

func main() {
	f, _ := os.Open("aoc-day9-input.txt")
	defer f.Close()

	i := newInput(f)

	fmt.Println("Expand", len(i.expandDiskMap()))
	fmt.Println("compress", len(i.compressDiskMap()))

	fmt.Println(i.checkSum())

	// fmt.Println(len(i.diskMap))
}
