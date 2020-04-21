package testcontroller

import ij "github.com/ElrondNetwork/elrond-vm-util/test-util/vmtestjson"

// VMScenarioExecutor describes a component that can run a VM scenario.
type VMScenarioExecutor interface {
	// Run executes the scenario and checks if it passed. Failure is signaled by returning an error.
	Run(*ij.Scenario) error
}

// ScenarioRunner is a component that can run json scenarios, using a provided executor.
type ScenarioRunner struct {
	Executor VMScenarioExecutor
	Parser   ij.Parser
}

// NewScenarioRunner creates new ScenarioRunner instance.
func NewScenarioRunner(executor VMScenarioExecutor, fileResolver ij.FileResolver) *ScenarioRunner {
	return &ScenarioRunner{
		Executor: executor,
		Parser: ij.Parser{
			FileResolver: fileResolver,
		},
	}
}
