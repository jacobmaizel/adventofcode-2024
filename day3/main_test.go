package main

import (
	"fmt"
	"testing"
)

func TestParseAndSum(t *testing.T) {
	inputs := []struct {
		input         string
		expectedTotal int // count of all complete mul statements in our input
	}{
		{"};//how():mul(422,702)'how()'from()-&when(551,888)from()#mul(694,437)", 599522},
		{"xmul(2,4)%&mul[3,7]!@^do_not_mul(5,5)+mul(32,64]then(mul(11,8)mul(8,5))", 161},
	}

	for i, tt := range inputs {

		tName := fmt.Sprintf("test idx %d", i)

		t.Run(tName, func(t *testing.T) {
			fp := New(tt.input)

			fp.ParseInputAndSum()
			res := fp.totalMultiplicationSum

			if res != tt.expectedTotal {
				t.Errorf("wrong summed total, got=%d want=%d", res, tt.expectedTotal)
			}
		})

	}
}

func TestParseMulExpression(t *testing.T) {
	inputs := []struct {
		input         string
		expectedLeft  int
		expectedRight int
		expectNil     bool
	}{
		{"mul(20,81)", 20, 81, false},
		{"mul(250,111)", 250, 111, true},
	}

	for i, tt := range inputs {

		tName := fmt.Sprintf("test idx %d", i)

		t.Run(tName, func(t *testing.T) {
			fp := New(tt.input)

			res := fp.ParseMulStatement()

			if res == nil && !tt.expectNil {
				t.Fatal("Did not expect res to be nil!")
			}

			if res.Left != tt.expectedLeft || res.Right != tt.expectedRight {
				t.Errorf("wrong params, got=%d,%d.... want=%d,%d", res.Left, res.Right, tt.expectedLeft, tt.expectedRight)
			}
		})

	}
}

func TestAdvance(t *testing.T) {
	inputs := []struct {
		input        string
		expectedPeek byte
	}{
		{"asdf123", 'f'},
	}

	for i, tt := range inputs {

		tName := fmt.Sprintf("test idx %d", i)

		t.Run(tName, func(t *testing.T) {
			fp := New(tt.input)
			fp.Advance()
			fp.Advance()

			if fp.peekCh != tt.expectedPeek {
				t.Errorf("peek not right, got=%s, want=%s", string(fp.peekCh), string(tt.expectedPeek))
			}
		})

	}
}

func BenchmarkParseAndSum(b *testing.B) {
	input := "xmul(2,4)%&mul[3,7]!@^do_not_mul(5,5)+mul(32,64]then(mul(11,8)mul(8,5))"
	// BenchmarkParseAndSum-12          6084108               183.9 ns/op
	// BenchmarkParseAndSum-12         559991629                1.900 ns/op           0 B/op without new call

	fp := New(input)
	for i := 0; i < b.N; i++ {
		fp.ParseInputAndSum()
	}
}
