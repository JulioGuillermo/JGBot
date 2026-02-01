package toolconf

import "JGBot/agent/toolconf/reactiontool"

var NativeTools = []ToolInitializerConf{
	&reactiontool.ReactionInitializerConf{},
}
