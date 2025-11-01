package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	args, err := validateArguments(os.Args)
	if err != nil {
		log.Fatalf("invalid arguments: %v", err)
	}

	sum, err := CountLines(args)
	if err != nil {
		log.Fatalf("failed to count lines: %v", err)
	}

	log.Printf("%d lines", sum)
}

func validateArguments(args []string) (*Params, error) {
	if len(args) < 2 {
		return nil, errors.New("there should be an argument for the path dummy")
	}

	path := args[1]
	_, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("invalid path: %w", err)
	}

	var fileTypes []string
	for i := 2; i < len(args); i++ {
		fileTypes = append(fileTypes, args[i])
	}

	cmdArgs := &Params{
		Path:      path,
		FileTypes: fileTypes,
	}

	return cmdArgs, nil
}

type Params struct {
	Path      string
	FileTypes []string
}

func (p *Params) ToString() string {
	fileTypes := strings.Join(p.FileTypes, ", ")
	return fmt.Sprintf("arguments - path: %s file types: %s", p.Path, fileTypes)
}

func CountLines(params *Params) (int, error) {
	filePaths, err := filePaths(params.Path, params.FileTypes)
	if err != nil {
		return -1, fmt.Errorf("file paths %v", err)
	}

	if len(filePaths) == 0 {
		return -1, fmt.Errorf("no files found for input %s ", params.ToString())
	}

	sum := 0
	for _, filePath := range filePaths {
		count, err := countPathLines(filePath)
		if err != nil {
			return -1, err
		}

		sum += count
	}

	return sum, nil
}

func isFilePath(path string, fileTypes []string) bool {
	// No filetypes default everything to valid
	if fileTypes == nil {
		return true
	}

	for _, fileType := range fileTypes {
		if strings.HasSuffix(path, fileType) {
			return true
		}
	}

	return false
}

func filePaths(root string, fileTypes []string) ([]string, error) {
	var paths []string
	err := filepath.WalkDir(root, func(path string, dir fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !dir.IsDir() {
			if isFilePath(path, fileTypes) {
				paths = append(paths, path)
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return paths, err
}

func countPathLines(path string) (int, error) {
	file, err := os.Open(path)
	if err != nil {
		return -1, err
	}
	defer file.Close()

	return countFileLines(file)
}

func countFileLines(file *os.File) (int, error) {
	scanner := bufio.NewScanner(file)

	count := 0
	for scanner.Scan() {
		count += 1
	}

	err := scanner.Err()
	if err != nil {
		return -1, err
	}

	return count, nil
}
