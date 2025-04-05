package scraper

import (
	"log"
	"os"
	"strings"
	"testing"
)

func setupTestScraper() *Scraper {
	logger := log.New(os.Stdout, "TEST: ", log.LstdFlags)
	baseURL := "https://example.com"
	return NewScraper(baseURL, logger)
}

func TestExtractDetailsContestPath(t *testing.T) {
	s := setupTestScraper()
	problemRoute := "contest/1234/problems/ABC"

	problemCode, contestCode, err := s.extractDetails(problemRoute)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if contestCode != 1234 {
		t.Errorf("expected contest code 1234, got %d", contestCode)
	}

	if problemCode != "ABC" {
		t.Errorf("expected problem code ABC, got %s", problemCode)
	}
}

func TestExtractDetailsProblemsetPath(t *testing.T) {
	s := setupTestScraper()
	problemRoute := "problemset/problem/5678/XYZ"

	problemCode, contestCode, err := s.extractDetails(problemRoute)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if contestCode != 5678 {
		t.Errorf("expected contest code 5678, got %d", contestCode)
	}

	if problemCode != "XYZ" {
		t.Errorf("expected problem code XYZ, got %s", problemCode)
	}
}

func TestScrapePage(t *testing.T) {
	s := setupTestScraper()

	testURL := "https://example.com/contest/4321/problems/DEF"

	data, err := s.ScrapePage(testURL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if data.ContestCode != 4321 {
		t.Errorf("expected contest code 4321, got %d", data.ContestCode)
	}

	if data.ProblemCode != "DEF" {
		t.Errorf("expected problem code DEF, got %s", data.ProblemCode)
	}
}

func TestInvalidPath(t *testing.T) {
	s := setupTestScraper()
	_, _, err := s.extractDetails("unknown/1234/something")

	if err == nil || !strings.Contains(err.Error(), "path not supported") {
		t.Errorf("expected path not supported error, got %v", err)
	}
}

