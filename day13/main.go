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

const (
	A_COST             = 3 // per button click
	B_COST             = 1
	MAX_BUTTON_PRESSES = 100 // cant press more than 100 for each button
)

func inputFile() *os.File {
	f, _ := os.Open("input.txt")
	return f
}

type input struct {
	games []*clawGame
}

const (
	PARSE_BUTTON_A = iota
	PARSE_BUTTON_B
	PARSE_PRIZE
)

func newInput(r io.Reader) *input {
	s := bufio.NewScanner(r)

	games := []*clawGame{}

	currGameLines := []string{}
	for {
		more := s.Scan()

		// fmt.Println(s.Err())
		line := strings.TrimSpace(s.Text())
		// fmt.Println("processing line", line)
		if line == "" || !more {
			// fmt.Println("Skipping empty line------")
			games = append(games, newClawGame(currGameLines))
			currGameLines = []string{}
			if !more {
				break
			}
			continue
		}
		currGameLines = append(currGameLines, line)

	}
	// for _, g := range games {
	// 	fmt.Printf("claw game: %+v\n", g)
	// }

	return &input{games}
}

type buttonStats struct {
	xMoves, yMoves int
}

func (b *buttonStats) String() string {
	return fmt.Sprintf("X: %d, Y: %d", b.xMoves, b.yMoves)
}

type position struct {
	x, y int
}

type clawGame struct {
	currentPosition position
	aButton         *buttonStats
	aPresses        int
	aCost           int
	bButton         *buttonStats
	bPresses        int
	bCost           int
	prize           position
}

func newClawGame(lines []string) *clawGame {
	// fmt.Println("new claw game with", lines)
	aRaw := lines[0][12:]
	aXVal, _ := strconv.Atoi(strings.Split(aRaw, ",")[0])
	aYValR := strings.Split(aRaw, "+")
	aYVal, _ := strconv.Atoi(aYValR[len(aYValR)-1])

	bRaw := lines[1][12:]
	bXVal, _ := strconv.Atoi(strings.Split(bRaw, ",")[0])
	bYValR := strings.Split(bRaw, "+")
	bYVal, _ := strconv.Atoi(bYValR[len(bYValR)-1])

	prizeRaw := lines[2][9:]
	prX, _ := strconv.Atoi(strings.Split(prizeRaw, ",")[0])
	pry := strings.Split(prizeRaw, "Y=")
	pryv, _ := strconv.Atoi(pry[len(pry)-1])

	abut := &buttonStats{xMoves: aXVal, yMoves: aYVal}
	bbut := &buttonStats{xMoves: bXVal, yMoves: bYVal}
	prize := position{x: prX, y: pryv}

	return &clawGame{aButton: abut, bButton: bbut, prize: prize}
}

func (c *clawGame) String() string {
	return fmt.Sprintf("Prize: %+v\nAButton: %s, Abutton presses: %d, Abutton cost: %d\nBButton: %s, Bbutton presses: %d, Bbutton cost: %d\n", c.prize, c.aButton, c.aPresses, c.aCost, c.bButton, c.bPresses, c.bCost)
}

/*
Problem:
Calculate the fewest tokens needed to win all possible prizes from the input list of games.
(NOTE: some might be unwinable!)

Some notes:
A button cost per press: 3 tokens
B Button ^: 1 token

cannot press a given button > 100 times per game in order to count the win.
it is possible that there is no valid combination that will allow us to win

some possible strategies:
prioritize smaller token cost button as much as possible

press that as many times as we can (up to 100)
check where we are at
if we are under the amount we are looking for (for both x and y)
  keep adding
*/

func (c *clawGame) checkEarlyInvalid() bool {
	maxX := c.aButton.xMoves*100 + c.bButton.xMoves*100
	maxY := c.aButton.yMoves*100 + c.bButton.yMoves*100
	if c.prize.x > maxX || c.prize.y > maxY {
		// fmt.Printf("game %+v is invalid!\n\n", c.prize)
		return true
	}
	return false
}

func (c *clawGame) pressButtons() {
	determinate := (c.aButton.xMoves * c.bButton.yMoves) - (c.aButton.yMoves * c.bButton.xMoves)
	if determinate == 0 {
		return
	}

	abuttonPresses := (float64(c.bButton.yMoves*c.prize.x) - float64(c.bButton.xMoves*c.prize.y)) / float64(determinate)
	bbuttonPresses := (float64(c.aButton.xMoves*c.prize.y) - float64(c.aButton.yMoves*c.prize.x)) / float64(determinate)

	if abuttonPresses < 0 || bbuttonPresses < 0 || math.Mod(abuttonPresses, 1) != 0 || math.Mod(bbuttonPresses, 1) != 0 {
		return
	}

	c.aPresses = int(abuttonPresses)
	c.bPresses = int(bbuttonPresses)
}

func (in *input) calcTokensP2() int {
	total := 0

	for _, c := range in.games {
		c.prize.x = c.prize.x + 10000000000000
		c.prize.y = c.prize.y + 10000000000000
		c.pressButtons()
		total += c.aPresses*A_COST + c.bPresses*B_COST
	}

	return total
}

// used as loop control, exits if win or invalid
func (c *clawGame) play(x, y, aPressLeft, bPressLeft, cost int, memo map[string]int) int {
	// check success at this point
	if x == c.prize.x && y == c.prize.y {
		return cost
	}

	// check sentinel / failure mode -> no more presses left, and we are not at goal!
	if aPressLeft == 0 && bPressLeft == 0 {
		return math.MaxInt
	}

	// check if our current state is in the cache, if so use that value
	key := fmt.Sprintf("%d,%d,%d,%d", x, y, aPressLeft, bPressLeft)
	if val, found := memo[key]; found {
		return val
	}

	// current state not seen before, calculate fresh
	pressACost := math.MaxInt
	pressBCost := math.MaxInt

	if aPressLeft > 0 {
		pressACost = c.play(x+c.aButton.xMoves, y+c.aButton.yMoves, aPressLeft-1, bPressLeft, cost+A_COST, memo)
	}

	if bPressLeft > 0 {
		pressBCost = c.play(x+c.bButton.xMoves, y+c.bButton.yMoves, aPressLeft, bPressLeft-1, cost+B_COST, memo)
	}

	res := min(pressACost, pressBCost)
	memo[key] = res
	return res
}

func main() {
	r := inputFile()
	in := newInput(r)
	res := in.calcTokensP2()

	fmt.Println(res)
}
