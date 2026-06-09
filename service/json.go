package service

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// jsonOutput represents the full machine-readable output.
type jsonOutput struct {
	Root    string      `json:"root"`
	Tree    bool        `json:"tree"`
	Entries []jsonEntry `json:"entries"`
	Summary jsonSummary `json:"summary"`
}

// jsonEntry represents one file-system entry in JSON output.
type jsonEntry struct {
	Name     string      `json:"name"`
	Path     string      `json:"path"`
	Type     string      `json:"type"`
	Size     int64       `json:"size"`
	Children []jsonEntry `json:"children,omitempty"`
}

// jsonSummary represents the filtered file and directory counts.
type jsonSummary struct {
	Directories int `json:"directories"`
	Files       int `json:"files"`
}

// printJSONDirectory prints directory contents as JSON.
func printJSONDirectory(options directoryOptions, output io.Writer) error {
	var (
		entries []jsonEntry
		summary jsonSummary
		err     error
	)

	if options.Flags.ShowTreeView {
		entries, summary, err = collectJSONTreeEntries(options, 0)
	} else {
		entries, summary, err = collectJSONFlatEntries(options)
	}
	if err != nil {
		return err
	}

	result := jsonOutput{
		Root:    options.Directory,
		Tree:    options.Flags.ShowTreeView,
		Entries: entries,
		Summary: summary,
	}

	encoder := json.NewEncoder(output)
	encoder.SetIndent("", "  ")
	return encoder.Encode(result)
}

// collectJSONFlatEntries returns visible direct children for flat JSON output.
func collectJSONFlatEntries(options directoryOptions) (entries []jsonEntry, summary jsonSummary, err error) {
	files, d, err := readDirectory(options.Directory)
	if err != nil {
		return nil, jsonSummary{}, fmt.Errorf("read directory %q: %w", options.Directory, err)
	}
	defer func() {
		closeErr := closeDirectory(d)
		if err == nil && closeErr != nil {
			err = closeErr
		}
	}()

	sortSlice(files)
	for _, file := range files {
		if !shouldShowFile(file, options) {
			continue
		}

		entries = append(entries, newJSONEntry(options.Directory, file))
		addJSONSummaryCount(&summary, file)
	}

	return entries, summary, nil
}

// collectJSONTreeEntries returns visible children recursively for tree JSON output.
func collectJSONTreeEntries(options directoryOptions, depth int) (entries []jsonEntry, summary jsonSummary, err error) {
	if reachedMaxDepth(options.Flags, depth) {
		return nil, jsonSummary{}, nil
	}

	files, d, err := readDirectory(options.Directory)
	if err != nil {
		return nil, jsonSummary{}, fmt.Errorf("read directory %q: %w", options.Directory, err)
	}
	defer func() {
		closeErr := closeDirectory(d)
		if err == nil && closeErr != nil {
			err = closeErr
		}
	}()

	sortSlice(files)
	for _, file := range files {
		if !shouldShowFile(file, options) {
			continue
		}

		entry := newJSONEntry(options.Directory, file)
		addJSONSummaryCount(&summary, file)

		if file.IsDir() {
			childOptions := options
			childOptions.Directory = filepath.Join(options.Directory, file.Name())

			children, childSummary, err := collectJSONTreeEntries(childOptions, depth+1)
			if err != nil {
				return nil, jsonSummary{}, err
			}
			entry.Children = children
			summary.Directories += childSummary.Directories
			summary.Files += childSummary.Files
		}

		entries = append(entries, entry)
	}

	return entries, summary, nil
}

// newJSONEntry creates a JSON entry from file metadata.
func newJSONEntry(parent string, file os.FileInfo) jsonEntry {
	entryType := "file"
	if file.IsDir() {
		entryType = "directory"
	}

	return jsonEntry{
		Name: file.Name(),
		Path: filepath.Join(parent, file.Name()),
		Type: entryType,
		Size: file.Size(),
	}
}

// addJSONSummaryCount increments summary counts for a file-system entry.
func addJSONSummaryCount(summary *jsonSummary, file os.FileInfo) {
	if file.IsDir() {
		summary.Directories++
		return
	}
	summary.Files++
}
