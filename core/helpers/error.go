package helpers

import "errors"

const (
	errorReadYamlFile   = "yaml file read error: %v"
	errorDecodeYamlFile = "yaml file decode error: %v"
)

var (
	errorPathIsEmpty     = errors.New("path is required")
	errorInvalidPath     = errors.New("path is invalid")
	errorYamlFileIsEmpty = errors.New("yaml file is empty")
)
