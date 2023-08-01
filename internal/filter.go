package internal

import (
	"fmt"
	"regexp"
)

func PrintFilteredLine(line string, filter *Filter) {
	dateMatch, _ := regexp.MatchString(filter.Date, line)
	nameMatch, _ := regexp.MatchString(filter.Name, line)

	if nameMatch || dateMatch {
		fmt.Println(line)
	}
}
