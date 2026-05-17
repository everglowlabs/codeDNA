package scanner

import (
	"path/filepath"

	gitignore "github.com/sabhiram/go-gitignore"
)

type Filter struct {
	ignore *gitignore.GitIgnore
}

func NewFilter(root string) *Filter {
	ignorePath := filepath.Join(root, ".gitignore")
	gi, err := gitignore.CompileIgnoreFile(ignorePath)
	if err != nil {
		// If no .gitignore, just ignore .git directory by default
		gi = gitignore.CompileIgnoreLines(".git")
	}
	return &Filter{ignore: gi}
}

func (f *Filter) ShouldIgnore(path string) bool {
	return f.ignore.MatchesPath(path)
}

func (f *Filter) IsSupported(path string) bool {
	ext := filepath.Ext(path)
	switch ext {
	case ".go", ".ts", ".tsx", ".py", ".js", ".rs", ".java":
		return true
	}
	return false
}
