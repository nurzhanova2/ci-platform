package model

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Step struct {
	Name     string   `json:"name"`
	Image    string   `json:"image"`
	Commands []string `json:"commands"`
}

type Pipeline struct {
	Steps []Step `json:"steps"`
}

func LoadPipeline(path string) (*Pipeline, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var pipeline Pipeline
	err = yaml.Unmarshal(data, &pipeline)
	if err != nil {
		return nil, err
	}

	return &pipeline, nil
}