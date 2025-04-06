package execution

import (
	"log"
	"os"
	"path/filepath"
	"testing"
)

func TestExecutionEngine_RunCPP(t *testing.T) {
	tmpDir := t.TempDir()

	// Write main.cpp
	code := `#include <iostream>
using namespace std;
int main() {
    int a, b;
    cin >> a >> b;
    cout << a + b << endl;
    return 0;
}`
	err := os.WriteFile(filepath.Join(tmpDir, "main.cpp"), []byte(code), 0o644)
	if err != nil {
		t.Fatalf("failed to write source file: %v", err)
	}

	// Setup testcases
	testcasesDir := filepath.Join(tmpDir, "testcases")
	err = os.Mkdir(testcasesDir, 0o755)
	if err != nil {
		t.Fatalf("failed to create testcases dir: %v", err)
	}

	// in0 / out0
	os.WriteFile(filepath.Join(testcasesDir, "in0"), []byte("3 5\n"), 0o644)
	os.WriteFile(filepath.Join(testcasesDir, "out0"), []byte("8\n"), 0o644)

	engine := NewEngine(
		tmpDir,                 // rootDir
		testcasesDir,           // testCasesDir
		"g++ main.cpp -o main", // build command
		"./main",               // execution command
		"in",                   // input prefix
		"out",                  // output prefix
		log.New(os.Stdout, "TEST: ", log.LstdFlags),
	)

	results, err := engine.Execute()
	if err != nil {
		t.Fatalf("execution failed: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}

	if !results[0].Ok {
		t.Errorf("expected test to pass, but it failed. Output: %s", results[0].ProgramOutput)
	}
}
