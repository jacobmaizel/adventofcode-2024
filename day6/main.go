package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

type State struct {
	pos Position
	dir byte
}

func (s *State) String() string {
	return fmt.Sprintf("%s -> %s", s.pos, string(s.dir))
}

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

func (pos Position) String() string {
	return fmt.Sprintf("row=%d, col=%d", pos.row, pos.col)
}

type Map struct {
	currGuardPos       *Position
	origGuardPos       *Position
	origGuardDir       byte
	distinctPositions  map[Position]byte // set of unique positions
	samePosAndDirCount int               // how many times we have been in the same position facing the same way
	grid               [][]byte
	origGrid           [][]byte
	rows               int
	cols               int
	guardLeavingGrid   bool
	stuckInLoop        bool // flag to determine if we are in a loop
	visited            map[State]bool
}

func (m *Map) reset(previousObstPos *Position, prevVal byte) {
	// m.distinctPositions = make(map[Position]byte)
	m.samePosAndDirCount = 0
	m.guardLeavingGrid = false
	m.stuckInLoop = false
	m.currGuardPos = &Position{row: m.origGuardPos.row, col: m.origGuardPos.col}
	m.setGridCoord(m.currGuardPos.row, m.currGuardPos.col, m.origGuardDir)
	if previousObstPos != nil {
		m.setGridCoord(previousObstPos.row, previousObstPos.col, prevVal)
	}

	m.visited = make(map[State]bool)

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
*/

var (
	UP    = byte('^')
	DOWN  = byte('v')
	LEFT  = byte('<')
	RIGHT = byte('>')
)

func NewMapFromReader(r io.Reader) *Map {
	m := &Map{distinctPositions: make(map[Position]byte), grid: [][]byte{}, visited: make(map[State]bool)}
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
				m.saveStoppingPoint()
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
	m := &Map{distinctPositions: make(map[Position]byte), grid: [][]byte{}, visited: make(map[State]bool)}
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
				m.saveStoppingPoint()
				m.origGuardDir = posVal
				m.origGuardPos = &Position{row: row, col: col}
			}
		}
	}
	m.cols = len(g[0]) - 1
	m.rows = len(g) - 1
	return m
}

func (m *Map) PrintGrid() {
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
func (m *Map) saveStoppingPoint() {
	newState := State{
		pos: *m.currGuardPos,
		dir: m.getGuardDir(),
	}

	if seen := m.visited[newState]; seen {
		m.stuckInLoop = true
		return

	}
	m.visited[newState] = true
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
		orig := m.getGuardDir()
		switch orig {
		case UP:
			m.guardUp()
		case RIGHT:
			m.guardRight()
		case DOWN:
			m.guardDown()
		case LEFT:
			m.guardLeft()
		}
		m.isGuardLeavingGrid(orig)
		m.turnGuard(orig)
		m.saveStoppingPoint()
	}
}

func (m *Map) simulateDifferentObstructionPositions() int {
	totalLoops := 0
	m.moveGuard()
	m.reset(nil, '0')
	for pos := range m.distinctPositions {
		if pos.row == m.origGuardPos.row && pos.col == m.origGuardPos.col {
			continue
		}
		prevVal := m.grid[pos.row][pos.col]
		prevPos := &Position{
			row: pos.row,
			col: pos.col,
		}
		m.setGridCoord(pos.row, pos.col, OBSTRUCTION)
		m.moveGuard()
		if m.stuckInLoop {
			totalLoops++
		}
		m.reset(prevPos, prevVal)
	}

	return totalLoops
}

func main() {
	file, _ := os.Open("day6-aoc-input.txt")
	defer file.Close()

	// m := NewMapFromReader(file)

	// m.moveGuard()

	// fmt.Printf("Day 6 Part 1 ans=%d\n", len(m.distinctPositions))

	p2 := NewMapFromReader(file)
	p2Res := p2.simulateDifferentObstructionPositions()

	fmt.Printf("Day 6 Part 2 ans=%d\n", p2Res)
}
