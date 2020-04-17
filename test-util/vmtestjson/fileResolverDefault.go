package vmtestjson

import (
	"io/ioutil"
	"path/filepath"
)

var _ FileResolver = (*DefaultFileResolver)(nil)

// DefaultFileResolver loads file contents for the test parser.
type DefaultFileResolver struct {
	contextPath              string
	contractPathReplacements map[string]string
}

// NewDefaultFileResolver yields a new DefaultFileResolver instance.
func NewDefaultFileResolver() *DefaultFileResolver {
	return &DefaultFileResolver{
		contextPath:              "",
		contractPathReplacements: make(map[string]string),
	}
}

// ReplacePath offers the possibility to swap a path with another withouot providing a new set of tests.
// It is very useful when testing multiple contracts against the same tests.
func (fr *DefaultFileResolver) ReplacePath(pathInTest, actualPath string) *DefaultFileResolver {
	fr.contractPathReplacements[pathInTest] = actualPath
	return fr
}

func (fr *DefaultFileResolver) SetContext(contextPath string) {
	fr.contextPath = contextPath
}

func (fr *DefaultFileResolver) ResolveFileValue(value string) ([]byte, error) {
	if len(value) == 0 {
		return []byte{}, nil
	}
	var fullPath string
	if replacement, shouldReplace := fr.contractPathReplacements[value]; shouldReplace {
		fullPath = replacement
	} else {
		testDirPath := filepath.Dir(fr.contextPath)
		fullPath = filepath.Join(testDirPath, value)
	}
	scCode, err := ioutil.ReadFile(fullPath)
	if err != nil {
		return []byte{}, err
	}

	return scCode, nil
}
