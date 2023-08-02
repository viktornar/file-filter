package internal

import "time"

var dateLayouts = []string{
	"2006-01-02",
	"2006/01/02",
	"01/02/2006",
	"Jan 02, 2006",
	"2006, Jun 02",
}

func IsValidDate(dateStr string) bool {
	for _, l := range dateLayouts {
		_, err := time.Parse(l, dateStr)
		if err == nil {
			return true
		}
	}

	return false
}

func GetDateLayouts() []string {
	return dateLayouts
}