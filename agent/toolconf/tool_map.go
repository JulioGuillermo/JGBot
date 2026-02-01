package toolconf

func GetToolMap() map[string]ToolInitializerConf {
	toolMap := make(map[string]ToolInitializerConf)
	for _, tool := range NativeTools {
		toolMap[tool.Name()] = tool
	}
	return toolMap
}
