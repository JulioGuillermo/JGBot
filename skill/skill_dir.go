package skill

import "os"

func readSkillDir() ([]string, error) {
	err := os.MkdirAll(SkillDir, 0755)
	if err != nil {
		return nil, err
	}

	elements, err := os.ReadDir(SkillDir)
	if err != nil {
		return nil, err
	}

	var plugs []string
	for _, element := range elements {
		if !element.IsDir() {
			continue
		}
		plugs = append(plugs, element.Name())
	}
	return plugs, nil
}
