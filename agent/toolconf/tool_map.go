package toolconf

import "JGBot/agent/toolconf/tools_conf"

func GetToolMap() map[string]tools_conf.ToolInitializerConf {
	toolMap := make(map[string]tools_conf.ToolInitializerConf)
	for _, tool := range tools_conf.NativeTools {
		toolMap[tool.Name()] = tool
	}
	return toolMap
}
