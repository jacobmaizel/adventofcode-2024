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

var _ = os.Open

type input struct {
	stones []int
	cache  map[stone]int
}

func newInput(r io.Reader) *input {
	s := bufio.NewScanner(r)
	s.Scan()

	line := strings.TrimSpace(s.Text())

	strItems := strings.Split(line, " ")

	st := []int{}

	for i := range strItems {
		v, _ := strconv.Atoi(strItems[i])
		st = append(st, v)
	}

	return &input{stones: st, cache: make(map[stone]int)}
}

func splitNum(num int) (int, int) {
	dig := len(fmt.Sprintf("%d", num))
	halfDig := dig / 2

	div := int(math.Pow10(halfDig))
	left := num / div
	right := num % div

	return left, right
}

type stone struct {
	blink int
	value int
}

func newStone(blink int, value int) stone {
	return stone{blink, value}
}

func (in *input) processStone(stone int, blink int) int {
	val := 0
	ns := newStone(blink, stone)
	if blink == 0 {
		return 1
	} else if v, ok := in.cache[ns]; ok {
		return v
	} else if stone == 0 {
		val = in.processStone(1, blink-1)
	} else if len(fmt.Sprintf("%d", stone))%2 == 0 {
		// even len numbers split in half, remove leading 0s from right half
		left, right := splitNum(stone)
		val = in.processStone(left, blink-1) + in.processStone(right, blink-1)
	} else {
		val = in.processStone(stone*2024, blink-1)
	}

	in.cache[ns] = val
	return val
}

func (i *input) simulateBlinks(blinks int) int {
	total := 0
	for _, stone := range i.stones {
		total += i.processStone(stone, blinks)
	}

	return total
}

func main() {
	f, _ := os.Open("input.txt")
	defer f.Close()
	i := newInput(f)
	res := i.simulateBlinks(75)
	fmt.Println("day11p2->", res)
}
