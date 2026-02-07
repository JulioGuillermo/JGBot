package cron

import (
	"JGBot/cron"
)

// CronSchedule is a struct that holds the arguments for a cron expression
type CronSchedule struct {
	// Every param is optional, by default is ignored if blank or *.
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
	Minute  string `json:"minute" description:"The minute of the hour, values range from 0-59, can be a number for an specific minute (e.g. '0'), a list of numbers separated by commas for multiple minutes (e.g. '0,15,30,45'), or every(n) for every n minutes (e.g. 'every(15)' for every 15 minutes)."`
	Hour    string `json:"hour" description:"The hour of the day, values range from 0-23, can be a number for an specific hour (e.g. '0'), a list of numbers separated by commas for multiple hours (e.g. '0,12,24'), or every(n) for every n hours (e.g. 'every(12)' for every 12 hours)."`
	Day     string `json:"day" description:"The day of the month, values range from 1-31, can be a number for an specific day (e.g. '1'), a list of numbers separated by commas for multiple days (e.g. '1,15,30'), or every(n) for every n days (e.g. 'every(15)' for every 15 days)."`
	Month   string `json:"month" description:"The month of the year, values range from 1-12, can be a number for an specific month (e.g. '1'), a list of numbers separated by commas for multiple months (e.g. '1,6,12'), or every(n) for every n months (e.g. 'every(6)' for every 6 months). Being 1 january, 2 february, etc. Can also be a name of the month (january, february, etc.), in case insensitive and accepting 3 first letters or full name."`
	Weekday string `json:"weekday" description:"The weekday of the week, values range from 0-6, can be a number for an specific weekday (e.g. '0'), a list of numbers separated by commas for multiple weekdays (e.g. '0,2,4'), or every(n) for every n weekdays (e.g. 'every(2)' for every 2 weekdays). Being 0 sunday, 1 monday, etc. Can also be a name of the weekday (sunday, monday, etc.), in case insensitive and accepting 3 first letters or full name."`
}

func (c *CronSchedule) ToCron() cron.CronArgs {
	return cron.CronArgs{
		Minute:  c.Minute,
		Hour:    c.Hour,
		Day:     c.Day,
		Month:   c.Month,
		Weekday: c.Weekday,
	}
}

func (c *CronSchedule) String() string {
	return c.ToCron().String()
}
