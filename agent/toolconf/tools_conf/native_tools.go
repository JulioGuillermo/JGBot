package tools_conf

import (
	"JGBot/agent/toolconf/cron"
	"JGBot/agent/toolconf/javascript"
	"JGBot/agent/toolconf/reactiontool"
	"JGBot/agent/toolconf/skill"
	"JGBot/agent/toolconf/timer"
)

var NativeTools = []ToolInitializerConf{
	&reactiontool.ReactionInitializerConf{},
	&javascript.JavaScriptInitializerConf{},
	&skill.SkillInitializerConf{},
	&cron.CronInitializerConf{},
	&timer.TimerInitializerConf{},
}
