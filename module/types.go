// Package module - module/types.go
package module

// Flags struct - represents flags for the application.
type Flags struct {
	ShowHidden       bool
	ShowTreeView     bool
	HideIcon         bool
	ShowReadableSize bool
	ShowVersion      bool
	ShowDirsOnly     bool
	ShowJSON         bool
	ShowLongFormat   bool
	HideSummary      bool
	RespectGitIgnore bool
	SortBy           string
	ReverseSort      bool
	DirsFirst        bool
	TreeDepth        int
	LimitTreeDepth   bool
}

// Options struct - Contains configuration options for directory listing.
type Options struct {
	Directory string
	Flags     Flags
}
