package ccparser

import (
	"log"
	"os"
	"strings"
	"testing"
)

func TestParse_ValidContestProblem(t *testing.T) {
	logger := log.New(os.Stdout, "[TEST] ", log.LstdFlags)
	parser := NewParser(logger)

	ccp := &CCProblem{
		Name: "A. Sum of Two Numbers",
		URL:  "https://codeforces.com/contest/1234/problem/A",
		Tests: []struct {
			Input  string "json:\"input\""
			Output string "json:\"output\""
		}{
			{Input: "1 2\n", Output: "3\n"},
			{Input: "4 5\n", Output: "9\n"},
		},
	}

	problem, err := parser.Parse(ccp)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if problem.ContestCode != 1234 {
		t.Errorf("Expected ContestCode 1234, got %d", problem.ContestCode)
	}

	if problem.ProblemCode != "A_Sum_of_Two_Numbers" {
		t.Errorf("Expected ProblemCode A_Sum_of_Two_Numbers, got %s", problem.ProblemCode)
	}

	if len(problem.TestCases) != 2 {
		t.Errorf("Expected 2 test cases, got %d", len(problem.TestCases))
	}
}

func TestParse_ValidProblemsetProblem(t *testing.T) {
	logger := log.New(os.Stdout, "[TEST] ", log.LstdFlags)
	parser := NewParser(logger)

	ccp := &CCProblem{
		Name: "B. Multiply",
		URL:  "https://codeforces.com/problemset/problem/5678/B",
		Tests: []struct {
			Input  string "json:\"input\""
			Output string "json:\"output\""
		}{
			{Input: "3 4\n", Output: "12\n"},
		},
	}

	problem, err := parser.Parse(ccp)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if problem.ContestCode != 5678 {
		t.Errorf("Expected ContestCode 5678, got %d", problem.ContestCode)
	}

	if !strings.HasPrefix(problem.ProblemCode, "B_") {
		t.Errorf("Expected ProblemCode to start with B_, got %s", problem.ProblemCode)
	}
}

func TestParse_InvalidURL(t *testing.T) {
	logger := log.New(os.Stdout, "[TEST] ", log.LstdFlags)
	parser := NewParser(logger)

	ccp := &CCProblem{
		Name: "A. Bad URL",
		URL:  "://bad_url", // malformed URL
	}

	_, err := parser.Parse(ccp)
	if err == nil {
		t.Fatal("Expected error for malformed URL, got nil")
	}
}

func TestParse_UnsupportedPath(t *testing.T) {
	logger := log.New(os.Stdout, "[TEST] ", log.LstdFlags)
	parser := NewParser(logger)

	ccp := &CCProblem{
		Name: "C. Weird Path",
		URL:  "https://codeforces.com/random/123/C",
	}

	_, err := parser.Parse(ccp)
	if err == nil {
		t.Fatal("Expected error for unsupported path, got nil")
	}
}

