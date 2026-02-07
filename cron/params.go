package cron

import (
	"regexp"
	"strings"
)

func GetNumberParam(input string) string {
	// input can be a 'number', a 'list of number separated by commas' with or without '[]', or 'every(number)' or 'every number'

	// match every
	re := regexp.MustCompile(`every\s*\((\d+)\)`)
	matches := re.FindStringSubmatch(input)
	if len(matches) > 1 {
		return "*/" + matches[1]
	}
	re = regexp.MustCompile(`every\s+(\d+)`)
	matches = re.FindStringSubmatch(input)
	if len(matches) > 1 {
		return "*/" + matches[1]
	}

	// match list of number
	re = regexp.MustCompile(`\[?\s*(\d+(,\s*\d+)*)\s*\]?`)
	matches = re.FindStringSubmatch(input)
	if len(matches) > 1 {
		return strings.ReplaceAll(matches[1], " ", "")
	}

	// match number
	re = regexp.MustCompile(`^\d+$`)
	matches = re.FindStringSubmatch(input)
	if len(matches) > 1 {
		return matches[1]
	}

	return "*"
}
