package tools

import (
	"github.com/tmc/langchaingo/tools"
)

func GetTools(customTools []tools.Tool) ([]tools.Tool, error) {
	agentsTools := append(
		customTools,
		tools.Calculator{},
	)

	// tool, err := serpapi.New()
	// if err != nil {
	// 	return nil, err
	// }
	// agentsTools = append(agentsTools, tool)

	return agentsTools, nil
}
