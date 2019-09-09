package testcontroller

import ij "github.com/ElrondNetwork/elrond-vm-util/test-util/ielejson"

// VMTestExecutor describes a component that can run a VM test.
type VMTestExecutor interface {
	// ProcessCode takes the code as it is represented in the test, and converts it to something the VM can execute.
	ProcessCode(testPath string, value string) (string, error)

	// Run executes the test and checks if it passed. Failure is signaled by returning an error.
	Run(*ij.Test) error
}
