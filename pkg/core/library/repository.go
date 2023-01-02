package library

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type (
	JSONRepository struct {
		dir string
	}
)

func NewRepository(dir string) *JSONRepository {
	return &JSONRepository{
		dir: dir,
	}
}

func (r *JSONRepository) CreateBuild(build *Build) error {
	return r.create(build, build.Reference, BuildFile)
}

func (r *JSONRepository) FindBuildByRef(ref string) (*Build, error) {
	file, err := r.findByReference(ref, BuildFile)
	if err != nil {
		return nil, err
	}

	build := Build{}
	if err := json.Unmarshal(file, &build); err != nil {
		return nil, fmt.Errorf("failed to unmarshal build config: %s", err)
	}

	return &build, nil
}

func (r *JSONRepository) CreatePackage(pack *Package) error {
	return r.create(pack, pack.Reference, PackgeFile)
}

func (r *JSONRepository) FindPackageByRef(ref string) (*Package, error) {
	file, err := r.findByReference(ref, PackgeFile)
	if err != nil {
		return nil, err
	}

	pack := Package{}
	if err := json.Unmarshal(file, &pack); err != nil {
		return nil, fmt.Errorf("failed to unmarshal package config: %s", err)
	}

	return &pack, nil
}

func (r *JSONRepository) getReferencePath(ref string) string {
	return filepath.Join(r.dir, ref)
}

func (r *JSONRepository) findByReference(ref string, file string) ([]byte, error) {
	f, err := os.ReadFile(filepath.Join(r.getReferencePath(ref), file))
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %s", err)
	}

	return f, nil
}

func (r *JSONRepository) create(obj interface{}, ref string, file string) error {
	data, err := json.Marshal(obj)
	if err != nil {
		return fmt.Errorf("failed to marshal object: %s", err)
	}

	// Local reference path
	path := r.getReferencePath(ref)

	// Create output directories
	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}

	if err := os.WriteFile(filepath.Join(path, file), data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %s", err)
	}

	return nil
}
