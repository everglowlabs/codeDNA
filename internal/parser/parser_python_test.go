package parser

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPythonParser(t *testing.T) {
	// Create a temp file for testing
	content := `
class DatabaseManager:
    def connect(self):
        pass

def calculate_dna_score(sequence):
    return 0.5
`
	tmpDir, err := os.MkdirTemp("", "parser_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	tmpFile := filepath.Join(tmpDir, "sample.py")
	if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	p := NewParser(".py")
	if p == nil {
		t.Fatal("Failed to create Python parser")
	}

	blocks, err := p.ExtractBlocks(tmpFile)
	if err != nil {
		t.Fatalf("Error extracting blocks: %v", err)
	}

	if len(blocks) == 0 {
		t.Error("Expected to find blocks, found 0")
	}

	foundClass := false
	foundFunc := false
	for _, b := range blocks {
		if b.Kind == "class" {
			foundClass = true
		}
		if b.Kind == "function" {
			foundFunc = true
		}
	}

	if !foundClass {
		t.Error("Expected to find class block")
	}
	if !foundFunc {
		t.Error("Expected to find function block")
	}
}
