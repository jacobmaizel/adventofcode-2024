package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

/*
Day 4 - Ceres Search

word search for ALL INSTANCES of XMAS (18 times)

- allows for horizontal, vertical, diagonal, backwards, and overlapping other words.


Part 2: X-MAS instead of XMAS

M.S
.A.
M.S


* Within the X, each MAS can be written forwards or backwards.


*/

const (
	XMAS     = "XMAS"
	MAS      = "MAS"
	NOTFOUND = 0
	FOUND    = 1
)

type InputGrid struct {
	data []string
	rows int
	cols int
}

func NewInputGrid(input []string) *InputGrid {
	im := &InputGrid{data: input, rows: len(input), cols: len(input[0])}
	return im
}

func (im *InputGrid) PrintGrid() {
	for x := range im.rows {
		for y := range im.cols {
			fmt.Printf("[%d,%d:%s] ", x, y, string(im.data[x][y]))
		}
		fmt.Print("\n")
	}
}

const (
	UP   = -1
	DOWN = 1
)

const (
	LEFT  = -1
	RIGHT = 1
)

var DIAG_DIRS = [][2]int{
	{LEFT, UP},
	{LEFT, DOWN},
	{RIGHT, DOWN},
	{RIGHT, UP},
}

func (ig *InputGrid) FindDiagonal(x, y int) int {
	if ig.data[x][y] != 'X' {
		return 0
	}
	total := 0
	for _, dir := range DIAG_DIRS {
		xDir := dir[0]
		yDir := dir[1]
		// var word []byte
		match := true

		for i := range len(XMAS) {
			newX := x + (i * xDir)
			newY := y + (i * yDir)

			if newX < 0 || newY < 0 || newX >= ig.rows || newY >= ig.cols {
				match = false
				break
			}

			if ig.data[newX][newY] != XMAS[i] {
				match = false
				break
			}
		}
		if match {
			total++
		}
	}

	return total
}

func (ig *InputGrid) searchLeftRight(x, y, dir int) int {
	for i := range len(XMAS) {
		point := ig.data[x][y+(i*dir)]
		if point != XMAS[i] {
			return NOTFOUND
		}
	}
	return FOUND
}

func (ig *InputGrid) FindHorizontal(x, y int) int {
	total := 0
	if y-len(XMAS) >= -1 {
		total += ig.searchLeftRight(x, y, LEFT)
	}

	if y+len(XMAS) <= ig.cols {
		total += ig.searchLeftRight(x, y, RIGHT)
	}

	return total
}

func (ig *InputGrid) searchUpDown(x, y, mod int) int {
	for i := range len(XMAS) {
		point := ig.data[x+(i*mod)][y]
		if point != XMAS[i] {
			return NOTFOUND
		}
	}
	return FOUND
}

func (ig *InputGrid) FindVertical(x, y int) int {
	total := 0

	if x-len(XMAS) >= -1 {
		total += ig.searchUpDown(x, y, UP)
	}

	if x+len(XMAS) <= ig.rows {
		total += ig.searchUpDown(x, y, DOWN)
	}

	return total
}

func (im *InputGrid) WordSearch() int {
	total := 0
	for x := range im.rows {
		for y := range im.cols {
			if im.data[x][y] != 'X' {
				continue
			}
			total += im.FindVertical(x, y)
			total += im.FindHorizontal(x, y)
			total += im.FindDiagonal(x, y)
		}
	}

	return total
}

var LeftTopToRightBot = [][]int{
	{LEFT, UP},
	{RIGHT, DOWN},
}

var LeftBotToRightTop = [][]int{
	{LEFT, DOWN},
	{RIGHT, UP},
}

func (ig *InputGrid) SearchCrossP2(x, y int) int {
	s1 := []byte{}
	for _, dir := range LeftTopToRightBot {

		newX, newY := x+dir[0], y+dir[1]

		if newX < 0 || newY < 0 || newX >= ig.rows || newY >= ig.cols {
			return NOTFOUND // dir loop
		}

		newPoint := ig.data[newX][newY]

		// if newPoint != 'M' && newPoint != 'S' {
		//   return NOTFOUND
		// }

		s1 = append(s1, newPoint)

	} // left to right
	s := string(s1)
	if s != "MS" && s != "SM" {
		return NOTFOUND
	}

	s2 := []byte{}
	for _, dir := range LeftBotToRightTop {

		newX, newY := x+dir[0], y+dir[1]

		if newX < 0 || newY < 0 || newX >= ig.rows || newY >= ig.cols {
			return NOTFOUND // dir loop
		}

		newPoint := ig.data[newX][newY]

		if newPoint != 'M' && newPoint != 'S' {
			return NOTFOUND
		}

		s2 = append(s2, newPoint)

	} // left to right
	s = string(s2)
	if s != "MS" && s != "SM" {
		return NOTFOUND
	}

	return FOUND
}

func (ig *InputGrid) XMASCrossSearch() int {
	total := 0
	for x := range ig.rows {
		for y := range ig.cols {
			currPoint := ig.data[x][y]
			if currPoint != 'A' {
				continue
			}

			total += ig.SearchCrossP2(x, y)

		} // col loop
	} // rows loop

	return total
}

func main() {
	fmt.Println("aoc day 4")
	file, _ := os.Open("aoc-day4-input.txt")
	defer file.Close()
	str, _ := io.ReadAll(file)
	lines := strings.Split(string(str), "\n")
	im := NewInputGrid(lines[:len(lines)-1])
	res := im.WordSearch()

	p2Res := im.XMASCrossSearch()

	fmt.Printf("D4P1 XMAS Count= %d\n", res)
	fmt.Printf("D4P2 X-MAS Count= %d\n", p2Res)
}
