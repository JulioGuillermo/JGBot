package skill

type Skill struct {
	Name        string
	Description string
	Content     string
	Dir         string
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
