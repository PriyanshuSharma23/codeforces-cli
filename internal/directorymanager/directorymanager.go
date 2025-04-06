package directorymanager

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/PriyanshuSharma23/codeforces-cli/internal/execution"
)

type DirectoryManager struct {
	logger   *log.Logger
	rootPath string
}

type Problem struct {
	ContestCode int
	ProblemCode string
}

func (p Problem) RelativeDir() string {
	return filepath.Join(strconv.Itoa(p.ContestCode), p.ProblemCode)
}

func NewDirectoryManager(root string, logger *log.Logger) *DirectoryManager {
	return &DirectoryManager{
		logger:   logger,
		rootPath: filepath.Clean(root),
	}
}

func (d *DirectoryManager) FullProblemPath(p Problem) string {
	return filepath.Join(d.rootPath, p.RelativeDir())
}

func (d *DirectoryManager) EnsureDir(p Problem) (string, error) {
	dir := d.FullProblemPath(p)
	d.logger.Printf("Ensuring directory: %s", dir)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}
	return dir, nil
}

func (d *DirectoryManager) WriteTestCases(p Problem, testCases []execution.TestCase, inputPrefix, outputPrefix string) error {
	dir := d.FullProblemPath(p)
	for i, tc := range testCases {
		inFile := filepath.Join(dir, fmt.Sprintf("%s%d", inputPrefix, i+1))
		outFile := filepath.Join(dir, fmt.Sprintf("%s%d", outputPrefix, i+1))

		if err := os.WriteFile(inFile, []byte(tc.Input), 0o644); err != nil {
			return fmt.Errorf("writing input file: %w", err)
		}
		if err := os.WriteFile(outFile, []byte(tc.Output), 0o644); err != nil {
			return fmt.Errorf("writing output file: %w", err)
		}
	}
	return nil
}

func (d *DirectoryManager) WriteMetadata(p Problem, metadata any) error {
	metaFile := filepath.Join(d.FullProblemPath(p), "problem.json")
	file, err := os.Create(metaFile)
	if err != nil {
		return fmt.Errorf("creating metadata file: %w", err)
	}
	defer file.Close()

	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")
	return enc.Encode(metadata)
}

func (d *DirectoryManager) WriteProgramFile(p Problem, filename, templateContent string) error {
	path := filepath.Join(d.FullProblemPath(p), filename)
	if _, err := os.Stat(path); err == nil {
		d.logger.Printf("Program file already exists: %s", path)
		return nil
	}
	return os.WriteFile(path, []byte(templateContent), 0o644)
}

func (d *DirectoryManager) LoadTemplate(templatePath string) (string, error) {
	// If the template path is empty, return an error
	if templatePath == "" {
		return "", fmt.Errorf("template path is empty")
	}

	// Read file contents
	content, err := os.ReadFile(templatePath)
	if err != nil {
		return "", fmt.Errorf("failed to read template file: %w", err)
	}

	return string(content), nil
}

func (d *DirectoryManager) GetCurrentProblemKey() (Problem, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		d.logger.Printf("ERROR: failed to get the current working dir: %s\n", err)
		return Problem{
			ContestCode: 0,
			ProblemCode: "",
		}, err // Return empty problem on error
	}

	parts := strings.Split(currentDir, string(filepath.Separator))

	if len(parts) < 2 {
		d.logger.Println("current working dir path segments less than two")
		return Problem{
			ContestCode: 0,
			ProblemCode: "",
		}, fmt.Errorf("current working dir path segments less than two") // Return empty problem if path is too short
	}

	problemCode := parts[len(parts)-1]
	contestCodeStr := parts[len(parts)-2]

	contestCode, err := strconv.Atoi(contestCodeStr)
	if err != nil {
		d.logger.Printf("failed to convert contest code %s to int: %s\n", contestCodeStr, err)
		return Problem{
			ContestCode: 0,
			ProblemCode: problemCode, // return the problem code even if contest code is invalid.
		}, err // Return problem code with contest code 0 if conversion fails
	}

	return Problem{
		ContestCode: contestCode,
		ProblemCode: problemCode,
	}, nil
}
