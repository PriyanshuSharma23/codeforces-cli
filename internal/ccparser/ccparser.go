package ccparser

import (
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"

	"github.com/PriyanshuSharma23/codeforces-cli/internal/execution"
)

type Problem struct {
	ContestCode  int
	ProblemCode  string
	TestCases    []execution.TestCase
	URL          string
	OriginalName string
}

type CCProblem struct {
	Name        string `json:"name"`
	Group       string `json:"group"`
	URL         string `json:"url"`
	Interactive bool   `json:"interactive"`
	MemoryLimit int    `json:"memoryLimit"`
	TimeLimit   int    `json:"timeLimit"`

	Tests []struct {
		Input  string `json:"input"`
		Output string `json:"output"`
	} `json:"tests"`

	TestType string `json:"testType"`

	Input struct {
		Type string `json:"type"`
	} `json:"input"`

	Output struct {
		Type string `json:"type"`
	} `json:"output"`

	Languages map[string]struct {
		MainClass string `json:"mainClass"`
		TaskClass string `json:"taskClass"`
	} `json:"languages"`

	Batch struct {
		ID   string `json:"id"`
		Size int    `json:"size"`
	} `json:"batch"`
}

type Parser struct {
	logger *log.Logger
}

func NewParser(logger *log.Logger) *Parser {
	return &Parser{logger}
}

func (s *Parser) Parse(ccp *CCProblem) (*Problem, error) {
	probURL, err := url.Parse(ccp.URL)
	if err != nil {
		s.logger.Printf("ERROR: failed to parse the problem url [%s]: %s\n", ccp.URL, err)
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	problemIndex, contestCode, err := s.extractDetails(probURL.Path[1:]) // remove the leading '/'
	if err != nil {
		s.logger.Printf("ERROR: failed to extract details from URL path [%s]: %s\n", probURL.Path, err)
		return nil, err
	}

	originalName := ccp.Name
	problemCode := normalizeProblemCode(problemIndex, originalName)

	testCases := make([]execution.TestCase, len(ccp.Tests))
	for i := range ccp.Tests {
		testCases[i] = execution.TestCase{
			Input:  strings.TrimSpace(ccp.Tests[i].Input),
			Output: strings.TrimSpace(ccp.Tests[i].Output),
		}
	}

	problem := Problem{
		ContestCode:  contestCode,
		ProblemCode:  problemCode,
		TestCases:    testCases,
		URL:          ccp.URL,
		OriginalName: originalName,
	}

	return &problem, nil
}

func (s *Parser) extractDetails(problemRoute string) (string, int, error) {
	parts := strings.Split(problemRoute, "/")

	if len(parts) < 4 {
		return "", 0, fmt.Errorf("invalid path format: %s", problemRoute)
	}

	var contestCode int
	var problemCode string

	switch parts[0] {
	case "contest":
		ccInt, err := strconv.Atoi(parts[1])
		if err != nil {
			return "", 0, fmt.Errorf("invalid contest code: %s", parts[1])
		}
		contestCode = ccInt
		problemCode = parts[3]

	case "problemset":
		ccInt, err := strconv.Atoi(parts[2])
		if err != nil {
			return "", 0, fmt.Errorf("invalid contest code: %s", parts[2])
		}
		contestCode = ccInt
		problemCode = parts[3]

	default:
		return "", 0, fmt.Errorf("path not supported for the url: %s", problemRoute)
	}

	return problemCode, contestCode, nil
}

func normalizeProblemCode(index string, name string) string {
	// Remove index prefix like "A. " or "B. "
	if idx := strings.Index(name, "."); idx != -1 {
		name = name[idx+1:]
	}
	// Slugify the name
	nameSlug := strings.Join(strings.Fields(name), "_")
	return fmt.Sprintf("%s_%s", index, nameSlug)
}

