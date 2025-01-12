package main

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
)

/*
Day 14 - Restroom Redoubt
Part 1 ---
Goal: simulate robot security movements for 100 seconds.
Answer: number of robots in each quadrant after 100 seconds
NOTE: robots exactly in the middle horizontally or vertically do not count
as being in any quadrant, so dont count them.

input: robots starting positions and their velocities (tile movements per second)

x values represent number of tiles from LEFT wall
y values represent number of tiles from the RIGHT wall

Notes:
robots can be on top of eachother (on same point)
robots will WRAP around the map if they run into the edge during movement

NOTE: space for part 1 will be 101 wide and 103 tall

will need Robot struct, with movement abilities robot.Move()
robots will need the ability to dteremine where they will end up when wraping around map


*/

func inputFile() *os.File {
	f, _ := os.Open("input.txt")
	return f
}

type robot struct {
	pos *position
	vel velocity
}

type velocity struct {
	dX, dY int
}

type position struct {
	x, y int
}

type input struct {
	robots                []robot
	gridHeight, gridWidth int
}

func newRobot(raw string) robot {
	parts := strings.Split(raw, " ")
	pointRaw := strings.Split(parts[0][2:], ",")
	xR, _ := strconv.Atoi(pointRaw[0])
	yR, _ := strconv.Atoi(pointRaw[1])
	p := position{xR, yR}
	velRaw := strings.Split(parts[1][2:], ",")
	dXR, _ := strconv.Atoi(velRaw[0])
	dYR, _ := strconv.Atoi(velRaw[1])
	v := velocity{dXR, dYR}
	return robot{&p, v}
}

func newInput(r io.Reader, gridHeight, gridWidth int) *input {
	s := bufio.NewScanner(r)
	robots := []robot{}
	for s.Scan() {
		t := strings.TrimSpace(s.Text())
		robots = append(robots, newRobot(t))
	}
	return &input{robots, gridHeight, gridWidth}
}

func (r *robot) simulateMovement(gridHeight, gridWidth, seconds int) {
	newX := (r.pos.x + (r.vel.dX * seconds)) % gridWidth
	newY := (r.pos.y + (r.vel.dY * seconds)) % gridHeight

	if newX < 0 {
		newX = gridWidth + newX
	}

	if newY < 0 {
		newY = gridHeight + newY
	}

	r.pos = &position{newX, newY}
}

func (in *input) quadrant(p position) int {
	hSplit := in.gridHeight / 2
	wSplit := in.gridWidth / 2

	if p.x < wSplit && p.y < hSplit {
		return 0
	} else if p.x < wSplit && p.y > hSplit {
		return 1
	} else if p.x > wSplit && p.y < hSplit {
		return 2
	} else {
		return 3
	}
}

func (in *input) moveRobotsP1(seconds int) int {
	// endingPositionsCounts := make(map[position]int)
	quadrants := make(map[int]int)

	for idx, r := range in.robots {
		// fmt.Println("Robot position before", r.pos)

		r.simulateMovement(in.gridHeight, in.gridWidth, seconds)

		// fmt.Println("Robot position after", r.pos)

		// update the robot in the input
		in.robots[idx] = r

		if r.pos.x == in.gridWidth/2 || r.pos.y == in.gridHeight/2 {
			// fmt.Println("Robot was in the middle hori or vertically, not counting", r.pos)
			continue
		}
		// endingPositionsCounts[*r.pos]++
		quadrants[in.quadrant(*r.pos)]++
	}

	t := 1
	for _, v := range quadrants {
		// fmt.Println(v)
		t *= v
	}
	return t
}

func (in *input) printGrid() {
	posLookup := make(map[position]int)

	// fmt.Println(len(in.robots))

	for _, rob := range in.robots {
		posLookup[*rob.pos]++
	}

	// fmt.Println(len(in.robots), len(posLookup))
	for y := range in.gridHeight {
		for x := range in.gridWidth {
			if v, ok := posLookup[position{x, y}]; ok {
				fmt.Print(v)
			} else {
				fmt.Print(".")
			}
		}
		fmt.Print("\n")
	}
	fmt.Print("\n\n")
}

func (in *input) P2XmasTree() {
	iters := 0
	for !in.foundRow() {
		in.moveRobotsP1(1)
		iters++
	}
	in.printGrid()
	fmt.Println(iters)
}

const minimumContiguous = 10

func (in *input) foundRow() bool {
	posLookup := make(map[position]int)

	for _, rob := range in.robots {
		posLookup[*rob.pos]++
	}

	currMaxContiguous := math.MinInt
	for y := range in.gridHeight {
		rowContig := 0
		for x := range in.gridWidth {
			if _, ok := posLookup[position{x, y}]; ok {
				rowContig++
			} else {
				currMaxContiguous = max(currMaxContiguous, rowContig)
				rowContig = 0
			}
		}
	}

	return currMaxContiguous >= minimumContiguous
}

func main() {
	f := inputFile()
	in := newInput(f, 103, 101)

	// res := in.moveRobotsP1(100)
	in.P2XmasTree()
	// fmt.Println(res)
}
