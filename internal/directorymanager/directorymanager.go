package directorymanager

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/PriyanshuSharma23/codeforces-cli/internal/execution"
)

type DirectoryManager struct {
	logger   *log.Logger
	rootPath string
}

func NewDirectoryManager(dirpath string, logger *log.Logger) *DirectoryManager {
	dirpath = filepath.Clean(dirpath)
	return &DirectoryManager{
		logger,
		dirpath,
	}
}

type Problem struct {
	ContestCode int
	ProblemCode string
}

func (p Problem) RelativeDir() string {
	return filepath.Join(strconv.Itoa(p.ContestCode), p.ProblemCode)
}

func (d *DirectoryManager) ProblemDirExists(p Problem) (string, error) {
	relativeDir := p.RelativeDir()
	d.logger.Printf("looking for directory: %s", relativeDir)

	fullDirPath := filepath.Join(d.rootPath, relativeDir)
	if stat, err := os.Stat(fullDirPath); err != nil {
		if os.IsNotExist(err) {
			d.logger.Printf("directory does not exist: %s", fullDirPath)
			return "", err
		}
		d.logger.Printf("error checking directory: %v", err)
		return "", err
	} else if !stat.IsDir() {
		d.logger.Printf("path exists but is not a directory: %s", fullDirPath)
		return "", fmt.Errorf("path exists but is not a directory: %s", fullDirPath)
	}
	return fullDirPath, nil
}

func (d *DirectoryManager) CreateProblemDir(p Problem) (string, error) {
	relativeDir := p.RelativeDir()
	d.logger.Printf("creating directory at %s", relativeDir)

	fullDirPath := filepath.Join(d.rootPath, relativeDir)
	if err := os.MkdirAll(fullDirPath, 0o755); err != nil {
		d.logger.Printf("failed to create directory: %v", err)
		return "", fmt.Errorf("could not create directory %s: %w", fullDirPath, err)
	}

	if stat, err := os.Stat(fullDirPath); err != nil {
		d.logger.Printf("error checking created directory: %v", err)
		return "", err
	} else if !stat.IsDir() {
		d.logger.Printf("created path is not a directory: %s", fullDirPath)
		return "", fmt.Errorf("created path is not a directory: %s", fullDirPath)
	}

	d.logger.Printf("directory successfully created: %s", fullDirPath)
	return fullDirPath, nil
}

func (d *DirectoryManager) Populate(p Problem, programFileName string, testCases []execution.TestCase, inputPrefix, outputPrefix string) error {
	problemDir := filepath.Join(d.rootPath, strconv.Itoa(p.ContestCode), p.ProblemCode)

	// Ensure problem directory exists
	if err := os.MkdirAll(problemDir, 0o755); err != nil {
		return fmt.Errorf("failed to create problem directory: %w", err)
	}

	// Write test cases with provided prefixes
	for i, tc := range testCases {
		inputFile := filepath.Join(problemDir, fmt.Sprintf("%s%d", inputPrefix, i+1))
		outputFile := filepath.Join(problemDir, fmt.Sprintf("%s%d", outputPrefix, i+1))

		if err := os.WriteFile(inputFile, []byte(tc.Input), 0o644); err != nil {
			return fmt.Errorf("failed to write input file %s: %w", inputFile, err)
		}

		if err := os.WriteFile(outputFile, []byte(tc.Output), 0o644); err != nil {
			return fmt.Errorf("failed to write output file %s: %w", outputFile, err)
		}
	}

	// Create the program file if it doesn't exist
	programPath := filepath.Join(problemDir, programFileName)
	if _, err := os.Stat(programPath); os.IsNotExist(err) {
		file, err := os.Create(programPath)
		if err != nil {
			return fmt.Errorf("failed to create program file %s: %w", programPath, err)
		}
		file.Close()
	}

	return nil
}
