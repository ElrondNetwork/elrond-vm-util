package testcontroller

import ij "github.com/ElrondNetwork/elrond-vm-util/test-util/vmtestjson"

// VMTestExecutor describes a component that can run a VM test.
type VMTestExecutor interface {
	// Run executes the test and checks if it passed. Failure is signaled by returning an error.
	Run(*ij.Test) error
}

type Runner struct {
	Executor VMTestExecutor
	Parser   ij.Parser
}

func NewRunner(executor VMTestExecutor, fileResolver ij.FileResolver) *Runner {
	return &Runner{
		Executor: executor,
		Parser: ij.Parser{
			FileResolver: fileResolver,
		},
	}
}
