package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

/*
Day 9
*/

type input struct {
	diskMap     []byte
	expandedMap []block
}

func newInput(r io.Reader) *input {
	v, _ := io.ReadAll(r)
	return &input{diskMap: v}
}

type file struct {
	// val byte // the value at that index from the original disk map
	id         int    // the original index that this came from in the input disk map'
	expanded   string // the complete file contents after expansion
	len        int    // length of expanded
	compressed bool   // has this file already been compressed?
}

type emptySpace struct {
	span int // how many contiguous spaces there is in this span
}

func newEmptySpace(span int) *emptySpace {
	return &emptySpace{span}
}

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

func (e emptySpace) String() string {
	return strings.Repeat(".", e.span)
}

func (f file) String() string {
	return f.expanded
}

var EMPTY_SPACE_BLOCK = emptySpace{}

func (i *input) expandDiskMap() {
	currFileId := 0

	var buff []block

	for diskIdx := range i.diskMap {
		diskmapVal := i.diskMap[diskIdx] // how many loops we need to do to expand

		switch diskIdx % 2 {
		case FILE_MODE:
			loops, _ := strconv.Atoi(string(diskmapVal))
			completeFileRaw := strconv.Itoa(currFileId)
			completeFileStr := strings.Repeat(completeFileRaw, loops)

			f := newFile(currFileId, completeFileStr)
			buff = append(buff, f)
			currFileId++

		case FREE_SPACE_MODE:
			// if the previous block is also an empty file,
			loops, _ := strconv.Atoi(string(diskmapVal))
			// current free space block.
			// TODO: check previous value, if its also a free space,
			// we need to combine the lengths.
			currEmptySpanLen := loops
			if es, prevIsEmptySpace := buff[len(buff)-1].(*emptySpace); prevIsEmptySpace {
				// our previous value is an empty space block, so we should combine them

				// NOTE: is this actually modifying anything?
				es.span += currEmptySpanLen
			} else {
				// previous was a file, so create a fresh empty space block.
				ns := newEmptySpace(currEmptySpanLen)
				buff = append(buff, ns)
			}
		}
	}

	i.expandedMap = buff
}

func newFile(id int, expanded string) *file {
	l := len(expanded)
	return &file{id: id, expanded: expanded, len: l}
}

func insertAtIdx[T any](li []T, idx int, val T) []T {
	return append(li[:idx], append([]T{val}, li[idx:]...)...)
}

func removeAtIndex[T any](s []T, index int) []T {
	return append(s[:index], s[index+1:]...)
}

func swap(l []block, i1, i2 int) []block {
	l[i1], l[i2] = l[i2], l[i1]
	return l
}

func (in *input) nextFileIdx(endIdx int) (int, error) {
	for endIdx > 0 {
		// fmt.Printf("nextFileIdx checking at %d\n", endIdx)
		if v, ok := in.expandedMap[endIdx].(*file); ok && !v.compressed {
			// fmt.Printf("Found a file at %d\n", endIdx)
			return endIdx, nil
		}
		endIdx--
	}

	return -1, fmt.Errorf("file was already compressed! we are done")
}

func (in *input) nextSpaceIdx(endIdx int, minSize int) (int, error) {
	start := 0
	for start <= endIdx {
		if v, ok := in.expandedMap[start].(*emptySpace); ok && v.span >= minSize {
			return start, nil
		}
		start++
	}

	return -1, fmt.Errorf("no space block found that fits that size")
}

func (i *input) compressDiskMap() {
	// fmt.Println("starting compressdiskmap")

	start := 0
	end := len(i.diskMap) - 1

	for start <= end {
		// fmt.Printf("at top of compress loop, start=%d, end=%d\n", start, end)

		fIdx, err := i.nextFileIdx(end)
		if err != nil {
			// fmt.Println("no uncompressed file found! exiting")
			return
		}
		currFile := i.expandedMap[fIdx].(*file)

		// fmt.Println("Got a file, about to look for space")

		// find a space that could fit that file
		spaceIdx, err := i.nextSpaceIdx(end, currFile.len)
		if err != nil {
			// fmt.Println("no valid space available that fits file w/ len", currFile.len, currFile.expanded)
			end = fIdx - 1
			start = 0
			continue
		}
		currSpace := i.expandedMap[spaceIdx].(*emptySpace)
		// fmt.Println("Got a space, about to look check for space left over")

		// do we have any left over space?
		spaceLeftOverAfterMove := currSpace.span - currFile.len

		// if we have space left over, we need to adjust the current space to the size of the file
		// and create a new space after our file that has the remainder.

		if spaceLeftOverAfterMove > 0 {
			currSpace.span = currFile.len
			remainderSpace := newEmptySpace(spaceLeftOverAfterMove)
			i.expandedMap = swap(i.expandedMap, spaceIdx, fIdx)
			i.expandedMap = insertAtIdx(i.expandedMap, spaceIdx+1, block(remainderSpace))
		} else {
			swap(i.expandedMap, spaceIdx, fIdx)
		}

		end = fIdx - 1
		start = 0
	}
}

func (i *input) checkSum() int {
	ans := 0

	var completeDisk strings.Builder

	for _, val := range i.expandedMap {
		completeDisk.WriteString(val.String())
	}

	comp := completeDisk.String()
	// fmt.Println(comp)

	// fmt.Println(comp)

	for idx, s := range comp {

		st, err := strconv.Atoi(string(s))
		if err != nil {
			continue
		}

		ans += idx * st

	}

	return ans
}

func main() {
	f, _ := os.Open("aoc-day9-input.txt")
	defer f.Close()

	// s := []int{1, 2, 3}

	// s = insertAtIdx(s, 1, 5)
	// s = insertAtIdx(s, 2, 10)

	// fmt.Println(s)

	i := newInput(f)

	i.expandDiskMap()
	i.compressDiskMap()
	// fmt.Println(i.expandedMap)

	fmt.Println(i.checkSum())

	fmt.Println(len(i.diskMap))
}
