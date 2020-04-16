package vmtestjson

// FileProvider resolves values starting with "file:"
type FileProvider interface {
	// ResolveFileValue converts a value prefixed with "file:" and replaes it with the file contents.
	ResolveFileValue(value string) (string, error)
}

// Parser performs parsing of both json tests (older) and scenarios (new).
type Parser struct {
	FileProvider FileProvider
}
