package testcontroller

import (
	"io/ioutil"
	"os"
	"path/filepath"

	ij "github.com/ElrondNetwork/elrond-vm-util/test-util/ielejson"
)

// RunSingleIeleTest parses and prepares test, then calls testCallback.
func RunSingleIeleTest(testFilePath string, testExecutor VMTestExecutor) error {
	var err error
	testFilePath, err = filepath.Abs(testFilePath)
	if err != nil {
		return err
	}

	// Open our jsonFile
	var jsonFile *os.File
	jsonFile, err = os.Open(testFilePath)
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

	top, parseErr := ij.ParseTopLevel(byteValue)
	if parseErr != nil {
		return parseErr
	}

	for _, test := range top {
		assembleErr := processCodeInTest(testFilePath, test, testExecutor)
		if assembleErr != nil {
			return assembleErr
		}
		testErr := testExecutor.Run(test)
		if testErr != nil {
			return testErr
		}
	}

	// toPath := strings.Replace(testFilePath, "iele-v2", "iele-v3", 1)
	// fmt.Println(toPath)
	// saveModifiedTest(toPath, top)

	return nil
}

func processCodeInTest(testFilePath string, test *ij.Test, testExecutor VMTestExecutor) error {
	testDirPath := filepath.Dir(testFilePath)

	var assErr error
	for _, acct := range test.Pre {
		acct.Code, assErr = testExecutor.ProcessCode(testDirPath, acct.Code)
		if assErr != nil {
			return assErr
		}
	}
	for _, acct := range test.PostState {
		acct.Code, assErr = testExecutor.ProcessCode(testDirPath, acct.Code)
		if assErr != nil {
			return assErr
		}
	}

	for _, block := range test.Blocks {
		for _, tx := range block.Transactions {
			if tx.IsCreate {
				tx.AssembledCode, assErr = testExecutor.ProcessCode(testDirPath, tx.ContractCode)
				if assErr != nil {
					return assErr
				}
			}
		}
	}

	return nil
}

// tool to modify tests
// use with extreme caution
func saveModifiedTest(toPath string, top []*ij.Test) {
	resultJSON := ij.ToJSONString(top)

	err := os.MkdirAll(filepath.Dir(toPath), os.ModePerm)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(toPath, []byte(resultJSON), 0644)
	if err != nil {
		panic(err)
	}
}
