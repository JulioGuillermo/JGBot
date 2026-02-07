package cron

import "fmt"

// CronArgs is a struct that holds the arguments for a cron expression
// Minute: must be in range 0-59
//
//	(can be a number for an specific minute, a list of numbers separated by commas for multiple minutes, or every(n) for every n minutes)
//
// Hour: must be in range 0-23
//
//	(can be a number for an specific hour, a list of numbers separated by commas for multiple hours, or every(n) for every n hours)
//
// Day: must be in range 1-31
//
//	(can be a number for an specific day, a list of numbers separated by commas for multiple days, or every(n) for every n days)
//
// Month: must be in range 1-12
//
//	(can be a number for an specific month, a list of numbers separated by commas for multiple months, or every(n) for every n months)
//	Being 1 january, 2 february, etc.
//	Can also be a name of the month (january, february, etc.), in case insensitive and accepting 3 first letters or full name
//
// Weekday: must be in range 0-6
//
//	(can be a number for an specific weekday, a list of numbers separated by commas for multiple weekdays, or every(n) for every n weekdays)
//	Being 0 sunday, 1 monday, etc.
//	Can also be a name of the weekday (sunday, monday, etc.), in case insensitive and accepting 3 first letters or full name
type CronArgs struct {
	Minute  string
	Hour    string
	Day     string
	Month   string
	Weekday string
}

func (c CronArgs) String() string {
	return fmt.Sprintf("Minute: %s, Hour: %s, Day: %s, Month: %s, Weekday: %s", c.Minute, c.Hour, c.Day, c.Month, c.Weekday)
}

func (c CronArgs) CronString() string {
	minute := GetNumberParam(c.Minute)
	hour := GetNumberParam(c.Hour)
	day := GetNumberParam(c.Day)
	month := GetNumberParam(ReplaceMonth(c.Month))
	weekday := GetNumberParam(ReplaceWeekday(c.Weekday))

	return minute + " " + hour + " " + day + " " + month + " " + weekday
}
