package toolargs

func GetValidArg(args string) string {
	// Extract raw content
	content := FromArgFormat(args)

	// return valid json
	content = NewToolArg(content).ToJSON()

	return NewToolArg(content).ToJSON()
}
