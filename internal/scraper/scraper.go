package scraper

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

type Scraper struct {
	baseURL string
	logger  *log.Logger
}

func NewScraper(baseURL string, logger *log.Logger) *Scraper {
	if baseURL[len(baseURL)-1] == '/' {
		baseURL = baseURL[:len(baseURL)-1]
	}
	return &Scraper{
		baseURL,
		logger,
	}
}

type ScrapedData struct {
	ProblemCode string
	ContestCode int
}

func (s *Scraper) ScrapePage(problemURL string) (*ScrapedData, error) {
	if len(problemURL)+1 <= len(s.baseURL) {
		return nil, fmt.Errorf("bad problemURL: %s", problemURL)
	}

	problemRoute := problemURL[len(s.baseURL)+1:]

	problemCode, contestCode, err := s.extractDetails(problemRoute)
	if err != nil {
		s.logger.Printf("failed to load page data: %s\n", err.Error())
		return nil, fmt.Errorf("failed to load page data for: %s", problemURL)
	}

	return &ScrapedData{
		ProblemCode: problemCode,
		ContestCode: contestCode,
	}, nil
}

func (s *Scraper) extractDetails(problemRoute string) (string, int, error) {
	parts := strings.Split(problemRoute, "/")

	var contestCode int
	var problemCode string

	switch parts[0] {
	case "contest":
		if len(parts) < 4 {
			return "", 0, fmt.Errorf("invalid contest path: %s", problemRoute)
		}

		ccInt, err := strconv.Atoi(parts[1])
		if err != nil {
			return "", 0, fmt.Errorf("invalid contest code: %s", parts[1])
		}

		contestCode = ccInt
		problemCode = parts[3]
	case "problemset":
		if len(parts) < 4 {
			return "", 0, fmt.Errorf("invalid problemset path: %s", problemRoute)
		}
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
