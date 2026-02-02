package skill

import (
	"JGBot/agent/tools"
	"JGBot/ctxs"
	"JGBot/skill"
	"JGBot/skill/skillexec"
	"context"
	"fmt"
	"strings"
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

func (c *SkillInitializerConf) listSkills() string {
	var sb strings.Builder
	sb.WriteString("# Available skills:\n")
	for _, skill := range skill.Skills {
		fmt.Fprintf(&sb, "- %s: %s\n", skill.Name, skill.Description)
		if skill.HasTool {
			fmt.Fprintf(&sb, "    This skill has a skill tool.\n")
		}
	}
	return sb.String()
}

func (c *SkillInitializerConf) readSkill(name string) string {
	skill, ok := skill.Skills[name]
	if !ok {
		return fmt.Sprintf("Skill %s not found", name)
	}
	if skill.HasTool {
		return fmt.Sprintf("Skill %s has a skill tool\n\n%s", name, skill.Content)
	}
	return skill.Content
}

func (c *SkillInitializerConf) execSkill(rCtx *ctxs.RespondCtx, name string, args skillexec.SkillArgs) (string, error) {
	sk, ok := skill.Skills[name]
	if !ok {
		return "", fmt.Errorf("Skill %s not found", name)
	}
	if !sk.HasTool {
		return "", fmt.Errorf("Skill %s has not skill tool", name)
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
				return c.listSkills(), nil
			case "read":
				return c.readSkill(args.Name), nil
			case "exec":
				return c.execSkill(rCtx, args.Name, args.Args)
			}

			return "Invalid action, please use 'list', 'read', or 'exec'", nil
		},
	}
}
