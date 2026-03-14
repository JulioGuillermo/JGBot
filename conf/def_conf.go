package conf

type Respond struct {
	Always *bool
	Match  *string
}

type Tool struct {
	Name    *string
	Enabled *bool
}

type Skill struct {
	Name        *string
	Enabled     *bool
	Description *string
}

type DefConf struct {
	Allowed *bool
	Respond *Respond

	HistorySize      *int
	Provider         *string
	SystemPromptFile *string
	AgentMaxIters    *int

	Tools  *[]Tool
	Skills *[]Skill
}
