package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"slices"
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

func NewInputFromParts(lines [][]int, rules []Rule) *Input {
	inp := &Input{
		printLines:  lines,
		rulesLookup: make(map[int][]Rule),
	}

	for _, r := range rules {
		inp.rulesLookup[r.before] = append(inp.rulesLookup[r.before], r)
	}

	return inp
}

func NewInputFromReader(read io.Reader) *Input {
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

// seperate validity check that does not have to worry about returning the rules
func (in *Input) rowValid(printRow []int) bool {
	seenPages := make(map[int]bool)
	for _, page := range printRow {
		for _, rule := range in.rulesLookup[page] {
			if ok := seenPages[rule.after]; ok {
				return false
			}
		}
		seenPages[page] = true
	}
	return true
}

func (in *Input) getAllRulesForRow(printRow []int) []Rule {
	allRules := []Rule{}
	for _, page := range printRow {
		for _, rule := range in.rulesLookup[page] {
			if slices.Contains(printRow, rule.before) && slices.Contains(printRow, rule.after) {
				allRules = append(allRules, rule)
			}
		}
	}

	return allRules
}

// part 2, sum correctly sorted updates' midpoints
func (in *Input) CalcP2() int {
	total := 0

	for i, printRow := range in.printLines {
		if valid := in.rowValid(printRow); !valid {

			allRules := in.getAllRulesForRow(printRow)

			res := in.topoSort(allRules)

			if !slices.Equal(res, in.printLines[i]) {
				mid := res[len(res)/2]
				total += mid
			}
		}
	}

	return total
}

func buildGraphWithIndegreeCount(rules []Rule) (map[int][]int, map[int]int) {
	graph := make(map[int][]int)
	inDegreeCount := make(map[int]int)

	for _, r := range rules {
		graph[r.before] = append(graph[r.before], r.after)
		inDegreeCount[r.after]++

		if _, exists := inDegreeCount[r.before]; !exists {
			inDegreeCount[r.before] = 0
		}
	}

	return graph, inDegreeCount
}

// kahns alg to topologically sort our DAG
func (in *Input) topoSort(rules []Rule) []int {
	depGraph, inDegree := buildGraphWithIndegreeCount(rules)

	result := []int{}
	queue := []int{}

	// prefill the queue with nodes that have no incoming edges
	for node, degree := range inDegree {
		if degree == 0 {
			queue = append(queue, node)
		}
	}

	// start
	for len(queue) > 0 {
		currNode := queue[0]
		queue = queue[1:]

		result = append(result, currNode)

		for _, neighbor := range depGraph[currNode] {
			inDegree[neighbor]--
			if inDegree[neighbor] == 0 {
				queue = append(queue, neighbor)
			}

		}
	}

	return result
}

func (in *Input) CalcP1() int {
	total := 0
	for _, printRow := range in.printLines {
		if valid := in.rowValid(printRow); valid {
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

	inp := NewInputFromReader(file)

	res := inp.CalcP1()

	fmt.Printf("Day 5 Part 1 ans=%d\n", res)

	res2 := inp.CalcP2()

	fmt.Printf("Day 5 Part 2 ans=%d\n", res2)
}
