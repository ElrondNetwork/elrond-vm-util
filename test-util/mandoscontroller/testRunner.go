package mandoscontroller

import mj "github.com/ElrondNetwork/elrond-vm-util/test-util/mandosjson"

// TestExecutor describes a component that can run a VM test.
type TestExecutor interface {
	// ExecuteTest executes the test and checks if it passed. Failure is signaled by returning an error.
	ExecuteTest(*mj.Test) error
}

// TestRunner is a component that can run tests, using a provided executor.
type TestRunner struct {
	Executor TestExecutor
	Parser   mj.Parser
}

// NewTestRunner creates new TestRunner instance.
func NewTestRunner(executor TestExecutor, fileResolver mj.FileResolver) *TestRunner {
	return &TestRunner{
		Executor: executor,
		Parser: mj.Parser{
			FileResolver: fileResolver,
		},
	}
}
