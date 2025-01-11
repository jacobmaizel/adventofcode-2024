package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// recursive with seen tracking, 4-dir expansion check for same crop type until
// expand if not seen one of our own type,
// or start new exploration if its a type we have not seen and its not same as current crop

type position struct {
	row int
	col int
}

func newPos(row, col int) position {
	return position{row, col}
}

type plotLookup map[byte][][]position

type input struct {
	grid [][]byte
	rows int
	cols int
}

var DIRS = [][2]int{
	{-1, 0}, // north
	{1, 0},  // south
	{0, -1}, // west
	{0, 1},  // east
}

func (in *input) String() string {
	var s strings.Builder

	for r := range in.rows {
		for c := range in.cols {
			s.WriteByte(in.grid[r][c])
		}
		s.WriteString("\n")
	}

	return s.String()
}

func (p plotLookup) String() string {
	var s strings.Builder

	for k, v := range p {
		s.WriteString(fmt.Sprintf("%s -> \n", string(k)))
		for i, l := range v {
			s.WriteString(fmt.Sprintf("%d = %v\n", i, l))
		}
		s.WriteString("\n")
	}
	return s.String()
}

func newInput(r io.Reader) *input {
	s := bufio.NewScanner(r)
	// plots plotLookup // ie. get all plots for crop A. (get a list of list of positions back)
	i := &input{}

	for s.Scan() {
		t := strings.TrimSpace(s.Text())
		i.grid = append(i.grid, []byte(t))
	}

	// fmt.Println(len(i.grid))

	i.rows = len(i.grid)
	i.cols = len(i.grid[0])

	return i
}

// global seen tracker
// if we hit a different crop than what we are currently looking for,
// and that crop has NOT been seen before, start an exploration there

// if we hit a point where everything surrounding our current point
// has been seen before, AND non of our surroudning squares are of the current byte
// we are done with this plot, add it to the plot lookup list for this byte and return

// recursion:
// add current position to seen + current exploration

// all of the surroudning same-crops that have not been seen yet, pass
func (in *input) getValPos(pos position) byte {
	return in.grid[pos.row][pos.col]
}

func (in *input) getVal(r, c int) byte {
	return in.grid[r][c]
}

// loops over every r c of the grid and when it finds a non seen value it starts a fresh exploration

func (in *input) bfs(startingPos position, seen map[position]bool) []position {
	queue := []position{startingPos}
	currentExploration := []position{}
	// currVal := in.getValPos(startingPos)

	seen[startingPos] = true
	for len(queue) > 0 {

		current := queue[0]
		currVal := in.getValPos(current)
		queue = queue[1:] // pop

		currentExploration = append(currentExploration, current)

		for _, dir := range DIRS {

			nextPos := newPos(dir[0]+current.row, dir[1]+current.col)

			if in.inBounds(nextPos.row, nextPos.col) && in.getValPos(nextPos) == currVal && !seen[nextPos] {
				seen[nextPos] = true
				queue = append(queue, nextPos)
			}

		}
	}
	return currentExploration
}

func (in *input) partOneQueueBFS(plots plotLookup) {
	seen := make(map[position]bool)
	for r := range in.rows {
		for c := range in.cols {
			pos := newPos(r, c)
			if !seen[pos] {

				exploration := in.bfs(pos, seen)

				v := in.getValPos(pos)
				// fmt.Println(exploration)
				plots[v] = append(plots[v], exploration)
			}
		}
	}
	// fmt.Println("Seen", len(seen), "grid size", in.rows*in.cols)
}

func (in *input) partOneRecursiveDFS(plots plotLookup) {
	seen := make(map[position]bool)
	for r := range in.rows {
		for c := range in.cols {
			pos := newPos(r, c)
			if alreadySeen := seen[pos]; !alreadySeen {
				exploration := []position{}
				in.recursiveExplore(pos, &exploration, seen)
				v := in.getValPos(pos)
				// fmt.Println(exploration)
				plots[v] = append(plots[v], exploration)
			}
		}
	}
	// fmt.Println("Seen", len(seen), "grid size", in.rows*in.cols)
}

// look in all dirs for the current
func (in *input) recursiveExplore(pos position, currentExploration *[]position, seen map[position]bool) {
	if ok := seen[pos]; ok {
		return
	}

	seen[pos] = true

	*currentExploration = append(*currentExploration, pos)

	currVal := in.getValPos(pos)
	// fmt.Println(string(currVal), currentExploration)

	for _, dir := range DIRS {
		nextPos := newPos(dir[0]+pos.row, dir[1]+pos.col)
		if !in.inBounds(nextPos.row, nextPos.col) {
			continue
		}

		nextVal := in.getValPos(nextPos)
		if nextVal == currVal {
			in.recursiveExplore(nextPos, currentExploration, seen)
		}

	}
}

func (in *input) inBounds(row, col int) bool {
	return row < in.rows && row >= 0 && col < in.cols && col >= 0
}

func (in *input) inBoundsPos(p position) bool {
	return p.row < in.rows && p.row >= 0 && p.col < in.cols && p.col >= 0
}

func (in *input) sameCrop(i, j int, crop byte) bool {
	return in.inBounds(i, j) && in.getVal(i, j) == crop
}

func (in *input) getCorners(r, c int) int {
	crop := in.getVal(r, c)

	NW := in.sameCrop(r-1, c-1, crop)
	N := in.sameCrop(r-1, c, crop)
	NE := in.sameCrop(r-1, c+1, crop)
	W := in.sameCrop(r, c-1, crop)
	E := in.sameCrop(r, c+1, crop)
	SW := in.sameCrop(r+1, c-1, crop)
	S := in.sameCrop(r+1, c, crop)
	SE := in.sameCrop(r+1, c+1, crop)

	return sum([]bool{
		N && W && !NW,
		N && E && !NE,
		S && W && !SW,
		S && E && !SE,
		!N && !W,
		!N && !E,
		!S && !W,
		!S && !E,
	})
}

// only counts unique external sides
// group-level counting
func (in *input) calcPart2FenceCost(plots plotLookup) int {
	totalFenceCost := 0

	for _, regions := range plots {
		for _, region := range regions {
			corners := 0
			for _, currPosition := range region {
				corners += in.getCorners(currPosition.row, currPosition.col)
			}
			totalFenceCost += corners * len(region)
		}
	}

	return totalFenceCost
}

func sum(values []bool) int {
	count := 0
	for _, v := range values {
		if v {
			count++
		}
	}
	return count
}

// total price for fencing needed
func (in *input) calcPart1FenceCost(plots plotLookup) int {
	totalFenceCost := 0

	// plots := make(plotLookup)

	for plotType, regions := range plots {
		for _, region := range regions {

			regionFenceCount := 0

			for _, coord := range region {
				for _, dir := range DIRS {
					nextPos := newPos(dir[0]+coord.row, dir[1]+coord.col)
					if !in.inBounds(nextPos.row, nextPos.col) {
						regionFenceCount++
						continue
					}

					if in.getValPos(nextPos) == plotType {
						continue
					}
					regionFenceCount++

				}
			}

			// fmt.Println("processed region", region, "with perim count", regionFenceCount)

			totalFenceCost += regionFenceCount * len(region)

		}
	}

	return totalFenceCost
}

func main() {
	f, _ := os.Open("input.txt")
	i := newInput(f)

	// i.gridSearch()
	plots := make(plotLookup)
	i.partOneQueueBFS(plots)
	p1 := i.calcPart2FenceCost(plots)

	fmt.Println(p1)
}
