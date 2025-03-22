package helpers

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
)

func ReadYAMLFile(path string, out interface{}) error {
	if path == "" {
		return errorPathIsEmpty
	}

	completePath, err := filepath.Abs(path)
	if err != nil {
		return errorInvalidPath
	}

	yamlFile, err := os.Open(completePath)
	if err != nil {
		return errors.New(fmt.Sprintf(errorReadYamlFile, err.Error()))
	}

	defer yamlFile.Close()

	if yamlFile == nil {
		return errorYamlFileIsEmpty
	}

	decoder := yaml.NewDecoder(yamlFile)
	err = decoder.Decode(out)
	if err != nil {
		return errors.New(fmt.Sprintf(errorDecodeYamlFile, err.Error()))
	}

	return nil
}
