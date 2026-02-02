package skill

import (
	"JGBot/log"
	"os"
	"path"
	"regexp"
)

func LoadSkill(name string) *Skill {
	dir := path.Join(SkillDir, name)
	if !checkSkillDir(dir) {
		return nil
	}

	content := getSkillContent(dir)
	if content == "" {
		return nil
	}

	skill := skillParse(content)
	if skill == nil {
		return nil
	}

	if skill.Description == "" || skill.Content == "" {
		return nil
	}

	skill.Dir = name
	if skill.Name == "" {
		skill.Name = name
	}

	skill.HasTool = hasSkillTool(dir)

	return skill
}

func checkSkillDir(dir string) bool {
	info, err := os.Stat(dir)
	if err != nil {
		log.Error("Skill dir error", "error", err)
		return false
	}
	return info.IsDir()
}

func getSkillContent(dir string) string {
	content, err := os.ReadFile(path.Join(dir, SkillFile))
	if err != nil {
		log.Error("Skill content error", "error", err)
		return ""
	}
	return string(content)
}

func hasSkillTool(dir string) bool {
	stat, err := os.Stat(path.Join(dir, SkillToolFile))
	return err == nil && !stat.IsDir()
}

func skillParse(content string) *Skill {
	re := regexp.MustCompile(`(?s)^---\s*name:\s*(?P<name>.*?)\s*description:\s*(?P<description>[^\n]*)\s*.*?---\s*(?P<content>.*)`)

	match := re.FindStringSubmatch(content)

	if match == nil {
		return nil
	}

	result := make(map[string]string)
	for i, name := range re.SubexpNames() {
		if i != 0 && name != "" {
			result[name] = match[i]
		}
	}

	return &Skill{
		Name:        result["name"],
		Description: result["description"],
		Content:     result["content"],
	}
}
