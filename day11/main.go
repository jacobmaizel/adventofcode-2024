package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

var _ = os.Open

type input struct {
	stones []string
}

func newInput(r io.Reader) *input {
	s := bufio.NewScanner(r)
	s.Scan()

	line := strings.TrimSpace(s.Text())

	return &input{strings.Split(line, " ")}
}

// func insert(l []string, val string, idx int) []string {
// 	return append(l[:idx], append([]string{val}, l[idx:]...)...)
// }

// has to be if else because order of rules matter, first applicable must run

func (in *input) processStone(stone string) []string {
	ans := []string{}

	if stone == "0" {
		ans = append(ans, "1")
	} else if len(stone)%2 == 0 {
		// even len numbers split in half, remove leading 0s from right half
		// fmt.Println("splitting stone=", stone)
		rightHalf := stone[len(stone)/2:]
		// fmt.Println("righthalf", rightHalf)
		if len(rightHalf) > 1 {
			ptr := 0
			for rightHalf[ptr] == '0' && ptr < len(rightHalf)-1 {
				ptr++
			}
			rightHalf = rightHalf[ptr:]
		}
		// fmt.Println("after trim 0s", rightHalf)
		ans = append(ans, stone[:len(stone)/2])
		ans = append(ans, rightHalf)

	} else {
		// mult stone val by 2024
		nVal, _ := strconv.Atoi(stone)
		ans = append(ans, strconv.Itoa(nVal*2024))
	}
	// fmt.Println("after adjustment", ans)

	return ans
}

func (in *input) blink() {
	n := []string{}

	for _, stone := range in.stones {
		v := in.processStone(stone)
		n = append(n, v...)
	}

	in.stones = n
}

func (i *input) simulateBlinks(blinks int) {
	for range blinks {
		i.blink()
	}
}

func main() {
	f, _ := os.Open("input.txt")
	i := newInput(f)
	i.simulateBlinks(20)
	fmt.Println("day11p1->", len(i.stones))
}
