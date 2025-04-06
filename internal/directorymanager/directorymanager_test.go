package directorymanager

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/PriyanshuSharma23/codeforces-cli/internal/execution"
)

func setupTestManager(t *testing.T) (*DirectoryManager, string) {
	t.Helper()
	root := t.TempDir()
	logger := log.New(os.Stdout, "[test] ", log.Lshortfile)
	return NewDirectoryManager(root, logger), root
}

func sampleProblem() Problem {
	return Problem{ContestCode: 1234, ProblemCode: "A"}
}

func TestFullProblemPath(t *testing.T) {
	dm, root := setupTestManager(t)
	p := sampleProblem()
	expected := filepath.Join(root, "1234", "A")
	if path := dm.FullProblemPath(p); path != expected {
		t.Errorf("expected path %q, got %q", expected, path)
	}
}

func TestEnsureDir(t *testing.T) {
	dm, _ := setupTestManager(t)
	p := sampleProblem()
	dir, err := dm.EnsureDir(p)
	if err != nil {
		t.Fatalf("EnsureDir failed: %v", err)
	}
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Errorf("directory was not created: %s", dir)
	}
}

func TestWriteTestCases(t *testing.T) {
	dm, _ := setupTestManager(t)
	p := sampleProblem()
	_, err := dm.EnsureDir(p)
	if err != nil {
		t.Fatalf("WriteTestCases failed: %v", err)
	}

	testCases := []execution.TestCase{
		{Input: "1 2", Output: "3"},
		{Input: "4 5", Output: "9"},
	}
	err = dm.WriteTestCases(p, testCases, "input", "output")
	if err != nil {
		t.Fatalf("WriteTestCases failed: %v", err)
	}

	for i := range testCases {
		inFile := filepath.Join(dm.FullProblemPath(p), fmt.Sprintf("input%d", i+1))
		outFile := filepath.Join(dm.FullProblemPath(p), fmt.Sprintf("output%d", i+1))

		checkFileContains(t, inFile, testCases[i].Input)
		checkFileContains(t, outFile, testCases[i].Output)
	}
}

func TestWriteMetadata(t *testing.T) {
	dm, _ := setupTestManager(t)
	p := sampleProblem()
	_, err := dm.EnsureDir(p)
	if err != nil {
		t.Fatalf("WriteTestCases failed: %v", err)
	}

	meta := map[string]string{"name": "Test Problem"}
	err = dm.WriteMetadata(p, meta)
	if err != nil {
		t.Fatalf("WriteMetadata failed: %v", err)
	}

	file := filepath.Join(dm.FullProblemPath(p), "problem.json")
	data, err := os.ReadFile(file)
	if err != nil {
		t.Fatalf("reading metadata file failed: %v", err)
	}

	var parsed map[string]string
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("metadata file is not valid JSON: %v", err)
	}
	if parsed["name"] != "Test Problem" {
		t.Errorf("unexpected metadata content: %v", parsed)
	}
}

func TestWriteProgramFile(t *testing.T) {
	dm, _ := setupTestManager(t)
	p := sampleProblem()
	_, err := dm.EnsureDir(p)
	if err != nil {
		t.Fatalf("WriteTestCases failed: %v", err)
	}

	code := `#include <iostream>`
	filename := "main.cpp"
	err = dm.WriteProgramFile(p, filename, code)
	if err != nil {
		t.Fatalf("WriteProgramFile failed: %v", err)
	}

	path := filepath.Join(dm.FullProblemPath(p), filename)
	checkFileContains(t, path, code)

	// Try again; should not overwrite
	err = dm.WriteProgramFile(p, filename, "new content")
	if err != nil {
		t.Fatalf("WriteProgramFile on existing file should succeed: %v", err)
	}

	data, _ := os.ReadFile(path)
	if string(data) != code {
		t.Errorf("existing file should not be overwritten")
	}
}

func checkFileContains(t *testing.T, path, expected string) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read file %s: %v", path, err)
	}
	if string(data) != expected {
		t.Errorf("file content mismatch: expected %q, got %q", expected, string(data))
	}
}
