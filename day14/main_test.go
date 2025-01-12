package main

import (
	"strings"
	"testing"
)

func Test_main(t *testing.T) {
	tests := []struct {
		name                                            string // description of this test case
		input                                           string
		gridHeight, gridWidth, seconds, expSafetyFactor int
	}{
		// TODO: Add test cases.
		{
			name: "web",

			input: `p=0,4 v=3,-3
              p=6,3 v=-1,-3
              p=10,3 v=-1,2
              p=2,0 v=2,-1
              p=0,0 v=1,3
              p=3,0 v=-2,-2
              p=7,6 v=-1,-3
              p=3,0 v=-1,-2
              p=9,3 v=2,3
              p=7,3 v=-1,2
              p=2,4 v=2,-3
              p=9,5 v=-3,-3`,
			gridHeight: 7, gridWidth: 11, seconds: 100, expSafetyFactor: 12,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := strings.NewReader(tt.input)
			in := newInput(f, tt.gridHeight, tt.gridWidth)

			// in.printGrid()

			res := in.moveRobotsP1(tt.seconds)

			if res != tt.expSafetyFactor {
				t.Errorf("got %v, want %v", res, tt.expSafetyFactor)
			}

			// in.printGrid()
		})
	}
}
