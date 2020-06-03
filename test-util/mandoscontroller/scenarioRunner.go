package mandoscontroller

import mj "github.com/ElrondNetwork/elrond-vm-util/test-util/mandosjson"

// ScenarioExecutor describes a component that can run a VM scenario.
type ScenarioExecutor interface {
	// Reset clears state/world.
	Reset()

	// ExecuteScenario executes the scenario and checks if it passed. Failure is signaled by returning an error.
	// The FileResolver helps with resolving external steps.
	// TODO: group into a "execution context" param.
	ExecuteScenario(*mj.Scenario, mj.FileResolver) error
}

// ScenarioRunner is a component that can run json scenarios, using a provided executor.
type ScenarioRunner struct {
	Executor ScenarioExecutor
	Parser   mj.Parser
}

// NewScenarioRunner creates new ScenarioRunner instance.
func NewScenarioRunner(executor ScenarioExecutor, fileResolver mj.FileResolver) *ScenarioRunner {
	return &ScenarioRunner{
		Executor: executor,
		Parser: mj.Parser{
			FileResolver: fileResolver,
		},
	}
}
