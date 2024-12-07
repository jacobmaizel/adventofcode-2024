package main

import "testing"

func TestAbsValHelper(t *testing.T) {
	tests := []struct {
		i1               int
		i2               int
		expectedDistance int
	}{
		{1, 2, 1},
		{5, 10, 5},
		{-100, 50, 50},
	}

	for _, tt := range tests {

		calcDistance := getDistanceBetweenInts(tt.i1, tt.i2)

		if calcDistance != tt.expectedDistance {
			t.Fatalf("Distance incorrect: got %d, wanted %d", calcDistance, tt.expectedDistance)
		}

	}
}
