package utils

import (
	"path/filepath"
	"strings"

	"github.com/mitchellh/go-homedir"
)

// NormalizePath the path
func NormalizePath(path string) string {
	if strings.HasPrefix(path, "~") {
		path, _ = homedir.Expand(filepath.FromSlash(path))
	}
	return path
}
