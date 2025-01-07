package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

/*
a good hiking trail is as long as possible and has an even, gradual, uphill slope.

a hiking trail is any path that starts at height 0, ends at height 9, and always increases by a height of exactly 1 at each step

Hiking trails never include diagonal steps - only up, down, left, or right (from the perspective of the map)

trailhead is any position that starts one or more hiking trails - here, these positions will always have height 0

NOTE: a trailhead's score is the number of 9-height positions reachable from that trailhead via a hiking trail

// approach 1: brute force, traverse each spot, if its a 0, start from it,
keep looking for next steps in each direction until we either hit a 9 or
cant go any farther.
ie. can explore up? explore up. move up. any directions from here a valid next
step? explore that. etc.

// approach 2: when bulding input keep track of starting spots
only start searches from there

// approach 3: recursive traversal at each valid next step

// approach 4: can we do some fancy Dynamic programming?
*/

type position struct {
	row int
	col int
}

func newPos(row, col int) position {
	return position{row, col}
}

type input struct {
	grid               [][]int
	trailheadPositions []position
	rows               int
	cols               int
}

func newInput(r io.Reader) *input {
	s := bufio.NewScanner(r)
	in := &input{}
	for s.Scan() {
		row := strings.TrimSpace(s.Text())
		r := make([]int, len(row))
		for i := range len(row) {
			if row[i] == '.' {
				r[i] = -1 // Use -1 to represent a period (or another unique value not in your number range)
			} else {
				r[i] = int(row[i] - '0') // Convert character digit to integer
				if r[i] == 0 {
					in.trailheadPositions = append(in.trailheadPositions, newPos(len(in.grid), i))
				}
			}
		}

		in.grid = append(in.grid, r)
	}

	in.rows = len(in.grid)
	in.cols = len(in.grid[0])

	return in
}

func (in *input) String() string {
	var s strings.Builder
	for r := range in.rows {
		for c := range in.cols {
			v := in.grid[r][c]
			if v == -1 {
				s.WriteByte('.')
			} else {
				s.WriteString(strconv.Itoa(v))
			}
		}
		s.WriteString("\n")
	}
	return s.String()
}

var DIRS = [][2]int{
	{0, 1},
	{1, 0},
	{0, -1},
	{-1, 0},
}

func (in *input) getGridVal(row, col int) int {
	return in.grid[row][col]
}

func (in *input) inBounds(row, col int) bool {
	return row < in.rows && row >= 0 && col < in.cols && col >= 0
}

// calcs a list of directions that are valid to take from a given position
// taking into account path rule that a next step has to be +1 our current pos val
func (in *input) nextSteps(currRow, currCol int) [][2]int {
	ans := [][2]int{}

	for _, dir := range DIRS {
		nR, nC := currRow+dir[0], currCol+dir[1]
		currVal := in.getGridVal(currRow, currCol)
		if in.inBounds(nR, nC) && in.getGridVal(nR, nC) == currVal+1 {
			ans = append(ans, dir)
		}
	}

	return ans
}

func (in *input) traverseTrail(currRow, currCol int, reached map[position]int) {
	currP := newPos(currRow, currCol)

	if in.getGridVal(currRow, currCol) == 9 {
		reached[currP]++
		return
	}

	ns := in.nextSteps(currRow, currCol)
	for _, dir := range ns {
		in.traverseTrail(currRow+dir[0], currCol+dir[1], reached)
	}
}

// how many distinct paths are there from each trailhead to a 9-height position?
func (in *input) calcRating() int {
	total := 0

	for _, starts := range in.trailheadPositions {
		reached := make(map[position]int)

		in.traverseTrail(starts.row, starts.col, reached)

		for _, k := range reached {
			total += k
		}
	}

	return total
}

func (in *input) calcScore() int {
	total := 0

	for _, starts := range in.trailheadPositions {
		reached := make(map[position]int)

		in.traverseTrail(starts.row, starts.col, reached)

		total += len(reached)
	}

	return total
}

func main() {
	f, _ := os.Open("input-day10.txt")

	in := newInput(f)

	fmt.Println("D10P1->\n", in.calcScore())
	fmt.Println("D10P2->\n", in.calcRating())
}
