package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

/*
Day 15 - warehouse woes
*/

func openInputFile() *os.File {
	f, _ := os.Open("input.txt")
	return f
}

type position struct {
	x, y int
}
type dir struct {
	x, y int
}

type box struct {
	left  *position
	right *position
}

const (
	LEFT  = byte('<')
	RIGHT = byte('>')
	UP    = byte('^')
	DOWN  = byte('v')
)

var DIRS = map[byte]dir{
	LEFT:  {-1, 0},
	RIGHT: {1, 0},
	UP:    {0, -1},
	DOWN:  {0, 1},
}

const (
	ROBOT     = byte('@')
	BOX_O     = byte('O')
	BOX_LEFT  = byte('[')
	BOX_RIGHT = byte(']')
	WALL      = byte('#')
	EMPTY     = byte('.')
)

type input struct {
	grid       [][]byte
	robot      *position
	movements  []byte
	rows, cols int
}

func (in *input) expandMap() {
	fmt.Println("Expanding map")
	newMap := [][]byte{}

	for r := range in.rows {
		newRow := []byte{}
		for c := range in.cols {
			currOriginalVal := in.grid[r][c]

			switch currOriginalVal {
			case WALL:
				newRow = append(newRow, []byte{WALL, WALL}...)
			case BOX_O:
				newRow = append(newRow, []byte{BOX_LEFT, BOX_RIGHT}...)
			case EMPTY:
				newRow = append(newRow, []byte{EMPTY, EMPTY}...)
			case ROBOT:
				newRow = append(newRow, []byte{ROBOT, EMPTY}...)
			}
		}
		newMap = append(newMap, newRow)
	}


	in.grid = newMap
	in.rows = len(in.grid)
	in.cols = len(in.grid[0])

	for r := range in.rows {
		for c := range in.cols {
			if in.grid[r][c] == ROBOT {
				in.robot = &position{r, c}
			}
		}
	}

}

func newInput(r io.Reader) *input {
	s := bufio.NewScanner(r)

	in := &input{}
	scanningMoves := false

	for s.Scan() {
		t := s.Text()
		if t == "" {
			scanningMoves = true
			continue
		}

		if scanningMoves {
			b := []byte(strings.TrimSpace(t))
			in.movements = append(in.movements, b...)
		} else {
			b := []byte(strings.TrimSpace(t))

			for idx, b := range b {
				if b == ROBOT {
					in.robot = &position{len(in.grid), idx}
				}
			}
			in.grid = append(in.grid, b)
		}
	}

	in.rows = len(in.grid)
	in.cols = len(in.grid[0])

	return in
}

func (in *input) handleMoveInDir(dirByte byte) {
	direc := DIRS[dirByte]
	trainPositions := []*position{in.robot}

	gridSearcher := &position{in.robot.x, in.robot.y}

	if in.grid[in.robot.x+direc.y][in.robot.y+direc.x] == WALL {
		return
	}


	gridSearcher.x += direc.y
	gridSearcher.y += direc.x
	for in.grid[gridSearcher.x][gridSearcher.y] == BOX_O {

		trainPositions = append(trainPositions, &position{gridSearcher.x, gridSearcher.y})

		gridSearcher.x += direc.y
		gridSearcher.y += direc.x
	}

	trainMoveCounter := 0

	for in.grid[gridSearcher.x][gridSearcher.y] == EMPTY {
		trainMoveCounter++
		gridSearcher.x += direc.y
		gridSearcher.y += direc.x
	}

	for _, t := range trainPositions {
		in.grid[t.x][t.y] = EMPTY
	}

	if trainMoveCounter > 0 {
		for _, trainCar := range trainPositions {
			trainCar.x += direc.y
			trainCar.y += direc.x
		}
	}

	afterMoveRobotPosition := trainPositions[0]

	if len(trainPositions) > 0 {
		boxPositions := trainPositions[1:]
		for _, b := range boxPositions {
			in.grid[b.x][b.y] = BOX_O
		}
	}

	in.grid[in.robot.x][in.robot.y] = EMPTY
	in.grid[afterMoveRobotPosition.x][afterMoveRobotPosition.y] = ROBOT

	in.robot.x = afterMoveRobotPosition.x
	in.robot.y = afterMoveRobotPosition.y
}

func (in *input) printGrid() {
	var b strings.Builder

	for row := range in.rows {
		for col := range in.cols {
			b.WriteByte(in.grid[row][col])
		}
		b.WriteByte('\n')
	}
	b.WriteByte('\n')

	fmt.Println(b.String())
}

func (in *input) bfsUp(start position) ([]*box, error) {
	queue := []*box{}
	res := []*box{}

	box := in.boxFromPos(start)
	if box != nil {
		queue = append(queue, box)
	}


	for len(queue) > 0 {
		currBox := queue[0]
		queue = queue[1:]

		itemAboveLeftSideOfCurrBox := in.grid[currBox.left.x-1][currBox.left.y]
		itemAboveRightSideOfCurrBox := in.grid[currBox.right.x-1][currBox.right.y]

		if itemAboveLeftSideOfCurrBox == WALL || itemAboveRightSideOfCurrBox == WALL {
			return nil, fmt.Errorf("wall in front of curr box! no op!%+v", currBox.String())
		}

		if itemAboveLeftSideOfCurrBox == BOX_LEFT && itemAboveRightSideOfCurrBox == BOX_RIGHT {
			queue = append(queue, in.boxFromPos(position{currBox.left.x - 1, currBox.left.y}))
		} else {
			if itemAboveLeftSideOfCurrBox == BOX_RIGHT {
				queue = append(queue, in.boxFromPos(position{currBox.left.x - 1, currBox.left.y}))
			}
			if itemAboveRightSideOfCurrBox == BOX_LEFT {
				queue = append(queue, in.boxFromPos(position{currBox.right.x - 1, currBox.right.y}))
			}
		}
		res = append(res, currBox)
	}

	return res, nil
}

func (in *input) bfsDown(start position) ([]*box, error) {
	queue := []*box{}
	res := []*box{}

	box := in.boxFromPos(start)
	if box != nil {
		queue = append(queue, box)
	}

	for len(queue) > 0 {
		currBox := queue[0]
		queue = queue[1:]

		itemBelowLeftSideOfBox := in.grid[currBox.left.x+1][currBox.left.y]
		itemBelowRightSideOfBox := in.grid[currBox.right.x+1][currBox.right.y]

		if itemBelowLeftSideOfBox == WALL || itemBelowRightSideOfBox == WALL {
			return nil, fmt.Errorf("wall in front of curr box! no op!%+v", currBox.String())
		}

		if itemBelowLeftSideOfBox == BOX_LEFT && itemBelowRightSideOfBox == BOX_RIGHT {
			queue = append(queue, in.boxFromPos(position{currBox.left.x + 1, currBox.left.y}))
		} else {
			if itemBelowLeftSideOfBox == BOX_RIGHT {
				queue = append(queue, in.boxFromPos(position{currBox.left.x + 1, currBox.left.y}))
			}
			if itemBelowRightSideOfBox == BOX_LEFT {
				queue = append(queue, in.boxFromPos(position{currBox.right.x + 1, currBox.right.y}))
			}
		}

		res = append(res, currBox)
	}

	return res, nil
}

func (in *input) boxFromPos(p position) *box {
	v := in.grid[p.x][p.y]
	switch v {
	case BOX_LEFT:
		return &box{left: &position{p.x, p.y}, right: &position{p.x, p.y + 1}}
	case BOX_RIGHT:
		return &box{left: &position{p.x, p.y - 1}, right: &position{p.x, p.y}}
	}

	return nil
}

func (in *input) p2Up(d dir) {
	gridSearcher := &position{in.robot.x, in.robot.y}

	if in.grid[in.robot.x+d.y][in.robot.y+d.x] == WALL {
		return
	}

	gridSearcher.x += d.y
	gridSearcher.y += d.x

	objectAbove := in.grid[gridSearcher.x][gridSearcher.y]

	if objectAbove == EMPTY {
		in.grid[in.robot.x][in.robot.y] = EMPTY

		in.robot.x += d.y
		in.robot.y += d.x

		in.grid[in.robot.x][in.robot.y] = ROBOT
		return
	}


	boxTrain, err := in.bfsUp(*gridSearcher)
	if err != nil {
		return
	}

	for _, b := range boxTrain {
		in.grid[b.left.x][b.left.y] = EMPTY
		in.grid[b.right.x][b.right.y] = EMPTY
	}

	in.grid[in.robot.x][in.robot.y] = EMPTY

	for _, b := range boxTrain {
		b.left.x += d.y
		b.left.y += d.x
		b.right.x += d.y
		b.right.y += d.x
	}

	in.robot.x += d.y
	in.robot.y += d.x

	for _, b := range boxTrain {
		in.grid[b.left.x][b.left.y] = BOX_LEFT
		in.grid[b.right.x][b.right.y] = BOX_RIGHT
	}

	in.grid[in.robot.x][in.robot.y] = ROBOT
}

func (in *input) p2Down(d dir) {

	gridSearcher := &position{in.robot.x, in.robot.y}

	if in.grid[in.robot.x+d.y][in.robot.y+d.x] == WALL {
		return
	}

	gridSearcher.x += d.y
	gridSearcher.y += d.x

	objectAbove := in.grid[gridSearcher.x][gridSearcher.y]

	if objectAbove == EMPTY {
		in.grid[in.robot.x][in.robot.y] = EMPTY

		in.robot.x += d.y
		in.robot.y += d.x

		in.grid[in.robot.x][in.robot.y] = ROBOT
		return
	}

	boxTrain, err := in.bfsDown(*gridSearcher)
	if err != nil {
		return
	}

	for _, b := range boxTrain {
		in.grid[b.left.x][b.left.y] = EMPTY
		in.grid[b.right.x][b.right.y] = EMPTY
	}

	in.grid[in.robot.x][in.robot.y] = EMPTY

	for _, b := range boxTrain {
		b.left.x += d.y
		b.left.y += d.x
		b.right.x += d.y
		b.right.y += d.x
	}

	in.robot.x += d.y
	in.robot.y += d.x

	for _, b := range boxTrain {
		in.grid[b.left.x][b.left.y] = BOX_LEFT
		in.grid[b.right.x][b.right.y] = BOX_RIGHT
	}

	in.grid[in.robot.x][in.robot.y] = ROBOT
}

func (b *box) String() string {
	return fmt.Sprintf("left:%+v, right:%+v", b.left, b.right)
}

func (in *input) p2Left(d dir) {

	boxTrain := []*box{}


	gridSearcher := &position{in.robot.x, in.robot.y}

	if in.grid[in.robot.x+d.y][in.robot.y+d.x] == WALL {
		return
	}

	gridSearcher.x += d.y
	gridSearcher.y += d.x

	for in.grid[gridSearcher.x][gridSearcher.y] == BOX_RIGHT {
		b := &box{left: &position{gridSearcher.x + d.y, gridSearcher.y + d.x}, right: &position{gridSearcher.x, gridSearcher.y}}

		boxTrain = append(boxTrain, b)

		gridSearcher.x += d.y * 2
		gridSearcher.y += d.x * 2
	}

	hasASpotToPush := in.grid[gridSearcher.x][gridSearcher.y] == EMPTY

	for _, b := range boxTrain {
		in.grid[b.left.x][b.left.y] = EMPTY
		in.grid[b.right.x][b.right.y] = EMPTY
	}

	in.grid[in.robot.x][in.robot.y] = EMPTY

	if hasASpotToPush {
		for _, b := range boxTrain {
			b.left.x += d.y
			b.left.y += d.x
			b.right.x += d.y
			b.right.y += d.x
		}

		in.robot.x += d.y
		in.robot.y += d.x
	}

	for _, b := range boxTrain {
		in.grid[b.left.x][b.left.y] = BOX_LEFT
		in.grid[b.right.x][b.right.y] = BOX_RIGHT
	}

	in.grid[in.robot.x][in.robot.y] = ROBOT
}

func (in *input) p2Right(d dir) {

	boxTrain := []*box{}


	gridSearcher := &position{in.robot.x, in.robot.y}

	if in.grid[in.robot.x+d.y][in.robot.y+d.x] == WALL {
		return
	}

	gridSearcher.x += d.y
	gridSearcher.y += d.x


	for in.grid[gridSearcher.x][gridSearcher.y] == BOX_LEFT {
		b := &box{left: &position{gridSearcher.x, gridSearcher.y}, right: &position{gridSearcher.x + d.y, gridSearcher.y + d.x}}


		boxTrain = append(boxTrain, b)

		gridSearcher.x += d.y * 2
		gridSearcher.y += d.x * 2
	}

	hasASpotToPush := in.grid[gridSearcher.x][gridSearcher.y] == EMPTY


	for _, b := range boxTrain {
		in.grid[b.left.x][b.left.y] = EMPTY
		in.grid[b.right.x][b.right.y] = EMPTY
	}

	in.grid[in.robot.x][in.robot.y] = EMPTY

	if hasASpotToPush {
		for _, b := range boxTrain {
			b.left.x += d.y
			b.left.y += d.x
			b.right.x += d.y
			b.right.y += d.x
		}

		in.robot.x += d.y
		in.robot.y += d.x
	}

	for _, b := range boxTrain {
		in.grid[b.left.x][b.left.y] = BOX_LEFT
		in.grid[b.right.x][b.right.y] = BOX_RIGHT
	}

	in.grid[in.robot.x][in.robot.y] = ROBOT
}

func (in *input) movep2(b byte) {
	d := DIRS[b]
	switch b {
	case LEFT:
		in.p2Left(d)
	case RIGHT:
		in.p2Right(d)
	case UP:
		in.p2Up(d)
	case DOWN:
		in.p2Down(d)
	}

}


func (in *input) Part2() int {
	total := 0

	for _, m := range in.movements {
		in.movep2(m)
	}

	for row := range in.rows {
		for col := range in.cols {
			if in.grid[row][col] == BOX_LEFT {
				total += (100 * row) + col
			}
		}
	}
	return total
}

func (in *input) Part1() int {
	total := 0

	for _, m := range in.movements {
		in.handleMoveInDir(m)
	}

	for row := range in.rows {
		for col := range in.cols {
			if in.grid[row][col] == BOX_O {
				total += 100*row + col
			}
		}
	}
	return total
}

func main() {
	f := openInputFile()
	in := newInput(f)

	in.expandMap()

	res := in.Part2()

	fmt.Println("P2", res, "inputs", len(in.movements))
}
