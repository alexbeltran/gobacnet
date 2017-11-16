package property

import (
	"fmt"
	"strings"
)

const numOfAdditionalSpaces = 15

func longestString(x map[string]uint32) int {
	max := 0
	for k, _ := range x {
		if len(k) > max {
			max = len(k)
		}
	}
	return max
}

func printRow(col1, col2 string, maxLen int) {
	spacing := strings.Repeat(" ", maxLen-len(col1)+numOfAdditionalSpaces)
	fmt.Printf("%s%s%s\n", col1, spacing, col2)
}

// PrintAll prints all of the properties within this package. This is only a
// subset of all properties.
func PrintAll() {
	max := longestString(enumMapping)

	printRow("Key", "Int", max)
	fmt.Println(strings.Repeat("-", max+numOfAdditionalSpaces+6))

	for k, id := range enumMapping {
		// Spacing
		printRow(k, fmt.Sprintf("%d", id), max)
	}
}
