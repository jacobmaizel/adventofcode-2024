package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

/*
Day 8
antenna has been reconfigured to emit a signal that makes people 0.1% more likely to ...
Each antenna is tuned to a specific frequency indicated by a single lowercase letter, uppercase letter, or digit

signal only applies its nefarious effect at specific antinodes based on the resonant frequencies of the antennas

antinode occurs at any point that is perfectly in line with two antennas of the same frequency - but only when one of the antennas is twice as far away as the other

for any pair of antennas with the same frequency, there are two antinodes, one on either side of them

antinodes can occur at locations that contain antennas

How many unique locations within the bounds of the map contain an antinode?


keep looping through all
until we find a non empty spot
create a new antenna pair with that being the start

keep looping until we find a matching antenna for the one we have
calculate anti node positions for both
and consider that pair completed

now if we are looping and


loop over all of the items
if its non empty space,
add it to a list of its own kind, each with their corresponding positions

ie. saw A at 5,6 8,10, and 10,2
etc

each of them represent a unique combination

calculate the anti node positions for all of the availabile pairs,
using a map set to track unique positions where they are found



*/

type antenna struct {
	pos        position
	signalChar byte
}

func newAntenna(signalChar byte, pos position) antenna {
	return antenna{pos, signalChar}
}

func (a *antenna) String() string {
	return fmt.Sprintf("%s=%s", string(a.signalChar), a.pos.String())
}

func (a *antennaPair) String() string {
	return fmt.Sprintf("%s -> %s", a.start.String(), a.end.String())
}

type position struct {
	row int
	col int
}

type antennaPair struct {
	start antenna
	end   antenna
}

func (p *position) String() string {
	return fmt.Sprintf("[%d,%d]", p.row, p.col)
}

func newPos(row, col int) position {
	return position{row, col}
}

type input struct {
	antinodePositions map[position]bool
	antennaLocations  map[byte][]antenna
	antennaPairs      []antennaPair
	grid              [][]byte
	rows              int
	cols              int
}

// distance from start is distance * 2
// do for both start -> end and end -> start
// save position in input
func (i *input) computePart1() {
	for _, p := range i.antennaPairs {
		// fmt.Printf("Starting AP bounds for pair=%s\n", p.String())
		dist1Row, dist1Col := distance(p.start.pos, p.end.pos)

		// fmt.Printf("dist1Row=%d, dist1Col=%d\n", dist1Row, dist1Col)

		// fmt.Printf("api1 = %d+%d*2, %d+%d*2\n", p.start.pos.row, dist1Row, p.start.pos.col, dist1Col)
		ap1 := newPos(p.start.pos.row+dist1Row*2, p.start.pos.col+dist1Col*2)
		// fmt.Printf("ap1: %s\n", ap1.String())

		if i.inBounds(&ap1) {
			// fmt.Println("ap1 was in bounds!")
			i.antinodePositions[ap1] = true
		}

		dist2Row, dist2Col := distance(p.end.pos, p.start.pos)
		ap2 := newPos(p.end.pos.row+dist2Row*2, p.end.pos.col+dist2Col*2)
		// fmt.Printf("ap2: %s\n", ap2.String())
		if i.inBounds(&ap2) {
			// fmt.Println("ap2 was in bounds!")
			i.antinodePositions[ap2] = true
		}
	}
}

// keep computing antinode positions until its out of bounds
// for both directions
func (i *input) computePart2() {
	for _, p := range i.antennaPairs {
		// fmt.Printf("Starting AP bounds for pair=%s\n", p.String())
		dist1Row, dist1Col := distance(p.start.pos, p.end.pos)

		// fmt.Printf("dist1Row=%d, dist1Col=%d\n", dist1Row, dist1Col)

		// fmt.Printf("api1 = %d+%d*2, %d+%d*2\n", p.start.pos.row, dist1Row, p.start.pos.col, dist1Col)
		i.antinodePositions[p.start.pos] = true
		i.antinodePositions[p.end.pos] = true
		nextPosRow := p.start.pos.row + dist1Row*2
		nextPosCol := p.start.pos.col + dist1Col*2
		ap1 := newPos(nextPosRow, nextPosCol)
		// fmt.Printf("ap1: %s\n", ap1.String())

		for i.inBounds(&ap1) {
			// fmt.Println("ap1 was inbounds")
			i.antinodePositions[ap1] = true

			nextPosRow += dist1Row
			nextPosCol += dist1Col

			ap1 = newPos(nextPosRow, nextPosCol)
			// fmt.Println("setting ap1 to new pos=", ap1)
		}

		// if i.inBounds(&ap1) {
		// 	// fmt.Println("ap1 was in bounds!")
		// }

		dist2Row, dist2Col := distance(p.end.pos, p.start.pos)
		nextPosRow = p.end.pos.row + dist2Row*2
		nextPosCol = p.end.pos.col + dist2Col*2

		ap2 := newPos(nextPosRow, nextPosCol)
		// fmt.Printf("ap2: %s\n", ap2.String())

		for i.inBounds(&ap2) {
			// fmt.Println("ap2 was inbounds")
			i.antinodePositions[ap2] = true

			nextPosRow += dist2Row
			nextPosCol += dist2Col

			ap2 = newPos(nextPosRow, nextPosCol)
			// fmt.Println("setting ap2 to new pos=", ap2)
		}
		// if i.inBounds(&ap2) {
		// 	// fmt.Println("ap2 was in bounds!")
		// 	i.antinodePositions[ap2] = true
		// }
	}
}

// as part of the initial reading, we can build the combos
// as well as their reversals at the same time
func newInput(r io.Reader) *input {
	i := &input{antinodePositions: make(map[position]bool), antennaLocations: make(map[byte][]antenna)}
	s := bufio.NewScanner(r)
	for s.Scan() {
		rawLine := strings.TrimSpace(s.Text())
		rawlineBytes := []byte(rawLine)
		for idx, a := range rawlineBytes {
			if a == EMPTY_SPACE {
				continue
			}
			x := newAntenna(a, newPos(len(i.grid), idx))
			i.createPairs(x)
			i.antennaLocations[a] = append(i.antennaLocations[a], x)
		}
		i.grid = append(i.grid, rawlineBytes)
	}
	i.rows = len(i.grid) - 1
	i.cols = len(i.grid[0]) - 1
	// fmt.Println("max rows=", i.rows, "max cols=", i.cols)
	return i
}

func newAntennaPair(start, end antenna) antennaPair {
	return antennaPair{start, end}
}

func (i *input) createPairs(newA antenna) []antennaPair {
	ans := []antennaPair{}
	list := i.antennaLocations[newA.signalChar]
	for _, a := range list {
		p := newAntennaPair(newA, a)
		ans = append(ans, p)
		i.antennaPairs = append(i.antennaPairs, p)
	}
	return ans
}

func (i *input) inBounds(p *position) bool {
	return p.row <= i.rows && p.row >= 0 && p.col <= i.cols && p.col >= 0
}

const EMPTY_SPACE byte = '.'

// returns a row, col difference between 2 positions
func distance(p1, p2 position) (int, int) {
	return p2.row - p1.row, p2.col - p1.col
}

func (i *input) String() string {
	var b strings.Builder
	for r := range i.grid {
		for c := range i.grid[r] {
			b.WriteString(string(i.grid[r][c]))
		}
		b.WriteString("\n")
	}
	return b.String()
}

func main() {
	f, _ := os.Open("aoc-day8-input.txt")
	defer f.Close()
	i := newInput(f)
	// i.computePart1()
	// fmt.Printf("Day 8 Part 1 ans=%d\n", len(i.antinodePositions))
	i.computePart2()
	fmt.Printf("Day 8 Part 2 ans=%d\n", len(i.antinodePositions))
	// fmt.Println(i.String())
}
