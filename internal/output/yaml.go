package output

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

func PrintYAML(config CompactConfig) error {
	yamlData, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	fmt.Print(string(yamlData))
	return nil
}
