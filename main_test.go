package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateArguments(t *testing.T) {
	args := []string{"./count", "test-data", ".txt", ".veeti"}
	params, err := validateArguments(args)
	require.NoError(t, err)

	assert.Equal(t, params.Path, args[1])
	assert.ElementsMatch(t, params.FileTypes, args[2:])
}

func TestCountLines(t *testing.T) {
	params := &Params{
		Path:      "test-data",
		FileTypes: []string{".veeti", ".txt"},
	}

	count, err := CountLines(params)
	require.NoError(t, err)

	assert.Equal(t, 10, count)
}

func TestIsFilePath(t *testing.T) {
	testCases := []struct {
		name      string
		path      string
		fileTypes []string
		expected  bool
	}{
		{
			name:      "empty file type returns all files",
			path:      "test-data.veeti",
			fileTypes: []string{".moi", ".veeti"},
			expected:  true,
		},
		{
			name:      "false when not in input",
			path:      "test-data.veeti",
			fileTypes: []string{".moi", ".iteev"},
			expected:  false,
		},
		{
			name:      "true when nil",
			path:      "test-data.veeti",
			fileTypes: nil,
			expected:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ok := isFilePath(tc.path, tc.fileTypes)
			assert.Equal(t, tc.expected, ok)
		})
	}
}

func TestFilePaths(t *testing.T) {
	testCases := []struct {
		name      string
		dir       string
		fileTypes []string
		expected  []string
	}{
		{
			name:      "empty file type returns all files",
			dir:       "test-data",
			fileTypes: nil,
			expected:  []string{"test-data/test-file.veeti", "test-data/test-dir/test-file-2.txt"},
		},
		{
			name:      "specified file type only returns specific ones",
			dir:       "test-data",
			fileTypes: []string{".txt"},
			expected:  []string{"test-data/test-dir/test-file-2.txt"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			paths, err := filePaths(tc.dir, tc.fileTypes)
			require.NoError(t, err)
			assert.ElementsMatch(t, tc.expected, paths)
		})
	}
}

func TestCountPathLines(t *testing.T) {
	path := "test-data/test-file.veeti"
	count, err := countPathLines(path)
	require.NoError(t, err)

	assert.Equal(t, 7, count)
}

func TestCountFileLines(t *testing.T) {
	path := "test-data/test-file.veeti"
	file, err := os.Open(path)
	require.NoError(t, err)
	defer file.Close()

	count, err := countFileLines(file)
	require.NoError(t, err)

	assert.Equal(t, 7, count)
}
