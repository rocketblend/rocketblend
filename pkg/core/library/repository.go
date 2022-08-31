package library

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type (
	JSONRepository struct {
	}
)

func NewRepository() *JSONRepository {
	return &JSONRepository{}
}

func (r *JSONRepository) FindBuildByPath(path string) (*Build, error) {
	file, err := r.findByPath(path, BuildFile)
	if err != nil {
		return nil, err
	}

	build := Build{}
	if err := json.Unmarshal(file, &build); err != nil {
		return nil, fmt.Errorf("failed to unmarshal build config: %s", err)
	}

	return &build, nil
}

func (r *JSONRepository) FindPackageByPath(path string) (*Package, error) {
	file, err := r.findByPath(path, PackgeFile)
	if err != nil {
		return nil, err
	}

	pack := Package{}
	if err := json.Unmarshal(file, &pack); err != nil {
		return nil, fmt.Errorf("failed to unmarshal package config: %s", err)
	}

	return &pack, nil
}

func (r *JSONRepository) findByPath(path string, file string) ([]byte, error) {
	f, err := os.ReadFile(filepath.Join(path, file))
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %s", err)
	}

	return f, nil
}
