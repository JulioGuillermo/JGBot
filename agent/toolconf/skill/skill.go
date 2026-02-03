package skill

import (
	"JGBot/agent/prompt"
	"JGBot/agent/tools"
	"JGBot/ctxs"
	"JGBot/skill"
	"JGBot/skill/skillexec"
	"context"
	"fmt"
)

type SkillArgs struct {
	Action string              `json:"action" description:"The action to execute. 'list' to list all the available skills, 'read' to read a skill, or 'exec' to execute a skill if it has a skill tool. (Note: The action is required)"`
	Name   string              `json:"name" description:"The name of the skill to read or execute (Note: Not required for 'list')."`
	Args   skillexec.SkillArgs `json:"args" description:"The arguments to pass to the skill (Note: Not required for 'list' or 'read', but required for 'exec')."`
}

type SkillInitializerConf struct{}

func (c *SkillInitializerConf) Name() string {
	return "skill"
}

func (c *SkillInitializerConf) listSkills(rCtx *ctxs.RespondCtx) string {
	return prompt.GetSkillsPrompt(rCtx.SessionConf)
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
