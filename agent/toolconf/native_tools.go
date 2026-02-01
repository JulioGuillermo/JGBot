package toolconf

import (
	"JGBot/agent/toolconf/javascript"
	"JGBot/agent/toolconf/reactiontool"
)

var NativeTools = []ToolInitializerConf{
	&reactiontool.ReactionInitializerConf{},
	&javascript.JavaScriptInitializerConf{},
}
