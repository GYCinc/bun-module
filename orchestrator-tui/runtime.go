package main

import (
	"fmt"
	"math/rand"
	"time"
)

type AgentConfig struct {
	ID                 string
	Name               string
	Focus              string
	Instructions       string
	OutputExpectations string
	Model              string
	Temperature        string
}

type Swarm struct {
	MainTask     string
	Orchestrator AgentConfig
	SubAgents    []AgentConfig
}

type AgentResult struct {
	ID     string
	Name   string
	Output string
}

func MockExecuteSubagent(spec AgentConfig, task string, updateDelay time.Duration) AgentResult {
	// Simulate work
	time.Sleep(updateDelay + time.Duration(rand.Intn(1000))*time.Millisecond)
	return AgentResult{
		ID:     spec.ID,
		Name:   spec.Name,
		Output: fmt.Sprintf("[%s via %s (T:%s)] result for task...", spec.Name, spec.Model, spec.Temperature),
	}
}

func MockSynthesize(spec AgentConfig, task string, results []AgentResult, updateDelay time.Duration) string {
	time.Sleep(updateDelay + time.Duration(rand.Intn(1000))*time.Millisecond)
	return fmt.Sprintf("Final synthesized answer by Orchestrator [%s via %s] based on %d sub-agents.", spec.Name, spec.Model, len(results))
}
