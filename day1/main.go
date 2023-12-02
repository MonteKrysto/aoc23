package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
)

type Match struct {
	Value string
	Index int
}

func main() {
	// Get args to decide which function to call
	// go run main.go a runs day 1, part 1
	// go run main.go b runs day 1, part 2
	args := os.Args
	aOrb := args[1]

	// Open the file with the calibrations
	file, err := os.Open("calibration.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	sum := 0
	scanner := bufio.NewScanner(file)

	// Read each line of the file and call the apporpriate function based on the arg value passed in
	for scanner.Scan() {
		line := scanner.Text()
		var calibrationValue int
		if aOrb == "a" {
			calibrationValue = getCalibrationValueByNumbersOnly(line)
		} else {
			calibrationValue = getRealCalibrationValue(line)
		}
		sum += calibrationValue
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading from file:", err)
		return
	}

	fmt.Println("Total calibration value:", sum)
}

// Find first and last digits
func getCalibrationValueByNumbersOnly(line string) int {
	// Get all occurences of digits in string
	re := regexp.MustCompile(`\d`)
	digits := re.FindAllString(line, -1)

	return calculateCalibrationValue(digits)
}

// This pain in the ass requirement looks for words that are numbers within strings
func getRealCalibrationValue(line string) int {
	matches := findAllMatches(line)

	// Sort matches by index
	sort.Slice(matches, func(i, j int) bool {
		return matches[i].Index < matches[j].Index
	})

	// Extract the Value field from each Match struct
	values := make([]string, len(matches))
	for i, match := range matches {
		values[i] = match.Value
	}

	return calculateCalibrationValue(values)
}

func findAllMatches(s string) []Match {
	var matches []Match

	// Declare a mapping of words to digits
	digitMap := map[string]string{
		"one": "1", "two": "2", "three": "3", "four": "4",
		"five": "5", "six": "6", "seven": "7", "eight": "8", "nine": "9",
	}

	for word, digit := range digitMap {
		re := regexp.MustCompile(word)
		indexes := re.FindAllStringIndex(s, -1)
		for _, index := range indexes {
			matches = append(matches, Match{Value: digit, Index: index[0]})
		}
	}

	// Also match standalone digits
	re := regexp.MustCompile(`\d`)
	digitIndexes := re.FindAllStringIndex(s, -1)
	for _, index := range digitIndexes {
		matches = append(matches, Match{Value: s[index[0] : index[0]+1], Index: index[0]})
	}

	return matches
}

func calculateCalibrationValue(values []string) int {
	// Get the first and last digits only
	if len(values) >= 2 {
		firstDigit, _ := strconv.Atoi(values[0])
		lastDigit, _ := strconv.Atoi(values[len(values)-1])
		combinedValue, _ := strconv.Atoi(strconv.Itoa(firstDigit) + strconv.Itoa(lastDigit))
		return combinedValue
	}

	// If there is only one digit multiply it by 11 to get two digits.  Ex: 9 will return 99
	if len(values) == 1 {
		singleDigit, _ := strconv.Atoi(values[0])
		return singleDigit * 11
	}

	return 0
}
