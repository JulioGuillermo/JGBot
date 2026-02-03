package tools_conf

import (
	"JGBot/agent/toolconf/javascript"
	"JGBot/agent/toolconf/reactiontool"
	"JGBot/agent/toolconf/skill"
)

var NativeTools = []ToolInitializerConf{
	&reactiontool.ReactionInitializerConf{},
	&javascript.JavaScriptInitializerConf{},
	&skill.SkillInitializerConf{},
}
