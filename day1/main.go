package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
)

func main() {
	fmt.Println("Day 1!")

	file, err := os.Open("input.txt")
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file)

	var l1 []int
	var l2 []int

	for scanner.Scan() {
		line := scanner.Text()
		values := strings.Split(line, "   ")
		if len(values) != 2 {
			log.Fatalf("Wrong len of values??, got %+v", values)
		}
		v1, err := strconv.ParseInt(values[0], 0, 0)
		if err != nil {
			log.Fatal("Failed to convert value from list 1 to int")
		}
		v2, err := strconv.ParseInt(values[1], 0, 0)
		if err != nil {
			log.Fatal("Failed to convert value from list 2 to int")
		}

		l1 = append(l1, int(v1))
		l2 = append(l2, int(v2))
	}

	slices.Sort(l1)
	slices.Sort(l2)

	// fmt.Printf("l1 len: %d, l2 len: %d", len(l1), len(l2))

	totalDistance := 0

	for i, list1Val := range l1 {
		distance := getDistanceBetweenInts(list1Val, l2[i])

		totalDistance += distance
	}

	// part 1 total distance
	fmt.Printf("Day 1 Part 1: Total Distance: %d\n", totalDistance)

	// Part 2: Get total similarity score
	// similarity score is val1 * freq of val1 in list2

	list2Freqs := make(map[int]int)

	for _, l2Val := range l2 {
		list2Freqs[l2Val]++
	}

	similarityScore := 0
	for _, list1Val := range l1 {
		similarityScore += list1Val * list2Freqs[list1Val]
	}

	fmt.Printf("Day 1 Part 2: Similarity Score: %d\n", similarityScore)
}

// total distance between two values, ensuring abs values throughout
func getDistanceBetweenInts(i1, i2 int) int {
	if i1 < 0 {
		i1 = i1 * -1
	}
	if i2 < 0 {
		i2 = i2 * -1
	}

	dist := i1 - i2

	if dist < 0 {
		dist = dist * -1
	}
	return dist
}
