package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

type Position struct {
	row int
	col int
}

func NewPos(row, col int) *Position {
	return &Position{
		row: row,
		col: col,
	}
}

type Map struct {
	currGuardPos      *Position
	distinctPositions map[Position]bool // set of unique positions
	grid              [][]byte
	rows              int
	cols              int
	outOfBounds       bool
}

const (
	OBSTRUCTION        = '#'
	EMPTY_SPACE        = '.'
	CUSTOM_OBSTRUCTION = '$' // what we will use to block the guard into a loop
)

/*
Part 1
- If there is something directly in front of you, turn right 90 degrees.
- Otherwise, take a step forward.

guard will eventually go out of the map, ie. facing right and hitting right boundary etc.

count all unique positions, INCLUDING start position

--- Part 2

count of different positions we can use to put an obstruction
that would get the guard stuck in a loop

part 2 idea:
brute force -> bunch of go routines each with its own version of the map with 
a unique spot for the new obstruction.
each goroutine will use its map and test for a cycle. 
if it finds one, report back on a channel (fan in)

cycle testing:
build a linked list and do fast+slow runners?
*/

var (
	UP    = '^'
	DOWN  = 'v'
	LEFT  = '<'
	RIGHT = '>'
)

func NewMapFromReader(r io.Reader) *Map {
	m := &Map{
		distinctPositions: make(map[Position]bool),
		grid:              [][]byte{},
	}

	s := bufio.NewScanner(r)

	for s.Scan() {
		rowString := s.Text()

		m.grid = append(m.grid, []byte(rowString))

		row := len(m.grid) - 1

		for col := range len(rowString) {

			currVal := rune(rowString[col])

			if currVal == UP || currVal == DOWN || currVal == LEFT || currVal == RIGHT {

				m.currGuardPos = &Position{
					row: row,
					col: col,
				}

				m.saveGuardPosition()
			}
		}
	}

	m.cols = len(m.grid[0]) - 1
	m.rows = len(m.grid) - 1

	return m
}

func NewMapFromGrid(g []string) *Map {
	m := &Map{
		distinctPositions: make(map[Position]bool),
		grid:              [][]byte{},
	}

	for _, l := range g {
		m.grid = append(m.grid, []byte(l))
	}

	for row, rowString := range g {
		for col := range len(rowString) {

			posVal := rune(rowString[col])

			if posVal == UP || posVal == DOWN || posVal == LEFT || posVal == RIGHT {
				m.currGuardPos = &Position{
					row: row,
					col: col,
				}
				m.saveGuardPosition()
			}

		}
	}

	m.cols = len(g[0]) - 1
	m.rows = len(g) - 1

	return m
}

func (m *Map) printGrid() {
	for r := range m.grid {
		for c := range m.grid[r] {
			fmt.Print(string(m.grid[r][c]))
		}
		fmt.Print("\n")
	}
	fmt.Print("\n")
}

func (m *Map) saveGuardPosition() {
	m.distinctPositions[*m.currGuardPos] = true
}

func (m *Map) getGuardPositionChar() byte {
	return m.grid[m.currGuardPos.row][m.currGuardPos.col]
}

func (m *Map) guardUp() {
	originalPos := NewPos(m.currGuardPos.row, m.currGuardPos.col)

	for m.currGuardPos.row > 0 && m.grid[m.currGuardPos.row-1][m.currGuardPos.col] != OBSTRUCTION {
		m.currGuardPos.row--
		m.saveGuardPosition()
	}

	m.isGuardLeavingGrid(UP)
	m.grid[m.currGuardPos.row][m.currGuardPos.col] = byte(RIGHT)
	m.grid[originalPos.row][originalPos.col] = EMPTY_SPACE
}

func (m *Map) guardRight() {
	originalPos := NewPos(m.currGuardPos.row, m.currGuardPos.col)

	for m.currGuardPos.col < m.cols && m.grid[m.currGuardPos.row][m.currGuardPos.col+1] != OBSTRUCTION {
		m.currGuardPos.col++
		m.saveGuardPosition()
	}

	m.isGuardLeavingGrid(RIGHT)
	m.grid[m.currGuardPos.row][m.currGuardPos.col] = byte(DOWN)
	m.grid[originalPos.row][originalPos.col] = EMPTY_SPACE
}

func (m *Map) guardDown() {
	originalPos := NewPos(m.currGuardPos.row, m.currGuardPos.col)

	for m.currGuardPos.row < m.rows && m.grid[m.currGuardPos.row+1][m.currGuardPos.col] != OBSTRUCTION {
		m.currGuardPos.row++
		m.saveGuardPosition()
	}

	m.isGuardLeavingGrid(DOWN)
	m.grid[m.currGuardPos.row][m.currGuardPos.col] = byte(LEFT)
	m.grid[originalPos.row][originalPos.col] = EMPTY_SPACE
}

func (m *Map) guardLeft() {
	originalPos := NewPos(m.currGuardPos.row, m.currGuardPos.col)

	for m.currGuardPos.col > 0 && m.grid[m.currGuardPos.row][m.currGuardPos.col-1] != OBSTRUCTION {
		m.currGuardPos.col--
		m.saveGuardPosition()
	}

	m.isGuardLeavingGrid(LEFT)
	m.grid[m.currGuardPos.row][m.currGuardPos.col] = byte(UP)
	m.grid[originalPos.row][originalPos.col] = EMPTY_SPACE
}

// bounds check for the guard
func (m *Map) isGuardLeavingGrid(movedDirection rune) {
	if m.currGuardPos.col == 0 && movedDirection == LEFT {
		m.outOfBounds = true
	} else if m.currGuardPos.col == m.cols && movedDirection == RIGHT {
		m.outOfBounds = true
	} else if m.currGuardPos.row == 0 && movedDirection == UP {
		m.outOfBounds = true
	} else if m.currGuardPos.row == m.rows && movedDirection == DOWN {
		m.outOfBounds = true
	}
}

// returns an error when the guard is out of bounds
func (m *Map) moveGuard() {
	for !m.outOfBounds {
		switch m.getGuardPositionChar() {
		case byte(UP):
			m.guardUp()
		case byte(RIGHT):
			m.guardRight()
		case byte(DOWN):
			m.guardDown()
		case byte(LEFT):
			m.guardLeft()
		}
	}
}

func main() {
	file, _ := os.Open("day6-aoc-input.txt")

	m := NewMapFromReader(file)

	m.moveGuard()

	fmt.Printf("Day 6 Part 1 ans=%d\n", len(m.distinctPositions))
}
