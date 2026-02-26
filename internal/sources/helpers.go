package sources

import (
	"gopkg.in/yaml.v3"
)

// ParseConfig Marshals and Unmarshals a config to a given model
func ParseConfig[T interface{}](config SourceConfigType) (T, error) {
	var cfg T

	encodedYaml, err := yaml.Marshal(config)
	if err != nil {
		return cfg, err
	}

	if err = yaml.Unmarshal(encodedYaml, &cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}
