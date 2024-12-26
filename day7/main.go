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
Day 7 Part 1
Calc total of all sums whos equations we can complete using the given + and * operators

***NOTE: Operators are always evaluated left-to-right, not according to precedence rules.
*/

const (
	MUL = Operator('*')
	ADD = Operator('+')
)

type Operator byte

func (o Operator) String() string {
	return fmt.Sprintf("%v", string(o))
}

type Statement struct {
	answer   int
	operands []int
}

type Input struct {
	statements []Statement
}

func NewStatement(answer int, ops []int) Statement {
	return Statement{answer: answer, operands: ops}
}

func NewInput(r io.Reader) *Input {
	i := &Input{statements: []Statement{}}
	s := bufio.NewScanner(r)
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		// fmt.Println("LINE", line)
		parts := strings.Split(line, ":")
		// fmt.Println("parts", parts)
		rawAns := parts[0]
		// fmt.Println("Raw ans", rawAns)
		rawOps := strings.Split(parts[1][1:], " ") // skip leading space
		rawIntOps := []int{}
		rawIntAns, _ := strconv.Atoi(rawAns)
		for _, ro := range rawOps {
			s, _ := strconv.Atoi(ro)
			rawIntOps = append(rawIntOps, s)
		}
		ni := NewStatement(rawIntAns, rawIntOps)
		// fmt.Println("NEW STATEMENT", ni)
		i.statements = append(i.statements, ni)
	}

	// fmt.Println("Created statements", i.statements)

	return i
}

// brute force generates list of lists of every permutation of ops for a statement
// num operands- 1 amt of operators needed
func genComb(items []Operator, size int) [][]Operator {
	if size == 0 {
		return [][]Operator{{}}
	}

	ans := [][]Operator{}

	for _, item := range items {
		smallerCombos := genComb(items, size-1)
		for _, combo := range smallerCombos {
			ans = append(ans, append([]Operator{item}, combo...))
		}
	}

	// fmt.Println("Generated combos", ans)

	return ans
}

func operate(oper1 int, op Operator, oper2 int) int {
	// fmt.Printf("%d %s %d", oper1, op, oper2)
	switch op {
	case MUL:
		return oper1 * oper2
	case ADD:
		return oper1 + oper2
	default:
		fmt.Println("Impossible operation")
		return 0
	}
}

// takes a set of ops to use, and returns if we have a match
func (s *Statement) checkCalcAnswerMatch(ops []Operator) bool {
	total := s.operands[0] // first op prefill

	for i := 1; i < len(s.operands); i++ {
		op := ops[i-1]
		total = operate(total, op, s.operands[i])

		// short circuit if we already are over our
		if total > s.answer {
			// fmt.Println("already over answer, short circuit", total)
			continue
		}

	}
	return total == s.answer
}

func (i *Input) RunPart1() int {
	// opCombos := genComb(OPS, len(i.statements)-1)
	total := 0

OUTER:
	for _, statement := range i.statements {
		// get all the combos and check each one, short circuit if we find an answer

		// fmt.Println("Generating combos for statement", statement)
		combos := genComb(OPS, len(statement.operands)-1)

		for _, combo := range combos {
			if ok := statement.checkCalcAnswerMatch(combo); ok {
				total += statement.answer
				continue OUTER
			}
		}

	}

	return total
}

var OPS []Operator = []Operator{MUL, ADD}

func main() {
	f, _ := os.Open("aoc-day7-input.txt")
	i := NewInput(f)

	p1 := i.RunPart1()
	fmt.Printf("Day 7 Part 1 ans=%d\n", p1)
}
