package scanner

import (
	"os"
	"path/filepath"
)

type Scanner struct {
	Root   string
	Filter *Filter
}

func NewScanner(root string) *Scanner {
	return &Scanner{
		Root:   root,
		Filter: NewFilter(root),
	}
}

func (s *Scanner) Scan(onFile func(path string, total int)) ([]string, error) {
	var files []string

	// First pass to count total supported files
	var total int
	_ = filepath.Walk(s.Root, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		rel, _ := filepath.Rel(s.Root, path)
		if s.Filter.ShouldIgnore(rel) || !s.Filter.IsSupported(rel) {
			return nil
		}
		total++
		return nil
	})

	err := filepath.Walk(s.Root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		rel, _ := filepath.Rel(s.Root, path)
		if s.Filter.ShouldIgnore(rel) || !s.Filter.IsSupported(rel) {
			return nil
		}

		files = append(files, path)
		if onFile != nil {
			onFile(rel, total)
		}
		return nil
	})

	return files, err
}
