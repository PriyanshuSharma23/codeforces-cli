package directorymanager

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

type DirectoryManager struct {
	logger           *log.Logger
	rootPath         string
	problemExtension string
}

func NewDirectoryManager(dirpath, problemExtension string, logger *log.Logger) *DirectoryManager {
	dirpath = filepath.Clean(dirpath)
	return &DirectoryManager{
		logger,
		dirpath,
		problemExtension,
	}
}

type Problem struct {
	ContestCode int
	ProblemCode string
}

func (p Problem) RelativePath(ext string) string {
	return filepath.Join(strconv.Itoa(p.ContestCode), fmt.Sprintf("%s.%s", p.ProblemCode, ext))
}

func (d *DirectoryManager) ProblemFileExists(p Problem) (string, error) {
	d.logger.Printf("looking for %s", p.RelativePath(d.problemExtension))

	fullDirPath := filepath.Join(d.rootPath, strconv.Itoa(p.ContestCode))
	fullPath := filepath.Join(d.rootPath, p.RelativePath(d.problemExtension))

	if _, err := os.Stat(fullDirPath); os.IsNotExist(err) {
		d.logger.Printf("directory for contest does not exist: %d", p.ContestCode)
		return "", err
	}

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		d.logger.Printf("file for problem code does not exist: %s", p.ProblemCode)
		return "", err
	}

	file, err := os.Open(fullPath)
	if err != nil {
		d.logger.Printf("failed to open problem file: %s", err)
		return "", err
	}
	defer file.Close()

	return fullPath, nil
}

func (d *DirectoryManager) CreateProblemFile(p Problem) (string, error) {
	relativePath := p.RelativePath(d.problemExtension)
	d.logger.Printf("creating file at %s", relativePath)

	fullDirPath := filepath.Join(d.rootPath, strconv.Itoa(p.ContestCode))
	fullPath := filepath.Join(d.rootPath, relativePath)

	if err := os.MkdirAll(fullDirPath, 0o755); err != nil {
		d.logger.Printf("failed to create contest directory: %v", err)
		return "", fmt.Errorf("could not create contest directory %s: %w", fullDirPath, err)
	}

	if _, err := os.Stat(fullPath); err == nil {
		d.logger.Printf("file already exists: %s", fullPath)
		return "", fmt.Errorf("file already exists: %s", fullPath)
	} else if !os.IsNotExist(err) {
		d.logger.Printf("error checking file existence: %v", err)
		return "", err
	}

	file, err := os.Create(fullPath)
	if err != nil {
		d.logger.Printf("failed to create file: %v", err)
		return "", fmt.Errorf("could not create file %s: %w", fullPath, err)
	}
	defer file.Close()

	d.logger.Printf("file successfully created: %s", fullPath)
	return fullPath, nil
}
