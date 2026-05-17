package parser

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRustParser(t *testing.T) {
	content := `
struct UserProfile {
    id: u64,
}

fn fetch_profile(id: u64) -> Result<UserProfile, String> {
    Ok(UserProfile { id })
}
`
	tmpDir, err := os.MkdirTemp("", "parser_test_rust")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	tmpFile := filepath.Join(tmpDir, "sample.rs")
	if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	p := NewParser(".rs")
	if p == nil {
		t.Fatal("Failed to create Rust parser")
	}

	blocks, err := p.ExtractBlocks(tmpFile)
	if err != nil {
		t.Fatalf("Error extracting Rust blocks: %v", err)
	}

	if len(blocks) == 0 {
		t.Error("Expected to find Rust blocks, found 0")
	}

	foundStruct := false
	foundFn := false
	for _, b := range blocks {
		if b.Kind == "class" {
			foundStruct = true
		}
		if b.Kind == "function" {
			foundFn = true
		}
	}

	if !foundStruct {
		t.Error("Expected to find Rust struct block (mapped to class)")
	}
	if !foundFn {
		t.Error("Expected to find Rust fn block (mapped to function)")
	}
}

func TestJavaParser(t *testing.T) {
	content := `
public class UserService {
    public void registerUser(String name) {
        System.out.println(name);
    }
}
`
	tmpDir, err := os.MkdirTemp("", "parser_test_java")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	tmpFile := filepath.Join(tmpDir, "sample.java")
	if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	p := NewParser(".java")
	if p == nil {
		t.Fatal("Failed to create Java parser")
	}

	blocks, err := p.ExtractBlocks(tmpFile)
	if err != nil {
		t.Fatalf("Error extracting Java blocks: %v", err)
	}

	if len(blocks) == 0 {
		t.Error("Expected to find Java blocks, found 0")
	}

	foundClass := false
	foundMethod := false
	for _, b := range blocks {
		if b.Kind == "class" {
			foundClass = true
		}
		if b.Kind == "method" {
			foundMethod = true
		}
	}

	if !foundClass {
		t.Error("Expected to find Java class block")
	}
	if !foundMethod {
		t.Error("Expected to find Java method block")
	}
}
