package execution

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

type Engine struct {
	root             string
	testCasesDir     string
	buildCommand     string // gcc -o main.exe main.cpp
	executionCommand string // ./main.exe
	inputPrefix      string
	outputPrefix     string
	logger           *log.Logger
}

func NewEngine(
	root,
	testCasesDir,
	buildCommand,
	executionCommand,
	inputPrefix,
	outputPrefix string,
	logger *log.Logger,
) *Engine {
	return &Engine{
		root:             root,
		testCasesDir:     testCasesDir,
		buildCommand:     buildCommand,
		executionCommand: executionCommand,
		inputPrefix:      inputPrefix,
		outputPrefix:     outputPrefix,
		logger:           logger,
	}
}

type Result struct {
	Ok             bool
	TestCase       int
	ExpectedOutput string
	ProgramOutput  string
}

type TestCase struct {
	Input  string
	Output string
}

func (e *Engine) Execute() ([]Result, error) {
	if e.buildCommand != "" {
		err := e.build()
		if err != nil {
			return nil, err
		}
	}

	testCases, err := e.readTestCases()
	if err != nil {
		return nil, err
	}

	results := make([]Result, 0, len(testCases))

	for k, v := range testCases {
		programOutput, ok, err := e.runTestCase(k, v)
		if err != nil {
			return nil, err
		}

		result := Result{
			Ok:             ok,
			TestCase:       k,
			ExpectedOutput: v.Output,
			ProgramOutput:  programOutput,
		}

		results = append(results, result)
	}

	return results, nil
}

func (e *Engine) build() error {
	e.logger.Println("Building program...")
	args := strings.Split(e.buildCommand, " ")

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = e.root
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("Build failed: %v\n", err)
		return err
	}

	return nil
}

func (e *Engine) readTestCases() (map[int]TestCase, error) {
	testCases := make(map[int]TestCase)

	entries, err := os.ReadDir(e.testCasesDir)
	if err != nil {
		e.logger.Printf("failed to read the testcases dir: %s\n", err)
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			e.logger.Printf("WARN: only files allowed in the test directory: %s", entry.Name())
			continue
		}

		fileName := entry.Name()

		var testCaseNumStr string
		var isInput bool

		if strings.HasPrefix(fileName, e.inputPrefix) {
			isInput = true
		} else if strings.HasPrefix(fileName, e.outputPrefix) {
			isInput = false
		} else {
			e.logger.Printf("WARN: invalid entry: %s", entry.Name())
			continue
		}

		if isInput {
			testCaseNumStr = fileName[len(e.inputPrefix):]
		} else {
			testCaseNumStr = fileName[len(e.outputPrefix):]
		}

		testCaseNum, err := strconv.Atoi(testCaseNumStr)
		if err != nil {
			e.logger.Printf("WARN: invalid trailing test case number: %s", fileName)
			continue
		}

		filePath := filepath.Join(e.testCasesDir, fileName)
		content, err := os.ReadFile(filePath)
		if err != nil {
			e.logger.Printf("WARN: failed to read file contents for: %s", fileName)
			continue
		}

		testCase := testCases[testCaseNum]

		if isInput {
			testCase.Input = string(content)
		} else {
			testCase.Output = string(content)
		}

		testCases[testCaseNum] = testCase
	}

	return testCases, nil
}

func (e *Engine) runTestCase(testNum int, t TestCase) (string, bool, error) {
	args := strings.Split(e.executionCommand, " ")

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = e.root

	inputReader := strings.NewReader(t.Input)
	cmd.Stdin = inputReader

	var out bytes.Buffer
	cmd.Stdout = &out

	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		e.logger.Printf("ERROR: Failed to execute test %d\n", testNum)
		return "", false, err
	}

	outBytes, err := io.ReadAll(&out)
	if err != nil {
		e.logger.Printf("ERROR: Failed to read output for the test: %d. %s\n", testNum, err)
		return "", false, err
	}

	outStr := string(outBytes)

	if outStr != t.Output {
		return outStr, false, nil
	}

	return outStr, true, nil
}
