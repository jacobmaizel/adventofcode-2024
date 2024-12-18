package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

/*
Day 5

part 1:
first section of input
X|Y pairs of rules
X must come before Y in any given update

part 2:
for each of the INCORRECTLY ordered updates,
use the page ordering rules to sort the page numbers in the right order
then sum their midpoints as we did before

p2 notes:
how to sort?
idea 1: move the offending page back 1 by 1 until its behind the before part of the rule?
    then the problem is do we mess up other rules that depend on that?

idea 2:
build some sort of dependency graph using the broke rules to build a plan
13,24,53,64,12,82
if rule 1 and rule 2 (13|82)
*/

type Rule struct {
	before int
	after  int
}

func NewRule(before, after int) Rule {
	n := Rule{before: before, after: after}
	return n
}

type Input struct {
	rulesLookup map[int][]Rule
	printLines  [][]int
}

func NewInput(read io.Reader) *Input {
	inp := &Input{
		rulesLookup: make(map[int][]Rule),
		printLines:  [][]int{},
	}

	scanner := bufio.NewScanner(read)

	scanningRules := true
	for scanner.Scan() {
		line := scanner.Text()

		if len(line) == 0 {
			scanningRules = false
			continue
		}

		if scanningRules {
			// scanning rules
			ruleSet := strings.Split(line, "|")
			before, _ := strconv.Atoi(ruleSet[0])
			after, _ := strconv.Atoi(ruleSet[1])

			r := NewRule(before, after)
			inp.rulesLookup[r.before] = append(inp.rulesLookup[r.before], r)

		} else {
			// scanning print lines
			pages := strings.Split(line, ",")
			l := []int{}
			for _, page := range pages {
				p, _ := strconv.Atoi(page)
				l = append(l, p)
			}
			inp.printLines = append(inp.printLines, l)
		}

	}

	return inp
}

func (in *Input) CheckValidity(printRow []int) ([]Rule, bool) {
	seenPages := make(map[int]bool)

	brokenRules := []Rule{}

	for _, page := range printRow {
		for _, beforeRule := range in.rulesLookup[page] {
			if ok := seenPages[beforeRule.after]; ok {
				brokenRules = append(brokenRules, beforeRule)
			}
		}
		seenPages[page] = true
	}

	if len(brokenRules) > 0 {
		return brokenRules, false
	}

	return nil, true
}

// part 2, sum correctly sorted updates' midpoints
func (in *Input) CalcP2() int {
	total := 0

	for _, printRow := range in.printLines {
		if brokenRules, valid := in.CheckValidity(printRow); !valid {

			// DEBUG
			fmt.Printf("Broken rules for row=%+v\n", printRow)
			for _, r := range brokenRules {
				fmt.Printf("(%d|%d)\n", r.before, r.after)
			}
			fmt.Println("---------------------------")

			sortedRow := in.SortUpdateRow(printRow)

			mid := sortedRow[len(sortedRow)/2]

			total += mid
		}
	}

	return total
}

func (in *Input) SortUpdateRow(printRow []int) []int {
	return printRow
}

func (in *Input) CalcP1() int {
	total := 0

	for _, printRow := range in.printLines {
		if _, valid := in.CheckValidity(printRow); valid {

			mid := printRow[len(printRow)/2]

			total += mid
		}
	}

	return total
}

func main() {
	fmt.Println("AOC - Day 5")
	file, _ := os.Open("aoc-day5-input.txt")

	defer file.Close()

	inp := NewInput(file)

	res := inp.CalcP1()

	fmt.Printf("Day 5 Part 1 ans=%d\n", res)

	res2 := inp.CalcP2()

	fmt.Printf("Day 5 Part 2 ans=%d\n", res2)
}
