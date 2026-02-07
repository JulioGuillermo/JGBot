package cron

import "regexp"

// Key are case insensitive, and can be 3 letters or full name
var weekdays = map[string]string{
	"(?i)sun[\\w]*": "0",
	"(?i)mon[\\w]*": "1",
	"(?i)tue[\\w]*": "2",
	"(?i)wed[\\w]*": "3",
	"(?i)thu[\\w]*": "4",
	"(?i)fri[\\w]*": "5",
	"(?i)sat[\\w]*": "6",
}

func ReplaceWeekday(input string) string {
	for k, v := range weekdays {
		re := regexp.MustCompile(k)
		input = re.ReplaceAllString(input, v)
	}
	return input
}
