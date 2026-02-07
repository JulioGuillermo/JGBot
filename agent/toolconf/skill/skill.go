package skill

import (
	"JGBot/agent/tools"
	"JGBot/ctxs"
	"JGBot/log"
	"JGBot/session/sessionconf/sc"
	"JGBot/skill"
	"JGBot/skill/skillexec"
	"context"
	"fmt"
	"strings"
)

type SkillArgs struct {
	Action string              `json:"action" description:"The action to execute. 'list' to list all the available skills, 'read' to read a skill, or 'exec' to execute a skill if it has a skill tool. (Note: The action is required). Important this action is for this tool only, not for the skill itself. The skill is executed using the 'exec' action and the 'name' of the skill and the 'args' of the skill, everything you need to pass to the skill tool on the exec action you have to pass it on this args parameter."`
	Name   string              `json:"name" description:"The name of the skill to read or execute (Note: Not required for 'list'). This is the name of the skill, and it is not passed to the skill tool on the exec action."`
	Args   skillexec.SkillArgs `json:"args" description:"The arguments to pass to the skill (Note: Not required for 'list' or 'read', but required for 'exec'). This is the arguments of the skill during the exec action, everything you need to pass to the skill tool on the exec action you have to pass it on this args parameter."`
}

type SkillInitializerConf struct{}

func (c *SkillInitializerConf) Name() string {
	return "skill"
}

func (c *SkillInitializerConf) listSkills(rCtx *ctxs.RespondCtx) string {
	return GetSkillsPrompt(rCtx.SessionConf)
}

func (c *SkillInitializerConf) readSkill(rCtx *ctxs.RespondCtx, name string) string {
	skillConf := rCtx.SessionConf.GetSkillConf(name)
	if skillConf == nil || !skillConf.Enabled {
		return fmt.Sprintf("Skill %s not available", name)
	}

	skill, ok := skill.Skills[skillConf.Name]
	if !ok {
		return fmt.Sprintf("Skill %s not found", name)
	}
	if skill.HasTool {
		return fmt.Sprintf("Skill %s has a skill tool\n\n%s", name, skill.Content)
	}
	return skill.Content
}

func (c *SkillInitializerConf) execSkill(rCtx *ctxs.RespondCtx, name string, args skillexec.SkillArgs) (string, error) {
	skillConf := rCtx.SessionConf.GetSkillConf(name)
	if skillConf == nil || !skillConf.Enabled {
		return "", fmt.Errorf("Skill %s not available", name)
	}

	sk, ok := skill.Skills[skillConf.Name]
	if !ok {
		return "", fmt.Errorf("Skill %s not found", skillConf.Name)
	}
	if !sk.HasTool {
		return "", fmt.Errorf("Skill %s has not skill tool", skillConf.Name)
	}
	return skillexec.ExecSkillTool(sk.Dir, args, rCtx)
}

func (c *SkillInitializerConf) ToolInitializer(rCtx *ctxs.RespondCtx) tools.Tool {
	return &tools.ToolAutoArgs[SkillArgs]{
		ToolName:        c.Name(),
		ToolDescription: "Allows you to list, read, or execute skills.",
		ToolFunc: func(ctx context.Context, args SkillArgs) (string, error) {
			switch args.Action {
			case "list":
				return c.listSkills(rCtx), nil
			case "read":
				return c.readSkill(rCtx, args.Name), nil
			case "exec":
				return c.execSkill(rCtx, args.Name, args.Args)
			}

			return "Invalid action, please use 'list', 'read', or 'exec'", nil
		},
	}
}

func GetSkillsPrompt(conf *sc.SessionConf) string {
	var sb strings.Builder
	sb.WriteString("**Available skills:**\n")
	sb.WriteString("You have access to the following skills, you can read them using the `skill` tool with the `read` action and the `name` of the skill or exec them using the `skill` tool with the `exec` action and the `name` of the skill and the `args` of the skill:\n")
	for _, skillConf := range conf.Skills {
		if !skillConf.Enabled {
			continue
		}

		skill, ok := skill.Skills[skillConf.Name]
		if !ok {
			log.Warn("Skill not found", "skill", skillConf.Name)
			continue
		}

		if skill.HasTool {
			fmt.Fprintf(&sb, "- %s: [SkillTool available to exec through the skill tool] %s\n", skill.Name, skill.Description)
		} else {
			fmt.Fprintf(&sb, "- %s: %s\n", skill.Name, skill.Description)
		}
	}
	return sb.String()
}
