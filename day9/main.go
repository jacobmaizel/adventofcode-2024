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
	input      []byte
	fileSystem []int
}

func newInput(r io.Reader) *input {
	v, _ := io.ReadAll(r)
	return &input{input: v}
}

const (
	FILE_MODE = iota
	FREE_SPACE_MODE
)

func (i *input) String() string {
	var s strings.Builder

	for _, v := range i.fileSystem {
		if v == -1 {
			s.WriteString(".")
		} else {
			con := strconv.Itoa(v)
			s.WriteString(con)
		}
	}

	return s.String()
}

func (i *input) buildFileSystem() {
	fileIdCounter := 0

	var fs []int

	for diskIdx := range i.input {
		diskmapVal := i.input[diskIdx] // how many loops we need to do to expand
		loops := int(diskmapVal - '0')
		switch diskIdx % 2 {
		case FILE_MODE:
			for range loops {
				fs = append(fs, fileIdCounter)
			}
			fileIdCounter++
		case FREE_SPACE_MODE:

			for range loops {
				fs = append(fs, -1)
			}
		}
	}

	i.fileSystem = fs
}

func (in *input) compressFileSystem() {
	end := len(in.fileSystem) - 1

	for end >= 0 {
		// start := 0
		// get next file
		for end >= 0 && in.fileSystem[end] == -1 {
			end--
		}
		if end < 0 {
			break
		}
		fileId := in.fileSystem[end]
		fileSize := in.getFileSize(fileId, end)
		nextSpaceStartIdx, spaceFound := in.nextSpaceIdx(fileSize, end-fileSize+1)

		if spaceFound {
			in.swapBlocks(fileSize, nextSpaceStartIdx, end)
		}

		end -= fileSize

	}
}

func swap(l []int, i1, i2 int) []int {
	l[i1], l[i2] = l[i2], l[i1]
	return l
}

func (in *input) getFileSize(fileId int, startIdx int) int {
	ctr := 0
	for startIdx >= 0 && in.fileSystem[startIdx] == fileId {
		startIdx--
		ctr++
	}
	return ctr
}

// starting always from 0, look for a block of spaces of adequate size for given file
func (in *input) nextSpaceIdx(size, end int) (int, bool) {
	start := 0

	for start < end {
		i := 0
		found := true
		for i < size {
			if in.fileSystem[start+i] != -1 {
				found = false
				break
			}
			i++
		}
		if found {
			return start, true
		}
		start += i
		for start < len(in.fileSystem) && in.fileSystem[start] != -1 {
			start++
		}
	}

	return start, false
}

func (in *input) swapBlocks(size int, start, end int) {
	for range size {
		in.fileSystem = swap(in.fileSystem, start, end)
		start++
		end--
	}
}

func (in *input) checkSum() int {
	ans := 0
	in.buildFileSystem()

	in.compressFileSystem()

	for idx, val := range in.fileSystem {
		if val == -1 {
			continue
		}
		ans += idx * val
	}

	return ans
}

func main() {
	f, _ := os.Open("aoc-day9-input.txt")

	i := newInput(f)

	i.buildFileSystem()

	i.compressFileSystem()
	fmt.Println(i.checkSum())
}
