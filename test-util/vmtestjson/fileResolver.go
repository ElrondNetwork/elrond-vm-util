package vmtestjson

// FileResolver resolves values starting with "file:"
type FileResolver interface {
	// SetContext sets directory where the test runs, to help resolve relative paths.
	SetContext(contextPath string)

	// ResolveFileValue converts a value prefixed with "file:" and replaces it with the file contents.
	ResolveFileValue(value string) ([]byte, error)
}
