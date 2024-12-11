package main

import (
	"bufio"
	"fmt"
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

Part 1:
- Report is considered safe when both of the following are true:
1. the levels are either ALL INCREASING or ALL DECREASING
2. any two adjacent levels differ by ATLEAST ONE and ATMOST THREE

Part 1 bench: BenchmarkReportSafety-12        226556745                5.184 ns/op
------------

Part 2:
- for each report, if the report is unsafe,
- remove 1 level at a time and re test with that new array until we get a successful run
- if we go through every single character and still have no success, it is truly unsafe

Part 2 bench: BenchmarkProcessReport-12       184762862                6.346 ns/op
*/

const (
	UNSAFE = iota
	SAFE
)

const (
	INVALID = iota
	INCREASING
	DECREASING
)

func main() {
	file, err := os.Open("aoc-day2-input.txt")
	if err != nil {
		log.Fatal("File failed to open")
	}

	scanner := bufio.NewScanner(file)

	safeReportCount := 0

	for scanner.Scan() {
		report := extractReportRow(scanner.Text())

		calc := processReport(report)
		safeReportCount += calc
	}
	fmt.Printf("d2p1: %d Reports Safe\n", safeReportCount)
}

func extractReportRow(in string) []int {
	report := strings.Split(in, " ")

	reportIntList := []int{}

	for _, r := range report {
		v, _ := strconv.Atoi(r)
		reportIntList = append(reportIntList, v)
	}

	return reportIntList
}

// handles checking report safety and re running safety check with removal if needed
func processReport(report []int) int {
	if reportSafetyCheck(report) == SAFE {
		// already a safe report
		return SAFE
	}

	// start the process of removing 1 level at a time and re checking report safety
	// when we get a safe, just return it
	for i := range report {
		if reportSafetyCheck(remove(i, report)) == SAFE {
			return SAFE
		}
	}

	// no removal would make this report safe
	return UNSAFE
}

func reportSafetyCheck(report []int) int {
	ptr1 := 0
	ptr2 := 1

	n := len(report)

	reportDirection := INVALID

	for ptr2 < n {

		val1 := report[ptr1]
		val2 := report[ptr2]

		currDirection := direction(val1, val2)

		if reportDirection == INVALID {
			reportDirection = currDirection
		}

		// if we change directions or if we have matching values, unsafe
		if currDirection != reportDirection || currDirection == INVALID {
			return UNSAFE
		}

		// if the diff between the two levels are out of bounds, unsafe
		if !validDifference(val1, val2) {
			return UNSAFE
		}

		// move ptrs
		ptr1++
		ptr2++
	}

	return SAFE
}

func remove(idx int, list []int) []int {
	res := []int{}
	res = append(res, list[:idx]...)
	res = append(res, list[idx+1:]...)

	return res
}

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
