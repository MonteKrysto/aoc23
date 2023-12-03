package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var defaultLimits = map[string]int{
	"red":   12,
	"green": 13,
	"blue":  14,
}

func SetDefaultLimits(newLimits map[string]int) {
	for color, value := range newLimits {
		if _, exists := defaultLimits[color]; exists {
			defaultLimits[color] = value
		}
	}
}

func main() {
	// Get args to decide which function to call
	// go run main.go a runs day 1, part 1
	// go run main.go b runs day 1, part 2
	args := os.Args
	aOrb := args[1]

	// Open the file with the calibrations
	file, err := os.Open("games.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	sum := 0
	for scanner.Scan() {
		line := scanner.Text()
		if aOrb == "a" {
			gameNumber, err := getPartsOfLine(line)

			if err != nil {
				fmt.Println("Error processing line:", err)
				continue
			}

			sum += gameNumber
		} else {
			power, err := getPartsOfLineB(line)
			if err != nil {
				fmt.Println("Error processing line:", err)
				continue
			}
			sum += power
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading from file:", err)
		return
	}

	fmt.Println("Total sum:", sum)
}

// Strip each line into parts

func getPartsOfLine(line string) (int, error) {
	// Check for Game number
	re := regexp.MustCompile(`Game (\d+):`)
	matches := re.FindStringSubmatch(line)
	if len(matches) < 2 {
		return 0, fmt.Errorf("no game number found in line: %s", line)
	}
	gameNumber, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, fmt.Errorf("invalid game number in line: %s", line)
	}

	// Break down line into separate games
	parts := strings.Split(line, ";")
	for i, part := range parts {
		group := strings.TrimSpace(part)
		if i == 0 {
			group = strings.TrimSpace(strings.SplitN(group, ":", 2)[1])
		}

		if err := parseGroupToCubeMatrix(group); err != nil {
			return 0, nil // should return err but i don't want to print out the error messages
		}
	}

	return gameNumber, nil
}

// Parse each game part

func parseGroupToCubeMatrix(group string) error {
	// Get each part of the game - 4 Blue, 3 Green, etc...
	items := strings.Split(group, ",")

	for _, item := range items {
		parts := strings.Fields(strings.TrimSpace(item))
		if len(parts) != 2 {
			return fmt.Errorf("invalid item format: '%s' in group '%s'", item, group)
		}

		// We know the first part is always the number so we can just use parts[0]
		// otherwise we would have to use a regex to grab the digit from the string
		count, err := strconv.Atoi(parts[0])
		if err != nil {
			return fmt.Errorf("invalid number '%s' in item: '%s'", parts[0], item)
		}

		// The color is always the second part so we can use parts[1]
		color := parts[1]

		// Check the color limit and color against the defaults we set in SetDefaultLimits
		limit, exists := defaultLimits[color]
		if !exists {
			return fmt.Errorf("invalid color '%s'", color)
		}
		if count > limit {
			return fmt.Errorf("count for color '%s' exceeds limit", color)
		}
	}
	return nil
}

func getPartsOfLineB(line string) (int, error) {
	re := regexp.MustCompile(`Game (\d+):`)
	matches := re.FindStringSubmatch(line)
	if len(matches) < 2 {
		return 0, fmt.Errorf("no game number found in line: %s", line)
	}

	parts := strings.Split(line, ";")
	minCubes := map[string]int{"red": 0, "green": 0, "blue": 0}

	for i, part := range parts {
		group := strings.TrimSpace(part)
		if i == 0 {
			group = strings.TrimSpace(strings.SplitN(group, ":", 2)[1])
		}

		if err := parseGroupToMinCubes(group, minCubes); err != nil {
			return 0, err
		}
	}

	// Calculate the power of the set of cubes
	power := 1
	for _, count := range minCubes {
		power *= count
	}

	return power, nil
}

func parseGroupToMinCubes(group string, minCubes map[string]int) error {
	items := strings.Split(group, ",")

	for _, item := range items {
		parts := strings.Fields(strings.TrimSpace(item))
		if len(parts) != 2 {
			return fmt.Errorf("invalid item format: '%s' in group '%s'", item, group)
		}
		count, err := strconv.Atoi(parts[0])
		if err != nil {
			return fmt.Errorf("invalid number '%s' in item: '%s'", parts[0], item)
		}
		color := parts[1]
		if count > minCubes[color] {
			minCubes[color] = count
		}
	}
	return nil
}
