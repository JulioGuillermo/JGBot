package toolargs

func ToArgFormat(args string) string {
	arg := ToolArgFromJSON(args)

	// not valid json
	if arg == nil {
		return ToolArgFromAny(args).ToJSON()
	}

	// valid json and format __arg1 (string)
	if arg.Arg != "" {
		return arg.ToJSON()
	}

	// valid json but __arg1 is object
	content := checkArgNotString(args)
	if content != "" {
		return NewToolArg(content).ToJSON()
	}

	// valid json but not format __arg1
	return NewToolArg(args).ToJSON()
}
