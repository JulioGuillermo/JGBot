package cron

import "regexp"

// Key are case insensitive, and can be 3 letters or full name
var months = map[string]string{
	"(?i)jan[\\w]*": "1",
	"(?i)feb[\\w]*": "2",
	"(?i)mar[\\w]*": "3",
	"(?i)apr[\\w]*": "4",
	"(?i)may[\\w]*": "5",
	"(?i)jun[\\w]*": "6",
	"(?i)jul[\\w]*": "7",
	"(?i)aug[\\w]*": "8",
	"(?i)sep[\\w]*": "9",
	"(?i)oct[\\w]*": "10",
	"(?i)nov[\\w]*": "11",
	"(?i)dec[\\w]*": "12",
}

func ReplaceMonth(input string) string {
	for k, v := range months {
		re := regexp.MustCompile(k)
		input = re.ReplaceAllString(input, v)
	}
	return input
}
