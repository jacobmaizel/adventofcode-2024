package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
)

/*
--------- DAY 3: Mull It Over ---------
- Memory is corrupted
- goal is to multiply numbers
- invalid sequences ie. mul(4*, mul(6,9!, ?(12,34), or mul ( 2 , 4 ) do NOTHING

- answer: scan input for uncorrupted mul instructions, and SUM all of their results.
*/

type FileParser struct {
	input                  string
	currIndex              int
	peekIndex              int
	peekCh                 byte
	ch                     byte
	totalMultiplicationSum int
}

func New(input string) *FileParser {
	fp := &FileParser{input: input}
	fp.Advance()
	return fp
}

func (i *FileParser) Advance() {
	if i.peekIndex >= len(i.input) {
		i.ch = 0
	} else {
		i.ch = i.input[i.peekIndex]
	}

	i.currIndex = i.peekIndex
	i.peekIndex++

	if i.peekIndex >= len(i.input) {
		i.peekCh = 0
	} else {
		i.peekCh = i.input[i.peekIndex]
	}
}

func (fp *FileParser) currCharIs(ch byte) bool {
	return fp.ch == ch
}

func (fp *FileParser) ParseParameters() {
}

func (fp *FileParser) ParseMulStatement() *MulExpression {
	word := fp.readWord()

	// fmt.Printf("read word: %s, peek ch is %s\n", word, string(fp.peekCh))

	if word == "mul" && fp.currCharIs('(') {
		// fmt.Println("Got mul and next is (")
		fp.Advance()

		param1 := fp.readNumber()

		if !fp.currCharIs(',') {
			// fmt.Printf("curr char was NOT a , it was %s\n", string(fp.ch))
			return nil
		}

		fp.Advance()
		param2 := fp.readNumber()

		if !fp.currCharIs(')') {
			// fmt.Printf("curr char was NOT a closing brace.. it was %s\n", string(fp.ch))
			return nil
		}

		p1Int, err := strconv.Atoi(param1)
		// fmt.Printf("param1: %d\n", p1Int)
		if err != nil {
			return nil
		}

		p2Int, err := strconv.Atoi(param2)
		if err != nil {
			return nil
		}

		m := &MulExpression{Left: p1Int, Right: p2Int}
		// fmt.Printf("Created mul exp: %+v\n", m)
		return m
	}

	return nil
}

func (fp *FileParser) ParseInputAndSum() {
	for fp.ch != 0 {

		if fp.ch == 'm' {
			// fmt.Printf("Got m at %d\n", fp.currIndex)
			mulStmt := fp.ParseMulStatement()
			if mulStmt == nil {
				// ended up being invalid, just skip over where it failed and continue.
				fp.Advance()
				continue
			}

			fp.totalMultiplicationSum += mulStmt.Calc()
		}

		fp.Advance()
	}
}

// keeps reading until we hit a non letter
func (fp *FileParser) readWord() string {
	startIdx := fp.currIndex
	for isLetter(fp.ch) {
		fp.Advance()
	}
	return fp.input[startIdx:fp.currIndex]
}

func (fp *FileParser) readNumber() string {
	startIdx := fp.currIndex
	for isDigit(fp.ch) {
		fp.Advance()
	}
	return fp.input[startIdx:fp.currIndex]
}

func (fp *FileParser) Peek() byte {
	return fp.input[fp.peekIndex]
}

// holds operators of a valid multiply statement
type MulExpression struct {
	Left  int
	Right int
}

func (me MulExpression) Calc() int {
	return me.Left * me.Right
}

func main() {
	fmt.Println("D3p1")
	file, _ := os.Open("aoc-day3-input.txt")
	data, _ := io.ReadAll(file)
	fp := New(string(data))

	fp.ParseInputAndSum()

	fmt.Printf("D2P1: Total sum of multiplations: %d", fp.totalMultiplicationSum)
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}
