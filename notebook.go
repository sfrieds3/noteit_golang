package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const notebookExt = ".md"

// Store manages notebooks on disk. Each notebook is a single markdown file
// named <notebook>.md inside Dir.
type Store struct {
	Dir string
}

// Path returns the file path for the named notebook. It rejects names that
// would resolve outside of the store directory.
func (s Store) Path(name string) (string, error) {
	if name == "" {
		return "", errors.New("notebook name cannot be empty")
	}
	if name == "." || name == ".." || strings.ContainsAny(name, `/\`) {
		return "", fmt.Errorf("invalid notebook name %q", name)
	}
	return filepath.Join(s.Dir, name+notebookExt), nil
}

// Ensure creates the named notebook (and the store directory) if it does not
// already exist, and returns its path.
func (s Store) Ensure(name string) (string, error) {
	path, err := s.Path(name)
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(s.Dir, 0o755); err != nil {
		return "", fmt.Errorf("creating notes directory: %w", err)
	}

	f, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o644)
	if errors.Is(err, fs.ErrExist) {
		return path, nil
	}
	if err != nil {
		return "", fmt.Errorf("creating notebook %q: %w", name, err)
	}
	if _, err := fmt.Fprintf(f, "# %s\n\n", name); err != nil {
		f.Close()
		return "", fmt.Errorf("writing notebook %q: %w", name, err)
	}
	return path, f.Close()
}

// Add appends a note as a markdown list item to the named notebook,
// creating the notebook first if needed.
func (s Store) Add(name, note string) (string, error) {
	path, err := s.Ensure(name)
	if err != nil {
		return "", err
	}

	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return "", fmt.Errorf("opening notebook %q: %w", name, err)
	}
	if _, err := fmt.Fprintf(f, "- %s\n", note); err != nil {
		f.Close()
		return "", fmt.Errorf("writing to notebook %q: %w", name, err)
	}
	return path, f.Close()
}

// Read returns the contents of the named notebook.
func (s Store) Read(name string) ([]byte, error) {
	path, err := s.Path(name)
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, fmt.Errorf("notebook %q does not exist", name)
	}
	return data, err
}

// List returns the names of all notebooks, sorted alphabetically.
func (s Store) List() ([]string, error) {
	entries, err := os.ReadDir(s.Dir)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var names []string
	for _, e := range entries {
		if e.IsDir() || filepath.Ext(e.Name()) != notebookExt {
			continue
		}
		names = append(names, strings.TrimSuffix(e.Name(), notebookExt))
	}
	sort.Strings(names)
	return names, nil
}
