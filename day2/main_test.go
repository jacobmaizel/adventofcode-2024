package main

import (
	"fmt"
	"reflect"
	"testing"
)

func TestReportSafetyHelper(t *testing.T) {
	input := []struct {
		report         []int
		expectedSafety int
	}{
		{
			[]int{69, 70, 71, 70, 71},
			0,
		},

		// from aoc website
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
			1,
		},
		{
			[]int{8, 6, 4, 4, 1},
			1,
		},
		{
			[]int{1, 3, 6, 7, 9},
			1,
		},
	}

	for _, tt := range input {
		tname := fmt.Sprintf("%+v", tt.report)

		t.Run(tname, func(t *testing.T) {
			val := processReport(tt.report)

			if val != tt.expectedSafety {
				t.Fatalf("report safety incorrect, got=%d, want=%d\n", val, tt.expectedSafety)
			}
		})

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

func TestRemoveByIndex(t *testing.T) {
	tests := []struct {
		idx          int
		list         []int
		expectedList []int
	}{
		{
			3,
			[]int{1, 2, 3, 4, 5},
			[]int{1, 2, 3, 5},
		},
		{
			7,
			[]int{38, 33, 30, 27, 25, 23, 22, 21},
			[]int{38, 33, 30, 27, 25, 23, 22},
		},
	}

	for _, tt := range tests {

		tname := fmt.Sprintf("idx:%d,list:%+v", tt.idx, tt.list)
		t.Run(tname, func(t *testing.T) {
			res := remove(tt.idx, tt.list)

			if !reflect.DeepEqual(res, tt.expectedList) {
				t.Fatalf("removal of %d gave wrong list: want: %+v, got: %+v\n", tt.idx, tt.expectedList, res)
			}
		})
	}
}

func BenchmarkProcessReport(b *testing.B) {
	list := []int{1, 3, 6, 7, 9}
	for i := 0; i < b.N; i++ {
		processReport(list)
	}
}

func BenchmarkReportSafety(b *testing.B) {
	list := []int{1, 3, 6, 7, 9}
	for i := 0; i < b.N; i++ {
		reportSafetyCheck(list)
	}
}
