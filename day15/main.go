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

data needed
input
position
robot position
box positions
movements
empty space

notes:
- while the robot moves, if there are boxes in the way, it will attempt to push the boxes
BUT if the movement would cause the robot or box to move into a wall, nothing moves
INCLUDING THE ROBOT!!

input is seperated by a \n
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
	ROBOT = byte('@')
	BOX   = byte('O')
	WALL  = byte('#')
	EMPTY = byte('.')
)

type input struct {
	grid       [][]byte
	robot      *position
	movements  []byte
	rows, cols int
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
			// scanning grid
			b := []byte(strings.TrimSpace(t))

			for idx, b := range b {
				if b == ROBOT {
					// fmt.Println("found robot", len(in.grid), idx)
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

/*
1. distance to a wall from our current point
2. find out how many rocks are in our way, and hold onto their positions.
3. reset their original positions to empty space
4. set reset our robots original position to empty space
5. re-place the boxes we removed up against the wall position (might not be the border of the mpa!)
6. replace our robot at the end of those boxes, wherever they stop.


  ////////// the robot does not push all the boxes as far as possible to the wall.
  // if there is a robot, then box, then a gap, then another box, the robot only pushes
  // the box in front of it 1 box.

  // ie. the robot only pushes continuous rows of boxes.
  // if the robot is directly in front of a row of 2 boxes, it will push them together
  // until it hits another box or a row.
*/

func (in *input) handleMoveInDir(dirByte byte) {
	direc := DIRS[dirByte]
	// set original robot position to empty,

	// find contiguous train of robot and boxes right next to him

	// after the end of the train, how many empty spaces are there until a wall or box?
	// shift each part of the train over by that many

	// our robot starting position is known
	// start by adding our robot starting position
	// keep searching in dir until we hit a wall
	// at each point, if we hit a box, append its position into our train array
	//
	trainPositions := []*position{in.robot}

	gridSearcher := &position{in.robot.x, in.robot.y}

	if in.grid[in.robot.x+direc.y][in.robot.y+direc.x] == WALL {
		// fmt.Println("wall in front, no moves, early return")
		return
	}

	// BUILD TRAIN

	// fmt.Println("Looking for box positions. grid searcher start", gridSearcher.x, gridSearcher.y)

	gridSearcher.x += direc.y
	gridSearcher.y += direc.x
	for in.grid[gridSearcher.x][gridSearcher.y] == BOX {
		// append a copy the current position into our train positions

		trainPositions = append(trainPositions, &position{gridSearcher.x, gridSearcher.y})

		gridSearcher.x += direc.y
		gridSearcher.y += direc.x
	}

	// fmt.Println("AFTER Looking for box positions. grid searcher pos", gridSearcher.x, gridSearcher.y)

	trainMoveCounter := 0

	// gridSearcher.x += direc.y
	// gridSearcher.y += direc.x
	for in.grid[gridSearcher.x][gridSearcher.y] == EMPTY {
		// fmt.Println("found empty space")
		trainMoveCounter++
		gridSearcher.x += direc.y
		gridSearcher.y += direc.x
	}

	// fmt.Println("Done finding empty space, found:", trainMoveCounter)

	// fmt.Println("TRAIN BEFORE MOVING")
	// for _, t := range trainPositions {
	// fmt.Printf("Train item:%+v\n", *t)
	// }

	for _, t := range trainPositions {
		in.grid[t.x][t.y] = EMPTY
	}

	if trainMoveCounter > 0 {
		// fmt.Println("MOVING TRAIN!!!!!!!!!!!")
		// for each of the positions in our train
		// reset the current positions value to empty on the map
		// then incr the position one by one by dir

		for _, trainCar := range trainPositions {
			trainCar.x += direc.y
			trainCar.y += direc.x
		}
	}
	// fmt.Println("TRAIN after MOVING")
	// for _, t := range trainPositions {
	// 	fmt.Printf("Train item:%+v\n", *t)
	// }

	// right at the end we can save the robots position back to our input (first pos in train)
	// set the robot icon on the map at its location
	// and set everything after the robot on the map to their respective box icons.
	afterMoveRobotPosition := trainPositions[0]

	if len(trainPositions) > 0 {
		boxPositions := trainPositions[1:]
		for _, b := range boxPositions {
			in.grid[b.x][b.y] = BOX
		}
	}

	in.grid[in.robot.x][in.robot.y] = EMPTY
	in.grid[afterMoveRobotPosition.x][afterMoveRobotPosition.y] = ROBOT

	in.robot.x = afterMoveRobotPosition.x
	in.robot.y = afterMoveRobotPosition.y
}

func (in *input) move(dirByte byte) {
	// d := DIRS[dirByte]

	// fmt.Printf("\n----BEFORE: moving in direction %v (%v)\n", d, string(dirByte))

	// in.printGrid()
	in.handleMoveInDir(dirByte)
	// in.printGrid()

	// fmt.Printf("\n----AFTER: moving in direction %v (%v)\n", d, string(dirByte))
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

func (in *input) Part1() int {
	total := 0

	for _, m := range in.movements {
		in.move(m)
	}

	for row := range in.rows {
		for col := range in.cols {
			if in.grid[row][col] == BOX {
				total += 100*row + col
			}
		}
	}
	return total
}

func main() {
	f := openInputFile()
	in := newInput(f)

	res := in.Part1()
	// in.printGrid()
	fmt.Println("P1", res)
}
