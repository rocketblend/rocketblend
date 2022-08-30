package library

import (
	"encoding/json"
	"fmt"
	"os"
)

type JSONRepository struct {
}

func NewRepository() *JSONRepository {
	return &JSONRepository{}
}

func (r *JSONRepository) FindBuildByPath(path string) (*Build, error) {
	file, err := r.findByPath(path)
	if err != nil {
		return nil, err
	}

	build := Build{}
	if err := json.Unmarshal(file, &build); err != nil {
		return nil, fmt.Errorf("failed to unmarshal build config: %s", err)
	}

	return &build, nil
}

func (r *JSONRepository) FindPackageByPath(path string) (*Build, error) {
	file, err := r.findByPath(path)
	if err != nil {
		return nil, err
	}

	build := Build{}
	if err := json.Unmarshal(file, &build); err != nil {
		return nil, fmt.Errorf("failed to unmarshal package config: %s", err)
	}

	return &build, nil
}

func (r *JSONRepository) findByPath(path string) ([]byte, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %s", err)
	}

	return file, nil
}
