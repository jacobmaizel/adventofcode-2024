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
	currGuardPos         *Position
	origGuardPos         *Position
	origGuardDir         byte
	distinctPositions    map[Position]byte // set of unique positions
	stoppingPointsLookup map[Position]byte // used to determine only stopping points and dirs, not each individual point.
	samePosAndDirCount   int               // how many times we have been in the same position facing the same way
	grid                 [][]byte
	origGrid             [][]byte
	rows                 int
	cols                 int
	guardLeavingGrid     bool
	stuckInLoop          bool // flag to determine if we are in a loop
}

func (m *Map) reset(previousObstPos *Position, prevVal byte) {
	m.distinctPositions = make(map[Position]byte)
	m.stoppingPointsLookup = make(map[Position]byte)
	m.samePosAndDirCount = 0
	m.guardLeavingGrid = false
	m.stuckInLoop = false
	m.currGuardPos = &Position{row: m.origGuardPos.row, col: m.origGuardPos.col}
	m.setGridCoord(m.currGuardPos.row, m.currGuardPos.col, m.origGuardDir)
	if previousObstPos != nil {
		m.setGridCoord(previousObstPos.row, previousObstPos.col, prevVal)
	}

	copy(m.origGrid, m.grid)

	// m.rows = 0
	// m.cols = 0
	// m.grid = [][]byte{}
}

const (
	OBSTRUCTION        = byte('#')
	EMPTY_SPACE        = byte('.')
	CUSTOM_OBSTRUCTION = byte('$') // what we will use to block the guard into a loop
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

we are in a loop when:
1. we hit a place that we have already been
2. We are facing the same direction that we did when we were last there
3. and we did this 4 times.
4. the only way we would be in the same spot, and facing the same direciton 4 times, is
if we were in a box.
*/

var (
	UP    = byte('^')
	DOWN  = byte('v')
	LEFT  = byte('<')
	RIGHT = byte('>')
)

func NewMapFromReader(r io.Reader) *Map {
	m := &Map{distinctPositions: make(map[Position]byte), grid: [][]byte{}, stoppingPointsLookup: make(map[Position]byte)}
	s := bufio.NewScanner(r)
	for s.Scan() {
		rowString := s.Text()
		m.grid = append(m.grid, []byte(rowString))
		m.origGrid = append(m.origGrid, []byte(rowString))
		row := len(m.grid) - 1
		for col := range len(rowString) {
			currVal := rowString[col]
			if currVal == UP || currVal == DOWN || currVal == LEFT || currVal == RIGHT {
				m.currGuardPos = &Position{row: row, col: col}
				m.saveGuardPosition()
				// m.saveStoppingPoint()
				m.origGuardDir = currVal
				m.origGuardPos = &Position{row: row, col: col}
			}
		}
	}
	m.cols = len(m.grid[0]) - 1
	m.rows = len(m.grid) - 1
	return m
}

func NewMapFromGrid(g []string) *Map {
	m := &Map{distinctPositions: make(map[Position]byte), grid: [][]byte{}, stoppingPointsLookup: make(map[Position]byte)}
	for _, l := range g {
		m.grid = append(m.grid, []byte(l))
		m.origGrid = append(m.origGrid, []byte(l))
	}
	for row, rowString := range g {
		for col := range len(rowString) {
			posVal := rowString[col]
			if posVal == UP || posVal == DOWN || posVal == LEFT || posVal == RIGHT {
				m.currGuardPos = &Position{row: row, col: col}
				m.saveGuardPosition()
				// m.saveStoppingPoint()
				m.origGuardDir = posVal
				m.origGuardPos = &Position{row: row, col: col}
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

// used each time the guard moves
func (m *Map) saveGuardPosition() {
	currGuardPositionChar := m.getGuardDir()
	m.distinctPositions[*m.currGuardPos] = currGuardPositionChar
}

// only used when the guard stops and turns. helpful for cycles!
// func (m *Map) saveStoppingPoint() {
// 	currGuardPositionChar := m.getGuardDir()
// 	m.stoppingPointsLookup[*m.currGuardPos] = currGuardPositionChar
// }

func (m *Map) checkGuardInLoop() {
	if dir, ok := m.stoppingPointsLookup[*m.currGuardPos]; ok && m.getGuardDir() == dir {
		m.samePosAndDirCount++
		if m.samePosAndDirCount == 4 {
			fmt.Printf("GOT A LOOP! current stopping poitns:%+v\n\n", m.stoppingPointsLookup)
			m.printGrid()
			m.stuckInLoop = true
		}
	} else {
		m.samePosAndDirCount = 0
	}
}

func (m *Map) setGridCoord(row, col int, val byte) {
	m.grid[row][col] = val
}

func (m *Map) getGuardDir() byte {
	// fmt.Println("getGuardPositionChar")
	return m.grid[m.currGuardPos.row][m.currGuardPos.col]
}

func (m *Map) guardUp() {
	// fmt.Println("guardUp")
	originalPos := NewPos(m.currGuardPos.row, m.currGuardPos.col)

	for m.currGuardPos.row > 0 && m.grid[m.currGuardPos.row-1][m.currGuardPos.col] != OBSTRUCTION {
		m.currGuardPos.row--
		m.saveGuardPosition()
	}

	if originalPos.row != m.currGuardPos.row || originalPos.col != m.currGuardPos.col {
		m.grid[originalPos.row][originalPos.col] = EMPTY_SPACE
	}
}

func (m *Map) guardRight() {
	// fmt.Println("guardRight")
	originalPos := NewPos(m.currGuardPos.row, m.currGuardPos.col)

	for m.currGuardPos.col < m.cols && m.grid[m.currGuardPos.row][m.currGuardPos.col+1] != OBSTRUCTION {
		m.currGuardPos.col++
		m.saveGuardPosition()
	}

	if originalPos.row != m.currGuardPos.row || originalPos.col != m.currGuardPos.col {
		m.grid[originalPos.row][originalPos.col] = EMPTY_SPACE
	}
}

func (m *Map) guardDown() {
	// fmt.Println("guardDown")
	originalPos := NewPos(m.currGuardPos.row, m.currGuardPos.col)

	for m.currGuardPos.row < m.rows && m.grid[m.currGuardPos.row+1][m.currGuardPos.col] != OBSTRUCTION {
		m.currGuardPos.row++
		m.saveGuardPosition()
	}

	if originalPos.row != m.currGuardPos.row || originalPos.col != m.currGuardPos.col {
		m.grid[originalPos.row][originalPos.col] = EMPTY_SPACE
	}
}

func (m *Map) guardLeft() {
	// fmt.Println("guardLeft")
	originalPos := NewPos(m.currGuardPos.row, m.currGuardPos.col)

	for m.currGuardPos.col > 0 && m.grid[m.currGuardPos.row][m.currGuardPos.col-1] != OBSTRUCTION {
		m.currGuardPos.col--
		m.saveGuardPosition()
	}

	if originalPos.row != m.currGuardPos.row || originalPos.col != m.currGuardPos.col {
		m.grid[originalPos.row][originalPos.col] = EMPTY_SPACE
	}
}

// bounds check for the guard
func (m *Map) isGuardLeavingGrid(movedDirection byte) {
	// fmt.Println("isGuardLeavingGrid, pos", m.currGuardPos)
	if m.currGuardPos.col == 0 && movedDirection == LEFT {
		m.guardLeavingGrid = true
	} else if m.currGuardPos.col == m.cols && movedDirection == RIGHT {
		m.guardLeavingGrid = true
	} else if m.currGuardPos.row == 0 && movedDirection == UP {
		m.guardLeavingGrid = true
	} else if m.currGuardPos.row == m.rows && movedDirection == DOWN {
		m.guardLeavingGrid = true
	}
}

func (m *Map) turnGuard(fromDir byte) {
	r := m.currGuardPos.row
	c := m.currGuardPos.col
	switch fromDir {
	case UP:
		m.setGridCoord(r, c, RIGHT)
	case RIGHT:
		m.setGridCoord(r, c, DOWN)
	case DOWN:
		m.setGridCoord(r, c, LEFT)
	case LEFT:
		m.setGridCoord(r, c, UP)
	}
}

// returns when guard is marked as out of bounds
func (m *Map) moveGuard() {
	for !m.guardLeavingGrid && !m.stuckInLoop {
		// fmt.Println("START---------------")
		// m.printGrid()
		orig := m.getGuardDir()
		// fmt.Println("Moving guard from dir", string(orig))
		// fmt.Printf("\nin move.. pos=%+v\n", m.currGuardPos)
		switch orig {
		case UP:
			m.guardUp()
		case RIGHT:
			m.guardRight()
		case DOWN:
			m.guardDown()
		case LEFT:
			m.guardLeft()
		default:
			fmt.Printf("How do we not have a valid dir char?%+v\n", m.currGuardPos)
			m.printGrid()
		}
		fmt.Println("after move switch, guard is now at", m.currGuardPos)
		m.isGuardLeavingGrid(orig)
		fmt.Println("Checked if guard left", m.guardLeavingGrid)
		fmt.Println("Saved stop point")
		m.turnGuard(orig)
		fmt.Println("turned guard")
		// m.saveStoppingPoint()
		fmt.Println("Saved stopping point")
		m.checkGuardInLoop()
		fmt.Println("Checked if guard is not in loop", m.stuckInLoop)

		// fmt.Println("AFTER-----------")
		// m.printGrid()
	}
	// fmt.Printf("Exiting move guard:: loop=%t, left grid=%t, -- ending position=%+v\n", m.stuckInLoop, m.guardLeavingGrid, m.currGuardPos)
}

// brute forces every single position on the map and tests each one to see if
// putting an obstruction there causes a loop
func (m *Map) simulateDifferentObstructionPositions() int {
	totalLoops := 0
	for r := range m.rows {
		for c := range m.cols {
			// m.reset(&Position{row: r, col: c})

			if r == m.origGuardPos.row && c == m.origGuardPos.col {
				continue
			}

			prevVal := m.grid[r][c]

			m.setGridCoord(r, c, OBSTRUCTION)

			m.moveGuard()

			if m.stuckInLoop {
				fmt.Println("loop for obst at ", r, c)
				// m.printGrid()
				totalLoops++
			}
			// m.printGrid()

			m.reset(&Position{row: r, col: c}, prevVal)

		}
		// only do first row to test
		// if r > 2 {
		// 	break
		// }
	}

	return totalLoops
}

func main() {
	file, _ := os.Open("day6-aoc-input.txt")
	defer file.Close()

	m := NewMapFromReader(file)

	m.moveGuard()

	fmt.Printf("Day 6 Part 1 ans=%d\n", len(m.distinctPositions))

	// p2 := NewMapFromReader(file)
	// p2Res := p2.simulateObstPlacmentLoopCount()

	// fmt.Printf("Day 6 Part 2 ans=%d\n", p2Res)
}
