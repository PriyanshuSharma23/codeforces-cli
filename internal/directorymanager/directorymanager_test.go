package directorymanager

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestEnsureProblemFileExists(t *testing.T) {
	// Setup: create a temporary root directory
	tempRoot, err := os.MkdirTemp("", "dm_test_root")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempRoot) // clean up

	// Setup: contest and problem values
	contestCode := 1234
	problemCode := "A"
	ext := "txt"
	problem := Problem{
		ContestCode: contestCode,
		ProblemCode: problemCode,
	}

	// Setup: create directory and file
	contestDir := filepath.Join(tempRoot, "1234")
	err = os.MkdirAll(contestDir, 0o755)
	if err != nil {
		t.Fatalf("Failed to create contest directory: %v", err)
	}

	problemFilePath := filepath.Join(contestDir, problemCode+"."+ext)
	content := "sample problem content"
	err = os.WriteFile(problemFilePath, []byte(content), 0o644)
	if err != nil {
		t.Fatalf("Failed to create problem file: %v", err)
	}

	// Instantiate the DirectoryManager
	logger := log.New(os.Stdout, "[TEST] ", log.LstdFlags)
	dm := NewDirectoryManager(tempRoot, ext, logger)

	// Call the method under test
	returnedPath, err := dm.ProblemFileExists(problem)
	if err != nil {
		t.Fatalf("EnsureProblemFileExists failed: %v", err)
	}

	// Assert: returned path should match full file path
	expectedPath := filepath.Join(tempRoot, "1234", "A.txt")
	if returnedPath != expectedPath {
		t.Errorf("Expected returned path to be %s, got %s", expectedPath, returnedPath)
	}
}

func TestEnsureProblemFileExists_MissingFile(t *testing.T) {
	// Setup: temporary root dir
	tempRoot, err := os.MkdirTemp("", "dm_test_missing_file")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempRoot)

	// Create contest directory but no problem file
	contestCode := 5678
	problemCode := "B"
	ext := "txt"
	problem := Problem{
		ContestCode: contestCode,
		ProblemCode: problemCode,
	}
	contestDir := filepath.Join(tempRoot, "5678")
	_ = os.MkdirAll(contestDir, 0o755)

	logger := log.New(os.Stdout, "[TEST] ", log.LstdFlags)
	dm := NewDirectoryManager(tempRoot, ext, logger)

	_, err = dm.ProblemFileExists(problem)
	if err == nil {
		t.Fatal("Expected error for missing problem file, got nil")
	}

	if !strings.Contains(err.Error(), "no such file or directory") {
		t.Errorf("Unexpected error message: %v", err)
	}
}

func TestEnsureProblemFileExists_MissingDirectory(t *testing.T) {
	tempRoot, err := os.MkdirTemp("", "dm_test_missing_dir")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempRoot)

	// No directory created
	problem := Problem{
		ContestCode: 9999,
		ProblemCode: "Z",
	}
	ext := "txt"
	logger := log.New(os.Stdout, "[TEST] ", log.LstdFlags)
	dm := NewDirectoryManager(tempRoot, ext, logger)

	_, err = dm.ProblemFileExists(problem)
	if err == nil {
		t.Fatal("Expected error for missing contest directory, got nil")
	}
	if !strings.Contains(err.Error(), "no such file or directory") {
		t.Errorf("Unexpected error message: %v", err)
	}
}

func TestCreateProblemFile_Success(t *testing.T) {
	// Setup
	tempRoot, err := os.MkdirTemp("", "dm_create_success")
	if err != nil {
		t.Fatalf("Failed to create temp root: %v", err)
	}
	defer os.RemoveAll(tempRoot)

	contestCode := 1111
	problemCode := "X"
	ext := "txt"
	problem := Problem{ContestCode: contestCode, ProblemCode: problemCode}

	logger := log.New(os.Stdout, "[TEST] ", log.LstdFlags)
	dm := NewDirectoryManager(tempRoot, ext, logger)

	// Act
	path, err := dm.CreateProblemFile(problem)
	if err != nil {
		t.Fatalf("CreateProblemFile failed: %v", err)
	}

	// Assert
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf("Expected file to exist at %s, but it does not", path)
	}
}

func TestCreateProblemFile_AlreadyExists(t *testing.T) {
	// Setup
	tempRoot, err := os.MkdirTemp("", "dm_create_exists")
	if err != nil {
		t.Fatalf("Failed to create temp root: %v", err)
	}
	defer os.RemoveAll(tempRoot)

	contestCode := 2222
	problemCode := "Y"
	ext := "txt"
	problem := Problem{ContestCode: contestCode, ProblemCode: problemCode}

	contestDir := filepath.Join(tempRoot, "2222")
	err = os.MkdirAll(contestDir, 0o755)
	if err != nil {
		t.Fatalf("Failed to create contest dir: %v", err)
	}

	// Create file manually
	filePath := filepath.Join(contestDir, "Y.txt")
	err = os.WriteFile(filePath, []byte("existing content"), 0o644)
	if err != nil {
		t.Fatalf("Failed to create initial file: %v", err)
	}

	logger := log.New(os.Stdout, "[TEST] ", log.LstdFlags)
	dm := NewDirectoryManager(tempRoot, ext, logger)

	// Act
	_, err = dm.CreateProblemFile(problem)
	if err == nil {
		t.Fatal("Expected error for existing file, got nil")
	}
	if !strings.Contains(err.Error(), "already exists") {
		t.Errorf("Expected 'already exists' error, got: %v", err)
	}
}

func TestCreateProblemFile_CreatesDirectory(t *testing.T) {
	// Setup
	tempRoot, err := os.MkdirTemp("", "dm_create_dir")
	if err != nil {
		t.Fatalf("Failed to create temp root: %v", err)
	}
	defer os.RemoveAll(tempRoot)

	// Directory will not exist beforehand
	contestCode := 3333
	problemCode := "Z"
	ext := "txt"
	problem := Problem{ContestCode: contestCode, ProblemCode: problemCode}

	logger := log.New(os.Stdout, "[TEST] ", log.LstdFlags)
	dm := NewDirectoryManager(tempRoot, ext, logger)

	// Act
	path, err := dm.CreateProblemFile(problem)
	if err != nil {
		t.Fatalf("CreateProblemFile failed: %v", err)
	}

	// Assert
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf("Expected file to exist at %s, but it does not", path)
	}
}
