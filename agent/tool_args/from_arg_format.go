package toolargs

func extractArgFormat(args string) (string, bool) {
	arg := ToolArgFromJSON(args)
	if arg != nil && arg.Arg != "" { // valid json and format __arg1, so extract arg
		return arg.Arg, true
	}
	// valid json but __arg1 is not string
	content := checkArgNotString(args)
	if content != "" {
		return content, true
	}
	// valid json but not format __arg1, so json input
	return args, false
}

func FromArgFormat(args string) string {
	ok := true
	for ok {
		args, ok = extractArgFormat(args)
	}
	return args
}
