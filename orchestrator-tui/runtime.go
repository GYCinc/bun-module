package main

import (
	"fmt"
	"math/rand"
	"time"
)

type SubAgentSpec struct {
	ID                 string
	Name               string
	Focus              string
	Instructions       string
	OutputExpectations string
}

type RoleBlueprint struct {
	Purpose   string
	MaxAgents int
}

func NewRoleBlueprint(purpose string, maxAgents int) *RoleBlueprint {
	if maxAgents < 1 {
		maxAgents = 1
	} else if maxAgents > 5 {
		maxAgents = 5
	}
	return &RoleBlueprint{
		Purpose:   purpose,
		MaxAgents: maxAgents,
	}
}

type AgentFocus struct {
	Name string
	Desc string
}

func (rb *RoleBlueprint) BuildSpecs(task string) []SubAgentSpec {
	var angles []AgentFocus

	switch rb.Purpose {
	case "researcher":
		angles = []AgentFocus{
			{"context", "Establish baseline context and key definitions"},
			{"evidence", "Find strongest supporting evidence and data"},
			{"gaps", "Identify unknowns, limitations, and open questions"},
			{"contrarian", "Stress-test with counterarguments and edge cases"},
			{"synthesis_view", "Extract actionable takeaways from all angles"},
		}
	case "coder":
		angles = []AgentFocus{
			{"implementation", "Build the minimal working solution"},
			{"edge_cases", "Identify and handle failure modes"},
			{"tests", "Propose concrete tests and validation"},
			{"performance", "Optimize for speed/cost/readability tradeoffs"},
			{"docs", "Write clear usage and integration notes"},
		}
	case "analyst":
		angles = []AgentFocus{
			{"quant", "Quantitative breakdown and key numbers"},
			{"qual", "Qualitative drivers and stakeholder view"},
			{"risks", "Downside risks and likelihood"},
			{"opps", "Upside opportunities and triggers"},
			{"recommendation", "Decide on next actions with rationale"},
		}
	default:
		angles = []AgentFocus{
			{"angle_1", "Cover the most important aspect"},
			{"angle_2", "Cover secondary considerations"},
			{"angle_3", "Flag risks and unknowns"},
			{"angle_4", "Suggest alternatives"},
			{"angle_5", "Summarize key implications"},
		}
	}

	limit := rb.MaxAgents
	if len(angles) < limit {
		limit = len(angles)
	}

	selected := angles[:limit]
	var specs []SubAgentSpec

	for i, angle := range selected {
		specs = append(specs, SubAgentSpec{
			ID:    fmt.Sprintf("sub_%d", i+1),
			Name:  angle.Name,
			Focus: angle.Desc,
			Instructions: fmt.Sprintf("You are sub-agent %d on a 1-orchestrator team. "+
				"Your unique responsibility: %s. "+
				"Work only on this angle for the task below. "+
				"Return your findings in a clear, self-contained form.", i+1, angle.Desc),
			OutputExpectations: "Output: your key findings, supporting bullets, and any " +
				"important assumptions specific to this angle.",
		})
	}

	return specs
}

type AgentResult struct {
	ID     string
	Name   string
	Output string
}

func MockExecuteSubagent(spec SubAgentSpec, task string, updateDelay time.Duration) AgentResult {
	// Simulate work
	time.Sleep(updateDelay + time.Duration(rand.Intn(1000))*time.Millisecond)
	return AgentResult{
		ID:     spec.ID,
		Name:   spec.Name,
		Output: fmt.Sprintf("[%s] result for task...", spec.Name),
	}
}

func MockSynthesize(task string, results []AgentResult, updateDelay time.Duration) string {
	time.Sleep(updateDelay + time.Duration(rand.Intn(1000))*time.Millisecond)
	return fmt.Sprintf("Final synthesized answer based on %d angles.", len(results))
}
