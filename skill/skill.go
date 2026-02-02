package skill

var Skills map[string]*Skill

type Skill struct {
	Name        string
	Description string
	Content     string
	Dir         string
	HasTool     bool
}

func GetSkills() ([]*Skill, error) {
	skillDirs, err := readSkillDir()
	if err != nil {
		return nil, err
	}

	skillList := make([]*Skill, 0)
	for _, skillDir := range skillDirs {
		skill := LoadSkill(skillDir)
		if skill != nil {
			skillList = append(skillList, skill)
		}
	}
	return skillList, nil
}

func InitSkills() error {
	skills, err := GetSkills()
	if err != nil {
		return err
	}

	Skills = make(map[string]*Skill)
	for _, skill := range skills {
		Skills[skill.Name] = skill
	}

	return nil
}
