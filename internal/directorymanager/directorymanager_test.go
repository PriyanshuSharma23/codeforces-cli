package directorymanager

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/PriyanshuSharma23/codeforces-cli/internal/execution"
)

func TestProblemDirExists_Success(t *testing.T) {
	tempRoot, err := os.MkdirTemp("", "dm_test_root")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempRoot)

	// Setup: Create contest/problem dir
	contestCode := 1234
	problemCode := "A"
	problem := Problem{ContestCode: contestCode, ProblemCode: problemCode}
	problemDir := filepath.Join(tempRoot, "1234", "A")

	err = os.MkdirAll(problemDir, 0o755)
	if err != nil {
		t.Fatalf("Failed to create problem directory: %v", err)
	}

	logger := log.New(os.Stdout, "[TEST] ", log.LstdFlags)
	dm := NewDirectoryManager(tempRoot, logger)

	returnedPath, err := dm.ProblemDirExists(problem)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if returnedPath != problemDir {
		t.Errorf("Expected path %s, got %s", problemDir, returnedPath)
	}
}

func TestProblemDirExists_MissingProblemDir(t *testing.T) {
	tempRoot, err := os.MkdirTemp("", "dm_test_missing")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempRoot)

	contestCode := 5678
	problemCode := "B"
	problem := Problem{ContestCode: contestCode, ProblemCode: problemCode}

	// Only contest dir, no problem dir
	os.MkdirAll(filepath.Join(tempRoot, "5678"), 0o755)

	logger := log.New(os.Stdout, "[TEST] ", log.LstdFlags)
	dm := NewDirectoryManager(tempRoot, logger)

	_, err = dm.ProblemDirExists(problem)
	if err == nil {
		t.Fatal("Expected error for missing problem directory, got nil")
	}
	if !strings.Contains(err.Error(), "no such file or directory") {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestProblemDirExists_MissingContestDir(t *testing.T) {
	tempRoot, err := os.MkdirTemp("", "dm_test_missing_contest")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempRoot)

	problem := Problem{ContestCode: 9999, ProblemCode: "Z"}

	logger := log.New(os.Stdout, "[TEST] ", log.LstdFlags)
	dm := NewDirectoryManager(tempRoot, logger)

	_, err = dm.ProblemDirExists(problem)
	if err == nil {
		t.Fatal("Expected error for missing contest directory, got nil")
	}
	if !strings.Contains(err.Error(), "no such file or directory") {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestCreateProblemDir_Success(t *testing.T) {
	tempRoot, err := os.MkdirTemp("", "dm_create_success")
	if err != nil {
		t.Fatalf("Failed to create temp root: %v", err)
	}
	defer os.RemoveAll(tempRoot)

	problem := Problem{ContestCode: 1111, ProblemCode: "X"}

	logger := log.New(os.Stdout, "[TEST] ", log.LstdFlags)
	dm := NewDirectoryManager(tempRoot, logger)

	path, err := dm.CreateProblemDir(problem)
	if err != nil {
		t.Fatalf("CreateProblemDir failed: %v", err)
	}

	if fi, err := os.Stat(path); err != nil || !fi.IsDir() {
		t.Errorf("Expected directory at %s to exist, got error: %v", path, err)
	}
}

func TestCreateProblemDir_AlreadyExists(t *testing.T) {
	tempRoot, err := os.MkdirTemp("", "dm_create_exists")
	if err != nil {
		t.Fatalf("Failed to create temp root: %v", err)
	}
	defer os.RemoveAll(tempRoot)

	problem := Problem{ContestCode: 2222, ProblemCode: "Y"}
	problemDir := filepath.Join(tempRoot, "2222", "Y")
	os.MkdirAll(problemDir, 0o755)

	logger := log.New(os.Stdout, "[TEST] ", log.LstdFlags)
	dm := NewDirectoryManager(tempRoot, logger)

	path, err := dm.CreateProblemDir(problem)
	if err != nil {
		t.Fatalf("Expected no error when dir already exists, got: %v", err)
	}

	if path != problemDir {
		t.Errorf("Expected path %s, got %s", problemDir, path)
	}
}

func TestPopulate_CreatesTestCasesAndPreservesExistingProgram(t *testing.T) {
	tempRoot := createTempDir(t)
	defer os.RemoveAll(tempRoot)

	problem := Problem{ContestCode: 1001, ProblemCode: "A"}
	programFileName := "main.cpp"
	inputPrefix := "in"
	outputPrefix := "out"

	problemDir := filepath.Join(tempRoot, "1001", "A")
	err := os.MkdirAll(problemDir, 0o755)
	if err != nil {
		t.Fatalf("Failed to create problem directory: %v", err)
	}

	// Create a pre-existing program file
	programPath := filepath.Join(problemDir, programFileName)
	err = os.WriteFile(programPath, []byte("// existing content"), 0o644)
	if err != nil {
		t.Fatalf("Failed to create existing program file: %v", err)
	}

	testCases := []execution.TestCase{
		{Input: "1 2\n", Output: "3\n"},
		{Input: "4 5\n", Output: "9\n"},
	}

	logger := log.New(os.Stdout, "[TEST] ", log.LstdFlags)
	dm := NewDirectoryManager(tempRoot, logger)

	err = dm.Populate(problem, programFileName, testCases, inputPrefix, outputPrefix)
	if err != nil {
		t.Fatalf("Populate failed: %v", err)
	}

	// Check if test cases exist
	for i := range testCases {
		inputPath := filepath.Join(problemDir, inputPrefix+strconv.Itoa(i+1))
		outputPath := filepath.Join(problemDir, outputPrefix+strconv.Itoa(i+1))

		assertFileExists(t, inputPath)
		assertFileExists(t, outputPath)
	}

	// Ensure program file content is unchanged
	data, err := os.ReadFile(programPath)
	if err != nil {
		t.Fatalf("Failed to read program file: %v", err)
	}
	if !strings.Contains(string(data), "existing content") {
		t.Errorf("Expected program file to retain existing content, got: %s", string(data))
	}
}

func TestPopulate_CreatesProgramIfMissing(t *testing.T) {
	tempRoot := createTempDir(t)
	defer os.RemoveAll(tempRoot)

	problem := Problem{ContestCode: 1002, ProblemCode: "B"}
	programFileName := "main.cpp"
	inputPrefix := "input"
	outputPrefix := "output"

	testCases := []execution.TestCase{
		{Input: "6 1\n", Output: "7\n"},
	}

	logger := log.New(os.Stdout, "[TEST] ", log.LstdFlags)
	dm := NewDirectoryManager(tempRoot, logger)

	err := dm.Populate(problem, programFileName, testCases, inputPrefix, outputPrefix)
	if err != nil {
		t.Fatalf("Populate failed: %v", err)
	}

	problemDir := filepath.Join(tempRoot, "1002", "B")
	programPath := filepath.Join(problemDir, programFileName)
	assertFileExists(t, programPath)

	inputPath := filepath.Join(problemDir, "input1")
	outputPath := filepath.Join(problemDir, "output1")
	assertFileExists(t, inputPath)
	assertFileExists(t, outputPath)
}

func createTempDir(t *testing.T) string {
	t.Helper()
	tempRoot, err := os.MkdirTemp("", "dm_test")
	if err != nil {
		t.Fatalf("Failed to create temp root dir: %v", err)
	}
	return tempRoot
}

func assertFileExists(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf("Expected file to exist at %s, but it does not", path)
	}
}
