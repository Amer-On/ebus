package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

func ReadConfig[T any](path string) (*T, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	target := new(T)

	if err := yaml.NewDecoder(f).Decode(target); err != nil {
		return nil, err
	}

	return target, nil
}
