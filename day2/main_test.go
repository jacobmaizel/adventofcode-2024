package main

import "testing"

func TestReportSafetyHelper(t *testing.T) {
	input := []struct {
		report         []int
		expectedSafety int
	}{
		// {
		// 	[]int{1, 2, 5, 6, 8},
		// 	1,
		// },
		// {
		// 	[]int{10, 7, 5, 4, 3, 2, 1},
		// 	1,
		// },
		// {
		// 	[]int{5, 2, 3, 4, 8},
		// 	0,
		// },

		// {
		// 	[]int{5, 4, 3, 2, 1},
		// 	1,
		// },
		// {
		// 	[]int{4, 4, 0, 4, 7, 0, 5, 0, 0, 5, 1, 0, 5, 3, 0, 5, 4, 0, 5, 3},
		// 	0,
		// },
		{
			[]int{7, 6, 4, 2, 1},
			1,
		},

		{
			[]int{1, 2, 7, 8, 9},
			0,
		},
		{
			[]int{9, 7, 6, 2, 1},
			0,
		},
		{
			[]int{1, 3, 2, 4, 5},
			0,
		},
		{
			[]int{8, 6, 4, 4, 1},
			0,
		},
		{
			[]int{1, 3, 6, 7, 9},
			1,
		},
	}

	for _, tt := range input {

		val := reportSafetyCheck(tt.report)

		if val != tt.expectedSafety {
			t.Fatalf("report safety incorrect, got=%d, want=%d\n", val, tt.expectedSafety)
		}

	}
}

func TestDifference(t *testing.T) {
	tests := []struct {
		val1     int
		val2     int
		expected bool
	}{
		{1, 3, true},
		{10, 5, false},
	}

	for _, tt := range tests {
		diff := validDifference(tt.val1, tt.val2)
		if tt.expected != diff {
			t.Fatalf("wrong diff, got=%t, want=%t for test: %d,%d", diff, tt.expected, tt.val1, tt.val2)
		}
	}
}
