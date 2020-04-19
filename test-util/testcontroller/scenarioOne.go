package testcontroller

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

// RunSingleJSONScenario parses and prepares test, then calls testCallback.
func (r *ScenarioRunner) RunSingleJSONScenario(contextPath string) error {
	var err error
	contextPath, err = filepath.Abs(contextPath)
	if err != nil {
		return err
	}

	// Open our jsonFile
	var jsonFile *os.File
	jsonFile, err = os.Open(contextPath)
	// if we os.Open returns an error then handle it
	if err != nil {
		return err
	}

	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return err
	}

	r.Parser.FileResolver.SetContext(contextPath)
	scenario, parseErr := r.Parser.ParseScenarioFile(byteValue)
	if parseErr != nil {
		return parseErr
	}

	return r.Executor.Run(scenario)
}
