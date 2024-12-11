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

const (
	DISABLED = iota
	ENABLED
)

type FileParser struct {
	input                  string
	currIndex              int
	peekIndex              int
	peekCh                 byte
	ch                     byte
	totalMultiplicationSum int
	mode                   int // enabled / disabled. enabled by default
}

func New(input string) *FileParser {
	fp := &FileParser{input: input, mode: ENABLED}
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

func (fp *FileParser) peekCharIs(ch byte) bool {
	return fp.peekCh == ch
}

func (fp *FileParser) currCharIs(ch byte) bool {
	return fp.ch == ch
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

		switch fp.ch {
		case 'm':
			// fmt.Printf("Got m at %d\n", fp.currIndex)
			mulStmt := fp.ParseMulStatement()
			if mulStmt == nil {
				// ended up being invalid, just skip over where it failed and continue.
				fp.Advance()
				continue
			}

			if fp.mode == ENABLED {
				// only saving the calc if we are currently enabled.
				fp.totalMultiplicationSum += mulStmt.Calc()
			}

		case 'd':
			// fmt.Printf("Parsing mod statement at %d\n", fp.currIndex)
			fp.ParseModifierStatement()
		}

		fp.Advance()
	}
}

// do() or don't(), enables / disables the parser from saving the calc of following muls
// assumes it is called when currChar is 'd'
func (fp *FileParser) ParseModifierStatement() {
	if !fp.peekCharIs('o') {
		return
	}

	word := fp.readWord()
	// fmt.Printf("Mod word: %s   curr token is: %s\n", word, string(fp.ch))

	switch word {
	case "do":
		if !fp.currCharIs('(') && fp.peekCharIs(')') {
			// fmt.Println("do did not have () after")
			return
		}
		// fmt.Printf("Flipping to enabled\n")
		fp.mode = ENABLED

	case "don":
		if fp.currCharIs('\'') {
			if fp.peekCharIs('t') {
				fp.Advance() // on 't'
				fp.Advance() // now should be on opening paren
				if fp.currCharIs('(') && fp.peekCharIs(')') {
					fp.Advance() // on closing paren
					fp.Advance() // now on item AFTER the closing paren

					// fmt.Printf("Flipping to disabled, withcurrent token now being: %s\n", string(fp.ch))
					fp.mode = DISABLED

				}
			}
		} else {
			// fmt.Printf("something went wrong in dont parsing, ending at char %s\n", string(fp.ch))
		}

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

	fmt.Printf("AOC D2: Total sum of multiplations: %d", fp.totalMultiplicationSum)
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}
