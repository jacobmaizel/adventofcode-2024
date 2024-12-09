package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

// Day 2: Red-Nosed Reports

/*
- each line is a "report"
- each report is a list of #s called "levels" that are seperated by spaces

Goal: Figure out which reports are "safe"

- Report is considered safe when both of the following are true:
1. the levels are either ALL INCREASING or ALL DECREASING
2. any two adjacent levels differ by ATLEAST ONE and ATMOST THREE
*/
func main() {
	file, err := os.Open("aoc-day2-input.txt")
	if err != nil {
		log.Fatal("File failed to open")
	}

	scanner := bufio.NewScanner(file)

	safeReportCount := 0

	for scanner.Scan() {
		report := strings.Split(scanner.Text(), " ")

		reportIntList := []int{}

		for _, r := range report {
			v, _ := strconv.Atoi(r)
			reportIntList = append(reportIntList, v)
		}

		safeReportCount += reportSafetyCheck(reportIntList)
	}
}

// checks
// 1. is the entire row all increasing or all decreasing?
// is the difference between each level > 0 and <=3 ?
// returns 1 if safe, 0 if unsafe
func reportSafetyCheck(report []int) int {
	ptr1 := 0
	ptr2 := 1

	n := len(report)

	reportDirection := direction(report[ptr1], report[ptr2])
	if reportDirection == INVALID {
		return 0
	}

	if !validDifference(report[ptr1], report[ptr2]) {
		return 0
	}

	ptr1++
	ptr2++

	for ptr2 < n {

		val1 := report[ptr1]
		val2 := report[ptr2]

		currDirection := direction(val1, val2)

		// if we change directions or if we have matching values, unsafe
		if currDirection != reportDirection || currDirection == INVALID {
			return 0
		}

		// if the diff between the two levels are out of bounds, unsafe
		if !validDifference(val1, val2) {
			return 0
		}

		// move ptrs
		ptr1++
		ptr2++
	}

	return 1
}

const (
	INVALID = iota
	INCREASING
	DECREASING
)

func direction(i1, i2 int) int {
	if i1 > i2 {
		return DECREASING
	} else if i2 > i1 {
		return INCREASING
	} else {
		// cant be the same!
		return INVALID
	}
}

func validDifference(i1, i2 int) bool {
	diff := i1 - i2
	if diff < 0 {
		diff = diff * -1
	}
	return diff >= 1 && diff <= 3
}
